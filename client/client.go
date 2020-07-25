package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/glvd/accipfs/config"
	httpapi "github.com/ipfs/go-ipfs-http-client"
	"github.com/multiformats/go-multiaddr"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/glvd/accipfs/core"
)

// DefaultClient ...
var DefaultClient core.API

type client struct {
	cfg *config.Config
	ds  *httpapi.HttpApi
	cli *http.Client
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
	DefaultClient = New(cfg)
}

// New ...
func New(cfg *config.Config) core.API {
	c := &http.Client{}
	c.Timeout = cfg.API.Timeout * time.Second
	ma, err := multiaddr.NewMultiaddr(cfg.IPFSAPIAddr())
	if err != nil {
		panic(err)
	}
	api, e := httpapi.NewApi(ma)
	if e != nil {
		panic(e)
	}
	return &client{
		ds:  api,
		cli: c,
		cfg: cfg,
	}
}

func (c *client) host() string {
	prefix := "http://"
	if c.cfg.UseTLS {
		prefix = "https://"
	}
	return strings.Join([]string{prefix, "127.0.0.1:", strconv.Itoa(c.cfg.API.Port), "/api/", c.cfg.API.Version}, "")
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

// Add ...
func Add(req *core.AddReq) (resp *core.AddResp, err error) {
	return DefaultClient.Add(req)
}

// UploadFile ...
func UploadFile(req *core.UploadReq) (resp *core.UploadResp, err error) {
	return DefaultClient.DataStoreAPI().UploadFile(req)
}
