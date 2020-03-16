package account

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/glvd/accipfs/config"
	"github.com/goextension/tool"
	"io/ioutil"
	"path/filepath"
)

// Account ...
type Account struct {
	Name     string
	Address  string
	KeyStore KeyStore
	Password string
}

// NewAccount ...
func NewAccount(cfg config.Config) (*Account, error) {
	var acc Account
	acc.Password = tool.GenerateRandomString(8)

	ks := keystore.NewKeyStore(config.KeyStoreDirETH(), keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount(acc.Password)
	if err != nil {
		return nil, err
	}

	acc.getName(&account)

	e := acc.loadKey(&account)
	if e != nil {
		return nil, e
	}
	return &acc, nil
}

func (acc *Account) getName(act *accounts.Account) {
	_, acc.Name = filepath.Split(act.URL.Path)
}

func (acc *Account) loadKey(act *accounts.Account) error {
	fileBytes, e := ioutil.ReadFile(act.URL.Path)
	if e != nil {
		return e
	}
	return json.Unmarshal(fileBytes, &acc.KeyStore)
}
