package service

import (
	"bufio"
	"context"
	"fmt"
	"github.com/glvd/accipfs/config"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
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
	cfg        *config.Config
	linker     *BustLinker
	rpcServer  *rpc.Server
	httpServer *http.Server
	route      *mux.Router
}

// NewRPCServer ...
func NewRPCServer(cfg *config.Config) (*Server, error) {
	rpcServer := rpc.NewServer()
	rpcServer.RegisterCodec(json2.NewCodec(), "application/json")
	rpcServer.RegisterCodec(json2.NewCodec(), "application/json;charset=UTF-8")

	acc, err := NewBustLinker(cfg)
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
		linker:    acc,
		route:     mux.NewRouter(),
	}, nil
}

// Start ...
func (s *Server) Start() error {
	s.route.Handle("/rpc", s.rpcServer)

	port := fmt.Sprintf(":%d", s.cfg.Port)
	s.httpServer = &http.Server{Addr: port, Handler: s.route}

	go s.linker.Start()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if s.linker.ipfsClient.IsReady() {
			wg.Done()
			return
		}
	}()
	wg.Add(1)
	go func() {
		if s.linker.ethClient.IsReady() {
			wg.Done()
			return
		}
	}()
	wg.Wait()

	var idError error
	for i := 0; i < 5; i++ {
		id, err := s.linker.localID()
		idError = err
		if err != nil {
			time.Sleep(3 * time.Second)
			continue
		}
		s.linker.id = id
		break
	}

	if idError != nil {
		return idError
	}
	fmt.Println(outputHead, "JSON RPC service listen and serving on port", port)
	s.httpServer.ListenAndServe()
	return nil
}

// Stop ...
func (s *Server) Stop() error {
	if err := s.httpServer.Shutdown(context.Background()); err != nil {
		return err
	}
	s.linker.Stop()
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
