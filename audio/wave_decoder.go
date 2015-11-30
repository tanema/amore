package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

type waveDecoder struct {
	decoderBase
}

// See http://www.topherlee.com/software/pcm-tut-wavformat.html.
func (decoder *waveDecoder) Decode(src *os.File) error {
	var riffMarker, waveMarker, fmtMarker, dataMarker [4]byte
	var fmtChunkSize int32

	//descriptor
	binary.Read(src, binary.BigEndian, &riffMarker)
	binary.Read(src, binary.LittleEndian, &decoder.fileSize)
	binary.Read(src, binary.BigEndian, &waveMarker)

	//fmt chunk
	binary.Read(src, binary.BigEndian, &fmtMarker)
	binary.Read(src, binary.LittleEndian, &fmtChunkSize)
	binary.Read(src, binary.LittleEndian, &decoder.audioFormat)
	binary.Read(src, binary.LittleEndian, &decoder.channels)
	binary.Read(src, binary.LittleEndian, &decoder.sampleRate)
	binary.Read(src, binary.LittleEndian, &decoder.byteRate)
	binary.Read(src, binary.LittleEndian, &decoder.blockAlign)
	binary.Read(src, binary.LittleEndian, &decoder.bitDepth)

	//data chunk
	binary.Read(src, binary.BigEndian, &dataMarker)
	//verify we have correct header data if we have the data marker we
	//definitely made it through, which is nice
	if !bytes.Equal(riffMarker[:], []byte("RIFF")) || //RIFF marker
		!bytes.Equal(waveMarker[:], []byte("WAVE")) || //WAVE marker
		!bytes.Equal(fmtMarker[:], []byte("fmt ")) || //fmt block
		fmtChunkSize != 16 || //fmt chunk size unknown bail out
		!bytes.Equal(dataMarker[:], []byte("data")) { //didnt find the data marker
		return errors.New("Not a wave file.")
	}
	binary.Read(src, binary.LittleEndian, &decoder.dataSize)
	//read data into fixed length array
	decoder.data = make([]byte, decoder.dataSize, decoder.dataSize)
	binary.Read(src, binary.LittleEndian, &decoder.data)

	//extra
	//calculate the duration form the data size and byterate
	decoder.duration = float32(decoder.dataSize) / float32(decoder.byteRate)

	switch channels, depth := decoder.channels, decoder.bitDepth; {
	case channels == 1 && depth == 8:
		decoder.format = Mono8
	case channels == 1 && depth == 16:
		decoder.format = Mono16
	case channels == 2 && depth == 8:
		decoder.format = Stereo8
	case channels == 2 && depth == 16:
		decoder.format = Stereo16
	default:
		return fmt.Errorf("unsupported format; num of channels=%d, bit rate=%d", channels, depth)
	}

	return nil
}
