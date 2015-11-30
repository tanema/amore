package audio

import (
	"os"
)

type decoderBase struct {
	fileSize    int32
	audioFormat int16
	format      Format
	channels    int16
	sampleRate  int32
	byteRate    int32
	blockAlign  int16
	bitDepth    int16
	dataSize    int32
	data        []byte
	duration    float32
}

func (decoder *decoderBase) Decode(src *os.File) error {
	panic("Decoder decode method not implemented")
	return nil
}

func (decoder *decoderBase) GetBuffer() *[]byte {
	return &decoder.data
}

func (decoder *decoderBase) Seek(s float32) bool {
	return true
}

func (decoder *decoderBase) Rewind() bool {
	return decoder.Seek(0)
}

func (decoder *decoderBase) IsSeekable() bool {
	return true
}

func (decoder *decoderBase) IsFinished() bool {
	return false
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

func (decoder *decoderBase) GetSize() int {
	return len(decoder.data)
}

func (decoder *decoderBase) durToByteOffset(offset float32) int64 {
	return 0
}

func (decoder *decoderBase) byteOffsetToDur(offset int64) float64 {
	return 0
}
