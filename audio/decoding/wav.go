package decoding

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type waveDecoder struct {
	decoderBase
	header         *riffHeader
	info           *riffChunkFmt
	firstSamplePos uint32
}

type riffHeader struct {
	Ftype       [4]byte
	ChunkSize   uint32
	ChunkFormat [4]byte
}

type riffChunkFmt struct {
	LengthOfHeader uint32
	AudioFormat    uint16 // 1 = PCM not compressed
	NumChannels    uint16
	SampleRate     uint32
	BytesPerSec    uint32
	BytesPerBloc   uint16
	BitsPerSample  uint16
}

func (decoder *waveDecoder) read() error {
	var err error

	decoder.header = &riffHeader{}
	if err = binary.Read(decoder.src, binary.LittleEndian, decoder.header); err != nil {
		return err
	}

	if !bytes.Equal(decoder.header.Ftype[:], []byte("RIFF")) ||
		!bytes.Equal(decoder.header.ChunkFormat[:], []byte("WAVE")) {
		return errors.New("Not a RIFF/WAVE file.")
	}

	var chunk [4]byte
	var chunkSize uint32
	for {
		// Read next chunkID
		err = binary.Read(decoder.src, binary.BigEndian, &chunk)
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		} else if err != nil {
			return err
		}

		// and it's size in bytes
		err = binary.Read(decoder.src, binary.LittleEndian, &chunkSize)
		if err == io.EOF {
			return io.ErrUnexpectedEOF
		} else if err != nil {
			return err
		}

		if bytes.Equal(chunk[:], []byte("fmt ")) {
			// seek 4 bytes back because riffChunkFmt reads the chunkSize again
			if _, err = decoder.src.Seek(-4, os.SEEK_CUR); err != nil {
				return err
			}

			decoder.info = &riffChunkFmt{}
			if err = binary.Read(decoder.src, binary.LittleEndian, decoder.info); err != nil {
				return err
			}
			if decoder.info.LengthOfHeader > 16 { // canonical format if chunklen == 16
				// Skip extra params
				if _, err = decoder.src.Seek(int64(decoder.info.LengthOfHeader-16), os.SEEK_CUR); err != nil {
					return err
				}
			}

			// Is audio supported ?
			if decoder.info.AudioFormat != 1 {
				return fmt.Errorf("Audio Format not supported")
			}
		} else if bytes.Equal(chunk[:], []byte("data")) {
			size, _ := decoder.src.Seek(0, os.SEEK_CUR)
			decoder.firstSamplePos = uint32(size)
			decoder.dataSize = int32(chunkSize)
			break
		} else {
			if _, err = decoder.src.Seek(int64(chunkSize), os.SEEK_CUR); err != nil {
				return err
			}
		}
	}

	if decoder.info == nil {
		return fmt.Errorf("unable to read the wav file format")
	}

	bytesPerSample := int32(decoder.info.BitsPerSample / 8)
	numSamples := decoder.dataSize / bytesPerSample

	decoder.channels = int16(decoder.info.NumChannels)
	decoder.sampleRate = int32(decoder.info.SampleRate)
	decoder.bitDepth = int16(decoder.info.BitsPerSample)
	decoder.format = getFormat(decoder.channels, decoder.bitDepth)
	decoder.duration = time.Duration(float64(numSamples)/float64(decoder.info.SampleRate)) * time.Second

	return nil
}
