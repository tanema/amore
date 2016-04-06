// Package decoding is used for converting familiar file types to data usable by
// OpenAL.
package decoding

import (
	"errors"
	"io"
	"time"

	"github.com/tanema/amore/audio/al"
	"github.com/tanema/amore/file"
)

type (
	// Decoder is an interface for all the decoders and future sound data objects
	Decoder interface {
		read() error
		GetBuffer() []byte
		Seek(int64) bool
		IsFinished() bool
		GetFormat() uint32
		GetChannels() int16
		GetBitDepth() int16
		GetSampleRate() int32
		GetData() []byte
		Decode() int
		GetDuration() time.Duration
		ByteOffsetToDur(offset int32) time.Duration
		DurToByteOffset(dur time.Duration) int32
	}
	// decoderBase is a base implementation of a few methods to keep tryings DRY
	decoderBase struct {
		src        io.ReadSeeker
		channels   int16
		sampleRate int32
		bitDepth   int16
		duration   time.Duration
		eof        bool
		buffer     []byte
		data       []byte
		format     uint32
	}
)

// The amount of bytes in each format
var formatBytes = map[uint32]int32{
	al.FormatMono8:    1,
	al.FormatMono16:   2,
	al.FormatStereo8:  2,
	al.FormatStereo16: 4,
}

// arbitrary buffer size, could be tuned in the future
const BUFFER_SIZE = 128 * 1024

// Decode will get the file at the path provided. It will then send it to the decoder
// that will handle its file type by the extention on the path. Supported formats
// are wav, ogg, and flac. If there is an error retrieving the file or decoding it,
// will return that error.
func Decode(filepath string) (Decoder, error) {
	src, err := file.NewFile(filepath)
	if err != nil {
		return nil, err
	}

	var decoder Decoder
	switch file.Ext(filepath) {
	case ".wav":
		decoder = &waveDecoder{decoderBase: decoderBase{src: src}}
	case ".ogg":
		decoder = &vorbisDecoder{decoderBase: decoderBase{src: src}}
	case ".flac":
		decoder = &flacDecoder{decoderBase: decoderBase{src: src}}
	default:
		src.Close()
		return nil, errors.New("unsupported audio file extention")
	}

	if err = decoder.read(); err != nil {
		src.Close()
		return nil, err
	}

	return decoder, nil
}

// getFormat will return the openal format for the channels and depth provided
func getFormat(channels, depth int16) uint32 {
	switch channels, depth := channels, depth; {
	case channels == 1 && depth == 8:
		return al.FormatMono8
	case channels == 1 && depth == 16:
		return al.FormatMono16
	case channels == 2 && depth == 8:
		return al.FormatStereo8
	case channels == 2 && depth == 16:
		return al.FormatStereo16
	default:
		return 0
	}
}

//Stubs for methods so functionality can be DRY
func (decoder *decoderBase) GetBuffer() []byte          { return decoder.buffer }
func (decoder *decoderBase) IsFinished() bool           { return decoder.eof }
func (decoder *decoderBase) GetSampleRate() int32       { return decoder.sampleRate }
func (decoder *decoderBase) GetChannels() int16         { return decoder.channels }
func (decoder *decoderBase) GetBitDepth() int16         { return decoder.bitDepth }
func (decoder *decoderBase) GetDuration() time.Duration { return decoder.duration }
func (decoder *decoderBase) GetFormat() uint32          { return decoder.format }
func (decoder *decoderBase) GetData() []byte            { return decoder.data }

// ByteOffsetToDur will translate byte count to time duration
func (decoder *decoderBase) ByteOffsetToDur(offset int32) time.Duration {
	return time.Duration(offset * formatBytes[decoder.GetFormat()] * int32(time.Second) / decoder.GetSampleRate())
}

// DurToByteOffset will translate time duration to a byte count
func (decoder *decoderBase) DurToByteOffset(dur time.Duration) int32 {
	return int32(dur) * int32(decoder.GetSampleRate()) / (formatBytes[decoder.GetFormat()] * int32(time.Second))
}

// Decode will read the next chunk into the buffer and return the amount of bytes read
func (decoder *decoderBase) Decode() int {
	buffer := make([]byte, BUFFER_SIZE)
	n, err := decoder.src.Read(buffer)
	decoder.eof = (err == io.EOF)
	decoder.buffer = buffer[:n]
	return n
}

// Seek will seek in the source data by count of bytes
func (decoder *decoderBase) Seek(s int64) bool {
	_, err := decoder.src.Seek(s, 0)
	decoder.eof = (err == io.EOF)
	return err == nil || decoder.eof
}
