package wave

import (
	"io"

	"github.com/bake/wave/riff"
	"github.com/pkg/errors"
)

// Reader reads samples from a WAVE file.
type Reader struct {
	rr  *riff.Reader
	fmt *Format
}

// NewReader reads the initial chunks from a WAVE file and returns a new reader.
func NewReader(r io.Reader) (*Reader, error) {
	t, rr, err := riff.NewReader(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not create new riff reader")
	}
	if t != "WAVE" {
		return nil, errors.Errorf("unexpected riff type %s", t)
	}
	if !rr.Next() {
		return nil, errors.Wrap(rr.Error(), "could not read format chunk")
	}
	id, _, data := rr.Chunk()
	if id != "fmt " {
		return nil, errors.Errorf("unexpected chunk id %s", id)
	}
	format, err := decodeFormat(data)
	if err != nil {
		return nil, errors.Wrap(err, "could not decode format chunk")
	}
	return &Reader{rr, format}, nil
}

// Sample returns the next sample from the wave file. Chunks that don't contain
// samples are skipped.
func (wavr *Reader) Sample() (int, error) {
	id, size, data := wavr.rr.Chunk()
	if id != "data" {
		body := make([]byte, size)
		if _, err := data.Read(body); err != nil && err != io.EOF {
			return 0, errors.Wrapf(err, "could not skip %s chunk", id)
		}
		wavr.rr.Next()
		return wavr.Sample()
	}
	s, err := wavr.sample(data)
	if err == io.EOF {
		if !wavr.rr.Next() {
			return 0, io.EOF
		}
		return wavr.Sample()
	}
	return s, err
}

// Samples reads the whole file and returns all samples.
func (wavr *Reader) Samples() ([]int, error) {
	var samples []int
	for {
		s, err := wavr.Sample()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "could not read sample")
		}
		samples = append(samples, s)
	}
	return samples, nil
}

func (wavr *Reader) sample(r io.Reader) (int, error) {
	s := make([]byte, wavr.fmt.BitsPerSample/8)
	if _, err := r.Read(s); err != nil {
		return 0, err
	}
	switch wavr.fmt.BitsPerSample {
	case 8:
		return int(uint8(s[0])), nil
	case 16:
		return int(int16(s[0]) | int16(s[1])<<8), nil
	case 24:
		return int(int32(s[0]) | int32(s[1])<<8 | int32(s[2])<<16), nil
	case 32:
		return int(int32(s[0]) | int32(s[1])<<8 | int32(s[2])<<16 | int32(s[3])<<24), nil
	default:
		return 0, errors.Errorf("unpexpected bps: %d", wavr.fmt.BitsPerSample)
	}
}
