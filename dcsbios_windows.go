package main

import "net"

func multicastInterface() (*net.Interface, error) {
	iface, err := net.InterfaceByName("Wi-Fi")
	return iface, err
}
