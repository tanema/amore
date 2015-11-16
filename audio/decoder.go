package audio

import (
	"errors"
	"os"

	"github.com/tanema/amore/file"
)

// See http://www.topherlee.com/software/pcm-tut-wavformat.html.
const (
	EXT_WAVE = ".wav"
)

const (
	//Indicates how many bytes of raw data should be generated at each call to Decode.
	DEFAULT_BUFFER_SIZE = 16384
	// Indicates the quality of the sound.
	DEFAULT_SAMPLE_RATE = 44100
	//	Default is stereo.
	DEFAULT_CHANNELS = 2
	// 16 bit audio is the default.
	DEFAULT_BIT_DEPTH = 16
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
		return nil, errors.New("unknow file extention")
	}

	err = decoder.Decode(src)
	if err != nil {
		return nil, err
	}
	return decoder, nil
}
