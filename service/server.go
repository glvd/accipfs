package service

import (
	"context"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"net/http"
)

// NodeServer ...
type NodeServer interface {
	Start() error
	Init() error
	Stop() error
}

// rpcServer ...
type rpcServer struct {
	cfg        *config.Config
	rpcServer  *rpc.Server
	httpServer *http.Server
	route      *mux.Router
}

// NewRPCServer ...
func newRPCServer(cfg *config.Config, linker *BustLinker) (*rpcServer, error) {
	serv := rpc.NewServer()
	serv.RegisterCodec(json2.NewCodec(), "application/json;charset=UTF-8")

	err := serv.RegisterService(linker, "")
	if err != nil {
		return nil, err
	}
	return &rpcServer{
		cfg:       cfg,
		rpcServer: serv,
		route:     mux.NewRouter(),
	}, nil
}

// Start ...
func (s *rpcServer) Start() error {
	s.route.Handle("/rpc", s.rpcServer)

	port := fmt.Sprintf(":%d", s.cfg.Port)
	s.httpServer = &http.Server{Addr: port, Handler: s.route}

	output("JSON RPC service listen and serving on port", port)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// Stop ...
func (s *rpcServer) Stop() error {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		return err
	}
	return nil
}
