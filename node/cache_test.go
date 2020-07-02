package node

import (
	"fmt"
	"testing"
	"time"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
)

var tmpData = `{"data_hash":"b6d69073de76cf7da1fddc752f515938488634e613f2d4144c93f50cd03
97c83","data_info":{"root_hash":"QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4"
,"media_info":{"no":"a3e465c1-bc35-11ea-aa2c-00155d639067","intro":"","alias":nu
ll,"role":null,"director":"","systematics":"","season":"","total_episode":"","ep
isode":"","producer":"","publisher":"","type":"","format":"","language":"","capt
ion":"","group":"","index":"","date":"","sharpness":"","series":"","tags":null,"
length":"","sample":null,"uncensored":false},"media_uri":"","media_hash":"","med
ia_index":"","info":{"thumb_hash":"","thumb_uri":"","poster_hash":"","poster_uri
":""},"info_uri":"","security":{"key":""},"last_update":1593674915,"version":[0,
0,0,0]},"addr_info":{"ID":"QmeqN54NuCYSuTxHGZbvG3YoKnLewaECxAiGZUQsXyuXdY","Addr
s":null,"IPFSAddrInfo":{"Addrs":null,"ID":""}}}`

func TestHashCache_Store(t *testing.T) {
	root := "QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4"
	cache := hashCacher(config.Default())
	for i := 0; i < 100000; i++ {
		ds1 := &core.DataInfoV1{
			RootHash: root,
			MediaInfo: core.MediaInfo{
				No:           basis.UUID(),
				Intro:        "",
				Alias:        nil,
				Role:         nil,
				Director:     "",
				Systematics:  "",
				Season:       "",
				TotalEpisode: "",
				Episode:      "",
				Producer:     "",
				Publisher:    "",
				Type:         "",
				Format:       "",
				Language:     "",
				Caption:      "",
				Group:        "",
				Index:        "",
				Date:         "",
				Sharpness:    "",
				Series:       "",
				Tags:         nil,
				Length:       "",
				Sample:       nil,
				Uncensored:   false,
			},
			MediaURI:   "",
			MediaHash:  "",
			MediaIndex: "",
			Info:       core.Info{},
			InfoURI:    "",
			Security:   core.Security{},
			LastUpdate: time.Now().Unix(),
			Version:    core.Version{},
		}
		info1 := newDataHashInfo(ds1)
		info1.AddrInfo.ID = "QmeqN54NuCYSuTxHGZbvG3YoKnLewaECxAiGZUQsXyuXdY"
		t.Logf("%+v", info1)
		err := cache.Store(info1.Hash(), info1)
		if err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 100000; i++ {
		ds2 := &core.DataInfoV1{}
		d1 := newDataHashInfo(ds2)
		err1 := cache.Load(root, d1)
		if err1 != nil {
			t.Fatal(err1)
		}
	}

}

func BenchmarkHashCache_Load(b *testing.B) {
	//root := "QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4"
	cache := hashCacher(config.Default())
	for i := 0; i < 100; i++ {
		ds2 := &core.DataInfoV1{}
		d1 := newDataHashInfo(ds2)
		err1 := cache.Load("ffffa7def937e0c5d5565598a2f6a99eb79eaf753a8d6e407ea79ddd72ed526c", d1)
		if err1 != nil {
			b.Fatal(err1)
		}
		//t.Logf("%+v", d1)
	}
	//err2 := cache.GC()
	//if err2 != nil {
	//	b.Fatal(err2)
	//}
}

func generateTestData() (key string, value string) {
	root := "QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4"
	uuid := basis.UUID()
	ds1 := &core.DataInfoV1{
		RootHash: root,
		MediaInfo: core.MediaInfo{
			No:           uuid,
			Intro:        "",
			Alias:        nil,
			Role:         nil,
			Director:     "",
			Systematics:  "",
			Season:       "",
			TotalEpisode: "",
			Episode:      "",
			Producer:     "",
			Publisher:    "",
			Type:         "",
			Format:       "",
			Language:     "",
			Caption:      "",
			Group:        "",
			Index:        "",
			Date:         "",
			Sharpness:    "",
			Series:       "",
			Tags:         nil,
			Length:       "",
			Sample:       nil,
			Uncensored:   false,
		},
		MediaURI:   "",
		MediaHash:  "",
		MediaIndex: "",
		Info:       core.Info{},
		InfoURI:    "",
		Security:   core.Security{},
		LastUpdate: time.Now().Unix(),
		Version:    core.Version{},
	}
	info1 := newDataHashInfo(ds1)
	info1.AddrInfo.ID = "QmeqN54NuCYSuTxHGZbvG3YoKnLewaECxAiGZUQsXyuXdY"
	encode, err := info1.Marshal()
	if err != nil {
		return "", ""
	}
	return uuid, string(encode)
}

func BenchmarkDatabase2(b *testing.B) {
	opt := badger.DefaultOptions("badger")
	db, err := badger.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var key string
	for i := 0; i < 10000; i++ {
		data, value := generateTestData()
		err = db.Update(
			func(txn *badger.Txn) error {
				return txn.Set([]byte(data), []byte(value))
			})
		if err != nil {
			continue
		}
		key = data
	}
	fmt.Println("last key", key)
}
func BenchmarkDatabase2Read(b *testing.B) {
	opt := badger.DefaultOptions("badger")
	db, err := badger.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	key := "a3e465c1-bc35-11ea-aa2c-00155d639067"

	for i := 0; i < 100; i++ {
		err := db.View(
			func(txn *badger.Txn) error {
				item, err := txn.Get([]byte(key))
				if err != nil {
					return err
				}
				return item.Value(func(val []byte) error {
					fmt.Println("getted", string(val))
					return nil
				})
			})
		if err != nil {
			continue
		}
	}
}
func BenchmarkDatabaseMap(b *testing.B) {
	v := map[string]string{}
	key := "a3e465c1-bc35-11ea-aa2c-00155d639067"
	b.StopTimer()
	for i := 0; i < 100000; i++ {
		data, value := generateTestData()
		v[data] = value
		key = data
	}
	fmt.Println("last key", key)
	b.StartTimer()
	for i := 0; i < 100; i++ {
		fmt.Println(v[key])
	}
}

func BenchmarkHashMap(b *testing.B) {
	m := hashmap.New()
	key := "a3e465c1-bc35-11ea-aa2c-00155d639067"
	b.StopTimer()
	for i := 0; i < 100000; i++ {
		data, value := generateTestData()
		m.Put(data, value)
		key = data
	}
	fmt.Println("last key", key)
	b.StartTimer()
	for i := 0; i < 100; i++ {
		fmt.Println(m.Get(key))
	}
}
