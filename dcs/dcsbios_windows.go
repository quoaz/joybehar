package dcs

import "net"

func multicastInterface() (*net.Interface, error) {
	iface, err := net.InterfaceByName("Ethernet")
	return iface, err
}
