package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/glvd/accipfs/config"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/core"
)

// DefaultClient ...
var DefaultClient core.API

type client struct {
	cfg *config.APIConfig
	cli *http.Client
}

// DataStoreAPI ...
func (c *client) DataStoreAPI() core.DataStoreAPI {
	return c
}

type jsonResp struct {
	Status  string `json:"status"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

func requestQuery(url string, req url.Values) string {
	if req == nil {
		return url
	}
	return url + "?" + req.Encode()
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
	r := &jsonResp{}
	err := decoder.Decode(r)
	if err != nil {
		return err
	}
	if r.Error != "" {
		return errors.New(r.Error)
	}
	return json.Unmarshal([]byte(r.Message), resp)
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
func (c *client) doGet(uri string, req url.Values, resp interface{}) error {
	request, err := http.NewRequest(http.MethodGet, requestQuery(c.RequestURL(uri), req), nil)
	if err != nil {
		return err
	}
	response, err := c.cli.Do(request)
	if err != nil {
		return err
	}
	return responseDecoder(response.Body, resp)
}
func (c *client) doPost(uri string, req, resp interface{}) error {
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
	err = c.doGet("ping", nil, resp)
	return
}

// ID ...
func ID(req *core.IDReq) (resp *core.IDResp, err error) {
	return DefaultClient.ID(req)
}

// ID ...
func (c *client) ID(req *core.IDReq) (resp *core.IDResp, err error) {
	resp = new(core.IDResp)
	err = c.doPost("id", req, resp)
	return
}

// Link ...
func (c *client) Link(req *core.NodeLinkReq) (resp *core.NodeLinkResp, err error) {
	resp = new(core.NodeLinkResp)
	err = c.doPost("/node/link", req, resp)
	return
}

// Add ...
func (c *client) Add(req *core.AddReq) (resp *core.AddResp, err error) {
	resp = new(core.AddResp)
	err = c.doPost("add", req, resp)
	return
}

// NodeAddrInfo ...
func NodeAddrInfo(req *core.AddrReq) (*core.AddrResp, error) {
	return DefaultClient.NodeAPI().NodeAddrInfo(req)
}

// NodeAddrInfo ...
func (c *client) NodeAddrInfo(req *core.AddrReq) (resp *core.AddrResp, err error) {
	resp = new(core.AddrResp)
	err = c.doPost("info", req, resp)
	return
}

// List ...
func List(req *core.NodeListReq) (resp *core.NodeListResp, err error) {
	return DefaultClient.NodeAPI().List(req)
}

// NodeLink ...
func NodeLink(req *core.NodeLinkReq) (resp *core.NodeLinkResp, err error) {
	return DefaultClient.NodeAPI().Link(req)
}

// Pins ...
func Pins(req *core.DataStoreReq) (*core.DataStoreResp, error) {
	return DefaultClient.DataStoreAPI().Pins(req)
}

// PinVideo ...
func PinVideo(url string, no string) error {
	logD("pin hash", "hash", no)
	b := new(bool)
	err := basis.RPCPost(url, "BustLinker.PinVideo", &no, b)
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
	if err := basis.RPCPost(url, "BustLinker.Peers", info, result); err != nil {
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
	if err := basis.RPCPost(url, "BustLinker.AddPeer", node, status); err != nil {
		logE("remote id error", "error", err.Error())
		return fmt.Errorf("remote id error: %w", err)
	}

	if !(*status) {
		return fmt.Errorf("connect failed: %s", basis.RPCAddress(node.Addrs()[0]))
	}
	return nil
}

// NodeAPI ...
func (c *client) NodeAPI() core.NodeAPI {
	return c
}

// Unlink ...
func (c *client) Unlink(req *core.NodeUnlinkReq) (resp *core.NodeUnlinkResp, err error) {
	resp = new(core.NodeUnlinkResp)
	err = c.doPost("node/unlink", req, resp)
	return
}

// List ...
func (c *client) List(req *core.NodeListReq) (resp *core.NodeListResp, err error) {
	resp = new(core.NodeListResp)
	err = c.doPost("node/list", req, resp)
	return
}

// DataStorePins ...
func DataStorePins(req *core.DataStoreReq) (resp *core.DataStoreResp, err error) {
	return DefaultClient.DataStoreAPI().Pins(req)
}

// Pins ...
func (c *client) Pins(req *core.DataStoreReq) (resp *core.DataStoreResp, err error) {
	resp = new(core.DataStoreResp)
	err = c.doPost("datastore/pins", req, resp)
	return
}
