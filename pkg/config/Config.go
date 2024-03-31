package config

type Config struct {
	//tap设备名称
	DeviceName string `json:"device_name"`
	//IPV4地址
	IPv4Address string `json:"ipv4_address"`
	//IPv6地址
	IPv6Address string `json:"ipv6_address"`
	//MTU
	MTU int `json:"mtu"`
	//表示本端的UDP地址和端口
	ListenAddress string `json:"listen"`
	//加密密匙
	Key string `json:"key"`
	//是否debug
	Verbose bool `json:"verbose"`
	//表示对端的UDP地址和端口
	PeerAddress string `json:"peer"`
}
