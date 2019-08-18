package riff_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/bake/wave/riff"
	"github.com/orcaman/writerseeker"
)

func TestWriter(t *testing.T) {
	chunks := []struct {
		id   string
		data []byte
	}{
		{"dat1", []byte{0x00, 0x01, 0x02, 0x03}},
		{"dat2", []byte{0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a}},
	}
	out := []byte{
		// R,    I,    F,    F,                     32,
		0x52, 0x49, 0x46, 0x46, 0x20, 0x00, 0x00, 0x00,
		// W,    A,    V,    E,    d,    a,    t,    1,
		0x57, 0x41, 0x56, 0x45, 0x64, 0x61, 0x74, 0x31,
		//                   4,    0,    1,    2,    3,
		0x04, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02, 0x03,
		// d,    a,    t,    2,                      7,
		0x64, 0x61, 0x74, 0x32, 0x07, 0x00, 0x00, 0x00,
		// 4,    5,    6,    0,    8,    9,   10,    0,
		0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x00,
	}

	ws := &writerseeker.WriterSeeker{}
	rw, err := riff.NewWriter(ws, "WAVE")
	if err != nil {
		t.Fatalf("could not create new riff writer: %v", err)
	}

	for _, c := range chunks {
		cw, err := rw.Chunk(c.id)
		if err != nil {
			t.Fatalf("could not create chunk %s: %v", c.id, err)
		}
		if _, err := cw.Write(c.data); err != nil {
			t.Fatalf("could not write to %s: %v", c.id, err)
		}
		if err := cw.Close(); err != nil {
			t.Fatalf("could not close %s: %v", c.id, err)
		}
	}
	if err := rw.Close(); err != nil {
		t.Fatalf("could not close riff: %v", err)
	}
	if err := ws.Close(); err != nil {
		t.Fatalf("could not close test file: %v", err)
	}

	defer ws.Close()
	body, _ := ioutil.ReadAll(ws.Reader())
	if fmt.Sprintf("% x", body) != fmt.Sprintf("% x", out) {
		t.Fatalf("expected body to be\n% x, got\n% x\n", out, body)
	}
}

func ExampleWriter() {
	ws := &writerseeker.WriterSeeker{}
	rw, err := riff.NewWriter(ws, "WAVE")
	if err != nil {
		log.Fatalf("could not create new riff writer: %v", err)
	}

	cw, err := rw.Chunk("foo1")
	if err != nil {
		log.Fatalf("could not create chunk: %v", err)
	}
	if _, err = cw.Write([]byte{0xff, 0xff}); err != nil {
		log.Fatalf("could not write data: %v", err)
	}
	if err := cw.Close(); err != nil {
		log.Fatalf("could not close chunk: %v", err)
	}

	// Could (or should) be deferred.
	if err := rw.Close(); err != nil {
		log.Fatalf("could not close reader: %v", err)
	}

	body, _ := ioutil.ReadAll(ws.Reader())
	fmt.Printf("% x\n", body)

	// Output:
	// 52 49 46 46 0e 00 00 00 57 41 56 45 66 6f 6f 31 02 00 00 00 ff ff
}
