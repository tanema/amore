// Package decoding is used for converting familiar file types to data usable by
// OpenAL.
package decoding

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/tanema/amore/audio/al"
	"github.com/tanema/amore/file"
)

type (
	// Decoder is a base implementation of a few methods to keep tryings DRY
	Decoder struct {
		src        io.ReadCloser
		codec      io.ReadSeeker
		bitDepth   int16
		eof        bool
		dataSize   int32
		currentPos int64
		byteRate   int32

		Channels   int16
		SampleRate int32
		Buffer     []byte
		Format     uint32
	}
)

var extHandlers = map[string]func(io.ReadCloser) (*Decoder, error){
	".wav":  newWaveDecoder,
	".ogg":  newVorbisDecoder,
	".flac": newFlacDecoder,
	".mp3":  newMp3Decoder,
}

// arbitrary buffer size, could be tuned in the future
const bufferSize = 128 * 1024

// Decode will get the file at the path provided. It will then send it to the decoder
// that will handle its file type by the extention on the path. Supported formats
// are wav, ogg, and flac. If there is an error retrieving the file or decoding it,
// will return that error.
func Decode(filepath string) (*Decoder, error) {
	src, err := file.NewFile(filepath)
	if err != nil {
		return nil, err
	}

	handler, hasHandler := extHandlers[file.Ext(filepath)]
	if !hasHandler {
		src.Close()
		return nil, errors.New("unsupported audio file extention")
	}

	decoder, err := handler(src)
	if err != nil {
		src.Close()
		return nil, err
	}

	return decoder, nil
}

func newDecoder(src io.ReadCloser, codec io.ReadSeeker, channels int16, sampleRate int32, bitDepth int16, dataSize int32) *Decoder {
	decoder := &Decoder{
		src:        src,
		codec:      codec,
		Channels:   channels,
		SampleRate: sampleRate,
		bitDepth:   bitDepth,
		dataSize:   dataSize,
		currentPos: 0,
		byteRate:   (sampleRate * int32(channels) * int32(bitDepth)) / 8,
	}

	switch {
	case channels == 1 && bitDepth == 8:
		decoder.Format = al.FormatMono8
	case channels == 1 && bitDepth == 16:
		decoder.Format = al.FormatMono16
	case channels == 2 && bitDepth == 8:
		decoder.Format = al.FormatStereo8
	case channels == 2 && bitDepth == 16:
		decoder.Format = al.FormatStereo16
	}

	return decoder
}

// IsFinished returns is the decoder is finished decoding
func (decoder *Decoder) IsFinished() bool {
	return decoder.eof
}

// Duration will return the total time duration of this clip
func (decoder *Decoder) Duration() time.Duration {
	return decoder.ByteOffsetToDur(decoder.dataSize)
}

// GetData returns the complete set of data
func (decoder *Decoder) GetData() []byte {
	data := make([]byte, decoder.dataSize)
	decoder.Seek(0)
	decoder.codec.Read(data)
	return data
}

// ByteOffsetToDur will translate byte count to time duration
func (decoder *Decoder) ByteOffsetToDur(offset int32) time.Duration {
	return time.Duration(offset/decoder.byteRate) * time.Second
}

// DurToByteOffset will translate time duration to a byte count
func (decoder *Decoder) DurToByteOffset(dur time.Duration) int32 {
	return int32(dur/time.Second) * decoder.byteRate
}

// Decode will read the next chunk into the buffer and return the amount of bytes read
func (decoder *Decoder) Decode() int {
	buffer := make([]byte, bufferSize)
	n, err := decoder.codec.Read(buffer)
	decoder.Buffer = buffer[:n]
	decoder.eof = (err == io.EOF)
	return n
}

// Seek will seek in the source data by count of bytes
func (decoder *Decoder) Seek(s int64) bool {
	decoder.currentPos = s
	_, err := decoder.codec.Seek(decoder.currentPos, os.SEEK_SET)
	decoder.eof = (err == io.EOF)
	return err == nil || decoder.eof
}
