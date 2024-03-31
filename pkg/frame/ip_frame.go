package frame

import (
	"encoding/hex"
	"fmt"
	"net"
)

type IPFrame struct {
	DestinationAddress net.HardwareAddr
	SourceAddress      net.HardwareAddr
	EtherType          uint16
	Payload            []byte
	//Except for the first byte, which is Type, the remaining bytes are Data.
	RAWData []byte
}

func (f *IPFrame) Bytes() []byte {
	return f.RAWData
}

func (f *IPFrame) Type() Type {
	return TypeIP
}

func (f *IPFrame) Raw() []byte {
	return f.Payload
}

func NewIPFrame(data []byte) *IPFrame {
	data[0] = byte(TypeIP)
	return &IPFrame{ //first byte
		DestinationAddress: net.HardwareAddr{data[1], data[2], data[3], data[4], data[5], data[6]},    //dst mac
		SourceAddress:      net.HardwareAddr{data[7], data[8], data[9], data[10], data[11], data[12]}, //src mac
		EtherType:          uint16(data[13])<<8 | uint16(data[14]),                                    //ethertype (0x0800)
		Payload:            data[15:],                                                                 //all bytes except for the first byte, which is Type.
		RAWData:            data,
	}
}

func (f *IPFrame) String() string {
	return fmt.Sprintf("Frame{type: %d, dst: %s, src: %s, ethertype: %s, payload: %s}", f.Type(), f.DestinationAddress, f.SourceAddress, fmt.Sprintf("%04X", f.EtherType), hex.EncodeToString(f.Payload))
}
