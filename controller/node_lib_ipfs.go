package controller

import (
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"github.com/ipfs/interface-go-ipfs-core/options"
	"os"
	"path/filepath"

	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	ipfsconfig "github.com/ipfs/go-ipfs-config"
	ipfscore "github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/plugin/loader"
	intercore "github.com/ipfs/interface-go-ipfs-core"
)

type nodeLibIPFS struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        *config.Config
	configRoot string
	api        intercore.CoreAPI
}

var _ core.ControllerService = &nodeLibIPFS{}

func newNodeLibIPFS(cfg *config.Config) *nodeLibIPFS {
	ctx, cancel := context.WithCancel(context.Background())
	root := filepath.Join(cfg.Path, ".ipfs")
	return &nodeLibIPFS{
		ctx:        ctx,
		cancel:     cancel,
		cfg:        cfg,
		configRoot: root,
	}
}

// Start ...
func (n *nodeLibIPFS) Start() error {
	node, err := createNode(n.ctx, n.configRoot)
	if err != nil {
		return err
	}
	n.api = node
	return nil
}

// Stop ...
func (n *nodeLibIPFS) Stop() error {
	if n.cancel != nil {
		n.cancel()
		n.cancel = nil
	}
	return nil
}

// Initialize ...
func (n nodeLibIPFS) Initialize() error {
	_ = os.Mkdir(n.configRoot, 0755)
	return n.spawnEphemeral(n.ctx)
}

// IsReady ...
func (n *nodeLibIPFS) IsReady() bool {
	return true
}

// MessageHandle ...
func (n *nodeLibIPFS) MessageHandle(f func(s string)) {

}

func setupPlugins(externalPluginsPath string) error {
	// Load any external plugins if available on externalPluginsPath
	plugins, err := loader.NewPluginLoader(filepath.Join(externalPluginsPath, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	// Load preloaded and external plugins
	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

// Spawns a node to be used just for this run (i.e. creates a tmp repo)
func (n *nodeLibIPFS) spawnEphemeral(ctx context.Context) error {
	if err := setupPlugins(""); err != nil {
		return err
	}

	// Create a Temporary Repo
	if err := n.createRepo(ctx); err != nil {
		return fmt.Errorf("failed to create temp repo: %s", err)
	}
	// Spawning an ephemeral IPFS node
	return nil
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

	// Create the repo with the config
	err = fsrepo.Init(n.configRoot, cfg)
	if err != nil {
		return fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return nil
}

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, repoPath string) (intercore.CoreAPI, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node

	nodeOptions := &ipfscore.BuildCfg{
		Online:  true,
		Routing: libp2p.NilRouterOption,
		//Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: repo,
	}

	node, err := ipfscore.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}

	// Attach the Core API to the constructed node
	return coreapi.NewCoreAPI(node)
}
