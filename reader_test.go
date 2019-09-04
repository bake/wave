package wave_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/bake/wave"
)

func TestReader(t *testing.T) {
	r := bytes.NewReader([]byte{
		// R,    I,    F,    F,                   2084,    W,    A,    V,    E,
		0x52, 0x49, 0x46, 0x46, 0x24, 0x08, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45,

		// f,    m,    t,    ␣,                     16,          1,          2,
		0x66, 0x6d, 0x74, 0x20, 0x10, 0x00, 0x00, 0x00, 0x01, 0x00, 0x02, 0x00,
		//               22050,                  88200,          4,         16,
		0x22, 0x56, 0x00, 0x00, 0x88, 0x58, 0x01, 0x00, 0x04, 0x00, 0x10, 0x00,

		// s,    l,    n,    t,                      4,
		0x73, 0x6c, 0x6e, 0x74, 0x04, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,

		// d,    a,    t,    a,                     28,          0,          0,
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
	out := []int{
		0, 0, 5924, -3298, 4924, 5180, -1770, -1768,
		-6348, -23005, -3524, -3548, -12783, 3354,
		0, 0, 5924, -3298, 4924, 5180, -1770, -1768,
	}
	wavr, err := wave.NewReader(r)
	if err != nil {
		t.Fatalf("could not create new wave reader: %v", err)
	}
	samples, err := wavr.Samples()
	if err != nil {
		t.Fatalf("could not read samples: %v", err)
	}
	if fmt.Sprint(samples) != fmt.Sprint(out) {
		t.Fatalf("expected samples to be\n%v, got\n%v", out, samples)
	}
}

func TestNewReader(t *testing.T) {
	tt := []struct {
		name string
		res  bool
		data []byte
	}{
		{
			name: "empty reader",
			res:  false,
			data: []byte{},
		},
		{
			name: "wrong riff id",
			res:  false,
			data: []byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "wrong riff type",
			res:  false,
			data: []byte{
				// R,    I,    F,    F,                      4,    0,    0,    0,    0,
				0x52, 0x49, 0x46, 0x46, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "no format chunk",
			res:  false,
			data: []byte{
				// R,    I,    F,    F,                      4,    W,    A,    V,    E,
				0x52, 0x49, 0x46, 0x46, 0x04, 0x00, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45,
			},
		},
		{
			name: "unexpected chunk id",
			res:  false,
			data: []byte{
				// R,    I,    F,    F,                     12,    W,    A,    V,    E,
				0x52, 0x49, 0x46, 0x46, 0x0c, 0x00, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45,
				// s,    l,    n,    t,                      4,
				0x73, 0x6c, 0x6e, 0x74, 0x04, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.data)
			_, err := wave.NewReader(r)
			if err != nil && tc.res {
				t.Fatalf("unexpected error: %v\n", err)
			}
			if err == nil && !tc.res {
				t.Fatal("expected an error")
			}
		})
	}
}

func ExampleReader() {
	r := bytes.NewReader([]byte{
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
	})
	wavr, err := wave.NewReader(r)
	if err != nil {
		log.Fatalf("could not create new wave reader: %v", err)
	}
	fmt.Printf("SampleRate: %d\n", wavr.Format.SampleRate)
	samples, err := wavr.Samples()
	if err != nil {
		log.Fatalf("could not read samples: %v", err)
	}
	fmt.Printf("Samples: %v\n", samples[:10])

	// Output:
	// SampleRate: 44100
	// Samples: [0 0 5924 -3298 4924 5180 -1770 -1768 -6348 -23005]
}
