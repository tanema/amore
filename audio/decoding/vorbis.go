package decoding

import (
	"io"

	"github.com/jfreymuth/oggvorbis"
)

func newVorbisDecoder(src io.ReadCloser) (*Decoder, error) {
	r, err := oggvorbis.NewReader(src)
	if err != nil {
		return nil, err
	}

	d := &vorbisReader{
		data:       make([]float32, r.Length()*2),
		totalBytes: int(r.Length()) * 4,
		source:     src,
		decoder:    r,
	}

	return newDecoder(
		src,
		d,
		int16(r.Channels()),
		int32(r.SampleRate()),
		16,
		int32(d.totalBytes),
	), nil
}

type vorbisReader struct {
	data       []float32
	totalBytes int
	readBytes  int
	pos        int
	source     io.Closer
	decoder    *oggvorbis.Reader
}

func (d *vorbisReader) readUntil(pos int) error {
	buffer := make([]float32, 8192)
	for d.readBytes < pos {
		n, err := d.decoder.Read(buffer)
		if n > 0 {
			if d.readBytes+n*2 > d.totalBytes {
				n = (d.totalBytes - d.readBytes) / 2
			}
			p := d.readBytes / 2
			for i := 0; i < n; i++ {
				d.data[p+i] = buffer[i]
			}
			d.readBytes += n * 2
		}
		if err == io.EOF {
			if err := d.source.Close(); err != nil {
				return err
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *vorbisReader) Read(b []byte) (int, error) {
	left := d.totalBytes - d.pos
	if left > len(b) {
		left = len(b)
	}
	if left < 0 {
		return 0, io.EOF
	}
	left = left / 2 * 2
	if err := d.readUntil(d.pos + left); err != nil {
		return 0, err
	}
	for i := 0; i < left/2; i++ {
		f := d.data[d.pos/2+i]
		s := int16(f * (1<<15 - 1))
		b[2*i] = uint8(s)
		b[2*i+1] = uint8(s >> 8)
	}
	d.pos += left
	if d.pos == d.totalBytes {
		return left, io.EOF
	}
	return left, nil
}

func (d *vorbisReader) Seek(offset int64, whence int) (int64, error) {
	next := int64(0)
	switch whence {
	case io.SeekStart:
		next = offset
	case io.SeekCurrent:
		next = int64(d.pos) + offset
	case io.SeekEnd:
		next = int64(d.totalBytes) + offset
	}
	// pos should be always even
	next = next / 2 * 2
	d.pos = int(next)
	if err := d.readUntil(d.pos); err != nil {
		return 0, err
	}
	return next, nil
}
