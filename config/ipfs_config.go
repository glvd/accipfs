package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// IPFSServerConfig ...
type IPFSServerConfig struct {
	API          API          `json:"API"`
	Addresses    Addresses    `json:"Addresses"`
	Bootstrap    []string     `json:"Bootstrap"`
	Datastore    Datastore    `json:"Datastore"`
	Discovery    Discovery    `json:"Discovery"`
	Experimental Experimental `json:"Experimental"`
	Gateway      Gateway      `json:"Gateway"`
	Identity     Identity     `json:"Identity"`
	Ipns         Ipns         `json:"Ipns"`
	Mounts       Mounts       `json:"Mounts"`
	Pubsub       Pubsub       `json:"Pubsub"`
	Reprovider   Reprovider   `json:"Reprovider"`
	Routing      Routing      `json:"Routing"`
	Swarm        Swarm        `json:"Swarm"`
}

// API ...
type API struct {
	HTTPHeaders APIHTTPHeaders `json:"HTTPHeaders"`
}

// APIHTTPHeaders ...
type APIHTTPHeaders struct {
}

// Addresses ...
type Addresses struct {
	API        string        `json:"API"`
	Announce   []interface{} `json:"Announce"`
	Gateway    string        `json:"Gateway"`
	NoAnnounce []interface{} `json:"NoAnnounce"`
	Swarm      []string      `json:"Swarm"`
}

// Datastore ...
type Datastore struct {
	BloomFilterSize    int64  `json:"BloomFilterSize"`
	GCPeriod           string `json:"GCPeriod"`
	HashOnRead         bool   `json:"HashOnRead"`
	Spec               Spec   `json:"Spec"`
	StorageGCWatermark int64  `json:"StorageGCWatermark"`
	StorageMax         string `json:"StorageMax"`
}

// Spec ...
type Spec struct {
	Child  Child  `json:"child"`
	Prefix string `json:"prefix"`
	Type   string `json:"type"`
}

// Child ...
type Child struct {
	Path       string `json:"path"`
	SyncWrites bool   `json:"syncWrites"`
	Truncate   bool   `json:"truncate"`
	Type       string `json:"type"`
}

// Discovery ...
type Discovery struct {
	Mdns Mdns `json:"MDNS"`
}

// Mdns ...
type Mdns struct {
	Enabled  bool  `json:"Enabled"`
	Interval int64 `json:"Interval"`
}

// Experimental ...
type Experimental struct {
	FilestoreEnabled     bool `json:"FilestoreEnabled"`
	Libp2PStreamMounting bool `json:"Libp2pStreamMounting"`
	P2PHTTPProxy         bool `json:"P2pHttpProxy"`
	PreferTLS            bool `json:"PreferTLS"`
	Quic                 bool `json:"QUIC"`
	ShardingEnabled      bool `json:"ShardingEnabled"`
	UrlstoreEnabled      bool `json:"UrlstoreEnabled"`
}

// Gateway ...
type Gateway struct {
	APICommands  []interface{}      `json:"APICommands"`
	HTTPHeaders  GatewayHTTPHeaders `json:"HTTPHeaders"`
	NoFetch      bool               `json:"NoFetch"`
	PathPrefixes []interface{}      `json:"PathPrefixes"`
	RootRedirect string             `json:"RootRedirect"`
	Writable     bool               `json:"Writable"`
}

// GatewayHTTPHeaders ...
type GatewayHTTPHeaders struct {
	AccessControlAllowHeaders []string `json:"Access-Control-Allow-Headers"`
	AccessControlAllowMethods []string `json:"Access-Control-Allow-Methods"`
	AccessControlAllowOrigin  []string `json:"Access-Control-Allow-Origin"`
}

// Identity ...
type Identity struct {
	PeerID  string `json:"PeerID"`
	PrivKey string `json:"PrivKey"`
}

// Ipns ...
type Ipns struct {
	RecordLifetime   string `json:"RecordLifetime"`
	RepublishPeriod  string `json:"RepublishPeriod"`
	ResolveCacheSize int64  `json:"ResolveCacheSize"`
}

// Mounts ...
type Mounts struct {
	FuseAllowOther bool   `json:"FuseAllowOther"`
	Ipfs           string `json:"IPFS"`
	Ipns           string `json:"IPNS"`
}

// Pubsub ...
type Pubsub struct {
	DisableSigning              bool   `json:"DisableSigning"`
	Router                      string `json:"Router"`
	StrictSignatureVerification bool   `json:"StrictSignatureVerification"`
}

// Reprovider ...
type Reprovider struct {
	Interval string `json:"Interval"`
	Strategy string `json:"Strategy"`
}

// Routing ...
type Routing struct {
	Type string `json:"Type"`
}

// Swarm ...
type Swarm struct {
	AddrFilters             interface{} `json:"AddrFilters"`
	ConnMgr                 ConnMgr     `json:"ConnMgr"`
	DisableBandwidthMetrics bool        `json:"DisableBandwidthMetrics"`
	DisableNatPortMap       bool        `json:"DisableNatPortMap"`
	DisableRelay            bool        `json:"DisableRelay"`
	EnableAutoNATService    bool        `json:"EnableAutoNATService"`
	EnableAutoRelay         bool        `json:"EnableAutoRelay"`
	EnableRelayHop          bool        `json:"EnableRelayHop"`
}

// ConnMgr ...
type ConnMgr struct {
	GracePeriod string `json:"GracePeriod"`
	HighWater   int64  `json:"HighWater"`
	LowWater    int64  `json:"LowWater"`
	Type        string `json:"Type"`
}

// LoadIPFSServerConfig ...
func LoadIPFSServerConfig(cfg Config) (*IPFSServerConfig, error) {
	var c IPFSServerConfig
	path := filepath.Join(cfg.Path, _dataDirIPFS, "config")
	open, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	dec := json.NewDecoder(open)
	err = dec.Decode(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
