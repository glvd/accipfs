package node

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/core"
	"github.com/portmapping/go-reuse"
	"go.uber.org/atomic"
	"math"
	"net"
	"sync"
)

const maxByteSize = 65520

type jsonNode struct {
	Addrs []core.Addr `json:"addrs"`
}

type node struct {
	ctx       context.Context
	cancel    context.CancelFunc
	api       core.API
	id        string
	callback  sync.Map
	session   *atomic.Uint32
	addrs     []core.Addr
	isRunning *atomic.Bool
	isAccept  bool
	conn      net.Conn
	isClosed  bool
	sendQueue chan *Queue
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
	n := defaultNode(conn)
	n.SetAPI(api)
	n.AppendAddr(core.Addr{
		Protocol: "tcp",
		IP:       ip,
		Port:     port,
	})
	return nodeRun(n)

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
	n := defaultNode(conn)
	n.SetAPI(api)
	n.AppendAddr(addr)
	return nodeRun(n)
}

func defaultNode(conn net.Conn) *node {
	ctx, fn := context.WithCancel(context.TODO())
	return &node{
		api:       nil,
		ctx:       ctx,
		cancel:    fn,
		id:        "", //id will get on running
		addrs:     nil,
		isRunning: atomic.NewBool(false),
		session:   atomic.NewUint32(0),
		isAccept:  false,
		conn:      conn,
		isClosed:  false,
		sendQueue: make(chan *Queue),
		callback:  sync.Map{},
		info:      nil,
	}
}

// AppendAddr ...
func (n *node) AppendAddr(addr core.Addr) {
	n.addrs = append(n.addrs, addr)
}

// SetAPI ...
func (n *node) SetAPI(api core.API) {
	n.api = api
}

func (n *node) recv(wg *sync.WaitGroup) {
	defer wg.Done()
	scan := dataScan(n.conn)
	for scan.Scan() {
		select {
		case <-n.ctx.Done():
			return
		default:
			exchange, err := ScanExchange(scan)
			if err != nil {
				continue
			}
			go n.doRecv(exchange)
		}
	}
	if scan.Err() != nil {
		log.Errorw("recv", "error", scan.Err())
		n.cancel()
	}
}

// Session ...
func (n *node) Session() uint32 {
	s := n.session.Load()
	if s != math.MaxUint32 {
		n.session.Inc()
	} else {
		n.session.Store(0)
	}
	return s
}

// RegisterCallback ...
func (n *node) RegisterCallback(queue *Queue) {
	s := n.Session()
	queue.SetSession(s)
	n.callback.Store(s, queue.Callback)
}

func (n *node) send(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-n.ctx.Done():
			return
		case q := <-n.sendQueue:
			if q.HasCallback() {
				n.RegisterCallback(q)
			}
			err := q.Exchange().Pack(n.conn)
			if err != nil {
				log.Errorw("recv", "error", err)
				n.cancel()
				return
			}
		}
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
	fmt.Println("node exit")
	n.isRunning.Store(false)
}

func (n *node) idRequest() string {
	ex := NewExchange(Request)
	q := NewQueue(ex, true)
	n.sendQueue <- q
	callback := q.WaitCallback()
	if callback == nil {
		return ""
	}
	return string(callback.Data)
}

// CallbackTrigger ...
func (n *node) CallbackTrigger(exchange *Exchange) {
	load, ok := n.callback.Load(exchange.Session)
	if ok {
		v, b := load.(func(exchange *Exchange))
		if b {
			v(exchange)
		}
		n.callback.Delete(exchange.Session)
	}
}

func (n *node) doRecv(exchange *Exchange) {
	switch exchange.Type {
	case Request:
		ex := NewExchange(Response)
		id, err := n.api.ID(&core.IDReq{})
		if err != nil {
			ex.Status = StatusFailed
			ex.Session = exchange.Session
			ex.SetData([]byte(err.Error()))
		} else {
			ex.Session = exchange.Session
			ex.SetData([]byte(id.Name))
		}
		q := NewQueue(ex, false)
		n.sendQueue <- q
	case Response:
		n.CallbackTrigger(exchange)
	default:
		return
	}
}

func (n *node) infoRequest() core.NodeInfo {
	return core.NodeInfo{}
}
