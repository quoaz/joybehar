package main

import "net"

func multicastInterface() (*net.Interface, error) {
	iface, err := net.InterfaceByName("en0")
	return iface, err
}
