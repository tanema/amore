/*
 * flac: 	 github.com/eaburns/flac
 * vorbis: github.com/runningwild/gorbis || github.com/mccoyst/vorbis
 * wav: 	 github.com/youpy/go-wav
 * mp3: 	 github.com/tcolgate/mp3
 */

package decoding

import (
	"errors"
	"io"
	"time"

	"github.com/tanema/amore/audio/al"
	"github.com/tanema/amore/file"
)

// ReadSeekCloser is an io.ReadSeeker and io.Closer.
type ReadSeekCloser interface {
	io.ReadSeeker
	io.Closer
}

type Decoder interface {
	read() error
	GetBuffer() []byte
	Seek(int64) bool
	Rewind() bool
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

func Decode(filepath string) (Decoder, error) {
	src, err := file.NewFile(filepath)
	if err != nil {
		return nil, err
	}

	var decoder Decoder
	base := decoderBase{src: src}

	switch file.Ext(filepath) {
	case ".wav":
		decoder = &waveDecoder{decoderBase: base}
	case ".ogg", ".oga":
		decoder = &vorbisDecoder{decoderBase: base}
	case ".mp3":
		decoder = &mp3Decoder{decoderBase: base}
	case ".flac":
		decoder = &flacDecoder{decoderBase: base}
	default:
		return nil, errors.New("unsupported audio file extention")
	}

	err = decoder.read()
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

type decoderBase struct {
	src        ReadSeekCloser
	channels   int16
	sampleRate int32
	bitDepth   int16
	duration   time.Duration
	eof        bool
	buffer     []byte
}

func (decoder *decoderBase) GetBuffer() []byte {
	return decoder.buffer
}

//Stubs for methods so functionality can be DRY
func (decoder *decoderBase) Seek(s int64) bool { return false }
func (decoder *decoderBase) Rewind() bool {
	return decoder.Seek(0)
}

func (decoder *decoderBase) IsFinished() bool {
	return decoder.eof
}

func (decoder *decoderBase) GetFormat() uint32 {
	return getFormat(decoder.GetChannels(), decoder.GetBitDepth())
}

func (decoder *decoderBase) GetSampleRate() int32       { return decoder.sampleRate }
func (decoder *decoderBase) GetChannels() int16         { return decoder.channels }
func (decoder *decoderBase) GetBitDepth() int16         { return decoder.bitDepth }
func (decoder *decoderBase) GetDuration() time.Duration { return decoder.duration }

func (decoder *decoderBase) ByteOffsetToDur(offset int32) time.Duration {
	return time.Duration(offset * formatBytes[decoder.GetFormat()] * int32(time.Second) / decoder.GetSampleRate())
}

func (decoder *decoderBase) DurToByteOffset(dur time.Duration) int32 {
	return int32(dur) * int32(decoder.GetSampleRate()) / (formatBytes[decoder.GetFormat()] * int32(time.Second))
}
