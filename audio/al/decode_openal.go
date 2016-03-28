// +build !js

package al

import (
	"errors"

	"github.com/tanema/amore/file"
)

func Decode(filepath string) (Decoder, error) {
	src, err := file.NewFile(filepath)
	if err != nil {
		return nil, err
	}

	var decoder Decoder
	switch file.Ext(filepath) {
	case ".wav":
		decoder = &waveDecoder{decoderBase: decoderBase{src: src}}
	case ".ogg":
		decoder = &vorbisDecoder{decoderBase: decoderBase{src: src}}
	case ".flac":
		decoder = &flacDecoder{decoderBase: decoderBase{src: src}}
	default:
		src.Close()
		return nil, errors.New("unsupported audio file extention")
	}

	if err = decoder.read(); err != nil {
		src.Close()
		return nil, err
	}

	return decoder, nil
}
