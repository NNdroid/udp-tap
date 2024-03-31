//go:build linux

package main

import (
	"context"
	"flag"
	"udp-tap/pkg/common"
	cfg "udp-tap/pkg/config"
	"udp-tap/pkg/log"
	"udp-tap/pkg/middle"
	"udp-tap/pkg/srv"
	"udp-tap/pkg/tap"
)

var config cfg.Config

func init() {
	flag.StringVar(&config.DeviceName, "d", "tun1", "device name")
	flag.StringVar(&config.IPv4Address, "c4", "10.0.0.1/24", "ipv4 address")
	flag.StringVar(&config.IPv6Address, "c6", "fd99:10::1/64", "ipv6 address")
	flag.IntVar(&config.MTU, "m", 1280, "mtu")
	flag.StringVar(&config.ListenAddress, "l", ":8000", "listen address")
	flag.StringVar(&config.Key, "k", "default", "auth key")
	flag.BoolVar(&config.Verbose, "v", false, "log debug")
	flag.StringVar(&config.PeerAddress, "peer", "192.168.177.26:8000", "peer address")
	flag.Parse()
}

func main() {
	common.PrintVersion()
	close := make(chan struct{})
	ctx := context.TODO()
	log.SetVerbose(config.Verbose)
	tunDevice, err := tap.New(config.DeviceName, config.IPv4Address, config.IPv6Address, config.MTU, ctx)
	if err != nil {
		log.Logger().Errorf("create tun device error: %v", err)
		return
	}
	log.Logger().Infoln("tun created")
	tunDevice.SetOffset(1) //增加一个offset为FrameType
	err = tunDevice.Open()
	if err != nil {
		log.Logger().Errorf("open tun device error: %v", err)
		return
	}
	log.Logger().Infoln("tun opened")
	tunDevice.Start()
	log.Logger().Infoln("tun started")

	peer, err := srv.NewLocalSrv(config.ListenAddress, config.PeerAddress, ctx)
	if err != nil {
		log.Logger().Errorf("create srv error: %v", err)
		return
	}
	peer.Start()
	log.Logger().Infoln("peer started")

	ware := middle.NewMiddle(peer, tunDevice, ctx)
	ware.Run()
	<-close
}
