package wave

import (
	"io"

	"github.com/bake/wave/riff"
	"github.com/pkg/errors"
)

// Writer writes samples to an io.Writer.
type Writer struct {
	rw  *riff.Writer
	cw  *riff.Writer
	fmt Format
}

// NewWriter creates a new WAVE Writer.
func NewWriter(ws io.WriteSeeker, format Format) (*Writer, error) {
	rw, err := riff.NewWriter(ws, "WAVE")
	if err != nil {
		return nil, errors.Wrap(err, "could not create new riff reader")
	}
	cw, err := rw.Chunk("fmt ")
	if err != nil {
		return nil, errors.Wrap(err, "could not create format chunk")
	}
	if err := format.encode(cw); err != nil {
		return nil, errors.Wrap(err, "could not encode format chunk")
	}
	if err := cw.Close(); err != nil {
		return nil, errors.Wrap(err, "could not close format chunk")
	}
	cw, err = rw.Chunk("data")
	if err != nil {
		return nil, errors.Wrap(err, "could not create data chunk")
	}
	return &Writer{rw, cw, format}, nil
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
	if _, err := wavw.cw.Write(p); err != nil {
		return errors.Wrap(err, "could not write sample")
	}
	return nil
}

// Samples writes a slice of samples.
func (wavw *Writer) Samples(samples []int) error {
	for _, s := range samples {
		if err := wavw.Sample(s); err != nil {
			return err
		}
	}
	return nil
}

// Close the underlying RIFF writer. The file writer needs to be closed
// separately.
func (wavw *Writer) Close() error {
	if err := wavw.cw.Close(); err != nil {
		return err
	}
	return wavw.rw.Close()
}
