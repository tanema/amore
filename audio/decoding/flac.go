package decoding

import (
	"io"

	"github.com/eaburns/flac"
)

func newFlacDecoder(src io.ReadCloser) (*Decoder, error) {
	r, err := flac.NewDecoder(src)
	if err != nil {
		return nil, err
	}

	totalBytes := int32(r.TotalSamples * int64(r.NChannels) * int64(r.BitsPerSample/8))
	d := &flacReader{
		data:       make([]byte, totalBytes),
		totalBytes: int(totalBytes),
		source:     src,
		decoder:    r,
	}

	return newDecoder(
		src,
		d,
		int16(r.NChannels),
		int32(r.SampleRate),
		int16(r.BitsPerSample),
		totalBytes,
	), nil
}

type flacReader struct {
	data       []byte
	totalBytes int
	readBytes  int
	pos        int
	source     io.Closer
	decoder    *flac.Decoder
}

func (d *flacReader) readUntil(pos int) error {
	for d.readBytes < pos {
		data, err := d.decoder.Next()
		if len(data) > 0 {
			p := d.readBytes
			for i := 0; i < len(data); i++ {
				d.data[p+i] = data[i]
			}
			d.readBytes += len(data)
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

func (d *flacReader) Read(b []byte) (int, error) {
	left := d.totalBytes - d.pos
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
	if d.pos == d.totalBytes {
		return left, io.EOF
	}
	return left, nil
}

func (d *flacReader) Seek(offset int64, whence int) (int64, error) {
	next := int64(0)
	switch whence {
	case io.SeekStart:
		next = offset
	case io.SeekCurrent:
		next = int64(d.pos) + offset
	case io.SeekEnd:
		next = int64(d.totalBytes) + offset
	}
	d.pos = int(next)
	if err := d.readUntil(d.pos); err != nil {
		return 0, err
	}
	return next, nil
}
