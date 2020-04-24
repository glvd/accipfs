package mdns

import (
	"fmt"
	"github.com/glvd/accipfs/config"
	"net"
	"strings"
)

// OptionConfig ...
type OptionConfig struct {
	//Zone              string
	NetInterface      *net.Interface
	IPv4Addr          *net.UDPAddr
	IPv6Addr          *net.UDPAddr
	WildcardAddrIPv4  *net.UDPAddr
	WildcardAddrIPv6  *net.UDPAddr
	LogEmptyResponses bool
	HostName          string
	instanceAddr      string
	serviceAddr       string
	enumAddr          string
	CustomPort        int
	Port              uint16
	TTL               uint32
	TXT               []string
	IPs               []net.IP
	Instance          string
	Service           string
	Enum              string
	Domain            string
}

// OptionConfigFunc ...
type OptionConfigFunc func(cfg *OptionConfig)

func serviceAddr(service, domain string) string {
	return fmt.Sprintf("%s.%s.", trimDot(service), trimDot(domain))
}
func instanceAddr(instance, service, domain string) string {
	return fmt.Sprintf("%s.%s.%s.", instance, trimDot(service), trimDot(domain))
}

func enumAddr(domain string) string {
	return fmt.Sprintf("_services._dns-sd._udp.%s.", trimDot(domain))
}

// trimDot is used to trim the dots from the start or end of a string
func trimDot(s string) string {
	return strings.Trim(s, ".")
}

// RegisterLocalIP ...
func (cfg *OptionConfig) RegisterLocalIP(c *config.Config) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for i := range addrs {
		if ipnet, ok := addrs[i].(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			var local bool
			if ipv4 := ipnet.IP.To4(); ipv4 != nil {
				if local = isLocalIP(ipv4); local {
					output("register mdns service ipv4 addr:", ipv4.String())
					cfg.IPs = append(cfg.IPs, ipv4)
				}
			}
			if local {
				if ipv6 := ipnet.IP.To16(); ipv6 != nil {
					output("register mdns service ipv6 addr:", ipv6.String())
					cfg.IPs = append(cfg.IPs, ipv6)
				}
			}
		}
	}
	cfg.Port = uint16(c.Port)
}

func isLocalIP(ip4 net.IP) bool {
	switch {
	case ip4[0] == 10:
		return true
	case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
		return true
	case ip4[0] == 192 && ip4[1] == 168:
		return true
	default:
		return false
	}
}
