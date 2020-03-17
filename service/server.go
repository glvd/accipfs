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
}

// NewServer ...
func NewServer(cfg config.Config) (*Server, error) {
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json.NewCodec(), "application/json")
	rpcServer.RegisterCodec(json.NewCodec(), "application/json;charset=UTF-8")
	acc := &nodeServerAccelerate{}
	err := rpcServer.RegisterService(acc, "accelerate")
	if err != nil {
		return nil, err
	}
	r := mux.NewRouter()
	r.Handle("/rpc", rpcServer)
	log.Println("JSON RPC service listen and serving on port 1234")
	if err := http.ListenAndServe(":1234", r); err != nil {
		log.Fatalf("Error serving: %s", err)
	}

	return &Server{
		cfg:       cfg,
		rpcServer: rpcServer,
	}, nil
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
