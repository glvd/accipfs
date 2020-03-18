package contract

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/contract/node"
	"github.com/glvd/accipfs/contract/token"
	"io/ioutil"
	"path/filepath"
)

const keyStore = `{"address":"945d35cd4a6549213e8d37feb5d708ec98906902","crypto":{"cipher":"aes-128-ctr","ciphertext":"649f5c7def3f345c39dc6f10e5438e179a5f06ff1d9ef2467ff7c84ec94f1a2a","cipherparams":{"iv":"0d66dfbc2c978ed1989e2fca05c16abe"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"547ed9895deda897adbe09058ebfb24fb5036695d490c2127da45c4f7ec9e4a8"},"mac":"db76804c69ceb8705de1a73ae0caf4761bd73c3d42aa43f801c03e7fdda6adff"},"id":"9aaeec2d-d639-425a-83f7-a0956dcc78a1","version":3}`

// instance ...
type instance struct {
	cfg       *config.Config
	cli       *ethclient.Client
	nodeAddr  common.Address
	tokenAddr common.Address
	key       *ecdsa.PrivateKey
}

// Contractor ...
type Contractor interface {
	Node(call NodeCall) error
	Token(call TokenCall) error
}

// NodeCall ...
type NodeCall func(node *node.AccelerateNode, opts *bind.TransactOpts) error

// TokenCall ...
type TokenCall func(token *token.DhToken, opts *bind.TransactOpts) error

// HexKey ...
func HexKey(cfg config.Config) *ecdsa.PrivateKey {
	privateKey, err := crypto.HexToECDSA(cfg.ETH.KeyHash)
	if err != nil {
		panic(err)
	}
	return privateKey
}

// FileKey ...
func FileKey(cfg *config.Config) *ecdsa.PrivateKey {
	newAccount, e := account.NewAccount(cfg)
	if e != nil {
		panic(e)
	}

	bys, e := ioutil.ReadFile(filepath.Join(cfg.Path, newAccount.Address))
	if e != nil {
		panic(e)
	}

	keys, err := keystore.DecryptKey(bys, cfg.ETH.Pass)
	if err != nil {
		panic(e)
	}
	return keys.PrivateKey
}

// Loader ...
func Loader(cfg *config.Config) Contractor {
	return &instance{
		cfg:       cfg,
		nodeAddr:  common.HexToAddress(cfg.ETH.NodeAddr),
		tokenAddr: common.HexToAddress(cfg.ETH.TokenAddr),
		key:       FileKey(cfg),
	}
}

//Node contract: Node init acceleratenode contract
func (c *instance) Node(call NodeCall) error {
	o := bind.NewKeyedTransactor(c.key)

	// gateway redirect to private chain
	client, err := ethclient.Dial(config.ETHAddr())
	if err != nil {
		return err
	}
	defer client.Close()
	instance, err := node.NewAccelerateNode(c.nodeAddr, client)
	if err != nil {
		return err
	}

	return call(instance, o)
}

//Token contract: Token init DHToken contract
func (c *instance) Token(call TokenCall) error {
	o := bind.NewKeyedTransactor(c.key)

	// gateway redirect to private chain
	client, err := ethclient.Dial(config.ETHAddr())
	if err != nil {
		return err
	}
	instance, err := token.NewDhToken(c.tokenAddr, client)
	if err != nil {
		return err
	}
	return call(instance, o)
}
