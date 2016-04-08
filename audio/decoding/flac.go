package decoding

import (
	"bytes"
	"io"

	"github.com/eaburns/flac"
)

type flacDecoder struct {
	decoderBase
}

// read will decode the file
func (decoder *flacDecoder) read() error {
	d, err := flac.NewDecoder(decoder.src)
	if err != nil {
		return err
	}

	decoder.channels = int16(d.NChannels)
	decoder.sampleRate = int32(d.SampleRate)
	decoder.bitDepth = int16(d.BitsPerSample)
	decoder.format = getFormat(decoder.channels, decoder.bitDepth)
	decoder.dataSize = int32(d.TotalSamples * int64(d.NChannels) * int64(d.BitsPerSample/8))
	data := make([]byte, 0, decoder.dataSize)
	for {
		frame, err := d.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		data = append(data, frame...)
	}

	decoder.src = bytes.NewReader(data)
	decoder.duration = decoder.ByteOffsetToDur(int32(len(data)) / formatBytes[decoder.GetFormat()])

	return err
}
