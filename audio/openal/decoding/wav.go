package decoding

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
)

type waveDecoder struct {
	io.Reader
	src            io.ReadCloser
	header         *riffHeader
	info           *riffChunkFmt
	firstSamplePos uint32
	dataSize       int32
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

func newWaveDecoder(src io.ReadCloser) (*Decoder, error) {
	decoder := &waveDecoder{
		src:    src,
		Reader: src,
		info:   &riffChunkFmt{},
		header: &riffHeader{},
	}

	if err := decoder.decode(); err != nil {
		return nil, err
	}

	return newDecoder(
		src,
		decoder,
		int16(decoder.info.NumChannels),
		int32(decoder.info.SampleRate),
		int16(decoder.info.BitsPerSample),
		decoder.dataSize,
	), nil
}

func (decoder *waveDecoder) decode() error {
	if err := binary.Read(decoder.src, binary.LittleEndian, decoder.header); err != nil {
		return err
	}

	if !bytes.Equal(decoder.header.Ftype[:], []byte("RIFF")) ||
		!bytes.Equal(decoder.header.ChunkFormat[:], []byte("WAVE")) {
		return errors.New("not a RIFF/WAVE file")
	}

	var chunk [4]byte
	var chunkSize uint32
	for {
		// Read next chunkID
		err := binary.Read(decoder.src, binary.BigEndian, &chunk)
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

		seeker := decoder.src.(io.Seeker)
		if bytes.Equal(chunk[:], []byte("fmt ")) {
			// seek 4 bytes back because riffChunkFmt reads the chunkSize again
			if _, err = seeker.Seek(-4, os.SEEK_CUR); err != nil {
				return err
			}

			if err = binary.Read(decoder.src, binary.LittleEndian, decoder.info); err != nil {
				return err
			}
			if decoder.info.LengthOfHeader > 16 { // canonical format if chunklen == 16
				// Skip extra params
				if _, err = seeker.Seek(int64(decoder.info.LengthOfHeader-16), os.SEEK_CUR); err != nil {
					return err
				}
			}

			// Is audio supported ?
			if decoder.info.AudioFormat != 1 {
				return fmt.Errorf("Audio Format not supported")
			}
		} else if bytes.Equal(chunk[:], []byte("data")) {
			size, _ := seeker.Seek(0, os.SEEK_CUR)
			decoder.firstSamplePos = uint32(size)
			decoder.dataSize = int32(chunkSize)
			break
		} else {
			if _, err = seeker.Seek(int64(chunkSize), os.SEEK_CUR); err != nil {
				return err
			}
		}
	}

	return nil
}

func (decoder *waveDecoder) Seek(offset int64, whence int) (int64, error) {
	seeker := decoder.src.(io.Seeker)
	switch whence {
	case io.SeekStart:
		offset += int64(decoder.firstSamplePos)
	case io.SeekCurrent:
	case io.SeekEnd:
		offset += int64(decoder.firstSamplePos) + int64(decoder.dataSize)
		whence = io.SeekStart
	}
	return seeker.Seek(offset, whence)
}
