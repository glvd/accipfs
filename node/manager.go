package node

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var _name = "bl.nodes"

type manager struct {
	cfg      *config.Config
	t        *time.Ticker
	ts       int64
	nodes    sync.Map
	expNodes sync.Map
}

// New ...
func New(cfg *config.Config) core.NodeManager {
	m := &manager{
		cfg: cfg,
		t:   time.NewTicker(cfg.Node.BackupSeconds),
	}
	return m
}

// Load ...
func (m *manager) Load() error {
	file := filepath.Join(m.cfg.Path, _name)
	stat, err := os.Stat(file)
	if err != nil || !os.IsNotExist(err) {
		return err
	}

	if stat.IsDir() {
		return fmt.Errorf("found file but is directory")
	}
	open, err := os.Open(file)
	if err != nil {
		return err
	}

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

// Push ...
func (m *manager) Push(node core.Node) {
	m.ts = time.Now().Unix()
	m.nodes.Store(node.ID(), node)
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
