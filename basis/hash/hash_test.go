package hash_test

import (
	"bytes"
	"github.com/glvd/accipfs/basis/hash"
	"testing"

	"github.com/glvd/accipfs/basis"
	"github.com/glvd/accipfs/core"
)

func TestSumStruct(t *testing.T) {
	id := basis.UUID()
	sum, err := hash.Sum(core.DataInfoV1{
		RootHash: "QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4",
		MediaInfo: core.MediaInfo{
			No:           id,
			Intro:        "",
			Alias:        []string{},
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
	})
	if err != nil {
		return
	}
	t.Logf("%x\n", sum)
	sum1, err := hash.Sum(core.DataInfoV1{
		RootHash: "QmVq9Du6jAgHBJfyXuhHkP9KHARxJ1RYoYPXTKdVkoN6F4",
		MediaInfo: core.MediaInfo{
			No:           id,
			Intro:        "",
			Alias:        []string{},
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
	})
	if err != nil {
		return
	}
	t.Logf("%x\n", sum1)
	if bytes.Compare(sum, sum1) != 0 {
		t.Fatalf("sum:%x,sum1:%x", sum, sum1)
	}
}
func TestSum(t *testing.T) {
	sum, err := hash.Sum(map[string]string{
		"vala": "a",
		"valb": "b",
		"valc": "c",
	})
	if err != nil {
		return
	}
	t.Logf("%x\n", sum)
	sum1, err := hash.Sum(map[string]string{
		"valb": "b",
		"vala": "a",
		"valc": "c",
	})
	if err != nil {
		return
	}
	t.Logf("%x\n", sum1)
	if bytes.Compare(sum, sum1) != 0 {
		t.Fatalf("sum:%x,sum1:%x", sum, sum1)
	}
}

func TestSumString(t *testing.T) {
	sum, err := hash.Sum("hello world")
	if err != nil {
		return
	}
	t.Logf("%x\n", sum)
	sum1, err := hash.Sum("hello world")
	if err != nil {
		return
	}
	t.Logf("%x\n", sum1)
	if bytes.Compare(sum, sum1) != 0 {
		t.Fatalf("sum:%x,sum1:%x", sum, sum1)
	}
}
