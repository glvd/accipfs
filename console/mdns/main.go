package main

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/mdns"
)

func main() {
	m, err := mdns.New(config.Default())
	if err != nil {
		panic(err)
	}

	s2, err := m.Server()
	if err != nil {
		panic(err)
	}
	s2.Start()
	defer s2.Stop()
}
