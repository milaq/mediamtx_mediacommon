# mediacommon

[![Test](https://github.com/bluenviron/mediacommon/workflows/test/badge.svg)](https://github.com/bluenviron/mediacommon/actions?query=workflow:test)
[![Lint](https://github.com/bluenviron/mediacommon/workflows/lint/badge.svg)](https://github.com/bluenviron/mediacommon/actions?query=workflow:lint)
[![Go Report Card](https://goreportcard.com/badge/github.com/bluenviron/mediacommon)](https://goreportcard.com/report/github.com/bluenviron/mediacommon)
[![CodeCov](https://codecov.io/gh/bluenviron/mediacommon/branch/main/graph/badge.svg)](https://app.codecov.io/gh/bluenviron/mediacommon/branch/main)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/bluenviron/mediacommon)](https://pkg.go.dev/github.com/bluenviron/mediacommon#pkg-index)

Definitions and functions shared between [gortsplib](https://github.com/bluenviron/gortsplib), [gohlslib](https://github.com/bluenviron/gohlslib) and [MediaMTX](https://github.com/bluenviron/mediamtx), in particular:

* [Codec utilities](https://pkg.go.dev/github.com/bluenviron/mediacommon/pkg/codecs)
* [Format utilities](https://pkg.go.dev/github.com/bluenviron/mediacommon/pkg/formats)
* [Bit reader and writer](https://pkg.go.dev/github.com/bluenviron/mediacommon/pkg/bits)

## Standards

* [ITU-T Rec. T-871, JPEG File Interchange Format](https://www.itu.int/rec/dologin_pub.asp?lang=e&id=T-REC-T.871-201105-I!!PDF-E&type=items)
* [ITU-T Rec. H.264 (08/2021)](https://www.itu.int/rec/dologin_pub.asp?lang=e&id=T-REC-H.264-202108-I!!PDF-E&type=items)
* [ITU-T Rec. H.265 (08/2021)](https://www.itu.int/rec/dologin_pub.asp?lang=e&id=T-REC-H.265-202108-I!!PDF-E&type=items)
* ISO 11172-3, Coding of moving pictures and associated audio
* ISO 13818-3, Generic Coding of Moving Pictures and Associated Audio: Audio
* ISO 14496-3, Coding of audio-visual objects, part 3, Audio
* [AV1 Bitstream & Decoding Process](https://aomediacodec.github.io/av1-spec/av1-spec.pdf)
* [AV1 Codec ISO Media File Format Binding](https://aomediacodec.github.io/av1-isobmff)
* [VP9 Bitstream & Decoding Process Specification v0.6](https://storage.googleapis.com/downloads.webmproject.org/docs/vp9/vp9-bitstream-specification-v0.6-20160331-draft.pdf)
* [VP9 Codec ISO Media FIle Format Binding](https://www.webmproject.org/vp9/mp4/)
* [RFC6716, Definition of the Opus Audio Codec](https://datatracker.ietf.org/doc/html/rfc6716)
* [Opus in MP4/ISOBMFF](https://opus-codec.org/docs/opus_in_isobmff.html)

## Related projects

* [MediaMTX](https://github.com/bluenviron/mediamtx)
* [gortsplib](https://github.com/bluenviron/gortsplib)
* [gohlslib](https://github.com/bluenviron/gohlslib)
