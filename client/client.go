package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// DefaultClient ...
var DefaultClient core.API

type client struct {
	cfg *config.Config
	//ds   *httpapi.HttpApi
	//node intercore.CoreAPI
	cli *http.Client
	ctx context.Context
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
	//c.Timeout = cfg.API.Timeout * time.Second
	//ma, err := multiaddr.NewMultiaddr(cfg.IPFSAPIAddr())
	//if err != nil {
	//	panic(err)
	//}
	//api, e := httpapi.NewApi(ma)
	//if e != nil {
	//	panic(e)
	//}
	//err := basis.SetupPlugins("")
	//if err != nil {
	//	panic(err)
	//}
	//node, err := basis.CreateNode(context.TODO(), filepath.Join(cfg.Path, ".ipfs"))
	//if err != nil {
	//	panic(err)
	//}
	return &client{
		//ds:   api,
		//node: node,
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

func (c *client) doGet(ctx context.Context, uri string, req url.Values, resp interface{}) error {
	request, err := http.NewRequest(http.MethodGet, requestQuery(c.RequestURL(uri), req), nil)
	if err != nil {
		return err
	}
	if ctx != nil {
		request.WithContext(ctx)
	}
	response, err := c.cli.Do(request)
	if err != nil {
		return err
	}
	return responseDecoder(response.Body, resp)
}
func (c *client) doPost(ctx context.Context, uri string, req, resp interface{}) error {
	reader, err := requestReader(req)
	if err != nil {
		return err
	}
	request, err := http.NewRequest(http.MethodPost, c.RequestURL(uri), reader)
	if err != nil {
		return err
	}
	if ctx != nil {
		request.WithContext(ctx)
	}
	response, err := c.cli.Do(request)
	if err != nil {
		return err
	}
	return responseDecoder(response.Body, resp)
}

// Ping ...
func Ping(ctx context.Context, req *core.PingReq) (resp *core.PingResp, err error) {
	return DefaultClient.Ping(ctx, req)
}

// Ping ...
func (c *client) Ping(ctx context.Context, req *core.PingReq) (resp *core.PingResp, err error) {
	resp = new(core.PingResp)
	err = c.doGet(ctx, "ping", nil, resp)
	return
}

// ID ...
func ID(ctx context.Context, req *core.IDReq) (resp *core.IDResp, err error) {
	return DefaultClient.ID(ctx, req)
}

// ID ...
func (c *client) ID(ctx context.Context, req *core.IDReq) (resp *core.IDResp, err error) {
	resp = new(core.IDResp)
	err = c.doPost(ctx, "id", req, resp)
	return
}

// Add ...
func (c *client) Add(ctx context.Context, req *core.NodeAddReq) (resp *core.NodeAddResp, err error) {
	resp = new(core.NodeAddResp)
	err = c.doPost(ctx, "add", req, resp)
	return
}

// NodeAddrInfo ...
func NodeAddrInfo(ctx context.Context, req *core.AddrReq) (*core.AddrResp, error) {
	return DefaultClient.NodeAPI().NodeAddrInfo(ctx, req)
}

// NodeAddrInfo ...
func (c *client) NodeAddrInfo(ctx context.Context, req *core.AddrReq) (resp *core.AddrResp, err error) {
	resp = new(core.AddrResp)
	err = c.doPost(ctx, "info", req, resp)
	return
}

// Add ...
func Add(ctx context.Context, req *core.NodeAddReq) (resp *core.NodeAddResp, err error) {
	return DefaultClient.Add(ctx, req)
}

// UploadFile ...
func UploadFile(ctx context.Context, req *core.UploadReq) (resp *core.UploadResp, err error) {
	return DefaultClient.DataStoreAPI().UploadFile(ctx, req)
}
