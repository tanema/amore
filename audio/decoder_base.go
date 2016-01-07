package audio

import (
	"io"
	"time"
)

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
}

func (decoder *decoderBase) readHeaders() error {
	panic("Decoder read headers method not implemented")
	return nil
}

func (decoder *decoderBase) decode() int32 {
	panic("Decoder decode method not implemented")
	return 0
}

func (decoder *decoderBase) getBuffer() []byte {
	return []byte{}
}

func (decoder *decoderBase) getData() []byte {
	data := make([]byte, decoder.dataSize)
	decoder.seek(int64(decoder.headerSize))
	decoder.src.Read(data)
	return data
}

func (decoder *decoderBase) seek(s int64) bool {
	_, err := decoder.src.Seek(s, 0)
	decoder.eof = (err == io.EOF)
	return err == nil || decoder.eof
}

func (decoder *decoderBase) rewind() bool {
	return decoder.seek(int64(decoder.headerSize))
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
