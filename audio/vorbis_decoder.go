package audio

import (
	"errors"
)

type vorbisDecoder struct {
	decoderBase
}

const (
	VORBIS_HEADER_SIZE = 44
)

func (decoder *vorbisDecoder) readHeaders() error {
	return errors.New("not implemented yet")
}
