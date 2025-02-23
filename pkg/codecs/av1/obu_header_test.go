package av1

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOBUHeaderUnmarshal(t *testing.T) {
	for _, ca := range []struct {
		name string
		byts []byte
		h    OBUHeader
	}{
		{
			"sequence header",
			[]byte{
				0x0a, 0x0e, 0x00, 0x00, 0x00, 0x4a, 0xab, 0xbf,
				0xc3, 0x77, 0x6b, 0xe4, 0x40, 0x40, 0x40, 0x41,
			},
			OBUHeader{
				Type:    OBUTypeSequenceHeader,
				HasSize: true,
			},
		},
	} {
		t.Run(ca.name, func(t *testing.T) {
			var h OBUHeader
			err := h.Unmarshal(ca.byts)
			require.NoError(t, err)
			require.Equal(t, ca.h, h)
		})
	}
}

func FuzzOBUHeaderUnmarshal(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		var h OBUHeader
		h.Unmarshal(b) //nolint:errcheck
	})
}
