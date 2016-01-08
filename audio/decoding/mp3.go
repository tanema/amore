package decoding

import (
	"errors"
)

type mp3Decoder struct {
	decoderBase
}

func (decoder *mp3Decoder) read() error {
	return errors.New("not implemented yet")
}

func (decoder *mp3Decoder) Decode() int {
	return 0
}

func (decoder *mp3Decoder) GetData() []byte {
	return []byte{}
}

func (decoder *mp3Decoder) Seek(s int64) bool {
	return false
}
