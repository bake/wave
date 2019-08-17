package wave

import (
	"io"

	"github.com/bake/wave/riff"
	"github.com/pkg/errors"
)

// Writer writes samples to an io.WriteCloser.
type Writer struct {
	io.WriteCloser
	fmt Format
}

// NewWriter creates a new WAVE Writer.
func NewWriter(ws io.WriteSeeker, fmt *Format) (*Writer, error) {
	rw, err := riff.NewWriter(ws, "WAVE")
	if err != nil {
		return nil, errors.Wrap(err, "could not create new riff reader")
	}
	cw, err := rw.Chunk("fmt ")
	if err != nil {
		return nil, errors.Wrap(err, "could not create format chunk")
	}
	if err := fmt.encode(cw); err != nil {
		return nil, errors.Wrap(err, "could not encode format chunk")
	}
	if err := cw.Close(); err != nil {
		return nil, errors.Wrap(err, "could not close format chunk")
	}
	cw, err = rw.Chunk("data")
	if err != nil {
		return nil, errors.Wrap(err, "could not create data chunk")
	}
	return &Writer{cw, *fmt}, nil
}

// Sample writes a sample.
func (wavw *Writer) Sample(s int) error {
	var p []byte
	switch wavw.fmt.BitsPerSample {
	case 8:
		p = []byte{byte(s)}
	case 16:
		p = []byte{byte(s), byte(s >> 8)}
	case 24:
		p = []byte{byte(s), byte(s >> 8), byte(s >> 16)}
	case 32:
		p = []byte{byte(s), byte(s >> 8), byte(s >> 16), byte(s >> 24)}
	}
	if _, err := wavw.Write(p); err != nil {
		return errors.Wrap(err, "could not write sample")
	}
	return nil
}
