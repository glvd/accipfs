package mdns

import (
	"fmt"
	"github.com/glvd/accipfs/config"
	"go.uber.org/atomic"
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
	Zone              string
	NetInterface      *net.Interface
	IPV4Addr          *net.UDPAddr
	IPV6Addr          *net.UDPAddr
	LogEmptyResponses bool
}

// OptionConfigFunc ...
type OptionConfigFunc func(cfg *OptionConfig)

// MulticastDNS ...
type MulticastDNS struct {
	cfg *OptionConfig
}

// Server ...
func (dns *MulticastDNS) Server() (s Server, err error) {
	// Create the listeners
	conn := make([]*net.UDPConn, ipmax)
	var udp4Err error
	if dns.cfg.IPV4Addr != nil {
		conn[udp4], udp4Err = net.ListenMulticastUDP("udp4", dns.cfg.NetInterface, dns.cfg.IPV4Addr)
	}
	var udp6Err error
	if dns.cfg.IPV6Addr != nil {
		conn[udp6], udp6Err = net.ListenMulticastUDP("udp6", dns.cfg.NetInterface, dns.cfg.IPV6Addr)
	}

	// Check if we have any listener
	if udp4Err != nil && udp6Err != nil {
		return nil, fmt.Errorf("no multicast listeners could be started")
	}

	return &server{
		cfg:  dns.cfg,
		conn: conn,
		stop: atomic.NewBool(false),
	}, nil
}

// Client ...
func (dns *MulticastDNS) Client() (c Client, err error) {
	// Create the listeners
	conn := make([]*net.UDPConn, ipmax)
	if dns.cfg.IPV4Addr != nil {
		conn[udp4], err = net.ListenMulticastUDP("udp4", nil, dns.cfg.IPV4Addr)
		if err != nil {
			return nil, err
		}
	}

	if dns.cfg.IPV6Addr != nil {
		conn[udp6], err = net.ListenMulticastUDP("udp6", nil, dns.cfg.IPV6Addr)
		if err != nil {
			return nil, err
		}
	}

	// Check if we have any listener
	if conn[udp4] == nil && conn[udp6] == nil {
		return nil, fmt.Errorf("no multicast listeners could be started")
	}

	return &client{
		cfg:  dns.cfg,
		conn: conn,
		stop: atomic.NewBool(false),
	}, nil
}

// New ...
func New(cfg *config.Config, opts ...OptionConfigFunc) (mdns *MulticastDNS, err error) {
	optionConfig := defaultConfig(cfg)
	for _, op := range opts {
		op(optionConfig)
	}

	return &MulticastDNS{
		cfg: optionConfig,
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
