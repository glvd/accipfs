package mdns

import (
	"fmt"
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/goextension/tool"
	"go.uber.org/atomic"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
	"net"
	"os"
)

const (
	mdnsIPv4Addr         = "224.0.0.251"
	mdnsWildcardIPv4Addr = "224.0.0.0"
	mdnsIPv6Addr         = "FF02::FB"
	mdnsWildcardIPv6Addr = "FF02::"
	mdnsPort             = 5353
)
const defaultTTL = 120
const defaultService = "_bustlinker._tcp"
const (
	udp4  = 0
	udp6  = 1
	ipmax = 2
)

// MulticastDNS ...
type MulticastDNS struct {
	cfg *OptionConfig
}

// Server ...
func (dns *MulticastDNS) Server() (s Server, err error) {
	// Create the listeners
	conn := make([]*net.UDPConn, ipmax)

	var udp4Err error
	if dns.cfg.IPv4Addr != nil {
		conn[udp4], udp4Err = net.ListenUDP("udp4", dns.cfg.WildcardAddrIPv4)
		if udp4Err != nil {
			conn[udp4] = &net.UDPConn{}
		}
	}
	var udp6Err error
	if dns.cfg.IPv6Addr != nil {
		conn[udp6], udp6Err = net.ListenUDP("udp6", dns.cfg.WildcardAddrIPv6)
		if udp6Err != nil {
			conn[udp6] = &net.UDPConn{}
		}
	}

	// Check if we have any listener
	if udp4Err != nil && udp6Err != nil {
		return nil, fmt.Errorf("no multicast listeners could be started")
	}
	p1 := ipv4.NewPacketConn(conn[udp4])
	p2 := ipv6.NewPacketConn(conn[udp6])
	p1.SetMulticastLoopback(true)
	p2.SetMulticastLoopback(true)

	if dns.cfg.NetInterface != nil {
		if err := p1.JoinGroup(dns.cfg.NetInterface, &net.UDPAddr{IP: net.ParseIP(mdnsIPv4Addr)}); err != nil {
			return nil, err
		}
		if err := p2.JoinGroup(dns.cfg.NetInterface, &net.UDPAddr{IP: net.ParseIP(mdnsIPv6Addr)}); err != nil {
			return nil, err
		}
	} else {
		ifaces, err := net.Interfaces()
		if err != nil {
			return nil, err
		}
		errCount1, errCount2 := 0, 0
		for _, iface := range ifaces {
			if err := p1.JoinGroup(&iface, &net.UDPAddr{IP: net.ParseIP(mdnsIPv4Addr)}); err != nil {
				errCount1++
			}
			if err := p2.JoinGroup(&iface, &net.UDPAddr{IP: net.ParseIP(mdnsIPv6Addr)}); err != nil {
				errCount2++
			}
		}
		if len(ifaces) == errCount1 && len(ifaces) == errCount2 {
			return nil, fmt.Errorf("failed to join multicast group on all interfaces")
		}
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
	if uudp4Err != nil && uudp6Err != nil {
		logE("failed to bind to port", "uudp6Err", uudp4Err, "uudp4Err", uudp6Err)
		return nil, fmt.Errorf("failed to bind to any unicast udp port")
	}
	conn := make([]*net.UDPConn, ipmax)
	var udp4Err error
	if dns.cfg.IPv4Addr != nil {
		conn[udp4], udp4Err = net.ListenMulticastUDP("udp4", nil, dns.cfg.IPv4Addr)
	}
	if udp4Err != nil {
		logE("failed to bind to port", "udp4Err", udp4Err)
		conn[udp4] = &net.UDPConn{}
	}
	var udp6Err error
	if dns.cfg.IPv6Addr != nil {
		conn[udp6], udp6Err = net.ListenMulticastUDP("udp6", nil, dns.cfg.IPv6Addr)
	}
	if udp6Err != nil {
		logE("failed to bind to port", "udp6Err", udp6Err)
		conn[udp6] = &net.UDPConn{}
	}
	// Check if we have any listener
	if udp4Err != nil && udp6Err != nil {
		logE("failed to bind to port", "udp6Err", udp6Err, "udp4Err", udp4Err)
		return nil, fmt.Errorf("failed to bind to any multicast udp port")
	}
	p1 := ipv4.NewPacketConn(conn[udp4])
	p2 := ipv6.NewPacketConn(conn[udp6])
	p1.SetMulticastLoopback(true)
	p2.SetMulticastLoopback(true)

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var errCount1, errCount2 int

	for _, iface := range ifaces {
		if err := p1.JoinGroup(&iface, &net.UDPAddr{IP: net.ParseIP(mdnsIPv4Addr)}); err != nil {
			errCount1++
		}
		if err := p2.JoinGroup(&iface, &net.UDPAddr{IP: net.ParseIP(mdnsIPv6Addr)}); err != nil {
			errCount2++
		}
	}

	if len(ifaces) == errCount1 && len(ifaces) == errCount2 {
		return nil, fmt.Errorf("failed to join multicast group on all interfaces")
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
	optionConfig.IPv4Addr = &net.UDPAddr{
		IP:   net.ParseIP(mdnsIPv4Addr),
		Port: optionConfig.CustomPort,
	}
	optionConfig.IPv6Addr = &net.UDPAddr{
		IP:   net.ParseIP(mdnsIPv6Addr),
		Port: optionConfig.CustomPort,
	}

	optionConfig.WildcardAddrIPv4 = &net.UDPAddr{
		IP:   net.ParseIP(mdnsWildcardIPv4Addr),
		Port: optionConfig.CustomPort,
	}
	optionConfig.WildcardAddrIPv6 = &net.UDPAddr{
		IP:   net.ParseIP(mdnsWildcardIPv6Addr),
		Port: optionConfig.CustomPort,
	}

	optionConfig.instanceAddr = instanceAddr(optionConfig.Instance, optionConfig.Service, optionConfig.Domain)
	optionConfig.serviceAddr = serviceAddr(optionConfig.Service, optionConfig.Domain)
	optionConfig.enumAddr = enumAddr(optionConfig.Domain)
	return &MulticastDNS{
		cfg: optionConfig,
	}, nil
}

func defaultConfig(cfg *config.Config) *OptionConfig {
	//ipv4Addr := &net.UDPAddr{
	//	IP:   net.ParseIP(mdnsIPv4Addr),
	//	Port: mdnsPort,
	//}
	//ipv6Addr := &net.UDPAddr{
	//	IP:   net.ParseIP(mdnsIPv6Addr),
	//	Port: mdnsPort,
	//}
	//crc32.NewIEEE().Sum(cfg.)
	loadAccount, err := account.LoadAccount(cfg)
	var name string
	if err != nil {
		logE("load account error", "error", err)
		name = tool.GenerateRandomString(8)
	} else {
		name = interceptAccountName(loadAccount.Name)
	}

	hostName, _ := os.Hostname()
	hostName = fmt.Sprintf("%s.", hostName)
	service := defaultService
	instance := name
	domain := "local."

	return &OptionConfig{
		NetInterface:      nil,
		IPv4Addr:          nil,
		IPv6Addr:          nil,
		LogEmptyResponses: false,
		HostName:          hostName,
		Instance:          instance,
		instanceAddr:      instanceAddr(instance, service, domain),
		Service:           service,
		serviceAddr:       serviceAddr(service, domain),
		Enum:              "",
		enumAddr:          enumAddr(domain),
		Domain:            domain,
		CustomPort:        mdnsPort,
		Port:              80,
		TTL:               defaultTTL,
		TXT:               []string{}, // TXT,
	}
}

func interceptAccountName(s string) string {
	if len(s) == 42 {
		return s[2:6] + s[38:]
	}
	return ""
}
