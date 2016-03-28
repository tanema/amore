package webal

type Buffer struct{}

func (b Buffer) Frequency() int32                                  { return 0 }
func (b Buffer) Bits() int32                                       { return 0 }
func (b Buffer) Channels() int32                                   { return 0 }
func (b Buffer) Size() int32                                       { return 0 }
func (b Buffer) BufferData(format uint32, data []byte, freq int32) {}
func (b Buffer) Valid() bool                                       { return false }
