package gfx

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func u16FromByte(b []byte) []uint16 {
	f := make([]uint16, len(b)/2)
	for i := 0; i < len(b)/2; i++ {
		f[i] = uint16(b[2*i+0]) | uint16(b[2*i+1])<<8
	}
	return f
}

func u16Bytes(values ...uint16) []byte {
	b := make([]byte, 2*len(values))
	for i, u := range values {
		b[2*i+0] = byte(u >> 0)
		b[2*i+1] = byte(u >> 8)
	}
	return b
}

func u32FromByte(b []byte) []uint32 {
	f := make([]uint32, len(b)/4)
	for i := 0; i < len(b)/4; i++ {
		f[i] = uint32(b[4*i+0]) | uint32(b[4*i+1])<<8 | uint32(b[4*i+2])<<16 | uint32(b[4*i+3])<<24
	}
	return f
}

func u32Bytes(values ...uint32) []byte {
	b := make([]byte, 4*len(values))
	for i, u := range values {
		b[4*i+0] = byte(u >> 0)
		b[4*i+1] = byte(u >> 8)
		b[4*i+2] = byte(u >> 16)
		b[4*i+3] = byte(u >> 24)
	}
	return b
}

func f32FromByte(b []byte) []float32 {
	f := make([]float32, len(b)/4)
	for i := 0; i < len(b)/4; i++ {
		f[i] = math.Float32frombits(uint32(b[4*i+0]) | uint32(b[4*i+1])<<8 | uint32(b[4*i+2])<<16 | uint32(b[4*i+3])<<24)
	}
	return f
}

func f32Bytes(values ...float32) []byte {
	b := make([]byte, 4*len(values))
	for i, v := range values {
		u := math.Float32bits(v)
		b[4*i+0] = byte(u >> 0)
		b[4*i+1] = byte(u >> 8)
		b[4*i+2] = byte(u >> 16)
		b[4*i+3] = byte(u >> 24)
	}
	return b
}

func v2Bytes(values ...mgl32.Vec2) []byte {
	b := make([]byte, 8*len(values))
	for i, v := range values {
		x := math.Float32bits(v[0])
		b[8*i+0] = byte(x >> 0)
		b[8*i+1] = byte(x >> 8)
		b[8*i+2] = byte(x >> 16)
		b[8*i+3] = byte(x >> 24)

		y := math.Float32bits(v[1])
		b[8*i+4] = byte(y >> 0)
		b[8*i+5] = byte(y >> 8)
		b[8*i+6] = byte(y >> 16)
		b[8*i+7] = byte(y >> 24)
	}
	return b
}

func colorBytes(values ...Color) []byte {
	by := make([]byte, 16*len(values))
	for i, v := range values {
		r := math.Float32bits(v[0])
		by[16*i+0] = byte(r >> 0)
		by[16*i+1] = byte(r >> 8)
		by[16*i+2] = byte(r >> 16)
		by[16*i+3] = byte(r >> 24)

		g := math.Float32bits(v[1])
		by[16*i+4] = byte(g >> 0)
		by[16*i+5] = byte(g >> 8)
		by[16*i+6] = byte(g >> 16)
		by[16*i+7] = byte(g >> 24)

		b := math.Float32bits(v[2])
		by[16*i+8] = byte(b >> 0)
		by[16*i+9] = byte(b >> 8)
		by[16*i+10] = byte(b >> 16)
		by[16*i+11] = byte(b >> 24)

		a := math.Float32bits(v[3])
		by[16*i+12] = byte(a >> 0)
		by[16*i+13] = byte(a >> 8)
		by[16*i+14] = byte(a >> 16)
		by[16*i+15] = byte(a >> 24)
	}
	return by
}
