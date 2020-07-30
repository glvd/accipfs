package controller

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/plugin/loader"
	ipfsversion "github.com/ipfs/go-ipfs"
	ipfsconfig "github.com/ipfs/go-ipfs-config"
	ipfscore "github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	intercore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"go.uber.org/atomic"
	"os"
	"path/filepath"
	"runtime"
)

type nodeLibIPFS struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        *config.Config
	isRunning  *atomic.Bool
	configRoot string
	intercore.CoreAPI
	repo    repo.Repo
	node    *ipfscore.IpfsNode
	plugins *loader.PluginLoader
}

var _ core.ControllerService = &nodeLibIPFS{}

func setupPlugins(externalPluginsPath string) (*loader.PluginLoader, error) {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return nil, fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return nil, fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return nil, fmt.Errorf("error initializing plugins: %s", err)
	}

	return plugins, nil
}

func newNodeLibIPFS(cfg *config.Config) *nodeLibIPFS {
	ctx, cancel := context.WithCancel(context.Background())
	root := filepath.Join(cfg.Path, ".ipfs")

	return &nodeLibIPFS{
		ctx:        ctx,
		cancel:     cancel,
		cfg:        cfg,
		configRoot: root,
		isRunning:  atomic.NewBool(false),
		//plugins:    plugins,
	}
}
func printVersion() {
	v := ipfsversion.CurrentVersionNumber
	if ipfsversion.CurrentCommit != "" {
		v += "-" + ipfsversion.CurrentCommit
	}
	fmt.Printf("Datastore version: %s\n", v)
	fmt.Printf("Repo version: %d\n", fsrepo.RepoVersion)
	fmt.Printf("System version: %s\n", runtime.GOARCH+"/"+runtime.GOOS)
	fmt.Printf("Golang version: %s\n", runtime.Version())
}

// Start ...
func (n *nodeLibIPFS) Start() (err error) {
	ipfsNode, err := startIPFSNode(n.ctx, n.configRoot)
	if err != nil {
		return err
	}
	n.CoreAPI = ipfsNode
	n.isRunning.Store(true)
	printVersion()
	fmt.Println("datastore is running")
	return nil
}

// Stop ...
func (n *nodeLibIPFS) Stop() error {
	if n.cancel != nil {
		n.cancel()
		n.cancel = nil
	}
	if n.node != nil {
		n.node.Close()
		n.node = nil
	}

	if n.repo != nil {
		n.repo.Close()
		n.repo = nil
	}

	n.isRunning.Store(false)
	return nil
}

// Initialize ...
func (n *nodeLibIPFS) Initialize() error {
	if _, err := setupPlugins(""); err != nil {
		return err
	}
	// Create a Repo
	if err := createRepo(n.ctx, n.configRoot); err != nil {
		return fmt.Errorf("failed to create temp repo: %s", err)
	}

	return nil
}

// IsReady ...
func (n *nodeLibIPFS) IsReady() bool {
	if n.isRunning.Load() {
		return true
	}
	return false
}

// MessageHandle ...
func (n *nodeLibIPFS) MessageHandle(f func(s string)) {

}

func createRepo(ctx context.Context, repoPath string) error {
	identity, err := ipfsconfig.CreateIdentity(os.Stdout, []options.KeyGenerateOption{options.Key.Type(options.Ed25519Key)})
	if err != nil {
		return err
	}

	// Create a config with default options and a 2048 bit key
	cfg, err := ipfsconfig.InitWithIdentity(identity)
	if err != nil {
		return err
	}
	cfg.Datastore.Spec = badgerSpec()
	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return nil
}

func badgerSpec() map[string]interface{} {
	return map[string]interface{}{
		"type":   "measure",
		"prefix": "badger.datastore",
		"child": map[string]interface{}{
			"type":       "badgerds",
			"path":       "badgerds",
			"syncWrites": false,
			"truncate":   true,
		},
	}
}

func startIPFSNode(ctx context.Context, path string) (intercore.CoreAPI, error) {
	/// --- Part I: Getting a IPFS node running
	fmt.Println("-- Getting an IPFS node running -- ")

	if _, err := setupPlugins(""); err != nil {
		return nil, err
	}
	// Spawning an ephemeral IPFS node
	return createNode(ctx, path)
}
