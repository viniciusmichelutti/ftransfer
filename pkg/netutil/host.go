package netutil

import (
	"net"
	"os"
	"strings"
)

func LocalHostname() string {
	h, err := os.Hostname()
	if err != nil || h == "" {
		return "ftransfer-host"
	}
	// Trim ".local" suffix that Bonjour appends on macOS.
	return strings.TrimSuffix(h, ".local")
}

func LocalIPv4s() []string {
	var out []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return out
	}
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagUp == 0 || ifc.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := ifc.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			ipnet, ok := a.(*net.IPNet)
			if !ok || ipnet.IP.IsLoopback() {
				continue
			}
			ip4 := ipnet.IP.To4()
			if ip4 == nil {
				continue
			}
			out = append(out, ip4.String())
		}
	}
	return out
}
