package controller

import (
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	"io/ioutil"
	"path/filepath"

	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	ipfsconfig "github.com/ipfs/go-ipfs-config"
	ipfscore "github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/plugin/loader"
	interfacecore "github.com/ipfs/interface-go-ipfs-core"
)

type nodeLibIPFS struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var _ core.ControllerService = &nodeLibIPFS{}

func newNodeLibIPFS(cfg *config.Config) *nodeLibIPFS {
	ctx, cancel := context.WithCancel(context.Background())

	return &nodeLibIPFS{
		ctx:    ctx,
		cancel: cancel,
	}
}

// Start ...
func (n nodeLibIPFS) Start() error {
	panic("implement me")
}

// Stop ...
func (n nodeLibIPFS) Stop() error {
	panic("implement me")
}

// Initialize ...
func (n nodeLibIPFS) Initialize() error {
	panic("implement me")
}

// IsReady ...
func (n nodeLibIPFS) IsReady() bool {
	panic("implement me")
}

// MessageHandle ...
func (n nodeLibIPFS) MessageHandle(f func(s string)) {
	panic("implement me")
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
func spawnEphemeral(ctx context.Context) (interfacecore.CoreAPI, error) {
	if err := setupPlugins(""); err != nil {
		return nil, err
	}

	// Create a Temporary Repo
	repoPath, err := createTempRepo(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp repo: %s", err)
	}

	// Spawning an ephemeral IPFS node
	return createNode(ctx, repoPath)
}

func createTempRepo(ctx context.Context) (string, error) {
	repoPath, err := ioutil.TempDir("", "ipfs-shell")
	if err != nil {
		return "", fmt.Errorf("failed to get temp dir: %s", err)
	}

	// Create a config with default options and a 2048 bit key
	cfg, err := ipfsconfig.Init(ioutil.Discard, 2048)
	if err != nil {
		return "", err
	}

	// Create the repo with the config
	err = fsrepo.Init(repoPath, cfg)
	if err != nil {
		return "", fmt.Errorf("failed to init ephemeral node: %s", err)
	}

	return repoPath, nil
}

// Creates an IPFS node and returns its coreAPI
func createNode(ctx context.Context, repoPath string) (interfacecore.CoreAPI, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, err
	}

	// Construct the node

	nodeOptions := &ipfscore.BuildCfg{
		Online:  true,
		Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
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
