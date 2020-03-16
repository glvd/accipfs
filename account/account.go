package account

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/glvd/accipfs/config"
	"path/filepath"
)

// NewAccount ...
func NewAccount(cfg config.Config) {

	ks := keystore.NewKeyStore(filepath.Join(cfg.Path, "key"), keystore.StandardScryptN, keystore.StandardScryptP)
	password := "secret"
	account, err := ks.NewAccount(password)
	if err != nil {
		log.Fatal(err)
	}

}
