package contract

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/contract/node"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/glvd/accipfs/contract/token"
)

const keyStore = `{"address":"945d35cd4a6549213e8d37feb5d708ec98906902","crypto":{"cipher":"aes-128-ctr","ciphertext":"649f5c7def3f345c39dc6f10e5438e179a5f06ff1d9ef2467ff7c84ec94f1a2a","cipherparams":{"iv":"0d66dfbc2c978ed1989e2fca05c16abe"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"547ed9895deda897adbe09058ebfb24fb5036695d490c2127da45c4f7ec9e4a8"},"mac":"db76804c69ceb8705de1a73ae0caf4761bd73c3d42aa43f801c03e7fdda6adff"},"id":"9aaeec2d-d639-425a-83f7-a0956dcc78a1","version":3}`
const nodeContractAddr = "0xbaEEB7a3AF34a365ACAa1f8464A3374B58ac9889"
const tokenContractAddr = "0x9064322CfeE623A447ba5aF0dA6AD3341c073535"

// instance ...
type instance struct {
	cfg      config.Config
	cli      *ethclient.Client
	keystore string
}

// Contractor ...
type Contractor interface {
	Node() (*node.AccelerateNode, *bind.TransactOpts, *ethclient.Client)
	Token() (*token.DhToken, *bind.TransactOpts, *ethclient.Client)
}

// NodeCall ...
type NodeCall func(node *node.AccelerateNode, opts *bind.TransactOpts)

// TokenCall ...
type TokenCall func(token *token.DhToken, opts *bind.TransactOpts)

// Loader ...
func Loader(cfg config.Config) Contractor {
	return &instance{
		cfg:      cfg,
		keystore: keyStore,
	}
}

//Node contract: Node init acceleratenode contract
func (c *instance) Node() (*node.AccelerateNode, *bind.TransactOpts, *ethclient.Client) {
	// TODO
	auth, err := bind.NewTransactor(strings.NewReader(c.keystore), "123")
	if err != nil {
		log.Fatal(err)
	}
	// gateway redirect to private chain
	// client, err := ethclient.Dial("http://gate.betabb.space:8545")
	client, err := ethclient.Dial(c.cfg.ETH.Addr)
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress(nodeContractAddr)
	instance, err := node.NewAccelerateNode(address, client)
	if err != nil {
		log.Fatal(err)
	}

	return instance, auth, client
}

//Token contract: Token init DHToken contract
func (c *instance) Token() (*token.DhToken, *bind.TransactOpts, *ethclient.Client) {
	// TODO
	auth, err := bind.NewTransactor(strings.NewReader(c.keystore), "123")
	if err != nil {
		log.Fatal(err)
	}
	// gateway redirect to private chain
	client, err := ethclient.Dial(c.cfg.ETH.Addr)
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress(tokenContractAddr)
	instance, err := token.NewDhToken(address, client)
	if err != nil {
		log.Fatal(err)
	}

	return instance, auth, client
}
