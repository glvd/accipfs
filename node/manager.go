package node

import (
	"bufio"
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
	nodes    sync.Map
	expNodes sync.Map
}

// New ...
func New(cfg *config.Config) core.NodeManager {
	m := manager{
		cfg:      cfg,
		t:        time.NewTicker(cfg.Node.BackupSeconds),
		nodes:    sync.Map{},
		expNodes: sync.Map{},
	}
	return &m
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
		line, prefix, err := reader.ReadLine()
		log.Debugw("load nodes", "line", string(line), "prefix", prefix)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

	}
}

func decodeNode(b []byte, node core.Node) error {

}
