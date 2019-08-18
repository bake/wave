package riff

import (
	"encoding/binary"
	"io"

	"github.com/pkg/errors"
)

const (
	riffTypeSize  = 4
	chunkTypeSize = 4
	sizeFieldSize = 4
)

// Writer extends an io.WriteSeeker by the ability to write in chunks.
type Writer struct {
	io.WriteSeeker
	start int64
	size  int64
}

// NewWriter creates a new RIFF writer and writes the initial RIFF chunk.
func NewWriter(ws io.WriteSeeker, riffType string) (*Writer, error) {
	if len(riffType) != riffTypeSize {
		return nil, errors.Errorf("riff type has to be %d bytes long", riffTypeSize)
	}
	w := &Writer{WriteSeeker: ws}
	cw, err := w.Chunk("RIFF")
	if err != nil {
		return nil, errors.Wrap(err, "could not create riff chunk")
	}
	if _, err := cw.Write([]byte(riffType)); err != nil {
		return nil, errors.Wrap(err, "could not write riff type")
	}
	return cw, nil
}

// Chunk creates a new chunk that has to be closed before creating a second one
// or closing the RIFF writer.
func (w *Writer) Chunk(chunkType string) (*Writer, error) {
	if len(chunkType) != chunkTypeSize {
		return nil, errors.Errorf("chunk type has to be %d bytes long", chunkTypeSize)
	}
	start, err := w.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, errors.Wrap(err, "could not get current position")
	}
	// The parent WriteSeeker might itself be a *Writer that counts written bytes.
	// When closed, chunks seek back to their starting position and overwrite the
	// initial size (an uint32), thus incrementing their parent chunks size by
	// additionon 4 bytes.
	// As a countermeasure and to not keep references to parent chunks, their
	// sizes are decremented by 4 on each creation of a new child.
	w.size -= sizeFieldSize
	cw := &Writer{WriteSeeker: w, start: start}
	header := append([]byte(chunkType), make([]byte, sizeFieldSize)...)
	if _, err := cw.WriteSeeker.Write(header); err != nil {
		return nil, errors.Wrap(err, "could not write chunk header")
	}
	return cw, nil
}

// Close seeks to the chunks beginning, writes its sice and seeks back to the
// writers end. The underlying io.WriteCloser has to be closed separately.
func (w *Writer) Close() error {
	size := w.size
	data := make([]byte, sizeFieldSize)
	binary.LittleEndian.PutUint32(data, uint32(size))
	if _, err := w.Seek(w.start+chunkTypeSize, io.SeekStart); err != nil {
		return errors.Wrap(err, "could not seek to beginning of chunk")
	}
	if _, err := w.WriteSeeker.Write(data); err != nil {
		return errors.Wrap(err, "could not write chunk size")
	}
	if _, err := w.Seek(size, io.SeekCurrent); err != nil {
		return errors.Wrap(err, "could not seek to end of chunk")
	}
	// Add an aditional byte if the data is not word aligned.
	if size%2 == 1 {
		if _, err := w.Write([]byte{0x00}); err != nil {
			return errors.Wrap(err, "could not write padding byte")
		}
	}
	return nil
}

// Write to the chunk.
func (w *Writer) Write(p []byte) (n int, err error) {
	n, err = w.WriteSeeker.Write(p)
	w.size += int64(n)
	return n, err
}
