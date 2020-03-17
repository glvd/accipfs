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
	Name string
}

// Accelerate ...
type Accelerate struct {
	self *account.Account
}

// NodeInfo ...
type NodeInfo struct {
	Name         string
	ContractAddr string
	DataAddr     string
	Version      string
}

// NodeList ...
type NodeList map[string]NodeInfo

// NewServerAccelerate ...
func NewServerAccelerate(cfg *config.Config) (*Accelerate, error) {
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

// ExchangeNode ...
func (n *Accelerate) ExchangeNode(r *http.Request, list NodeList, result NodeList) error {
	return nil
}
