package account

import (
	"encoding/base64"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/glvd/accipfs/config"
	"github.com/goextension/tool"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

// LoadAccount ...
func LoadAccount(cfg config.Config) (*Account, error) {
	r := strings.NewReader(cfg.Account)
	dec := base64.NewDecoder(base64.StdEncoding, r)
	target, err := ioutil.ReadAll(dec)
	if err != nil {
		return nil, err
	}
	var acc Account
	err = json.Unmarshal(target, &acc)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

// SaveAccountToConfig ...
func saveAccountToConfig(account *Account, cfg config.Config) error {
	bytes, err := json.Marshal(account)
	if err != nil {
		return err
	}
	acc := base64.StdEncoding.EncodeToString(bytes)
	cfg.Account = acc
	return config.SaveConfig(&cfg)
}

// Check ...
func (acc *Account) Check() error {
	path := filepath.Join(config.KeyStoreDirETH(), acc.Name)
	_, e := os.Stat(path)
	if e != nil && os.IsNotExist(e) {
		bytes, e := json.Marshal(acc.KeyStore)
		if e != nil {
			return e
		}
		e = ioutil.WriteFile(path, bytes, 0755)
		if e != nil {
			return e
		}
		return nil
	}
	return nil
}

// Save ...
func (acc *Account) Save(cfg config.Config) error {
	if err := acc.Check(); err != nil {
		return err
	}
	return saveAccountToConfig(acc, cfg)
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
