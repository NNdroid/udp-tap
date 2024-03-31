package tap

import (
	"context"
	"encoding/hex"
	"fmt"
	"golang.org/x/sys/unix"
	"net"
	"net/netip"
	"os"
	"syscall"
	"udp-tap/pkg/common"
	"udp-tap/pkg/frame"
	"udp-tap/pkg/log"
	"udp-tap/pkg/xchan"
	"unsafe"
)

type NativeTun struct {
	name              string
	file              *os.File
	ipv4AddressString string
	ipv4Address       netip.Addr
	ipv4Cidr          *net.IPNet
	ipv6AddressString string
	ipv6Address       netip.Addr
	ipv6Cidr          *net.IPNet
	mtu               int
	ctx               context.Context
	cancel            context.CancelFunc
	receiveQueue      <-chan frame.Frame
	receiveChan       chan<- frame.Frame
	transportQueue    chan<- frame.Frame
	transportChan     <-chan frame.Frame
	offset            int
}

// New
// name:  tun name, eg:  tun0, tun1, tun2, etc.
// ipv4Address:  ipv4 address, eg: 10.0.0.1/24, 10.0.0.1/32, 10.0.0.1/128, 10.0.0.1/0, etc.
// ipv6Address:  ipv6 address,  eg: 2001:db8::1/64, 2001:db8::1/128, 2001:db8::1/0, etc.
func New(name string, ipv4Address string, ipv6Address string, mtu int, ctx context.Context) (*NativeTun, error) {
	var ip4, ip6 net.IP
	var cidr4, cidr6 *net.IPNet
	var err error
	ip4, cidr4, err = net.ParseCIDR(ipv4Address)
	if err != nil {
		return nil, err
	}
	ip6, cidr6, err = net.ParseCIDR(ipv6Address)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	receiveUnboundedChan := xchan.NewUnboundedChan[frame.Frame](ctx, 3000)
	transportUnboundedChan := xchan.NewUnboundedChan[frame.Frame](ctx, 3000)
	var obj = &NativeTun{
		name:              name,
		ipv4AddressString: ipv4Address,
		ipv4Address:       netip.AddrFrom4([4]byte(ip4.To4())),
		ipv4Cidr:          cidr4,
		ipv6AddressString: ipv6Address,
		ipv6Address:       netip.AddrFrom16([16]byte(ip6.To16())),
		ipv6Cidr:          cidr6,
		mtu:               mtu,
		ctx:               ctx,
		cancel:            cancel,
		receiveQueue:      receiveUnboundedChan.Out,
		receiveChan:       receiveUnboundedChan.In,
		transportQueue:    transportUnboundedChan.In,
		transportChan:     transportUnboundedChan.Out,
		offset:            0,
	}
	return obj, nil
}

func (tun *NativeTun) SetOffset(offset int) {
	tun.offset = offset
}

func (tun *NativeTun) Start() {
	go tun.loopReceive()
	go tun.loopTransport()
}

func (tun *NativeTun) loopReceive() {
	var buf = make([]byte, 65535) //TODO 自动大小 +tun.offset
	var buffer []byte
	for {
		//log.Logger().Infof("try to read tun device, 2222\n")
		select {
		case <-tun.ctx.Done():
			log.Logger().Infof("loopReceive: active close")
			return
		default:
			//log.Logger().Infof("try to read tun device\n")
			n, err := tun.file.Read(buf[tun.offset:])
			if err != nil {
				log.Logger().Infof("loopReceive: err to write data: %s\n", err)
				return
			}
			if n == 0 {
				continue
			}
			buffer = make([]byte, tun.offset+n)
			copy(buffer, buf[:tun.offset+n])
			log.Logger().Infof("tun: read %s\n", hex.EncodeToString(buffer))
			tun.receiveChan <- frame.NewIPFrame(buffer)
		}
	}
}

func (tun *NativeTun) loopTransport() {
	for {
		select {
		case <-tun.ctx.Done():
			log.Logger().Infoln("loopTransport: active close")
			return
		case data := <-tun.transportChan:
			log.Logger().Infof("tun: write %s\n", hex.EncodeToString(data.Bytes()[tun.offset:]))
			_, err := tun.file.Write(data.Bytes()[tun.offset:]) //RAWData[tun.offset:]
			if err != nil {
				log.Logger().Infof("err to write data: %s\n", err)
				return
			}
		}
	}
}

