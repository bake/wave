package wave_test

import (
	"io/ioutil"
	"testing"

	"github.com/bake/wave"
	"github.com/orcaman/writerseeker"
)

func TestWriter(t *testing.T) {
	samples := []int{
		0, 0,
		5924, -3298, 4924, 5180, -1770, -1768,
		-6348, -23005, -3524, -3548, -12783, 3354,
	}

	ws := &writerseeker.WriterSeeker{}
	wavw, err := wave.NewWriter(ws)
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

	b, _ := ioutil.ReadAll(ws.Reader())
	t.Errorf("% x\n", b)
}
