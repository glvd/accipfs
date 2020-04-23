package service

import (
	"net/http"

	"github.com/glvd/accipfs/config"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

// rpcHandle ...
type rpcHandle struct {
	cfg       *config.Config
	rpcServer *rpc.Server
	linker    *BustLinker
}

func newRPCHandle(cfg *config.Config, linker *BustLinker) (*rpcHandle, error) {
	serv := rpc.NewServer()
	serv.RegisterCodec(json2.NewCodec(), "application/json")
	serv.RegisterCodec(json2.NewCodec(), "application/json;charset=UTF-8")

	err := serv.RegisterService(linker, "")
	if err != nil {
		return nil, err
	}
	return &rpcHandle{
		cfg:       cfg,
		rpcServer: serv,
		linker:    linker,
	}, nil
}

// Handler ...
func (s *rpcHandle) Handler() (string, http.Handler) {
	return "/rpc", s.rpcServer
}