// Open
// Open tun device, and set tun device name, ipv4 address and ipv6 address.
func (tun *NativeTun) Open() error {
	var err error
	tun.file, err = os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return err
	}
	//if err := syscall.SetNonblock(int(tun.file.Fd()), true); err != nil {
	//	return err
	//}
	ifr, err := unix.NewIfreq(tun.name)
	if err != nil {
		return err
	}
	ifr.SetUint16(unix.IFF_TAP | unix.IFF_NO_PI) //| unix.IFF_NO_PI
	var errno syscall.Errno
	_, _, errno = unix.Syscall(unix.SYS_IOCTL, tun.file.Fd(), uintptr(unix.TUNSETIFF), uintptr(unsafe.Pointer(ifr)))
	if errno != 0 {
		return fmt.Errorf("err to syscall: %s", errno)
	}
	//_, _, errno = unix.Syscall(unix.SYS_IOCTL, f.Fd(), unix.TUNSETPERSIST, uintptr(1))
	//if errno != 0 {
	//	return fmt.Errorf("err to syscall: %s", errno)
	//}
	err = tun.setMTU()
	if err != nil {
		return err
	}
	err = tun.setIpv4()
	if err != nil {
		return err
	}
	err = tun.setIpv6()
	if err != nil {
		return err
	}
	err = tun.deviceUp()
	if err != nil {
		return err
	}
	return nil
}

// deviceUp set tun device up
func (tun *NativeTun) deviceUp() error {
	var err error
	err = common.ExecCmd("/bin/ip", "link", "set", "dev", tun.name, "up")
	if err != nil {
		return err
	}
	return nil
}

// deviceDown set tun device down
func (tun *NativeTun) deviceDown() error {
	var err error
	err = common.ExecCmd("/bin/ip", "link", "set", "dev", tun.name, "down")
	if err != nil {
		return err
	}
	return nil
}

// setIpv4 set ipv4 address for tun device.
func (tun *NativeTun) setIpv4() error {
	var err error
	err = common.ExecCmd("/bin/ip", "addr", "add", tun.ipv4AddressString, "dev", tun.name)
	if err != nil {
		return err
	}
	return nil
}

// setIpv6 set ipv6 address for tun device.
func (tun *NativeTun) setIpv6() error {
	var err error
	err = common.ExecCmd("/bin/ip", "-6", "addr", "add", tun.ipv6AddressString, "dev", tun.name)
	if err != nil {
		return err
	}
	return nil
}

func (tun *NativeTun) setMTU() error {
	var err error
	err = common.ExecCmd("/bin/ip", "link", "set", "dev", tun.name, "mtu", fmt.Sprintf("%d", tun.mtu))
	if err != nil {
		return err
	}
	return nil
}

// Close close tun device
func (tun *NativeTun) Close() error {
	err := tun.deviceDown()
	if err != nil {
		return err
	}
	return tun.file.Close()
}

// Read read data from tun device.
//func (tun *NativeTun) Read(b []byte) (int, error) {
//	data := <-tun.receiveQueue
//	return copy(b, data), nil
//}

// ReceiveQueue read data from tun device.
func (tun *NativeTun) ReceiveQueue() <-chan frame.Frame {
	return tun.receiveQueue
}

// TransportQueue write data to tun device.
func (tun *NativeTun) TransportQueue() chan<- frame.Frame {
	return tun.transportQueue
}

// Write write data to tun device.
//func (tun *NativeTun) Write(b []byte) (int, error) {
//	tun.transportQueue <- b
//	return len(b), nil
//}

// Name return tun device name. eg: tun0, tun1, tun2, etc.
func (tun *NativeTun) Name() string {
	return tun.name
}

func (tun *NativeTun) MTU() int {
	return tun.mtu
}

func (tun *NativeTun) Ipv4Address() netip.Addr {
	return tun.ipv4Address
}

func (tun *NativeTun) Ipv6Address() netip.Addr {
	return tun.ipv6Address
}

//func openTun(name string) {
//
//	log.Println("read loop")
//	var buffer = make([]byte, 1500)
//	for {
//		n, err := f.Read(buffer)
//		if err != nil {
//			if err == unix.EAGAIN || err == syscall.EWOULDBLOCK {
//				log.Println("No data available at the moment")
//				continue
//			} else {
//				log.Println("Error reading from file descriptor:", err)
//				return
//			}
//		} else {
//			log.Println("reading from file descriptor:", hex.EncodeToString(buffer[:n]))
//		}
//	}
//}
