package wave_test

import (
	"testing"

	"github.com/bake/wave"
	"github.com/orcaman/writerseeker"
)

func TestWriter(t *testing.T) {
	// TODO: Compare written data.
	samples := []int{
		0, 0,
		5924, -3298, 4924, 5180, -1770, -1768,
		-6348, -23005, -3524, -3548, -12783, 3354,
	}

	ws := &writerseeker.WriterSeeker{}
	fmt := &wave.Format{
		AudioFormat:   0x1,
		NumChans:      0x2,
		SampleRate:    0xac44,
		ByteRate:      0x2b110,
		BlockAlign:    0x4,
		BitsPerSample: 0x10,
	}
	wavw, err := wave.NewWriter(ws, fmt)
	if err != nil {
		t.Fatalf("could not create wave writer: %v", err)
	}
	for _, s := range samples {
		if err := wavw.Sample(s); err != nil {
			t.Fatalf("could not write sample %d: %v", s, err)
		}
	}
	if err := wavw.Close(); err != nil {
		t.Fatalf("could not close wave writer: %v", err)
	}
}
