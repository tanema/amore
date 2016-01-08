package decoding

import (
	"io"

	"github.com/eaburns/flac"
)

type flacDecoder struct {
	decoderBase
	handle   *flac.Decoder
	dataSize int32
}

func (decoder *flacDecoder) read() error {
	var err error
	decoder.handle, err = flac.NewDecoder(decoder.src)
	decoder.channels = int16(decoder.handle.NChannels)
	decoder.sampleRate = int32(decoder.handle.SampleRate)
	decoder.bitDepth = int16(decoder.handle.BitsPerSample)
	decoder.dataSize = int32(decoder.handle.TotalSamples) * int32(decoder.channels) * int32(decoder.bitDepth/8)
	return err
}

func (decoder *flacDecoder) Decode() int {
	buffer, err := decoder.handle.Next()
	decoder.eof = (err == io.EOF)
	decoder.buffer = buffer
	return len(decoder.buffer)
}

func (decoder *flacDecoder) GetData() []byte {
	data := make([]byte, decoder.dataSize)
	for {
		buffer, err := decoder.handle.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return []byte{}
		}
		data = append(data, buffer...)
	}
	return data
}

func (decoder *flacDecoder) Seek(s int64) bool {
	return false
}
