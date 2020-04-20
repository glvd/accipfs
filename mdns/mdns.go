package mdns

import (
	"fmt"
	"github.com/glvd/accipfs/config"
	"go.uber.org/atomic"
	"net"
	"os"
)

const (
	mdnsIPV4Addr = "224.0.0.251"
	mdnsIPV6Addr = "FF02::FB"
	mdnsPort     = 5353
)
const defaultTTL = 120

const (
	udp4  = 0
	udp6  = 1
	ipmax = 2
)

// OptionConfig ...
type OptionConfig struct {
	//Zone              string
	NetInterface      *net.Interface
	IPV4Addr          *net.UDPAddr
	IPV6Addr          *net.UDPAddr
	LogEmptyResponses bool
	HostName          string
	InstanceAddr      string
	ServiceAddr       string
	EnumAddr          string
	Port              uint16
	TTL               uint32
	TXT               []string
	IPs               []net.IP
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
	uniConn := make([]*net.UDPConn, ipmax)
	var uudp4Err error
	uniConn[udp4], uudp4Err = net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})

	var uudp6Err error
	uniConn[udp6], uudp6Err = net.ListenUDP("udp6", &net.UDPAddr{IP: net.IPv6zero, Port: 0})
	if uudp4Err == nil && uudp6Err == nil {
		logE("failed to bind to port", "uudp6Err", uudp4Err, "uudp4Err", uudp6Err)
		return nil, fmt.Errorf("failed to bind to any unicast udp port")
	}
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
		logE("failed to bind to port", "udp6Err", udp6Err, "udp4Err", udp4Err)
		return nil, fmt.Errorf("failed to bind to any multicast udp port")
	}

	return &client{
		cfg:     dns.cfg,
		uniConn: uniConn,
		conn:    conn,
		stop:    atomic.NewBool(false),
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

	hostName, _ := os.Hostname()
	hostName = fmt.Sprintf("%s.", hostName)
	service := "_http._tcp"
	instance := "hostname"
	domain := "local."
	return &OptionConfig{
		//Zone:              "",
		NetInterface:      nil,
		IPV4Addr:          ipv4Addr,
		IPV6Addr:          ipv6Addr,
		LogEmptyResponses: false,
		HostName:          hostName,
		InstanceAddr:      instanceAddr(instance, service, domain),
		ServiceAddr:       serviceAddr(service, domain),
		EnumAddr:          enumAddr(domain),
		Port:              80,
		TTL:               defaultTTL,
		TXT:               nil,
		IPs:               nil,
	}
}