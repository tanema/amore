package decoding

import (
	"errors"
)

type modDecoder struct {
	decoderBase
}

func (decoder *modDecoder) read() error {
	return errors.New("not implemented yet")
}
