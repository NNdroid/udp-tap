package frame

import (
	"udp-tap/pkg/log"
)

type Frame interface {
	Bytes() []byte
	Type() Type
	Raw() []byte
	String() string
}

type Type uint8

const (
	TypeIP Type = 1 + iota
	TypeCommand
)

func Parse(data []byte) Frame {
	t := Type(data[0])
	switch t {
	case TypeIP:
		return NewIPFrame(data)
	case TypeCommand:
		//return ParseCommand(data)
	default:
		log.Logger().Errorf("unknown frame type %d", t)
	}
	return nil
}
