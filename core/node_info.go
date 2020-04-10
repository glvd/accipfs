package core

// NodeType ...
type NodeType int

const (
	// NodeUndefined ...
	NodeUndefined NodeType = -1
	// NodeAccount ...
	NodeAccount NodeType = 0x01
	// NodeAccelerate ...
	NodeAccelerate NodeType = 0x02
)

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
