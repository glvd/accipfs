package service

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"net"

	"github.com/panjf2000/ants/v2"
	"github.com/portmapping/go-reuse"
)

type linkListener struct {
	listener net.Listener
	protocol string
	bindPort int
	port     int
	cb       func(conn net.Conn) (core.Node, error)
}

// newLinkListener listen other client connections
func newLinkListener(cfg *config.Config, cb func(conn net.Conn) (core.Node, error)) core.Listener {
	l := &linkListener{
		protocol: "tcp",
		bindPort: cfg.Node.BindPort,
		port:     cfg.Node.Port,
		cb:       cb,
	}
	return l
}

// Stop ...
func (h *linkListener) Stop() error {
	if h.listener != nil {
		return h.listener.Close()
	}
	return nil
}

// Listen ...
func (h *linkListener) Listen() (err error) {
	local := &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: h.port,
	}
	h.listener, err = reuse.ListenTCP(h.protocol, local)
	if err != nil {
		return err
	}
	for {
		conn, err := h.listener.Accept()
		if err != nil {
			continue
		}
		if h.cb != nil {
			h.cb(conn)
			continue
		}
		//no callback closed
		err = conn.Close()
		if err != nil {
			log.Errorw("accept conn error", "err", err)
		}
	}
}
func mustPool(size int, f func(interface{})) *ants.PoolWithFunc {
	pf, err := ants.NewPoolWithFunc(size, f, ants.WithNonblocking(false))
	if err != nil {
		panic(err)
	}
	return pf
}
