package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/glvd/accipfs/config"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
)

// DefaultClient ...
var DefaultClient = ""

type client struct {
	cfg *config.APIConfig
	cli *http.Client
}

func requestReader(req interface{}) (io.Reader, error) {
	marshal, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(marshal), nil
}

func responseDecoder(rc io.ReadCloser, resp interface{}) error {
	decoder := json.NewDecoder(rc)
	return decoder.Decode(resp)
}

// New ...
func New(cfg *config.APIConfig) core.API {
	c := &http.Client{}
	c.Timeout = cfg.Timeout
	return &client{
		cli: c,
		cfg: cfg,
	}
}

func (c *client) host() string {
	prefix := "http://"
	if c.cfg.UseTLS {
		prefix = "https://"
	}
	return strings.Join([]string{prefix, "127.0.0.1:", strconv.Itoa(c.cfg.Port), "/api/", c.cfg.Version}, "")
}

// RequestURL ...
func (c *client) RequestURL(uri string) string {
	if uri[0] == '/' {
		uri = uri[1:]
	}
	return strings.Join([]string{c.host(), uri}, "/")
}

// Ping ...
func (c *client) Ping(req *core.PingReq) (*core.PingResp, error) {
	reader, err := requestReader(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, c.RequestURL("ping"), reader)
	if err != nil {
		return nil, err
	}
	resp, err := c.cli.Do(request)
	if err != nil {
		return nil, err
	}
	result := new(core.PingResp)
	err = responseDecoder(resp.Body, result)
	return result, err
}

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

// ConnectTo ...
func ConnectTo(url string, req *core.ConnectToReq) (*core.ConnectToResp, error) {
	remoteNode := new(core.ConnectToResp)
	if err := general.RPCPost(url, "BustLinker.ConnectTo", req, remoteNode); err != nil {
		return nil, err
	}
	return remoteNode, nil
}

// Pins ...
func Pins(address core.Addr) ([]string, error) {
	logD("ping info", "addr", address.IP, "port", address.Port)
	pingAddr := strings.Join([]string{address.IP.String(), strconv.Itoa(address.Port)}, ":")
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
	//if len(*result) == 0 {
	//	return nil, fmt.Errorf("no data response")
	//}
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
		return fmt.Errorf("connect failed: %s", general.RPCAddress(info.Addr[0]))
	}
	return nil
}
