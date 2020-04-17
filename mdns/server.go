package mdns

import (
	"github.com/glvd/accipfs/config"
)

const (
	mdnsIPV4Addr = "224.0.0.251"
	mdnsIPV6Addr = "FF02::FB"
	mdnsPort     = 5353
)

// Config ...
type Config struct {
}

// ConfigFunc ...
type ConfigFunc func(cfg *Config)

// New ...
func New(config *config.Config, opts ...ConfigFunc) {

}
