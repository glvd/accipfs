package core

// NodeType ...
type NodeType int

const (
	// NodeUndefined ...
	NodeUndefined NodeType = -1
	// NodeAccount ...
	NodeAccount NodeType = iota
)

// NodeInfo ...
type NodeInfo struct {
	Name       string
	Schema     string
	RemoteAddr string
	NodeType   string
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
