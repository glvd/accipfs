package account

import (
	"encoding/json"
	"github.com/glvd/accipfs/config"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestNewAccount(t *testing.T) {
	config.WorkDir = "D:\\workspace\\pvt"
	if err := config.LoadConfig(); err != nil {
		t.Fatal(err)
		return
	}
	acc, err := NewAccount(config.Global())
	if err != nil {
		t.Fatal(err)
		return
	}
	marshal, e := json.Marshal(acc)
	if e != nil {
		return
	}
	ioutil.WriteFile(filepath.Join(config.WorkDir, "tmp"), marshal, 0755)
}
