package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/core"
	"github.com/portmapping/go-reuse"
	"go.uber.org/atomic"
	"net"
	"sync"
	"time"
)

const maxByteSize = 65520

type jsonNode struct {
	Addrs []core.Addr `json:"addrs"`
}

type node struct {
	api       core.API
	id        string
	addrs     []core.Addr
	isRunning *atomic.Bool
	session   *atomic.Uint64
	isAccept  bool
	conn      net.Conn
	isClosed  bool
	sendData  chan []byte
	callback  sync.Map
	info      *core.NodeInfo
}

var _ core.Node = &node{}

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
func AcceptNode(conn net.Conn, api core.API) (core.Node, error) {
	addr := conn.RemoteAddr()
	ip, port := basis.SplitIP(addr.String())

	return nodeRun(&node{
		id:        "", //id will get on running
		api:       api,
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

// ConnectNode ...
func ConnectNode(addr core.Addr, bind int, api core.API) (core.Node, error) {
	conn, err := reuse.DialTCP(addr.Protocol, &net.TCPAddr{
		IP:   net.IPv4zero,
		Port: bind,
	}, addr.TCP())
	if err != nil {
		return nil, err
	}
	return nodeRun(&node{
		id:        "", //id will get on running
		api:       api,
		isRunning: atomic.NewBool(false),
		sendData:  make(chan []byte),
		addrs:     []core.Addr{addr},
		conn:      conn,
	})
}

func defaultNode(conn net.Conn) *node {
	return &node{
		api:       nil,
		id:        "", //id will get on running
		addrs:     nil,
		isRunning: atomic.NewBool(false),
		session:   atomic.NewUint64(0),
		isAccept:  false,
		conn:      conn,
		isClosed:  false,
		sendData:  make(chan []byte),
		callback:  sync.Map{},
		info:      nil,
	}
}

// SetAPI ...
func (n *node) SetAPI(api core.API) {
	n.api = api
}

func (n *node) recv(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		tmp := make([]byte, maxByteSize)
		_, err := n.conn.Read(tmp)
		if err != nil {
			return
		}
		fmt.Println("recv", string(tmp))
		indexByte := bytes.IndexByte(tmp, 0)
		go n.doRecv(tmp[:indexByte])
	}
}

func (n *node) send(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		tmp := make([]byte, maxByteSize)
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
	return n.idRequest()
}

// Info ...
func (n *node) Info() core.NodeInfo {
	if n.info != nil {
		return *n.info
	}
	return n.infoRequest()
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
	go n.recv(wg)
	go n.send(wg)
	wg.Wait()
}

func (n *node) idRequest() string {
	ex := &Exchange{
		Type: Request,
		Data: nil,
	}

	resp := make(chan []byte)
	n.sendData <- ex.JSON()
	t := time.NewTimer(5 * time.Second)
	select {
	case id := <-resp:
		n.id = string(id)
	case <-t.C:
		return n.id
	}
	return n.id
}

func (n *node) doRecv(r []byte) {
	var ed Exchange
	err := json.Unmarshal(r, &ed)
	if err != nil {
		fmt.Println("failed", err)
		return
	}
	switch ed.Type {
	case Request:
		ex := &Exchange{Type: Response}
		id, err := n.api.ID(&core.IDReq{})
		if err != nil {
			ex.Status = StatusFailed
			ex.Data = []byte(err.Error())
		} else {
			ex.Data = []byte(id.Name)
		}
		n.sendData <- ex.JSON()
	default:
		n.cb(&ed)
	}
}

func (n *node) cb(ed *Exchange) {
	v, b := n.callback.Load(ed.Type)
	if b {
		cb, b := v.(chan []byte)
		if b {
			cb <- ed.Data
		}
	}
}

func (n *node) infoRequest() core.NodeInfo {
	return core.NodeInfo{}
}
