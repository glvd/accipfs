package account

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/glvd/accipfs/config"
	"github.com/goextension/tool"
	"io/ioutil"
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
	var target []byte
	r := strings.NewReader(cfg.Account)
	dec := base64.NewDecoder(base64.StdEncoding, r)
	read, err := dec.Read(target)
	if err != nil {
		return nil, err
	}
	if read == 0 {
		return nil, errors.New("read account with size:0")
	}
	var acc Account
	err = json.Unmarshal(target, &acc)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}

// SaveAccountToConfig ...
func SaveAccountToConfig(cfg config.Config, account *Account) error {
	bytes, err := json.Marshal(account)
	if err != nil {
		return err
	}
	acc := base64.StdEncoding.EncodeToString(bytes)
	cfg.Account = acc
	return config.SaveConfig(&cfg)
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
