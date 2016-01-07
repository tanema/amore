package audio

import (
	"errors"
	"io"
	"time"

	"github.com/tanema/amore/audio/al"
	"github.com/tanema/amore/file"
)

const (
	// See http://www.topherlee.com/software/pcm-tut-wavformat.html.
	EXT_WAVE = ".wav"
)

// ReadSeekCloser is an io.ReadSeeker and io.Closer.
type ReadSeekCloser interface {
	io.ReadSeeker
	io.Closer
}

type Decoder interface {
	readHeaders() error
	getBuffer() []byte
	seek(int64) bool
	rewind() bool
	isFinished() bool
	getFormat() uint32
	getSize() int32
	getChannels() int16
	getBitDepth() int16
	getSampleRate() int32
	getData() []byte
	decode() int32
	byteOffsetToDur(offset int32) time.Duration
	durToByteOffset(dur time.Duration) int32
}

func decode(filepath string) (Decoder, error) {
	src, err := file.NewFile(filepath)
	if err != nil {
		return nil, err
	}

	var decoder Decoder

	switch file.Ext(filepath) {
	case EXT_WAVE:
		decoder = &waveDecoder{decoderBase: decoderBase{src: src}}
	default:
		return nil, errors.New("unsupported audio file extention")
	}

	err = decoder.readHeaders()
	if err != nil {
		return nil, err
	}

	return decoder, nil
}

var formatBytes = map[uint32]int32{
	al.FormatMono8:    1,
	al.FormatMono16:   2,
	al.FormatStereo8:  2,
	al.FormatStereo16: 4,
}

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
