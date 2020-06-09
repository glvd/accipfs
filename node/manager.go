package node

import (
	"bufio"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/robfig/cron/v3"
	"io"
	"os"
	"path/filepath"
	"sync"
)

var _name = "bl.nodes"

type manager struct {
	nodes sync.Map
	cfg   *config.Config
	c     *cron.Cron
}

// LoadNodesFromFile ...
func LoadNodesFromFile(path string) {

}

// New ...
func New(cfg *config.Config) core.NodeManager {
	var m manager
	m.c = cron.New(cron.WithSeconds())
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
