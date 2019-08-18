package wave_test

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/bake/wave"
	"github.com/orcaman/writerseeker"
)

func TestWriter(t *testing.T) {
	// TODO: Compare written data.
	format := wave.Format{
		AudioFormat:   1,
		NumChans:      2,
		SampleRate:    44100,
		ByteRate:      176400,
		BlockAlign:    4,
		BitsPerSample: 16,
	}
	samples := []int{
		0, 0, 5924, -3298, 4924, 5180, -1770, -1768,
		-6348, -23005, -3524, -3548, -12783, 3354,
		0, 0, 5924, -3298, 4924, 5180, -1770, -1768,
	}
	out := []byte{
		// R,    I,    F,    F,                     76,    W,    A,    V,    E,
		0x52, 0x49, 0x46, 0x46, 0x50, 0x00, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45,

		// f,    m,    t,    ␣,                     16,          1,          2,
		0x66, 0x6d, 0x74, 0x20, 0x10, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00,
		//               44100,                 176400,          4,         16,
		0x44, 0xac, 0x00, 0x00, 0x10, 0xb1, 0x02, 0x00, 0x04, 0x00, 0x10, 0x00,

		// d,    a,    t,    a,                     44,          0,          0,
		0x64, 0x61, 0x74, 0x61, 0x2c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		//    5924,      -3298,       4924,       5180,      -1770,      -1768,
		0x24, 0x17, 0x1e, 0xf3, 0x3c, 0x13, 0x3c, 0x14, 0x16, 0xf9, 0x18, 0xf9,
		//   -6348,     -23005,      -3524,      -3548,     -12783,       3354,
		0x34, 0xe7, 0x23, 0xa6, 0x3c, 0xf2, 0x24, 0xf2, 0x11, 0xce, 0x1a, 0x0d,
		//       0,          0,       5924,      -3298,       4924,       5180,
		0x00, 0x00, 0x00, 0x00, 0x24, 0x17, 0x1e, 0xf3, 0x3c, 0x13, 0x3c, 0x14,
		//   -1770,      -1768,
		0x16, 0xf9, 0x18, 0xf9,
	}

	ws := &writerseeker.WriterSeeker{}
	wavw, err := wave.NewWriter(ws, format)
	if err != nil {
		t.Fatalf("could not create wave writer: %v", err)
	}
	for _, s := range samples {
		if err := wavw.Sample(s); err != nil {
			t.Fatalf("could not write sample %d: %v", s, err)
		}
	}
	wavw.Close()

	body, _ := ioutil.ReadAll(ws.Reader())
	if fmt.Sprintf("% x", body) != fmt.Sprintf("% x", out) {
		t.Fatalf("expected body to be\n% x, got\n% x\n", out, body)
	}
}
