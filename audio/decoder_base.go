package audio

import (
	"os"
)

type decoderBase struct {
	fileSize         int32
	formatDataLength int32
	sampleRate       int32
	byteRate         int32
	dataChunkSize    int32
	audioFormat      int16
	channels         int16
	byteSampleRate   int16
	bitsPerSample    int16
	data             []byte
	duration         float32
}

func (decoder *decoderBase) Decode(src *os.File) error {
	panic("Decoder decode method not implemented")
	return nil
}

func (decoder *decoderBase) GetBuffer() *[]byte {
	return &decoder.data
}

func (decoder *decoderBase) Seek(s float32) bool {
	return false
}

func (decoder *decoderBase) Rewind() bool {
	return false
}

func (decoder *decoderBase) IsSeekable() bool {
	return false
}

func (decoder *decoderBase) IsFinished() bool {
	return false
}

func (decoder *decoderBase) GetChannels() int16 {
	return decoder.channels
}

func (decoder *decoderBase) GetBitDepth() int16 {
	return decoder.bitsPerSample
}

func (decoder *decoderBase) GetSampleRate() int32 {
	return decoder.sampleRate
}

func (decoder *decoderBase) GetSize() int {
	return len(decoder.data)
}
