package node

import (
	"encoding/json"
	"fmt"
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
	sendData  chan []byte
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

// AcceptNode ...
func AcceptNode(conn net.Conn, ctrl *controller.Controller) (core.Node, error) {
	addr := conn.RemoteAddr()
	ip, port := basis.SplitIP(addr.String())

	return nodeRun(&node{
		id:        "", //id will get on running
		c:         ctrl,
		sendData:  make(chan []byte),
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

// ConnectToNode ...
func ConnectToNode(addr core.Addr, bind int, ctrl *controller.Controller) (core.Node, error) {
	conn, err := reuse.DialTCP(addr.Protocol, &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: bind,
	}, addr.TCP())
	if err != nil {
		return nil, err
	}
	return nodeRun(&node{
		id:        "", //id will get on running
		c:         ctrl,
		isRunning: atomic.NewBool(false),
		sendData:  make(chan []byte),
		addrs:     []core.Addr{addr},
		conn:      conn,
	})
}

func (n *node) recv(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		tmp := make([]byte, maxByteSize)
		read, err := n.conn.Read(tmp)
		if err != nil {
			return
		}
		log.Debugw("recv", "read", read)
		go n.doRecv(tmp)
	}
}

func (n *node) send(wg *sync.WaitGroup) {
	defer wg.Done()
	tmp := make([]byte, maxByteSize)
	for {
		copy(tmp, <-n.sendData)
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
	n.idRequest()
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
	//recvData := make(chan []byte)
	go n.recv(wg)
	go n.send(wg)
	wg.Wait()
}

func (n *node) idRequest() {
	ex := &Exchange{
		Type: RequestID,
		Data: nil,
	}
	n.sendData <- ex.JSON()
}

func (n *node) doRecv(r []byte) {
	var ed Exchange
	err := json.Unmarshal(r, &ed)
	if err != nil {
		return
	}
	switch ed.Type {
	case RequestID:
		id := n.ID()
		ex := &Exchange{
			Type: ResponseID,
			Data: []byte(id),
		}
		n.sendData <- ex.JSON()
	case ResponseID:
		n.id = string(ed.Data)
	default:
		fmt.Println("wrong type")
	}
}
