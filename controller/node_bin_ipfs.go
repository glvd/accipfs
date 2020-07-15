package controller

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/goextension/io"
	files "github.com/ipfs/go-ipfs-files"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	iface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"github.com/ipfs/interface-go-ipfs-core/path"
	"github.com/multiformats/go-multiaddr"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var _ core.ControllerService = &nodeBinIPFS{}

type nodeBinIPFS struct {
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *config.Config
	name   string
	cmd    *exec.Cmd
	msg    func(string)
	api    *httpapi.HttpApi
}

// MessageHandle ...
func (n *nodeBinIPFS) MessageHandle(f func(s string)) {
	n.msg = f
}

// Msg ...
func (n *nodeBinIPFS) Msg(s string) {
	if n.msg != nil {
		n.msg(s)
	}
}

// Start ...
func (n *nodeBinIPFS) Start() error {
	n.cmd = exec.CommandContext(n.ctx, n.name, "daemon", "--routing", "none")
	output("ipfs cmd: ", n.cmd.Args)
	pipe, err2 := n.cmd.StderrPipe()
	if err2 != nil {
		return err2
	}
	stdoutPipe, err2 := n.cmd.StdoutPipe()
	if err2 != nil {
		return err2
	}
	m := io.MultiReader(pipe, stdoutPipe)
	if n.cfg.IPFS.LogOutput {
		go basis.PipeReader(n.ctx, m, n.msg)
	}
	//else {
	//	go basis.PipeDummy(n.ctx, module, m)
	//}
	err := n.cmd.Start()
	if err != nil {
		return err
	}

	return nil
}

// Stop ...
func (n *nodeBinIPFS) Stop() error {
	if n.cmd != nil {
		n.cancel()
		n.cmd = nil
	}
	return nil
}

// Initialize ...
func (n *nodeBinIPFS) Initialize() error {
	_, err := os.Stat(config.DataDirIPFS())
	if err != nil && os.IsNotExist(err) {
		_ = os.MkdirAll(config.DataDirIPFS(), 0755)

		//os.Setenv("IPFS_PATH", filepath.Join(n.cfg.Path, config.DataDirIPFS()))
		cmd := exec.Command(n.name, "init", "--profile", "badgerds")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("init:%w", err)
		}

		version := n.getVersion()
		logI("ipfs init", "log", string(out), "version", version)
		n.Msg(string(out))
		if version[1] < 5 {
			cmd = exec.Command(n.name, "config", "Swarm.EnableAutoNATService", "--bool", "true")
			out, err = cmd.CombinedOutput()
			if err != nil {
				return fmt.Errorf("config(nat):%w", err)
			}
		}
		n.Msg(string(out))
		logI("ipfs init config set", "log", string(out))
		cmd = exec.Command(n.name, "config", "Swarm.EnableRelayHop", "--bool", "true")
		out, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("config(relay):%w", err)
		}
		n.Msg(string(out))
		logI("ipfs init config set", "log", string(out))
		logI("ipfs init end")
	}
	return nil
}

func newNodeBinIPFS(cfg *config.Config) *nodeBinIPFS {
	path := filepath.Join(cfg.Path, "bin", basis.BinName(cfg.IPFS.Name))
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &nodeBinIPFS{
		ctx:    ctx,
		cancel: cancelFunc,
		cfg:    cfg,
		name:   path,
	}
}

func (n *nodeBinIPFS) getVersion() (ver [3]int) {
	cmd := exec.Command(n.name, "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	split := strings.Split(string(out), " ")
	if len(split) < 3 {
		return
	}
	verS := strings.Split(split[2], ".")
	for i := range verS {
		parseInt, err := strconv.ParseInt(verS[i], 32, 10)
		if err != nil {
			return
		}
		ver[i] = int(parseInt)
	}
	return
}

// IsReady ...
func (n *nodeBinIPFS) IsReady() bool {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	_, err := n.ID(ctx)
	if err != nil {
		return false
	}

	return true
}

// APIContext ...
func (n *nodeBinIPFS) API() *httpapi.HttpApi {
	if n.api == nil {
		ma, err := multiaddr.NewMultiaddr(config.IPFSAddr())
		if err != nil {
			panic(err)
		}
		api, e := httpapi.NewApi(ma)
		if e != nil {
			panic(e)
		}
		n.api = api
	}
	return n.api
}

// ID get self serviceNode info
func (n *nodeBinIPFS) ID(ctx context.Context) (pid *core.DataStoreInfo, e error) {
	pid = &core.DataStoreInfo{}
	e = n.API().Request("id").Exec(ctx, pid)
	if e != nil {
		return nil, e
	}
	return pid, nil
}

// PinAdd ...
func (n *nodeBinIPFS) FileAdd(ctx context.Context, filename string, option options.UnixfsAddOption) (hash string, e error) {
	stat, e := os.Stat(filename)
	if e != nil {
		return "", e
	}
	var node files.Node
	//var err error
	if !stat.IsDir() {
		file, e := os.Open(filename)
		if e != nil {
			return "", e
		}
		node = files.NewReaderFile(file)
	} else {
		sf, e := files.NewSerialFile(filename, false, stat)
		if e != nil {
			return "", e
		}
		node = sf
	}
	resolved, e := n.API().Unixfs().Add(ctx, node, option)
	if e != nil {
		return "", e
	}
	return resolved.Cid().String(), nil
}

// PinAdd ...
func (n *nodeBinIPFS) PinAdd(ctx context.Context, hash string) (e error) {
	p := path.New(hash)
	return n.API().Pin().Add(ctx, p, options.Pin.Recursive(true))
}

// PinLS ...
func (n *nodeBinIPFS) PinLS(ctx context.Context) (pins <-chan iface.Pin, e error) {
	return n.API().Pin().Ls(ctx, options.Pin.Ls.Recursive())
}

// PinRm ...
func (n *nodeBinIPFS) PinRm(ctx context.Context, hash string) (e error) {
	p := path.New(hash)
	return n.API().Pin().Rm(ctx, p)
}

// SwarmPeers ...
func (n *nodeBinIPFS) SwarmPeers(ctx context.Context) ([]iface.ConnectionInfo, error) {
	return n.API().Swarm().Peers(ctx)
}
