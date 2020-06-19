package node

import (
	"bytes"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/controller"
	"github.com/glvd/accipfs/core"
	"github.com/portmapping/go-reuse"
	"go.uber.org/atomic"
	"net"
	"sync"
)

const maxByteSize = 65520

type jsonNode struct {
	Addrs []core.Addr `json:"addrs"`
}

type node struct {
	c         *controller.Controller
	id        string
	addrs     []core.Addr
	isRunning *atomic.Bool
	isAccept  bool
	conn      net.Conn
	isClosed  bool
}

// IsClosed ...
func (n *node) IsClosed() bool {
	return n.isClosed
}

// Closed ...
func (n *node) Closed(f func(core.Node) bool) {
	if f != nil {
		n.isClosed = f(n)
	}
}

// IsConnecting ...
func (n *node) IsConnecting() bool {
	return true
}

var _ core.Node = &node{}

// Close ...
func (n *node) Close() (err error) {
	if n.conn != nil {
		err = n.conn.Close()
		n.conn = nil
	}
	return
}

// Verify ...
func (n *node) Verify() bool {
	return true
}

// ConnectToNode ...
func ConnectToNode(addr core.Addr, bind int, ctrl *controller.Controller) (core.Node, error) {
	tcp, err := reuse.DialTCP(addr.Protocol, &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: bind,
	}, addr.TCP())
	if err != nil {
		return nil, err
	}
	return &node{
		id:    "",
		c:     ctrl,
		addrs: []core.Addr{addr},
		conn:  tcp,
	}, nil
}

func (n *node) recv(wg *sync.WaitGroup, b chan<- []byte) {
	defer wg.Done()
	for {
		tmp := make([]byte, maxByteSize)
		read, err := n.conn.Read(tmp)
		if err != nil {
			return
		}
		log.Debugw("recv", "read", read)
		b <- tmp

	}
}

func (n *node) send(wg *sync.WaitGroup, b <-chan []byte) {
	defer wg.Done()
	tmp := make([]byte, maxByteSize)
	for {
		copy(tmp, <-b)
		write, err := n.conn.Write(tmp)
		if err != nil {
			return
		}
		log.Debugw("send", "write", write)
	}
}

func nodeRun(node *node) (core.Node, error) {
	go node.running()
	return node, nil
}

// AcceptNode ...
func AcceptNode(conn net.Conn, ctrl *controller.Controller) (core.Node, error) {
	addr := conn.RemoteAddr()
	ip, port := basis.SplitIP(addr.String())

	return nodeRun(&node{
		id:        "", //todo
		c:         ctrl,
		isRunning: atomic.NewBool(false),
		isAccept:  true,
		addrs: []core.Addr{
			{
				Protocol: addr.Network(),
				IP:       ip,
				Port:     port,
			},
		},
		conn: conn,
	})
}

// Addrs ...
func (n node) Addrs() []core.Addr {
	return n.addrs
}

// ID ...
func (n *node) ID() string {
	if n.id != "" {
		return n.id
	}
	if n.isAccept {
		id, err := n.c.LocalAPI().ID(&core.IDReq{})
		if err != nil {
			return ""
		}
		n.id = id.Name
	}

}

// Info ...
func (n *node) Info() core.NodeInfo {
	panic("implement me")
}

// Ping ...
func (n *node) Ping() error {
	panic("implement me")
}

func (n *node) running() {
	if n.isRunning.Load() {
		return
	}
	n.isRunning.Store(true)
	defer func() {
		if n.conn != nil {
			n.conn.Close()
		}
	}()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	cache := make(map[int]bytes.Buffer)

	recvData := make(chan []byte)
	sendData := make(chan []byte)
	go n.recv(wg, recvData)
	go n.send(wg, sendData)
	for {

		select {
		case snd <- buf:
		case buf1 := <-rec:
			fmt.Println("rec", string(buf1))
		}
	}

	wg.Wait()
}
