package core

import "fmt"

// AddressInfo ...
type AddressInfo struct {
	Address string
	Schema  string
	URI     string
	Port    int
}

// URL ...
func (a *AddressInfo) URL() string {
	return fmt.Sprintf("http://%s:%d/rpc", a.Address, a.Port)
}
