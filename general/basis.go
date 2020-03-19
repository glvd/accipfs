package general

import (
	"os"
	"strconv"
	"strings"
)

// CurrentDir ...
func CurrentDir() string {
	dir, err := os.Getwd()
	if err == nil {
		return dir
	}
	return "."
}

// SplitIP ...
func SplitIP(addr string) (ip string, port int) {
	if addr == "" {
		return
	}
	s := strings.Split(addr, ":")
	if len(s) < 2 {
		return
	}
	ip = s[0]
	port, _ = strconv.Atoi(s[1])
	return
}
