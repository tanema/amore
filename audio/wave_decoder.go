package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type waveDecoder struct {
	decoderBase
}

// See http://www.topherlee.com/software/pcm-tut-wavformat.html.
const (
	WAV_HEADER_SIZE = 44
)

func (decoder *waveDecoder) readHeaders() error {
	var riffMarker, waveMarker, fmtMarker, dataMarker [4]byte
	var fmtChunkSize int32

	//descriptor
	binary.Read(decoder.src, binary.BigEndian, &riffMarker)
	binary.Read(decoder.src, binary.LittleEndian, &decoder.fileSize)
	binary.Read(decoder.src, binary.BigEndian, &waveMarker)

	//fmt chunk
	binary.Read(decoder.src, binary.BigEndian, &fmtMarker)
	binary.Read(decoder.src, binary.LittleEndian, &fmtChunkSize)
	binary.Read(decoder.src, binary.LittleEndian, &decoder.audioFormat)
	binary.Read(decoder.src, binary.LittleEndian, &decoder.channels)
	binary.Read(decoder.src, binary.LittleEndian, &decoder.sampleRate)
	binary.Read(decoder.src, binary.LittleEndian, &decoder.byteRate)
	binary.Read(decoder.src, binary.LittleEndian, &decoder.blockAlign)
	binary.Read(decoder.src, binary.LittleEndian, &decoder.bitDepth)

	//data chunk
	binary.Read(decoder.src, binary.BigEndian, &dataMarker)
	//verify we have correct header data if we have the data marker we
	//definitely made it through, which is nice
	if !bytes.Equal(riffMarker[:], []byte("RIFF")) || //RIFF marker
		!bytes.Equal(waveMarker[:], []byte("WAVE")) || //WAVE marker
		!bytes.Equal(fmtMarker[:], []byte("fmt ")) || //fmt block
		fmtChunkSize != 16 || //fmt chunk size unknown bail out
		!bytes.Equal(dataMarker[:], []byte("data")) { //didnt find the data marker
		return errors.New("Not a wave file.")
	}
	decoder.headerSize = WAV_HEADER_SIZE
	decoder.dataSize = decoder.fileSize - decoder.headerSize
	//we could read file the file like this
	//binary.Read(decoder.src, binary.LittleEndian, &decoder.dataSize)
	//and subtract 8 for the data marker and the data size but why read and calculate
	//when we can just calculate

	//extra
	//calculate the duration form the data size and byterate
	decoder.duration = float32(decoder.dataSize) / float32(decoder.sampleRate)

	decoder.format = getFormat(decoder.channels, decoder.bitDepth)
	if decoder.format == 0 {
		return fmt.Errorf("unsupported format; num of channels=%d, bit rate=%d", decoder.channels, decoder.bitDepth)
	}

	return nil
}
