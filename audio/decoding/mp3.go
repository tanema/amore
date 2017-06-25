package decoding

import (
	"io"

	"github.com/hajimehoshi/go-mp3"
)

func newMp3Decoder(src io.ReadCloser) (*Decoder, error) {
	r, err := mp3.Decode(src)
	if err != nil {
		return nil, err
	}

	d := &mp3Reader{
		data:    []byte{},
		source:  src,
		decoder: r,
	}

	return newDecoder(
		src,
		d,
		2,
		int32(r.SampleRate()),
		16,
		int32(r.Length())*2,
	), nil
}

type mp3Reader struct {
	data      []byte
	readBytes int
	pos       int
	source    io.Closer
	decoder   *mp3.Decoded
}

func (d *mp3Reader) readUntil(pos int) error {
	for len(d.data) <= pos {
		buf := make([]uint8, 8192)
		n, err := d.decoder.Read(buf)
		d.data = append(d.data, buf[:n]...)
		if err != nil {
			if err == io.EOF {
				return io.EOF
			}
			return err
		}
	}
	return nil
}

func (d *mp3Reader) Read(b []byte) (int, error) {
	left := int(d.decoder.Length()) - d.pos
	if left > len(b) {
		left = len(b)
	}
	if left <= 0 {
		return 0, io.EOF
	}
	if err := d.readUntil(d.pos + left); err != nil {
		return 0, err
	}
	copy(b, d.data[d.pos:d.pos+left])
	d.pos += left
	if d.pos == int(d.decoder.Length()) {
		return left, io.EOF
	}
	return left, nil
}

func (d *mp3Reader) Seek(offset int64, whence int) (int64, error) {
	next := int64(0)
	switch whence {
	case io.SeekStart:
		next = offset
	case io.SeekCurrent:
		next = int64(d.pos) + offset
	case io.SeekEnd:
		next = int64(d.decoder.Length()) + offset
	}
	d.pos = int(next)
	if err := d.readUntil(d.pos); err != nil {
		return 0, err
	}
	return next, nil
}
