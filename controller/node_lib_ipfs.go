package controller

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/plugin/loader"
	ipfscore "github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/jbenet/goprocess"
	"go.uber.org/atomic"
	"runtime"
	"sort"

	"os"
	"path/filepath"

	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	ipfsversion "github.com/ipfs/go-ipfs"
	ipfsconfig "github.com/ipfs/go-ipfs-config"
	utilmain "github.com/ipfs/go-ipfs/cmd/ipfs/util"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	mprome "github.com/ipfs/go-metrics-prometheus"
	intercore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

type nodeLibIPFS struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        *config.Config
	isRunning  *atomic.Bool
	configRoot string
	intercore.CoreAPI
	repo    repo.Repo
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
	plugins, err := setupPlugins("")
	if err != nil {
		panic(err)
	}

	return &nodeLibIPFS{
		ctx:        ctx,
		cancel:     cancel,
		cfg:        cfg,
		configRoot: root,
		isRunning:  atomic.NewBool(false),
		plugins:    plugins,
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

// printSwarmAddrs prints the addresses of the host
func printSwarmAddrs(node *ipfscore.IpfsNode) {
	if !node.IsOnline {
		fmt.Println("Swarm not listening, running in offline mode.")
		return
	}

	var lisAddrs []string
	ifaceAddrs, err := node.PeerHost.Network().InterfaceListenAddresses()
	if err != nil {
		fmt.Printf("failed to read listening addresses: %s", err)
	}
	for _, addr := range ifaceAddrs {
		lisAddrs = append(lisAddrs, addr.String())
	}
	sort.Strings(lisAddrs)
	for _, addr := range lisAddrs {
		fmt.Printf("Swarm listening on %s\n", addr)
	}

	var addrs []string
	for _, addr := range node.PeerHost.Addrs() {
		addrs = append(addrs, addr.String())
	}
	sort.Strings(addrs)
	for _, addr := range addrs {
		fmt.Printf("Swarm announcing %s\n", addr)
	}

}

// Start ...
func (n *nodeLibIPFS) Start() (_err error) {
	err := mprome.Inject()
	if err != nil {
		log.Errorf("Injecting prometheus handler for metrics failed with message: %s\n", err.Error())
	}
	// let the user know we're going.
	fmt.Printf("Initializing daemon...\n")

	defer func() {
		if _err != nil {
			// Print an extra line before any errors. This could go
			// in the commands lib but doesn't really make sense for
			// all commands.
			fmt.Println(_err)
		}
	}()

	printVersion()

	if true {
		if _, _, err := utilmain.ManageFdLimit(); err != nil {
			log.Errorf("setting file descriptor limit: %s", err)
		}
	}

	if !fsrepo.IsInitialized(n.configRoot) {
		if err := n.Initialize(); err != nil {
			return err
		}
	}

	repo, err := basis.OpenRepo(n.configRoot)
	if err != nil {
		return err
	}
	n.repo = repo

	node, err := basis.CreateNode(n.ctx, repo)
	if err != nil {
		return err
	}
	printSwarmAddrs(node)
	// Attach the Core API to the constructed node

	_err = n.plugins.Start(node)
	if _err != nil {
		return _err
	}
	node.Process.AddChild(goprocess.WithTeardown(n.plugins.Close))
	api, _err := coreapi.NewCoreAPI(node)
	if _err != nil {
		return _err
	}
	n.CoreAPI = api
	n.isRunning.Store(true)
	log.Infow("datastore is ready")
	return nil
}

// Stop ...
func (n *nodeLibIPFS) Stop() error {
	if n.cancel != nil {
		n.cancel()
		n.cancel = nil
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
	log.Infow("datastore init")
	_ = os.Mkdir(n.configRoot, 0755)
	//if err := basis.SetupPlugins(""); err != nil {
	//	return err
	//}
	// Create a Temporary Repo
	if err := n.createRepo(n.ctx); err != nil {
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

func (n *nodeLibIPFS) createRepo(ctx context.Context) error {
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
	err = fsrepo.Init(n.configRoot, cfg)
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
