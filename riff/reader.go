package riff

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

const riffID = "RIFF"

// Reader reads a RIFF file chunk by chunk.
type Reader struct {
	r     io.Reader
	chunk struct {
		id   string
		size int64
		data io.Reader
		err  error
	}
}

// NewReader reads the initial RIFF header and returns a chunk reader and its
// type.
func NewReader(r io.Reader) (rr *Reader, riffType string, err error) {
	rr = &Reader{r: r}
	if !rr.Next() {
		if rr.Error() == nil {
			return nil, "", errors.Wrap(io.EOF, "unecpected EOF")
		}
		return nil, "", errors.Wrap(rr.Error(), "could not read RIFF chunk")
	}
	id, _, data := rr.Chunk()
	if id != riffID {
		return nil, "", errors.Errorf("unexpected chunk id %s", id)
	}
	t := make([]byte, 4)
	if _, err := data.Read(t); err != nil {
		return nil, "", errors.Wrap(err, "could not read RIFF type")
	}
	return rr, string(t), nil
}

// Next returns true until the underlying reader returns an error like EOF. The
// caller is responsible to read or seek to the end of the chunk before calling
// Next again.
func (rr *Reader) Next() bool {
	header := make([]byte, 8)
	_, err := rr.r.Read(header)
	if err == io.EOF {
		rr.chunk.err = io.EOF
		return false
	}
	if err != nil {
		rr.chunk.err = errors.Wrap(err, "could not read chunk header")
		return false
	}
	rr.chunk.id = string(header[:4])
	rr.chunk.size = int64(binary.LittleEndian.Uint32(header[4:]))
	rr.chunk.data = io.LimitReader(rr.r, rr.chunk.size)
	return rr.chunk.err == nil
}

// Chunk returns the current chunk. This function can be called multiple times.
func (rr *Reader) Chunk() (id string, size int64, data io.Reader) {
	return rr.chunk.id, rr.chunk.size, rr.chunk.data
}

// Err returns the first non-EOF error.
func (rr Reader) Error() error {
	if rr.chunk.err == io.EOF {
		return nil
	}
	return rr.chunk.err
}
