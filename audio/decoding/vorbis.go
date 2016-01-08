package decoding

import (
	"errors"

	"github.com/mccoyst/vorbis"
)

type vorbisDecoder struct {
	decoderBase
	handle *vorbis.Vorbis
}

func (decoder *vorbisDecoder) read() error {
	return errors.New("not implemented")
}

func (decoder *vorbisDecoder) Decode() int {
	return 0
}

func (decoder *vorbisDecoder) GetData() []byte {
	return []byte{}
}

func (decoder *vorbisDecoder) Seek(s int64) bool {
	return false
}
