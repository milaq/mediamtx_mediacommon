package mpegts

import (
	"context"
	"fmt"
	"io"

	"github.com/asticode/go-astits"

	"github.com/bluenviron/mediacommon/pkg/codecs/h264"
	"github.com/bluenviron/mediacommon/pkg/codecs/mpeg1audio"
	"github.com/bluenviron/mediacommon/pkg/codecs/mpeg4audio"
)

const (
	streamIDVideo = 224
	streamIDAudio = 192
)

func leadingTrack(tracks []*Track) *Track {
	for _, track := range tracks {
		if track.Codec.IsVideo() {
			return track
		}
	}
	return tracks[0]
}

// Writer is a MPEG-TS writer.
type Writer struct {
	nextPID uint16
	mux     *astits.Muxer
}

// NewWriter allocates a Writer.
func NewWriter(
	bw io.Writer,
	tracks []*Track,
) *Writer {
	w := &Writer{
		nextPID: 256,
	}

	w.mux = astits.NewMuxer(
		context.Background(),
		bw)

	for _, track := range tracks {
		if track.PID == 0 {
			track.PID = w.nextPID
			w.nextPID++
		}
		es, _ := track.marshal()

		err := w.mux.AddElementaryStream(*es)
		if err != nil {
			panic(err) // TODO: return error instead of panicking
		}
	}

	w.mux.SetPCRPID(leadingTrack(tracks).PID)

	// WriteTables() is not necessary
	// since it's called automatically when WriteData() is called with
	// * PID == PCRPID
	// * AdaptationField != nil
	// * RandomAccessIndicator = true

	return w
}

// WriteH26x writes a H26x access unit.
func (w *Writer) WriteH26x(
	track *Track,
	pts int64,
	dts int64,
	idrPresent bool,
	au [][]byte,
) error {
	enc, err := h264.AnnexBMarshal(au)
	if err != nil {
		return err
	}

	var af *astits.PacketAdaptationField

	if idrPresent {
		af = &astits.PacketAdaptationField{}
		af.RandomAccessIndicator = true
	}

	oh := &astits.PESOptionalHeader{
		MarkerBits: 2,
	}

	if dts == pts {
		oh.PTSDTSIndicator = astits.PTSDTSIndicatorOnlyPTS
		oh.PTS = &astits.ClockReference{Base: pts}
	} else {
		oh.PTSDTSIndicator = astits.PTSDTSIndicatorBothPresent
		oh.DTS = &astits.ClockReference{Base: dts}
		oh.PTS = &astits.ClockReference{Base: pts}
	}

	_, err = w.mux.WriteData(&astits.MuxerData{
		PID:             track.PID,
		AdaptationField: af,
		PES: &astits.PESData{
			Header: &astits.PESHeader{
				OptionalHeader: oh,
				StreamID:       streamIDVideo,
			},
			Data: enc,
		},
	})
	return err
}

// WriteOpus writes Opus packets.
func (w *Writer) WriteOpus(
	track *Track,
	pts int64,
	packets [][]byte,
) error {
	n := 0
	for _, packet := range packets {
		au := opusAccessUnit{
			ControlHeader: opusControlHeader{
				PayloadSize: len(packet),
			},
			Packet: packet,
		}
		n += au.marshalSize()
	}

	enc := make([]byte, n)
	n = 0
	for _, packet := range packets {
		au := opusAccessUnit{
			ControlHeader: opusControlHeader{
				PayloadSize: len(packet),
			},
			Packet: packet,
		}
		mn, err := au.marshalTo(enc[n:])
		if err != nil {
			return err
		}
		n += mn
	}

	_, err := w.mux.WriteData(&astits.MuxerData{
		PID: track.PID,
		AdaptationField: &astits.PacketAdaptationField{
			RandomAccessIndicator: true,
		},
		PES: &astits.PESData{
			Header: &astits.PESHeader{
				OptionalHeader: &astits.PESOptionalHeader{
					MarkerBits:      2,
					PTSDTSIndicator: astits.PTSDTSIndicatorOnlyPTS,
					PTS:             &astits.ClockReference{Base: pts},
				},
				StreamID: streamIDAudio,
			},
			Data: enc,
		},
	})
	return err
}

// WriteMPEG4Audio writes MPEG-4 Audio access units.
func (w *Writer) WriteMPEG4Audio(
	track *Track,
	pts int64,
	aus [][]byte,
) error {
	aacCodec := track.Codec.(*CodecMPEG4Audio)

	pkts := make(mpeg4audio.ADTSPackets, len(aus))

	for i, au := range aus {
		pkts[i] = &mpeg4audio.ADTSPacket{
			Type:         aacCodec.Config.Type,
			SampleRate:   aacCodec.SampleRate,
			ChannelCount: aacCodec.Config.ChannelCount,
			AU:           au,
		}
	}

	enc, err := pkts.Marshal()
	if err != nil {
		return err
	}

	_, err = w.mux.WriteData(&astits.MuxerData{
		PID: track.PID,
		AdaptationField: &astits.PacketAdaptationField{
			RandomAccessIndicator: true,
		},
		PES: &astits.PESData{
			Header: &astits.PESHeader{
				OptionalHeader: &astits.PESOptionalHeader{
					MarkerBits:      2,
					PTSDTSIndicator: astits.PTSDTSIndicatorOnlyPTS,
					PTS:             &astits.ClockReference{Base: pts},
				},
				StreamID: streamIDAudio,
			},
			Data: enc,
		},
	})
	return err
}

// WriteMPEG1Audio writes MPEG-1 Audio packets.
func (w *Writer) WriteMPEG1Audio(
	track *Track,
	pts int64,
	frames [][]byte,
) error {
	if !track.mp3Checked {
		var h mpeg1audio.FrameHeader
		err := h.Unmarshal(frames[0])
		if err != nil {
			return err
		}

		if h.MPEG2 {
			return fmt.Errorf("Only MPEG-1 audio is supported")
		}

		track.mp3Checked = true
	}

	n := 0
	for _, frame := range frames {
		n += len(frame)
	}

	enc := make([]byte, n)
	n = 0
	for _, frame := range frames {
		n += copy(enc[n:], frame)
	}

	_, err := w.mux.WriteData(&astits.MuxerData{
		PID: track.PID,
		AdaptationField: &astits.PacketAdaptationField{
			RandomAccessIndicator: true,
		},
		PES: &astits.PESData{
			Header: &astits.PESHeader{
				OptionalHeader: &astits.PESOptionalHeader{
					MarkerBits:      2,
					PTSDTSIndicator: astits.PTSDTSIndicatorOnlyPTS,
					PTS:             &astits.ClockReference{Base: pts},
				},
				StreamID: streamIDAudio,
			},
			Data: enc,
		},
	})
	return err
}
