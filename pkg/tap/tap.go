package tap

import (
	"net/netip"
	"udp-tap/pkg/frame"
)

type Tun interface {
	Open() error
	Close() error
	//Read(b []byte) (int, error)
	//Write(b []byte) (int, error)
	Name() string
	Ipv4Address() netip.Addr
	Ipv6Address() netip.Addr
	Start()
	MTU() int
	ReceiveQueue() <-chan frame.Frame
	TransportQueue() chan<- frame.Frame
}
