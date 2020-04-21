package mdns

import (
	"fmt"
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
