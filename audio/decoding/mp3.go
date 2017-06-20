package decoding

import (
	"bytes"
	"io"

	"github.com/hajimehoshi/go-mp3"
	"github.com/tanema/amore/file"
)

type mp3Decoder struct {
	decoderBase
}

// read will decode the file
func (decoder *mp3Decoder) read() error {
	src, err := file.NewFile(decoder.src_path)
	if err != nil {
		return err
	}
	defer src.Close()

	d, err := mp3.Decode(src)
	if err != nil {
		return err
	}

	decoder.channels = int16(2)
	decoder.sampleRate = int32(d.SampleRate())
	decoder.bitDepth = int16(16)
	decoder.format = getFormat(decoder.channels, decoder.bitDepth)

	data := []byte{}
	for {
		tmp_buffer := make([]byte, 4097)
		n, err := d.Read(tmp_buffer)
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
