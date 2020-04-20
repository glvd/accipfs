package mdns

import (
	"fmt"
	"github.com/glvd/accipfs/config"
	"net"
)

const (
	mdnsIPV4Addr = "224.0.0.251"
	mdnsIPV6Addr = "FF02::FB"
	mdnsPort     = 5353
)

const (
	udp4  = 0
	udp6  = 1
	ipmax = 2
)

// OptionConfig ...
type OptionConfig struct {
	IPV4Addr *net.UDPAddr
	IPV6Addr *net.UDPAddr
}

// OptionConfigFunc ...
type OptionConfigFunc func(cfg *OptionConfig)

// MulticastDNS ...
type MulticastDNS struct {
	conn []*net.UDPConn
}

// New ...
func New(cfg *config.Config, opts ...OptionConfigFunc) (mdns *MulticastDNS, err error) {
	optionConfig := defaultConfig(cfg)
	for _, op := range opts {
		op(optionConfig)
	}

	// Create the listeners
	conn := make([]*net.UDPConn, ipmax)
	if optionConfig.IPV4Addr != nil {
		conn[udp4], err = net.ListenMulticastUDP("udp4", config.Iface, optionConfig.IPV4Addr)
		if err != nil {
			return nil, err
		}
	}
	if optionConfig.IPV6Addr != nil {
		conn[udp6], err = net.ListenMulticastUDP("udp6", config.Iface, optionConfig.IPV6Addr)
		if err != nil {
			return nil, err
		}
	}

	// Check if we have any listener
	if conn[udp4] == nil && conn[udp6] == nil {
		return nil, fmt.Errorf("no multicast listeners could be started")
	}
	return &MulticastDNS{
		conn: conn,
	}, nil
}

func defaultConfig(cfg *config.Config) *OptionConfig {
	ipv4Addr := &net.UDPAddr{
		IP:   net.ParseIP(mdnsIPV4Addr),
		Port: mdnsPort,
	}
	ipv6Addr := &net.UDPAddr{
		IP:   net.ParseIP(mdnsIPV6Addr),
		Port: mdnsPort,
	}
	return &OptionConfig{
		IPV4Addr: ipv4Addr,
		IPV6Addr: ipv6Addr,
	}
}
