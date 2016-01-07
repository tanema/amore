package decoding

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type vorbisDecoder struct {
	decoderBase
}

func (decoder *vorbisDecoder) read() error {
	decoder.bitDepth = 16

	var version, header_type, page_segments int8
	var capture_patter, sn, sequence_number, checksum int32
	var granual_position int64

	binary.Read(decoder.src, binary.LittleEndian, &capture_patter)
	binary.Read(decoder.src, binary.LittleEndian, &version)
	binary.Read(decoder.src, binary.LittleEndian, &header_type)
	binary.Read(decoder.src, binary.LittleEndian, &granual_position)
	binary.Read(decoder.src, binary.LittleEndian, &sn)
	binary.Read(decoder.src, binary.LittleEndian, &sequence_number)
	binary.Read(decoder.src, binary.LittleEndian, &checksum)
	binary.Read(decoder.src, binary.LittleEndian, &page_segments)

	fmt.Println("capture_pattern:", capture_patter)
	fmt.Println("version:", version)
	fmt.Println("header_type:", header_type)
	fmt.Println("granual_position:", granual_position)
	fmt.Println("sn:", sn)
	fmt.Println("sequence_number:", sequence_number)
	fmt.Println("checksum:", checksum)
	fmt.Println("page_segments:", capture_patter)

	return errors.New("not implemented yet")
}
