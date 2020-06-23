package node

import (
	"context"
	"fmt"
	"math"
	"net"
	"sync"
	"time"

	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/core"
	"github.com/portmapping/go-reuse"
	"go.uber.org/atomic"
)

const maxByteSize = 65520

type jsonNode struct {
	Addrs []core.Addr `json:"addrs"`
}

//temp Data
type nodeLocal struct {
	id *string
}

type node struct {
	ctx       context.Context
	cancel    context.CancelFunc
	closeCB   func(node core.Node)
	local     *nodeLocal
	api       core.API
	callback  sync.Map
	connector bool
	isTimeout *atomic.Bool
	timeout   *time.Timer
	session   *atomic.Uint32
	addrs     []core.Addr
	isRunning *atomic.Bool
	isAccept  bool
	conn      net.Conn
	isClosed  bool
	sendQueue chan *Queue
	info      *core.NodeInfo
	heartBeat *time.Ticker
}

var _ core.Node = &node{}
var heartBeatTimer = 15 * time.Second

// IsClosed ...
func (n *node) IsClosed() bool {
	return n.isClosed
}

// Closed ...
func (n *node) Closed(f func(core.Node)) {
	if f != nil {
		n.closeCB = f
	}
}

// IsConnecting ...
func (n *node) IsConnecting() bool {
	return true
}

// Close ...
func (n *node) Close() (err error) {
	defer func() {
		n.isClosed = true
		if n.closeCB != nil {
			n.closeCB(n)
		}
	}()
	if n.cancel != nil {
		n.cancel()
		n.cancel = nil
	}

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
	n.connector = true
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
		heartBeat: time.NewTicker(heartBeatTimer),
		local:     &nodeLocal{},
		isTimeout: atomic.NewBool(false),
		timeout:   time.NewTimer(24 * time.Hour),
		addrs:     nil,
		isRunning: atomic.NewBool(false),
		session:   atomic.NewUint32(math.MaxUint32 - 5),
		isAccept:  false,
		conn:      conn,
		isClosed:  false,
		sendQueue: make(chan *Queue),
		callback:  sync.Map{},
		info:      nil,
	}
}

// Session ...
func (n *node) Session() uint32 {
	s := n.session.Load()
	log.Infow("session", "num", s)
	if s != math.MaxUint32 {
		n.session.Inc()
	} else {
		n.session.Store(1)
	}
	return s
}

// RegisterCallback ...
func (n *node) RegisterCallback(queue *Queue) {
	s := n.Session()
	queue.SetSession(s)
	n.callback.Store(s, queue.Callback)
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
	defer func() {
		wg.Done()
		if e := recover(); e != nil {
			_ = n.Close()
			log.Errorw("panic", "error", e)
		}
	}()
	scan := dataScan(n.conn)
	for scan.Scan() {
		select {
		case <-n.ctx.Done():
			return
		default:
			exchange, err := ScanExchange(scan)
			if err != nil {
				panic(err)
			}
			go n.doRecv(exchange)
		}
	}
	if scan.Err() != nil {
		panic(scan.Err())
	}
}

func (n *node) send(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		if e := recover(); e != nil {
			_ = n.Close()
			log.Errorw("panic", "error", e)
		}
	}()

	for {
		select {
		case <-n.ctx.Done():
			return
		case <-n.heartBeat.C:
			if !n.connector {
				continue
			}
			n.timeout.Reset(heartBeatTimer)
			err := NewRequestExchange(TypeDetailPing).Pack(n.conn)
			if err != nil {
				panic(err)
			}
		case q := <-n.sendQueue:
			if q.HasCallback() {
				n.RegisterCallback(q)
			}
			err := q.Exchange().Pack(n.conn)
			if err != nil {
				panic(err)
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
	if n.local.id == nil {
		id := n.idRequest()
		n.local.id = &id
	}
	return *n.local.id
}

// Info ...
func (n *node) Info() core.NodeInfo {
	if n.info == nil {
		n.info = n.infoRequest()
	}
	return *n.info
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
		_ = n.Close()
	}()
	go n.beatChecker()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go n.recv(wg)
	go n.send(wg)
	wg.Wait()
	fmt.Println("node exit")
	n.isRunning.Store(false)
}

func (n *node) idRequest() string {
	ex := NewRequestExchange(TypeDetailID)
	if cb, b := n.SendExchange(ex); b {
		callback := cb.WaitCallback()
		if callback != nil {
			return string(callback.Data)
		}
	}
	return ""
}

// CallbackTrigger ...
func (n *node) CallbackTrigger(exchange *Exchange) {
	if exchange.Session == 0 {
		return
	}
	load, ok := n.callback.Load(exchange.Session)
	if ok {
		v, b := load.(func(exchange *Exchange))
		if b {
			v(exchange)
		}
		n.callback.Delete(exchange.Session)
	}
}

var recvReqFunc = map[TypeDetail]func(api core.API) (exchange *Exchange, err error){
	TypeDetailID:   recvRequestID,
	TypeDetailPing: recvRequestPing,
}

func recvRequestID(api core.API) (exchange *Exchange, err error) {
	exchange = NewResponseExchange(TypeDetailID)
	var idResp *core.IDResp
	idResp, err = api.ID(&core.IDReq{})
	if err != nil {
		exchange.Status = StatusFailed
		exchange.SetData([]byte(err.Error()))
	} else {
		exchange.SetData([]byte(idResp.Name))
	}
	return exchange, nil
}
func recvRequestPing(api core.API) (exchange *Exchange, err error) {
	exchange = NewResponseExchange(TypeDetailPing)
	return exchange, nil
}
func (n *node) recvRequest(exchange *Exchange) {
	f, b := recvReqFunc[exchange.TypeDetail]
	if !b {
		return
	}
	//ignore error
	ex, _ := f(n.api)
	ex.Session = exchange.Session
	NewQueue(ex).Send(n.sendQueue)
}

func (n *node) recvResponse(exchange *Exchange) {
	if exchange.TypeDetail == TypeDetailPing {
		if !n.isTimeout.Load() {
			n.timeout.Stop()
		}
		return
	}
	n.CallbackTrigger(exchange)
}

func (n *node) doRecv(exchange *Exchange) {
	switch exchange.Type {
	case TypeRequest:
		n.recvRequest(exchange)
	case TypeResponse:
		n.recvResponse(exchange)
	default:
		return
	}
}

func (n *node) infoRequest() *core.NodeInfo {
	return &core.NodeInfo{}
}

func (n *node) pingRequest() bool {
	ex := newExchange(TypeRequest, TypeDetailPing)
	if cb, b := n.SendExchange(ex); b {
		callback := cb.WaitCallback()
		if callback != nil {
			return true
		}
	}
	return false
}

// SendQueue ...
func (n *node) SendExchange(ex *Exchange) (cb CallbackWaiter, b bool) {
	queue := NewQueue(ex, func(option *QueueOption) {
		option.Callback = true
	})
	b = queue.Send(n.sendQueue)
	return queue, b
}

// SendQueue ...
func (n *node) SendQueue(queue *Queue) bool {
	return queue.Send(n.sendQueue)
}

func (n *node) beatChecker() {
	defer func() {
		n.Close()
		if e := recover(); e != nil {
			log.Errorw("beatChecker timeout", "error", e)
		}
	}()
	select {
	case <-n.timeout.C:
		n.isTimeout.Store(true)
		panic("heart beat timeout")
	case <-n.ctx.Done():
		return
	}
}
