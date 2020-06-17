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
	"time"

	"github.com/glvd/accipfs/core"
	"github.com/glvd/accipfs/general"
)

// DefaultClient ...
var DefaultClient core.API

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

// InitGlobalClient ...
func InitGlobalClient(cfg *config.Config) {
	DefaultClient = New(&cfg.API)
}

// New ...
func New(cfg *config.APIConfig) core.API {
	c := &http.Client{}
	c.Timeout = cfg.Timeout * time.Minute
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

func (c *client) do(uri string, req, resp interface{}) error {
	reader, err := requestReader(req)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, c.RequestURL(uri), reader)
	if err != nil {
		return err
	}
	response, err := c.cli.Do(request)
	if err != nil {
		return err
	}
	return responseDecoder(response.Body, resp)
}

// Ping ...
func Ping(req *core.PingReq) (resp *core.PingResp, err error) {
	return DefaultClient.Ping(req)
}

// Ping ...
func (c *client) Ping(req *core.PingReq) (resp *core.PingResp, err error) {
	resp = new(core.PingResp)
	err = c.do("ping", req, resp)
	return
}

// ID ...
func ID(req *core.IDReq) (resp *core.IDResp, err error) {
	return DefaultClient.ID(req)
}

// ID ...
func (c *client) ID(req *core.IDReq) (resp *core.IDResp, err error) {
	resp = new(core.IDResp)
	err = c.do("id", req, resp)
	return
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
func AddPeer(url string, node core.Node) error {
	status := new(bool)
	if err := general.RPCPost(url, "BustLinker.AddPeer", node, status); err != nil {
		logE("remote id error", "error", err.Error())
		return fmt.Errorf("remote id error: %w", err)
	}

	if !(*status) {
		return fmt.Errorf("connect failed: %s", general.RPCAddress(node.Addrs()[0]))
	}
	return nil
}
