package riff_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/bake/wave/riff"
)

func exampleInt8WaveReader() io.ReadSeeker {
	return bytes.NewReader([]byte{})
}

func exampleInt16WaveReader() io.ReadSeeker {
	// This example is borrowed from http://soundfile.sapp.org/doc/WaveFormat/.
	return bytes.NewReader([]byte{
		// R,    I,    F,    F,                   2084,    W,    A,    V,    E,
		0x52, 0x49, 0x46, 0x46, 0x24, 0x08, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45,

		// f,    m,    t,    ␣,                     16,          1,          2,
		0x66, 0x6d, 0x74, 0x20, 0x10, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00,
		//               22050,                  88200,          4,         16,
		0x22, 0x56, 0x00, 0x00, 0x88, 0x58, 0x01, 0x00, 0x04, 0x00, 0x10, 0x00,

		// s,    l,    n,    t,                      4,
		0x73, 0x6c, 0x6e, 0x74, 0x04, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,

		// d,    a,    t,    a,                     27,          0,          0,
		0x64, 0x61, 0x74, 0x61, 0x1c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		//    5924,      -3298,       4924,       5180,      -1770,      -1768,
		0x24, 0x17, 0x1e, 0xf3, 0x3c, 0x13, 0x3c, 0x14, 0x16, 0xf9, 0x18, 0xf9,
		//   -6348,     -23005,      -3524,      -3548,     -12783,       3354,
		0x34, 0xe7, 0x23, 0xa6, 0x3c, 0xf2, 0x24, 0xf2, 0x11, 0xce, 0x1a, 0x0d,

		// d,    a,    t,    a,                     16,          0,          0,
		0x64, 0x61, 0x74, 0x61, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		//    5924,      -3298,       4924,       5180,      -1770,      -1768,
		0x24, 0x17, 0x1e, 0xf3, 0x3c, 0x13, 0x3c, 0x14, 0x16, 0xf9, 0x18, 0xf9,
	})
}

func TestReader(t *testing.T) {
	r := exampleInt16WaveReader()
	rr, riffType, err := riff.NewReader(r)
	if err != nil {
		t.Fatalf("could not create riff reader: %v", err)
	}
	if riffType != "WAVE" {
		t.Fatalf("expected RIFF type to be \"WAVE\", got \"%s\"", riffType)
	}

	chunks := []struct {
		id   string
		size int64
	}{{"fmt ", 16}, {"slnt", 4}, {"data", 28}, {"data", 16}}

	for i := 0; rr.Next(); i++ {
		id, size, _ := rr.Chunk()
		if id != chunks[i].id {
			t.Fatalf("expected id to be \"%s\", got \"%s\"", chunks[i].id, id)
		}
		if size != chunks[i].size {
			t.Fatalf("expected size to be %d, got %d", chunks[i].size, size)
		}
		if _, err := r.Seek(size, io.SeekCurrent); err != nil {
			t.Fatalf("could not seek chunk: %v", err)
		}
	}
	if err := rr.Error(); err != nil {
		t.Fatalf("could not read wav: %v", err)
	}
}

func ExampleReader() {
	r := exampleInt16WaveReader()
	rr, riffType, err := riff.NewReader(r)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("type: %s\n", riffType)
	for rr.Next() {
		id, size, data := rr.Chunk()
		body := make([]byte, size)
		if _, err := data.Read(body); err != nil {
			log.Fatal(err)
		}
		switch id {
		case "fmt ", "slnt":
			fmt.Printf("%s: % x\n", id, body)
		default:
			fmt.Printf("%s: ...\n", id)
		}
	}
	if err := rr.Error(); err != nil {
		log.Fatal(err)
	}
	// Output:
	// type: WAVE
	// fmt : 01 00 02 00 22 56 00 00 88 58 01 00 04 00 10 00
	// slnt: ff ff ff ff
	// data: ...
	// data: ...
}
