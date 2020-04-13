package core

const (
	// NodeUndefined ...
	NodeUndefined NodeType = -1
	// NodeAccount ...
	NodeAccount NodeType = 0x01
	// NodeAccelerate ...
	NodeAccelerate NodeType = 0x02
)

// Version ...
const Version = "0.0.1"

// NodeType ...
type NodeType int

// DataStoreNode ...
type DataStoreNode struct {
	ID              string   `json:"ID"`
	PublicKey       string   `json:"PublicKey"`
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ProtocolVersion string   `json:"ProtocolVersion"`
}

// ContractNode ...
type ContractNode struct {
	Enode      string    `json:"enode"`
	Enr        string    `json:"enr"`
	ID         string    `json:"id"`
	IP         string    `json:"ip"`
	ListenAddr string    `json:"listenAddr"`
	Name       string    `json:"name"`
	Ports      Ports     `json:"ports"`
	Protocols  Protocols `json:"protocols"`
}

// Ports ...
type Ports struct {
	Discovery int64 `json:"discovery"`
	Listener  int64 `json:"listener"`
}

// Protocols ...
type Protocols struct {
	Eth Eth `json:"eth"`
}

// Eth ...
type Eth struct {
	Config     Config `json:"config"`
	Difficulty int64  `json:"difficulty"`
	Genesis    string `json:"genesis"`
	Head       string `json:"head"`
	Network    int64  `json:"network"`
}

// Config ...
type Config struct {
	ByzantiumBlock      int64  `json:"byzantiumBlock"`
	ChainID             int64  `json:"chainId"`
	Clique              Clique `json:"clique"`
	ConstantinopleBlock int64  `json:"constantinopleBlock"`
	Eip150Block         int64  `json:"eip150Block"`
	Eip150Hash          string `json:"eip150Hash"`
	Eip155Block         int64  `json:"eip155Block"`
	Eip158Block         int64  `json:"eip158Block"`
	HomesteadBlock      int64  `json:"homesteadBlock"`
}

// Clique ...
type Clique struct {
	Epoch  int64 `json:"epoch"`
	Period int64 `json:"period"`
}

// NodeInfo ...
type NodeInfo struct {
	Name       string
	Schema     string
	RemoteAddr string
	NodeType   NodeType
	Port       int
	Contract   ContractNode
	DataStore  DataStoreNode
	Version    string
}

// Address ...
func (n *NodeInfo) Address() *AddressInfo {
	return &AddressInfo{
		Address: n.RemoteAddr,
		Schema:  n.Schema,
		Port:    n.Port,
	}
}
