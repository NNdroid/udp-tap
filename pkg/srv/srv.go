package srv

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"golang.org/x/net/context"
	"net"
	"net/netip"
	"udp-tap/pkg/frame"
	"udp-tap/pkg/log"
	"udp-tap/pkg/xchan"
)

type LocalSrv struct {
	Listen         *net.UDPAddr
	Peer           *net.UDPAddr
	conn           *net.UDPConn
	listenIPPort   netip.AddrPort
	peerIPPort     netip.AddrPort
	receiveQueue   <-chan frame.Frame
	receiveChan    chan<- frame.Frame
	transportQueue chan<- frame.Frame
	transportChan  <-chan frame.Frame
	ctx            context.Context
	cancel         context.CancelFunc

	key    []byte
	nonce  []byte
	aesgcm cipher.AEAD
}

func NewLocalSrv(listen, peer string, parentCtx context.Context) (*LocalSrv, error) {
	listenUdpAddr, err := net.ResolveUDPAddr("udp", listen)
	if err != nil {
		return nil, err
	}
	peerUdpAddr, err := net.ResolveUDPAddr("udp", peer)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(parentCtx)
	receiveUnboundedChan := xchan.NewUnboundedChan[frame.Frame](ctx, 3000)
	transportUnboundedChan := xchan.NewUnboundedChan[frame.Frame](ctx, 3000)
	loSrv := &LocalSrv{
		Listen:         listenUdpAddr,
		Peer:           peerUdpAddr,
		listenIPPort:   listenUdpAddr.AddrPort(),
		peerIPPort:     peerUdpAddr.AddrPort(),
		receiveQueue:   receiveUnboundedChan.Out,
		receiveChan:    receiveUnboundedChan.In,
		transportQueue: transportUnboundedChan.In,
		transportChan:  transportUnboundedChan.Out,
		ctx:            ctx,
		cancel:         cancel,
	}
	loSrv.key, _ = hex.DecodeString("6368616e676520746869732070617373776f726420746f206120736563726574")
	loSrv.nonce, _ = hex.DecodeString("64a9433eae7ccceee2fc0eda")
	block, err := aes.NewCipher(loSrv.key)
	if err != nil {
		return nil, err
	}
	loSrv.aesgcm, err = cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return loSrv, nil
}

func (s *LocalSrv) loopReceive() {
	var buf [65535]byte
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			n, addr, err := s.conn.ReadFromUDP(buf[:])
			if err != nil {
				continue
			}
			if addr.AddrPort().Compare(s.peerIPPort) == 0 {
				log.Logger().Debugf("receive from peer %s, %d, %s", addr.String(), n, hex.EncodeToString(buf[:n]))
				plainData, err := s.aesgcm.Open(nil, s.nonce, buf[:n], nil)
				if err != nil {
					log.Logger().Errorf("decrypt failed, %s", err)
					continue
				}
				frm := frame.Parse(plainData)
				if frm == nil {
					log.Logger().Errorf("parse frame failed")
					continue
				}
				s.receiveChan <- frm
			}
		}
	}
}

func (s *LocalSrv) loopTransport() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case data := <-s.transportChan:
			log.Logger().Debugf("send to peer %s, %d, %s", s.Peer.String(), len(data.Bytes()), hex.EncodeToString(data.Bytes()))
			cipherData := s.aesgcm.Seal(nil, s.nonce, data.Bytes(), nil)
			_, err := s.conn.WriteToUDP(cipherData, s.Peer)
			if err != nil {
				log.Logger().Errorf("send to peer %s, %s", s.Peer.String(), err)
			}
		}
	}
}

func (s *LocalSrv) Start() error {
	conn, err := net.ListenUDP("udp", s.Listen)
	if err != nil {
		return err
	}
	s.conn = conn
	go s.loopReceive()
	go s.loopTransport()
	return nil
}

func (s *LocalSrv) Stop() error {
	s.cancel()
	return s.conn.Close()
}

// ReceiveQueue read data from tun device.
func (s *LocalSrv) ReceiveQueue() <-chan frame.Frame {
	return s.receiveQueue
}

// TransportQueue write data to tun device.
func (s *LocalSrv) TransportQueue() chan<- frame.Frame {
	return s.transportQueue
}
