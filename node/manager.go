package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/atomic"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type manager struct {
	cfg          *config.Config
	exchangePool *ants.PoolWithFunc
	t            *time.Ticker
	currentTS    int64
	ts           int64
	currentNodes *atomic.Uint64
	nodes        sync.Map
	expNodes     sync.Map
	path         string
	expPath      string
}

var _nodes = "bl.nodes"
var _expNodes = "exp.nodes"
var _ core.NodeManager = &manager{}

// New ...
func New(cfg *config.Config) core.NodeManager {
	m := &manager{
		cfg:     cfg,
		path:    filepath.Join(cfg.Path, _nodes),
		expPath: filepath.Join(cfg.Path, _expNodes),
		t:       time.NewTicker(cfg.Node.BackupSeconds),
	}
	m.exchangePool = mustPool(ants.DefaultAntsPoolSize, m.handleConn)
	go m.loop()

	return m
}

// Store ...
func (m *manager) Store() error {
	err := os.Remove(m.path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	file, err := os.OpenFile(m.path, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	m.nodes.Range(func(key, value interface{}) bool {
		n, b := value.(core.Node)
		if !b {
			return true
		}
		nodeData, err := encodeNode(n)
		if err != nil {
			return false
		}
		_, err = writer.Write(nodeData)
		if err != nil {
			return false
		}
		_, err = writer.WriteString(basis.NewLine)
		if err != nil {
			return false
		}
		return true
	})
	return writer.Flush()
}

// Load ...
func (m *manager) Load() error {
	stat, err := os.Stat(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if stat.IsDir() {
		return fmt.Errorf("found file but is directory")
	}
	open, err := os.Open(m.path)
	if err != nil {
		return err
	}
	defer open.Close()
	reader := bufio.NewReader(open)
	for {
		n, prefix, err := reader.ReadLine()
		log.Debugw("load nodes", "line", string(n), "prefix", prefix)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		err = decodeNode(m, n)
		if err != nil {
			log.Errorw("decode failed", "error", err, "data", string(n))
			continue
		}
	}
}

// StateExamination ...
func (m *manager) StateExamination(id string, f func(node core.Node) bool) {
	if f == nil {
		return
	}
	node, ok := m.nodes.Load(id)
	if ok {
		if f(node.(core.Node)) {
			m.nodes.Delete(id)
			m.expNodes.Store(id, node)
		}
	}

	exp, ok := m.expNodes.Load(id)
	if ok {
		if f(exp.(core.Node)) {
			m.expNodes.Delete(id)
			m.nodes.Store(id, exp)
		}
	}
}

// Range ...
func (m *manager) Range(f func(key string, node core.Node) bool) {
	m.nodes.Range(func(key, value interface{}) bool {
		k, b1 := key.(string)
		n, b2 := value.(core.Node)
		if !b1 || !b2 {
			return true
		}
		if f != nil {
			return f(k, n)
		}
		return false
	})
}

// Push ...
func (m *manager) Push(node core.Node) {
	m.ts = time.Now().Unix()
	m.nodes.Store(node.ID(), node)
}

// save nodes
func (m *manager) loop() {
	for {
		<-m.t.C
		if m.ts != m.currentTS {
			if err := m.Store(); err != nil {
				continue
			}
			m.currentTS = m.ts
		}
	}
}

// HandleConn ...
func (m *manager) HandleConn(i interface{}) {

}

func decodeNode(m core.NodeManager, b []byte) error {
	nodes := map[string]jsonNode{}
	err := json.Unmarshal(b, &nodes)
	if err != nil {
		return err
	}
	for id, nodes := range nodes {
		m.Push(&node{
			id:    id,
			addrs: nodes.Addrs,
		})
	}
	return nil
}

func encodeNode(node core.Node) ([]byte, error) {
	n := map[string]jsonNode{
		node.ID(): {Addrs: node.Addrs()},
	}
	return json.Marshal(n)
}

func mustPool(size int, f func(interface{})) *ants.PoolWithFunc {
	pf, err := ants.NewPoolWithFunc(size, f, ants.WithNonblocking(false))
	if err != nil {
		panic(err)
	}
	return pf
}
