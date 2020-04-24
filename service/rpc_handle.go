package service

import (
	"github.com/gin-gonic/gin"
	"github.com/glvd/accipfs/config"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
)

// rpcHandle ...
type rpcHandle struct {
	cfg       *config.Config
	rpcServer *rpc.Server
	//linker    *BustLinker
}

func newRPCHandle(cfg *config.Config, handle interface{}) (*rpcHandle, error) {
	serv := rpc.NewServer()
	serv.RegisterCodec(json2.NewCodec(), "application/json")
	serv.RegisterCodec(json2.NewCodec(), "application/json;charset=UTF-8")

	err := serv.RegisterService(handle, "")
	if err != nil {
		return nil, err
	}
	return &rpcHandle{
		cfg:       cfg,
		rpcServer: serv,
		//linker:    handle,
	}, nil
}

// Handler ...
func (s *rpcHandle) ginHandler() (string, gin.HandlerFunc) {
	return "/rpc", s.gin
}

// RPCHandle ...
func (s *rpcHandle) gin(ctx *gin.Context) {
	s.rpcServer.ServeHTTP(ctx.Writer, ctx.Request)
}
