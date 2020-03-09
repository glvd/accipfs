package service

import (
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/dhcrypto"
	"go.uber.org/atomic"
	"time"
)

// HandleInfo ...
type HandleInfo struct {
	ServiceName string
	Data        interface{}
	Callback    HandleCallback
}

// HandleCallback ...
type HandleCallback func(src interface{})

// Node ...
type Node interface {
	Start()
}

var dateKey = time.Date(2019, time.November, 11, 10, 20, 10, 300, time.Local)

type serviceNode struct {
	lock *atomic.Bool
}

func nodeInstance() *serviceNode {
	return &serviceNode{lock: atomic.NewBool(false)}
}
func decodeNodes(cfg config.Config, nodes []string) []string {
	// init contract
	var decodeNodes []string
	decoder := dhcrypto.NewCipherDecode([]byte(cfg.PrivateKey), dateKey)
	if len(nodes) == 0 {
		return decodeNodes
	}
	for _, node := range nodes {
		decoded, err := decoder.Decode(node)
		if err != nil {
			continue
		}
		decodeNodes = append(decodeNodes, string(decoded))
	}
	return decodeNodes
}

func encodeNodes(cfg config.Config, nodes []string) []string {
	var encodedNodes []string
	encoder := dhcrypto.NewCipherEncoder([]byte(cfg.PublicKey), 10, dateKey)
	for _, node := range nodes {
		encoded, err := encoder.Encode(node)
		if err != nil {
			continue
		}
		encodedNodes = append(encodedNodes, string(encoded))
	}
	return encodedNodes
}
