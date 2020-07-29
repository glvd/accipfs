package basis

import (
	"context"
	"fmt"
	ipfscore "github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipfs/go-ipfs/repo/fsrepo"
	migrate "github.com/ipfs/go-ipfs/repo/fsrepo/migrations"
)

// OpenRepo ...
func OpenRepo(repoPath string) (repo.Repo, error) {
	// Open the repo
	repo, err := fsrepo.Open(repoPath)
	switch err {
	default:
		return nil, err
	case fsrepo.ErrNeedMigration:
		err = migrate.RunMigration(fsrepo.RepoVersion)
		if err != nil {
			fmt.Println("The migrations of fs-repo failed:")
			fmt.Printf("  %s\n", err)
			fmt.Println("If you think this is a bug, please file an issue and include this whole log output.")
			fmt.Println("  https://github.com/ipfs/fs-repo-migrations")
			return nil, err
		}
		repo, err = fsrepo.Open(repoPath)
		if err != nil {
			return nil, err
		}
	}
	return repo, nil
}

// CreateNode Creates an IPFS node and returns its coreAPI
func CreateNode(ctx context.Context, r repo.Repo) (*ipfscore.IpfsNode, error) {

	// Construct the node
	nodeOptions := &ipfscore.BuildCfg{
		Online:                      true,
		Routing:                     libp2p.NilRouterOption,
		Permanent:                   true, // It is temporary way to signify that node is permanent
		DisableEncryptedConnections: false,
		ExtraOpts: map[string]bool{
			"pubsub": true,
			"ipnsps": true,
		},
		//Routing: libp2p.DHTOption, // This option sets the node to be a full DHT node (both fetching and storing DHT Records)
		// Routing: libp2p.DHTClientOption, // This option sets the node to be a client DHT node (only fetching records)
		Repo: r,
	}

	node, err := ipfscore.NewNode(ctx, nodeOptions)
	if err != nil {
		return nil, err
	}
	return node, nil
}
