package service

import (
	"github.com/glvd/accipfs/account"
	"github.com/glvd/accipfs/config"
	"net/http"
)

// Empty ...
type Empty struct {
}

// Account ...
type Account struct {
	Name         string
	ContractAddr string
	DataAddr     string
}

// Accelerate ...
type Accelerate struct {
	self *account.Account
}

// NewAccelerate ...
func NewAccelerate(cfg config.Config) (*Accelerate, error) {
	account, err := account.LoadAccount(cfg)
	if err != nil {
		return nil, err
	}
	return &Accelerate{
		self: account,
	}, nil
}

// Ping ...
func (n *Accelerate) Ping(r *http.Request, s *Empty, result *string) error {
	*result = "pong pong pong"
	return nil
}

// ID ...
func (n *Accelerate) ID(r *http.Request, s *Empty, result *Account) error {
	result.Name = n.self.Name
	return nil
}
