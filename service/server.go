package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// NodeServer ...
type NodeServer interface {
	Start() error
	Init() error
	Stop() error
	Node() (Node, error)
}

// Server ...
type Server struct {
	cfg       config.Config
	rpcServer *rpc.Server
	route     *mux.Router
}

// NewRPCServer ...
func NewRPCServer(cfg config.Config) (*Server, error) {
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")

	acc, err := NewServerAccelerate(cfg)
	if err != nil {
		return nil, err
	}

	err = rpcServer.RegisterService(acc, "")
	if err != nil {
		return nil, err
	}
	return &Server{
		cfg:       cfg,
		rpcServer: rpcServer,
		route:     mux.NewRouter(),
	}, nil
}

// Start ...
func (s *Server) Start() error {
	s.route.Handle("/rpc", s.rpcServer)
	port := fmt.Sprintf(":%d", s.cfg.Port)
	log.Println("JSON RPC service listen and serving on port", port)
	if err := http.ListenAndServe(port, s.route); err != nil {
		log.Fatalf("Error serving: %s", err)
		return err
	}
	return nil
}

func screenOutput(ctx context.Context, reader io.Reader) (e error) {
	r := bufio.NewReader(reader)
	var lines []byte
END:
	for {
		select {
		case <-ctx.Done():
			return
		default:
			lines, _, e = r.ReadLine()
			if e != nil || io.EOF == e {
				break END
			}
			if strings.TrimSpace(string(lines)) != "" {
				fmt.Println(outputHead, string(lines))
			}
		}
	}

	return nil
}
