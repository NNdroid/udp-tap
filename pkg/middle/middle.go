package middle

import (
	"golang.org/x/net/context"
	"udp-tap/pkg/log"
	"udp-tap/pkg/srv"
	"udp-tap/pkg/tap"
)

type Middle struct {
	srv       *srv.LocalSrv
	tunDevice tap.Tun
	ctx       context.Context
}

func NewMiddle(srv *srv.LocalSrv, tunDevice tap.Tun, ctx context.Context) *Middle {
	return &Middle{
		srv:       srv,
		tunDevice: tunDevice,
		ctx:       ctx,
	}
}

func (m *Middle) loopReceive() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case data := <-m.tunDevice.ReceiveQueue():
			log.Logger().Debugf("tun: read %s\v", data.String())
			m.srv.TransportQueue() <- data
		}
	}
}

func (m *Middle) loopTransport() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case data := <-m.srv.ReceiveQueue():
			log.Logger().Debugf("srv: write %s\v", data.String())
			m.tunDevice.TransportQueue() <- data
		}
	}
}

func (m *Middle) Run() {
	go m.loopReceive()
	go m.loopTransport()
}
