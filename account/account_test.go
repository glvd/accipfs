package account

import (
	"github.com/glvd/accipfs/config"
	"testing"
)

func TestNewAccount(t *testing.T) {
	config.WorkDir = "D:\\workspace\\pvt"
	if err := config.LoadConfig(); err != nil {
		t.Fatal(err)
		return
	}
	cfg := config.Global()

	acc, err := NewAccount(&cfg)
	if err != nil {
		t.Fatal(err)
		return
	}

	if err := acc.Save(&cfg); err != nil {
		t.Fatal(err)
		return
	}

}
