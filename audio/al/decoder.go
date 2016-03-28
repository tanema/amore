// Package decoding is used for converting familiar file types to data usable by
// OpenAL.

package al

import (
	"io"
	"time"
)

type (
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

var formatBytes = map[uint32]int32{
	FormatMono8:    1,
	FormatMono16:   2,
	FormatStereo8:  2,
	FormatStereo16: 4,
}

const BUFFER_SIZE = 128 * 1024

func getFormat(channels, depth int16) uint32 {
	switch channels, depth := channels, depth; {
	case channels == 1 && depth == 8:
		return FormatMono8
	case channels == 1 && depth == 16:
		return FormatMono16
	case channels == 2 && depth == 8:
		return FormatStereo8
	case channels == 2 && depth == 16:
		return FormatStereo16
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

func (decoder *decoderBase) ByteOffsetToDur(offset int32) time.Duration {
	return time.Duration(offset * formatBytes[decoder.GetFormat()] * int32(time.Second) / decoder.GetSampleRate())
}

func (decoder *decoderBase) DurToByteOffset(dur time.Duration) int32 {
	return int32(dur) * int32(decoder.GetSampleRate()) / (formatBytes[decoder.GetFormat()] * int32(time.Second))
}

func (decoder *decoderBase) Decode() int {
	buffer := make([]byte, BUFFER_SIZE)
	n, err := decoder.src.Read(buffer)
	decoder.eof = (err == io.EOF)
	decoder.buffer = buffer[:n]
	return n
}

func (decoder *decoderBase) Seek(s int64) bool {
	_, err := decoder.src.Seek(s, 0)
	decoder.eof = (err == io.EOF)
	return err == nil || decoder.eof
}
