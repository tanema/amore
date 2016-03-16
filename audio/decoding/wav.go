// See http://www.topherlee.com/software/pcm-tut-wavformat.html.
package decoding

import (
	"bytes"
	"encoding/binary"
	"errors"
	"time"
)

type waveDecoder struct {
	decoderBase
	dataSize int32
}

const (
	WAV_HEADER_SIZE = 44
)

func (decoder *waveDecoder) read() error {
	var riffMarker, waveMarker, fmtMarker, dataMarker [4]byte
	var fmtChunkSize, fileSize, byteRate int32
	var audioFormat, blockAlign int16

	//descriptor
	binary.Read(decoder.src, binary.BigEndian, &riffMarker)
	binary.Read(decoder.src, binary.LittleEndian, &fileSize)
	binary.Read(decoder.src, binary.BigEndian, &waveMarker)

	//fmt chunk
	binary.Read(decoder.src, binary.BigEndian, &fmtMarker)
	binary.Read(decoder.src, binary.LittleEndian, &fmtChunkSize)
	binary.Read(decoder.src, binary.LittleEndian, &audioFormat)
	binary.Read(decoder.src, binary.LittleEndian, &decoder.channels)
	binary.Read(decoder.src, binary.LittleEndian, &decoder.sampleRate)
	binary.Read(decoder.src, binary.LittleEndian, &byteRate)
	binary.Read(decoder.src, binary.LittleEndian, &blockAlign)
	binary.Read(decoder.src, binary.LittleEndian, &decoder.bitDepth)

	decoder.format = getFormat(decoder.channels, decoder.bitDepth)

	//data chunk
	binary.Read(decoder.src, binary.BigEndian, &dataMarker)
	//verify we have correct header data if we have the data marker we
	//definitely made it through, which is nice
	if !bytes.Equal(riffMarker[:], []byte("RIFF")) || //RIFF marker
		!bytes.Equal(waveMarker[:], []byte("WAVE")) || //WAVE marker
		!bytes.Equal(fmtMarker[:], []byte("fmt ")) || //fmt block
		fmtChunkSize != 16 || //fmt chunk size unknown bail out
		!bytes.Equal(dataMarker[:], []byte("data")) { //didnt find the data marker
		return errors.New("Not a RIFF/WAVE file.")
	}
	decoder.dataSize = fileSize - WAV_HEADER_SIZE
	//we could read file the file like this
	//binary.Read(decoder.src, binary.LittleEndian, &decoder.dataSize)
	//and subtract 8 for the data marker and the data size but why read and calculate
	//when we can just calculate

	//extra
	//calculate the duration form the data size and samplerate
	decoder.duration = time.Duration((float32(decoder.dataSize) / float32(decoder.sampleRate)) * float32(time.Second))

	return nil
}

func (decoder *waveDecoder) GetData() []byte {
	data := make([]byte, decoder.dataSize)
	decoder.Seek(0)
	decoder.src.Read(data)
	return data
}
