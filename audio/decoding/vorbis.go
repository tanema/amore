package decoding

/*
 * This is not the best, it creates a byte array of all the pcm data which is not good,
 * but it makes it a lot easier to seek and handle the data. Streaming will still have a
 * smaller memory footprint while buffering to openal. I am not sure how nice this will
 * be in a mobile platform
 */

import (
	"bytes"
	"io"

	"github.com/hajimehoshi/go-vorbis"
)

type vorbisDecoder struct {
	decoderBase
}

// read will decode the file
func (decoder *vorbisDecoder) read() error {
	var err error
	var channels, sampleRate int
	var reader io.ReadCloser
	reader, channels, sampleRate, err = vorbis.Decode(decoder.src)
	defer reader.Close()

	decoder.channels = int16(channels)
	decoder.sampleRate = int32(sampleRate)
	decoder.bitDepth = 16
	decoder.format = getFormat(decoder.channels, decoder.bitDepth)

	data := []byte{}
	for {
		tmp_buffer := make([]byte, 4097)
		n, err := reader.Read(tmp_buffer)
		if err == io.EOF {
			break
		}
		data = append(data, tmp_buffer[:n]...)
	}

	decoder.dataSize = int32(len(data))
	decoder.src = bytes.NewReader(data)
	decoder.duration = decoder.ByteOffsetToDur(int32(len(data)) / formatBytes[decoder.GetFormat()])

	return err
}
