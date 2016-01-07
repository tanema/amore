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
	EXT_WAVE   = ".wav"
	EXT_VORBIS = ".ogg"
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
	decode() int
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
	case EXT_VORBIS:
		decoder = &vorbisDecoder{decoderBase: decoderBase{src: src}}
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

type decoderBase struct {
	src         ReadSeekCloser
	fileSize    int32
	audioFormat int16
	format      uint32
	channels    int16
	sampleRate  int32
	byteRate    int32
	blockAlign  int16
	bitDepth    int16
	dataSize    int32
	duration    float32
	headerSize  int32
	eof         bool
	buffer      []byte
}

func (decoder *decoderBase) readHeaders() error {
	panic("Decoder read headers method not implemented")
	return nil
}

func (decoder *decoderBase) decode() int {
	buffer := make([]byte, 128*1024)
	n, err := decoder.src.Read(buffer)
	decoder.eof = (err == io.EOF)
	decoder.buffer = buffer[:n]
	return n
}

func (decoder *decoderBase) getBuffer() []byte {
	return decoder.buffer
}

func (decoder *decoderBase) getData() []byte {
	data := make([]byte, decoder.dataSize)
	decoder.rewind()
	decoder.src.Read(data)
	return data
}

func (decoder *decoderBase) seek(s int64) bool {
	_, err := decoder.src.Seek(int64(decoder.headerSize)+s, 0)
	decoder.eof = (err == io.EOF)
	return err == nil || decoder.eof
}

func (decoder *decoderBase) rewind() bool {
	return decoder.seek(0)
}

func (decoder *decoderBase) isFinished() bool {
	return decoder.eof
}

func (decoder *decoderBase) getFormat() uint32 {
	return decoder.format
}

func (decoder *decoderBase) getChannels() int16 {
	return decoder.channels
}

func (decoder *decoderBase) getBitDepth() int16 {
	return decoder.bitDepth
}

func (decoder *decoderBase) getSampleRate() int32 {
	return decoder.sampleRate
}

func (decoder *decoderBase) getSize() int32 {
	return decoder.dataSize
}

func (decoder *decoderBase) byteOffsetToDur(offset int32) time.Duration {
	return time.Duration(offset * formatBytes[decoder.format] * int32(time.Second) / decoder.sampleRate)
}

func (decoder *decoderBase) durToByteOffset(dur time.Duration) int32 {
	return int32(dur) * int32(decoder.sampleRate) / (formatBytes[decoder.format] * int32(time.Second))
}
