package openal

// A buffer represents a chunk of PCM audio data that could be buffered to an audio
// source. A single buffer could be shared between multiple sources.
type Buffer uint32

// Frequency returns the frequency of the buffer data in Hertz (Hz).
func (b Buffer) Frequency() int32 {
	return alGetBufferi(b, paramFreq)
}

// Bits return the number of bits used to represent a sample.
func (b Buffer) Bits() int32 {
	return alGetBufferi(b, paramBits)
}

// Channels return the number of the audio channels.
func (b Buffer) Channels() int32 {
	return alGetBufferi(b, paramChannels)
}

// Size returns the size of the data.
func (b Buffer) Size() int32 {
	return alGetBufferi(b, paramSize)
}

// BufferData buffers PCM data to the current buffer.
func (b Buffer) BufferData(format uint32, data []byte, freq int32) {
	alBufferData(b, format, data, freq)
}

// Valid returns true if the buffer exists and is valid.
func (b Buffer) Valid() bool {
	return alIsBuffer(b)
}
