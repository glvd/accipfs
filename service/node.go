package service

import (
	"bug.vlavr.com/godcong/dhcrypto"
	"github.com/glvd/accipfs/config"
	"go.uber.org/atomic"
	"net"
	"strings"
	"time"
)

// Node ...
type Node interface {
	Run()
}

var dateKey = time.Date(2019, time.November, 11, 10, 20, 10, 300, time.Local)

type serviceNode struct {
	lock *atomic.Bool
}

func nodeInstance() *serviceNode {
	return &serviceNode{lock: atomic.NewBool(false)}
}

func decodeNodes(cfg *config.Config, nodes []string) []string {
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

func encodeNodes(cfg *config.Config, nodes []string) []string {
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

func getAccessibleEthNodes(addresses []string, port string, to time.Duration) []string {
	var accessible []string
	for _, address := range addresses {
		strs := strings.Split(address, "@")
		if len(strs) < 2 {
			continue
		}
		url := strs[1]
		ip := strings.Split(url, ":")[0]

		conn, e := net.DialTimeout("tcp", ip+":"+port, to)
		if e == nil {
			addr := strs[0] + "@" + ip + ":" + port
			accessible = append(accessible, addr)
			_ = conn.Close()
		}
	}
	return accessible
}
