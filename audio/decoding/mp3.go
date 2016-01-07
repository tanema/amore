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
