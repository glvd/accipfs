module github.com/glvd/accipfs

go 1.13

require (
	github.com/DataDog/zstd v1.4.5 // indirect
	github.com/dgraph-io/badger/v2 v2.0.3
	github.com/dgraph-io/ristretto v0.0.3 // indirect
	github.com/dgryski/go-farm v0.0.0-20200201041132-a6ae2369ad13 // indirect
	github.com/dustin/go-humanize v1.0.0
	github.com/elastic/gosigar v0.10.5 // indirect
	github.com/ethereum/go-ethereum v1.9.11
	github.com/gin-gonic/gin v1.6.2
	github.com/godcong/scdt v0.0.20
	github.com/goextension/extmap v0.0.1
	github.com/goextension/io v0.0.0-20191016080154-50dbafac3df3
	github.com/goextension/tool v0.0.2
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/rpc v1.2.0
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-ds-badger2 v0.1.0
	github.com/ipfs/go-ds-flatfs v0.4.4
	github.com/ipfs/go-ds-leveldb v0.4.2
	github.com/ipfs/go-ipfs v0.6.0
	github.com/ipfs/go-ipfs-config v0.9.0
	github.com/ipfs/go-ipfs-files v0.0.8
	github.com/ipfs/go-ipfs-http-client v0.0.5
	github.com/ipfs/go-ipld-format v0.2.0
	github.com/ipfs/go-ipld-git v0.0.3
	github.com/ipfs/go-log v1.0.4
	github.com/ipfs/go-log/v2 v2.1.1 // indirect
	github.com/ipfs/go-metrics-prometheus v0.0.2
	github.com/ipfs/interface-go-ipfs-core v0.3.0
	github.com/jbenet/goprocess v0.1.4
	github.com/libp2p/go-libp2p v0.9.6
	github.com/libp2p/go-libp2p-core v0.5.7
	github.com/libp2p/go-libp2p-peerstore v0.2.6
	github.com/libp2p/go-openssl v0.0.6 // indirect
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/miekg/dns v1.1.29
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-multiaddr v0.2.2
	github.com/multiformats/go-multiaddr-dns v0.2.0
	github.com/multiformats/go-multiaddr-net v0.1.5
	github.com/multiformats/go-multihash v0.0.14
	github.com/opentracing/opentracing-go v1.2.0
	github.com/panjf2000/ants/v2 v2.4.1
	github.com/polydawn/refmt v0.0.0-20190807091052-3d65705ee9f1 // indirect
	github.com/portmapping/go-reuse v0.0.3
	github.com/prometheus/tsdb v0.7.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.3.2
	github.com/syndtr/goleveldb v1.0.1-0.20190923125748-758128399b1d
	github.com/whyrusleeping/cbor-gen v0.0.0-20200706173030-3bb387cdd4d1 // indirect
	go.opencensus.io v0.22.4 // indirect
	go.uber.org/atomic v1.6.0
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/sys v0.0.0-20200727154430-2d971f7391a4 // indirect
	golang.org/x/tools v0.0.0-20200702044944-0cc1aa72b347 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	honnef.co/go/tools v0.0.1-2020.1.4 // indirect
)

replace github.com/ipfs/go-ipfs-http-client v0.0.5 => github.com/godcong/go-ipfs-http-client v0.0.11
