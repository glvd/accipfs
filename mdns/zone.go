package mdns

import (
	"fmt"
	"strings"
)

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
