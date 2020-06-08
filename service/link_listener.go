package service

import (
	"github.com/glvd/accipfs/config"
	"net"

	"github.com/panjf2000/ants/v2"
	"github.com/portmapping/go-reuse"
)

type tcpListener struct {
	protocol string
	bindPort int
	port     int
	connBack func(conn net.Conn)
	pool     *ants.Pool
}

// NewLinkListener ...
func NewLinkListener(cfg *config.Config) Listener {
	l := &tcpListener{
		protocol: "tcp",
		bindPort: cfg.Node.BindPort,
		port:     cfg.Node.Port,
	}
	pool, err := ants.NewPool(5000)
	if err != nil {
		return nil
	}
	l.pool = pool
	return l
}

// Listen ...
func (h *tcpListener) Listen() error {
	local := &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: h.port,
	}
	tcp, err := reuse.ListenTCP(h.protocol, local)
	if err != nil {
		return err
	}
	for {
		conn, err := tcp.Accept()
		if err != nil {
			continue
		}
		if h.connBack != nil {
			h.pool.Submit(func() {
				h.connBack(conn)
			})
			continue
		}
		//no callback closed
		conn.Close()
	}
}

// Accept ...
func (h *tcpListener) Accept(f func(conn net.Conn)) {
	if f != nil {
		h.connBack = f
	}
}
