package node

import (
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"testing"
	"time"
)

func TestHashCache_Store(t *testing.T) {
	cache := newHashCacher(config.Default())
	ds1 := &core.DataInfoV1{
		RootHash: "QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4",
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
	ds2 := &core.DataInfoV1{}
	d1 := newDataHashInfo(ds2)
	err = cache.Load(ds1.RootHash, d1)
	if err != nil {
		t.Fatal(err)
	}
	err = cache.GC()
	t.Logf("%+v", d1)

}
