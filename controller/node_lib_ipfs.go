package controller

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/basis"
	migrate "github.com/ipfs/go-ipfs/repo/fsrepo/migrations"
	"go.uber.org/atomic"

	"os"
	"path/filepath"

	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	ipfsconfig "github.com/ipfs/go-ipfs-config"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	intercore "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/options"
)

type nodeLibIPFS struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        *config.Config
	isRunning  *atomic.Bool
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
		isRunning:  atomic.NewBool(false),
	}
}

// Start ...
func (n *nodeLibIPFS) Start() error {
	if err := basis.SetupPlugins(""); err != nil {
		return err
	}

	if !fsrepo.IsInitialized(n.configRoot) {
		if err := n.Initialize(); err != nil {
			return err
		}
	}
	_, err := fsrepo.Open(n.configRoot)
	switch err {
	default:
		return err
	case fsrepo.ErrNeedMigration:
		err = migrate.RunMigration(fsrepo.RepoVersion)
		if err != nil {
			fmt.Println("The migrations of fs-repo failed:")
			fmt.Printf("  %s\n", err)
			fmt.Println("If you think this is a bug, please file an issue and include this whole log output.")
			fmt.Println("  https://github.com/ipfs/fs-repo-migrations")
			return err
		}
	}
	// Spawning an ephemeral IPFS node
	node, err := basis.CreateNode(n.ctx, n.configRoot)
	if err != nil {
		return err
	}
	n.api = node
	n.isRunning.Store(true)
	fmt.Println("datastore is ready")
	return nil
}

// Stop ...
func (n *nodeLibIPFS) Stop() error {
	if n.cancel != nil {
		n.cancel()
		n.cancel = nil
	}
	n.isRunning.Store(false)
	return nil
}

// Initialize ...
func (n nodeLibIPFS) Initialize() error {
	_ = os.Mkdir(n.configRoot, 0755)
	if err := basis.SetupPlugins(""); err != nil {
		return err
	}
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

// Spawns a node to be used just for this run (i.e. creates a tmp repo)
func (n *nodeLibIPFS) spawnEphemeral(ctx context.Context) error {

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
