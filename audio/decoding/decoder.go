package decoding

import (
	"errors"
	"io"
	"time"

	"github.com/tanema/amore/audio/al"
	"github.com/tanema/amore/file"
)

const (
	EXT_WAVE   = ".wav"
	EXT_VORBIS = ".ogg"
	EXT_MP3    = ".mp3"
	EXT_MOD    = ".mod"
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
	GetSize() int32
	GetChannels() int16
	GetBitDepth() int16
	GetSampleRate() int32
	GetData() []byte
	Decode() int
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
	case EXT_WAVE:
		decoder = &waveDecoder{base}
	case EXT_VORBIS:
		decoder = &vorbisDecoder{base}
	case EXT_MP3:
		decoder = &mp3Decoder{base}
	case EXT_MOD:
		decoder = &modDecoder{base}
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

func (decoder *decoderBase) Read() error {
	panic("Decoder read headers method not implemented")
	return nil
}

func (decoder *decoderBase) Decode() int {
	buffer := make([]byte, 128*1024)
	n, err := decoder.src.Read(buffer)
	decoder.eof = (err == io.EOF)
	decoder.buffer = buffer[:n]
	return n
}

func (decoder *decoderBase) GetBuffer() []byte {
	return decoder.buffer
}

func (decoder *decoderBase) GetData() []byte {
	data := make([]byte, decoder.dataSize)
	decoder.Rewind()
	decoder.src.Read(data)
	return data
}

func (decoder *decoderBase) Seek(s int64) bool {
	_, err := decoder.src.Seek(int64(decoder.headerSize)+s, 0)
	decoder.eof = (err == io.EOF)
	return err == nil || decoder.eof
}

func (decoder *decoderBase) Rewind() bool {
	return decoder.Seek(0)
}

func (decoder *decoderBase) IsFinished() bool {
	return decoder.eof
}

func (decoder *decoderBase) GetFormat() uint32 {
	return decoder.format
}

func (decoder *decoderBase) GetChannels() int16 {
	return decoder.channels
}

func (decoder *decoderBase) GetBitDepth() int16 {
	return decoder.bitDepth
}

func (decoder *decoderBase) GetSampleRate() int32 {
	return decoder.sampleRate
}

func (decoder *decoderBase) GetSize() int32 {
	return decoder.dataSize
}

func (decoder *decoderBase) ByteOffsetToDur(offset int32) time.Duration {
	return time.Duration(offset * formatBytes[decoder.format] * int32(time.Second) / decoder.sampleRate)
}

func (decoder *decoderBase) DurToByteOffset(dur time.Duration) int32 {
	return int32(dur) * int32(decoder.sampleRate) / (formatBytes[decoder.format] * int32(time.Second))
}
