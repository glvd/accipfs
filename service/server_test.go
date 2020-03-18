package service

import (
	"bytes"
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/goextension/log/zap"
	"github.com/gorilla/rpc/v2/json2"
	"io/ioutil"
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
	cfg := config.Global()
	eth := NewNodeServerETH(&cfg)
	t.Logf("%+v", eth)
	if err := eth.Init(); err != nil {
		t.Error(err)
		return
	}
	//if err := eth.Start(); err != nil {
	//	t.Error(err)
	//	return
	//}

	ipfs := NewNodeServerIPFS(&cfg)
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
	cfg := config.Global()
	acc, e := account.NewAccount(&cfg)
	if e != nil {
		t.Fatal(e)
	}
	if err := acc.Save(&cfg); err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", cfg)
	//server, e := NewRPCServer(&cfg)
	//if e != nil {
	//	t.Fatal(e)
	//}
	//go server.Start()
	url := "http://localhost:1234/rpc"

	m1, e := json2.EncodeClientRequest("Accelerate.Ping", &Empty{})
	if e != nil {
		return
	}
	r1, err := http.Post(url, "application/json", bytes.NewReader(m1))
	if err != nil {
		t.Fatal(err)
	}
	readAll, e := ioutil.ReadAll(r1.Body)
	if e != nil {
		return
	}
	t.Log(string(readAll))
	message, err := json2.EncodeClientRequest("Accelerate.ID", &Empty{})
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(message))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	all, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		t.Fatal(e)
	}
	t.Log(string(all))
	reply := new(core.NodeInfo)
	err = json2.DecodeClientResponse(bytes.NewReader(all), reply)
	if err != nil {
		t.Fatal(err)
	}
	log.Printf(" %+v\n", *reply)
}
