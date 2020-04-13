package client

import (
	"fmt"
	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
	"strconv"
	"strings"
)

// ID ...
func ID(url string) (*core.Node, error) {
	reply := new(core.Node)
	if err := general.RPCPost(url, "BustLinker.ID", core.DummyEmpty(), reply); err != nil {
		return nil, err
	}
	return reply, nil
}

// Ping ...
func Ping(url string) error {
	logD("ping info", "url", url)
	result := new(core.PingResp)
	if err := general.RPCPost(url, "BustLinker.Ping", &core.PingReq{}, result); err != nil {
		return err
	}
	if result.Resp != "pong" {
		return fmt.Errorf("get wrong response data:%+v", *result)
	}
	return nil
}

// Pins ...
func Pins(address core.NodeAddress) ([]string, error) {
	logD("ping info", "addr", address.Address, "port", address.Port)
	pingAddr := strings.Join([]string{address.Address, strconv.Itoa(address.Port)}, ":")
	url := fmt.Sprintf("http://%s/rpc", pingAddr)
	result := new([]string)
	if err := general.RPCPost(url, "BustLinker.Pins", core.DummyEmpty(), result); err != nil {
		return nil, err
	}
	return *result, nil
}

// PinVideo ...
func PinVideo(url string, no string) error {
	logD("pin hash", "hash", no)
	b := new(bool)
	err := general.RPCPost(url, "BustLinker.PinVideo", &no, b)
	if err != nil {
		return err
	}
	if *b {
		fmt.Printf("pin (%s) success\n", no)
	}
	return nil
}

// Peers ...
func Peers(url string, info *core.Node) ([]*core.Node, error) {
	//pingAddr := strings.Join([]string{info.RemoteAddr, strconv.Itoa(info.Port)}, ":")
	//url := fmt.Sprintf("http://%s/rpc", pingAddr)
	result := new([]*core.Node)
	if err := general.RPCPost(url, "BustLinker.Peers", info, result); err != nil {
		return nil, err
	}
	if len(*result) == 0 {
		return nil, fmt.Errorf("no data response")
	}
	return *result, nil
}

// AddPeer ...
func AddPeer(url string, info *core.Node) error {
	status := new(bool)
	if err := general.RPCPost(url, "BustLinker.AddPeer", info, status); err != nil {
		logE("remote id error", "error", err.Error())
		return fmt.Errorf("remote id error: %w", err)
	}

	if !(*status) {
		return fmt.Errorf("connect failed:%s", url)
	}
	return nil
}
