module github.com/glvd/accipfs

go 1.13

require (
	github.com/dgraph-io/badger/v2 v2.0.3
	github.com/elastic/gosigar v0.10.5 // indirect
	github.com/emirpasic/gods v1.12.0
	github.com/ethereum/go-ethereum v1.9.11
	github.com/gin-gonic/gin v1.6.2
	github.com/godcong/scdt v0.0.12
	github.com/goextension/extmap v0.0.1
	github.com/goextension/io v0.0.0-20191016080154-50dbafac3df3
	github.com/goextension/tool v0.0.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/rpc v1.2.0
	github.com/ipfs/bbloom v0.0.4 // indirect
	github.com/ipfs/go-cid v0.0.6 // indirect
	github.com/ipfs/go-ipfs-http-client v0.0.5
	github.com/ipfs/interface-go-ipfs-core v0.2.3
	github.com/libp2p/go-libp2p-core v0.6.0
	github.com/libp2p/go-openssl v0.0.6 // indirect
	github.com/miekg/dns v1.1.12
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-multiaddr v0.2.2
	github.com/multiformats/go-multiaddr-net v0.1.0
	github.com/multiformats/go-multihash v0.0.14 // indirect
	github.com/panjf2000/ants/v2 v2.4.1
	github.com/portmapping/go-reuse v0.0.3
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.3.2
	go.uber.org/atomic v1.6.0
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/net v0.0.0-20200226121028-0de0cce0169b
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	golang.org/x/tools v0.0.0-20200702044944-0cc1aa72b347 // indirect
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
)

replace github.com/libp2p/go-libp2p-core v0.6.0 => git.5gnode.cn/chain/go-libp2p-core v0.0.0-20200706033803-a09882e801a0
