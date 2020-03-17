package service

import (
	"bytes"
	"github.com/glvd/accipfs/config"
	"github.com/goextension/log/zap"
	"github.com/gorilla/rpc/v2/json"
	"log"
	"net/http"
	"testing"
)

func init() {
	zap.InitZapFileSugar()
}
func TestNodeServerETH(t *testing.T) {
	config.WorkDir = "D:\\workspace\\pvt"
	err := config.SaveConfig(config.Default())
	if err != nil {
		t.Error(err)
		return
	}
	config.Initialize()
	eth := NewNodeServerETH(config.Global())
	t.Logf("%+v", eth)
	if err := eth.Init(); err != nil {
		t.Error(err)
		return
	}
	//if err := eth.Start(); err != nil {
	//	t.Error(err)
	//	return
	//}

	ipfs := NewNodeServerIPFS(config.Global())
	t.Logf("%+v", ipfs)
	if err := ipfs.Init(); err != nil {
		t.Error(err)
		return
	}
	//if err := ipfs.Start(); err != nil {
	//	t.Error(err)
	//	return
	//}

	//time.Sleep(1 * time.Minute)

	//if err := eth.Stop(); err != nil {
	//	t.Error(err)
	//	return
	//}
	//
	//if err := ipfs.Stop(); err != nil {
	//	t.Error(err)
	//	return
	//}
	t.Log("done")
}

func TestNewServer(t *testing.T) {
	config.WorkDir = "d:\\workspace\\pvt"
	config.Initialize()
	server, e := NewServer(config.Global())
	if e != nil {
		t.Fatal(e)
	}
	go server.Start()
	url := "http://localhost:1234/rpc"

	message, err := json.EncodeClientRequest("Accelerate.Ping", &Empty{})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(message))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	reply := new(string)
	err = json.DecodeClientResponse(resp.Body, reply)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf(" %s\n", *reply)
}
