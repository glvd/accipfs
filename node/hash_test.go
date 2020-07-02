package node

import (
	"fmt"
	"testing"
	"time"

	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"github.com/xujiajun/nutsdb"
)

func TestHashCache_Store(t *testing.T) {
	root := "QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4"
	cache := newHashCacher(config.Default())
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
		//t.Logf("%+v", d1)
	}
	err2 := cache.GC()
	if err2 != nil {
		t.Fatal(err2)
	}

}

func BenchmarkHashCache_Load(b *testing.B) {
	//root := "QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4"
	cache := newHashCacher(config.Default())
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

func BenchmarkDatabase(b *testing.B) {
	bucket := "test"
	// Open the database located in the /tmp/nutsdb directory.
	// It will be created if it doesn't exist.
	opt := nutsdb.DefaultOptions
	opt.Dir = "/tmp/nutsdb"
	db, err := nutsdb.Open(opt)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	var uuid string
	root := "QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4"
	for i := 0; i < 100000; i++ {
		if uuid == "" {
			uuid = basis.UUID()
			fmt.Println("first uuid:", uuid)
		} else {
			uuid = basis.UUID()
		}

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
		encode, err := info1.Encode()
		if err != nil {
			continue
		}
		err = db.Update(
			func(tx *nutsdb.Tx) error {
				return tx.Put(bucket, []byte(uuid), []byte(encode), 0)
			})
		if err != nil {
			continue
		}
	}
	fmt.Println("last uuid", uuid)
	for i := 0; i < 100; i++ {
		err := db.View(
			func(tx *nutsdb.Tx) error {
				entry, err := tx.Get(bucket, []byte(uuid))
				if err != nil {
					return err
				}
				fmt.Println("entry", string(entry.Encode()))
				return nil
			})
		if err != nil {
			continue
		}
	}
}
