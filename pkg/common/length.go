package common

import "encoding/binary"

const DataLengthHeaderSize = 4

func ReadUint32(data []byte) uint32 {
	return binary.LittleEndian.Uint32(data)
}

func WriteUint32(data []byte, value uint32) {
	binary.LittleEndian.PutUint32(data, value)
}
