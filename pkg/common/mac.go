package common

import "net"

// GetMACAddress returns the MAC address of the specified interface.
func GetMACAddress(ifaceName string) (net.HardwareAddr, error) {
	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return nil, err
	}
	return iface.HardwareAddr, nil
}

// IsUnicast returns true if the MAC address is a unicast address.
func IsUnicast(mac net.HardwareAddr) bool {
	return (mac[0]&0x1) == 0 && !IsMulticast(mac) && !IsBroadcast(mac)
}

// IsMulticast returns true if the MAC address is a multicast address.
func IsMulticast(mac net.HardwareAddr) bool {
	return (mac[0] | 0xfe) == 0xff
}

// IsBroadcast returns true if the MAC address is a broadcast address.
func IsBroadcast(mac net.HardwareAddr) bool {
	return mac[0] == 0xff && mac[1] == 0xff && mac[2] == 0xff && mac[3] == 0xff && mac[4] == 0xff && mac[5] == 0xff
}
