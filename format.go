package wave

import (
	"encoding/binary"
	"io"
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
func decodeFormat(r io.Reader) (Format, error) {
	var dst Format
	err := binary.Read(r, binary.LittleEndian, &dst)
	return dst, err
}

// encode a format struct into an io.Writer.
func (f *Format) encode(w io.Writer) error {
	return binary.Write(w, binary.LittleEndian, f)
}
