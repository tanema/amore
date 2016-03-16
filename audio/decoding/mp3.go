package decoding

import (
	"fmt"
	"io"

	"github.com/tcolgate/mp3"
)

type mp3Decoder struct {
	decoderBase
	data []byte
}

func (decoder *mp3Decoder) read() error {
	d := mp3.NewDecoder(decoder.src)
	var frame mp3.Frame
	for {
		err := d.Decode(&frame)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return err
		}
	}

	header := frame.Header()
	decoder.duration = frame.Duration()
	decoder.sampleRate = int32(header.SampleRate())
	decoder.bitDepth = 16
	switch header.ChannelMode() {
	case mp3.Stereo, mp3.DualChannel, mp3.JointStereo:
		decoder.channels = int16(2)
	case mp3.SingleChannel:
		decoder.channels = int16(1)
	}
	decoder.format = getFormat(decoder.channels, decoder.bitDepth)

	return nil
}
