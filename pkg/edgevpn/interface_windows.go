//go:build windows
// +build windows

package edgevpn

import (
	"net"
	"os/exec"
	"fmt"
	"log"

	"github.com/songgao/water"
)

func (e *EdgeVPN) prepareInterface() error {
	err := netsh("interface", "ip", "set", "address", "name=", e.config.InterfaceName, "static", e.config.InterfaceAddress)
	if err != nil {
		log.Println(err)
	}
	err = netsh("interface", "ipv4", "set", "subinterface", e.config.InterfaceName, "mtu=", fmt.Sprintf("%d", e.config.InterfaceMTU))
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (e *EdgeVPN) createInterface() (*water.Interface, error) {
	// TUN on Windows requires address and network to be set on device creation stage
	// We also set network to 0.0.0.0/0 so we able to reach networks behind the node
	// https://github.com/songgao/water/blob/master/params_windows.go
	// https://gitlab.com/openconnect/openconnect/-/blob/master/tun-win32.c
	ip, _, err := net.ParseCIDR(e.config.InterfaceAddress)
	if err != nil {
		return nil, err
	}
	network := net.IPNet{
		IP:   ip,
		Mask: net.IPv4Mask(0,0,0,0),
	}
	config := water.Config{
		DeviceType: e.config.DeviceType,
		PlatformSpecificParams: water.PlatformSpecificParams{
			ComponentID:   "tap0901",
			InterfaceName: e.config.InterfaceName,
			Network:       network.String(),
		},
	}

	return water.New(config)
}

func netsh(args ...string) (err error) {
	cmd := exec.Command("netsh", args...)
	err = cmd.Run()
	return
}
