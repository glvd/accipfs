package node

import (
	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/config"
	"github.com/glvd/accipfs/core"
	"testing"
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
		Version:    core.Version{},
	}
	t.Log(ds1)
	cache.Store("QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4", ds1)
	var d1 core.DataInfoV1
	err := cache.Load("QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4", &d1)
	if err != nil {
		return
	}
	cache.GC()
	t.Log(d1)

}
