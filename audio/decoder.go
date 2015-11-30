package audio

import (
	"errors"
	"os"

	"github.com/tanema/amore/file"
)

const (
	// See http://www.topherlee.com/software/pcm-tut-wavformat.html.
	EXT_WAVE = ".wav"
)

type Decoder interface {
	Decode(*os.File) error
	GetBuffer() *[]byte
	Seek(float32) bool
	Rewind() bool
	IsSeekable() bool
	IsFinished() bool
	GetSize() int
	GetChannels() int16
	GetBitDepth() int16
	GetSampleRate() int32
	durToByteOffset(float32) int64
	byteOffsetToDur(offset int64) float64
}

func decode(filepath string) (Decoder, error) {
	src, err := file.NewFile(filepath)
	if err != nil {
		return nil, err
	}

	var decoder Decoder

	switch file.Ext(filepath) {
	case EXT_WAVE:
		decoder = &waveDecoder{}
	default:
		return nil, errors.New("unsupported audio file extention")
	}

	err = decoder.Decode(src)
	if err != nil {
		return nil, err
	}
	return decoder, nil
}
