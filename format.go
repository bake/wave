package wave

import (
	"encoding/binary"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Format holds configuration about the WAVE.
type Format struct {
	AudioFormat   uint16 // 1 if PCM is used.
	NumChans      uint16 // Number of channels (1 = mono, 2 = stereo, ...)
	SampleRate    uint32 // Samples per second (44100, ...).
	ByteRate      uint32 // Average bytes per second.
	BlockAlign    uint16 // Bytes per sample.
	BitsPerSample uint16 // Bits per sample.
}

// decodeFormat decodes a chunk in a format chunk.
func decodeFormat(r io.Reader) (*Format, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "could not read format chunk")
	}
	return &Format{
		AudioFormat:   binary.LittleEndian.Uint16(data[:2]),
		NumChans:      binary.LittleEndian.Uint16(data[2:4]),
		SampleRate:    binary.LittleEndian.Uint32(data[4:8]),
		ByteRate:      binary.LittleEndian.Uint32(data[8:12]),
		BlockAlign:    binary.LittleEndian.Uint16(data[12:14]),
		BitsPerSample: binary.LittleEndian.Uint16(data[14:16]),
	}, nil
}

// Encode encodes a format struct into an io.Writer.
func (f *Format) Encode(w io.Writer) error {
	var err error
	write := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Write(w, binary.LittleEndian, data)
	}
	write(f.AudioFormat)
	write(f.NumChans)
	write(f.SampleRate)
	write(f.ByteRate)
	write(f.BlockAlign)
	write(f.BitsPerSample)
	return err
}
