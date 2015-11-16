package audio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

type waveDecoder struct {
	decoderBase
}

// See http://www.topherlee.com/software/pcm-tut-wavformat.html.
const WavHeaderSize = 44

func (decoder *waveDecoder) Decode(src *os.File) error {
	new_decoder := waveDecoder{}
	var riffMarker, waveMarker, fmtMarker, dataMarker [4]byte
	binary.Read(src, binary.LittleEndian, &riffMarker)
	binary.Read(src, binary.LittleEndian, &new_decoder.fileSize)
	binary.Read(src, binary.LittleEndian, &waveMarker)
	binary.Read(src, binary.LittleEndian, &fmtMarker)
	binary.Read(src, binary.LittleEndian, &new_decoder.formatDataLength)
	binary.Read(src, binary.LittleEndian, &new_decoder.audioFormat)
	binary.Read(src, binary.LittleEndian, &new_decoder.channels)
	binary.Read(src, binary.LittleEndian, &new_decoder.sampleRate)
	binary.Read(src, binary.LittleEndian, &new_decoder.byteRate)
	binary.Read(src, binary.LittleEndian, &new_decoder.byteSampleRate)
	binary.Read(src, binary.LittleEndian, &new_decoder.bitsPerSample)
	binary.Read(src, binary.LittleEndian, &dataMarker)

	//verify we have correct header data if we have the data marker we
	//definitely made it through, which is nice
	if !bytes.Equal(riffMarker[:], []byte("RIFF")) ||
		!bytes.Equal(waveMarker[:], []byte("WAVE")) ||
		!bytes.Equal(fmtMarker[:], []byte("fmt ")) ||
		!bytes.Equal(dataMarker[:], []byte("data")) {
		return errors.New("Not a wave file.")
	}

	binary.Read(src, binary.LittleEndian, &new_decoder.dataChunkSize)
	new_decoder.data = make([]byte, new_decoder.dataChunkSize, new_decoder.dataChunkSize)
	binary.Read(src, binary.LittleEndian, &new_decoder.data)
	new_decoder.duration = float32(new_decoder.dataChunkSize) / float32(new_decoder.byteRate)
	return nil
}
