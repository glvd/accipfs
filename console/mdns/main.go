package main

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/mdns"
	"time"
)

func main() {
	m, err := mdns.New(config.Default(), func(cfg *mdns.OptionConfig) {
		cfg.Service = "_foobar._tcp"
	})
	if err != nil {
		panic(err)
	}

	s2, err := m.Server()
	if err != nil {
		panic(err)
	}
	s2.Start()
	time.Sleep(5 * time.Minute)
	defer s2.Stop()
}
