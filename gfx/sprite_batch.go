package gfx

type SpriteBatch struct {
}

func NewSpriteBatch(text iTexture, size int) *SpriteBatch {
	return NewSpriteBatchExt(text, size, USAGE_DYNAMIC)
}

func NewSpriteBatchExt(texture iTexture, size int, usage Usage) *SpriteBatch {
	return &SpriteBatch{}
}
