package lib

import (
	"net"
	"os"
)

func GetIpAddresses() (Ips []string, err error) {
	var (
		addrs []net.Addr
	)
	addrs, err = net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil || ipNet.IP.To16 != nil {
					Ips = append(Ips, ipNet.IP.String())
				}
			}
		}
	}
	return Ips, err
}

func GetIpAddresses2() (Ips []string, err error) {
	var (
		hostname string
		addrs    []string
	)
	hostname, err = os.Hostname()
	if err == nil {
		addrs, err = net.LookupHost(hostname)
		if err == nil {
			for _, addr := range addrs {
				Ips = append(Ips, addr)
			}
		}
	}

	return Ips, err
}
