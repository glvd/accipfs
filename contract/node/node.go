// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package node

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AccelerateNodeABI is the input ABI used to generate the binding from.
const AccelerateNodeABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"nodes\",\"type\":\"string[]\"}],\"name\":\"addEthNodes\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nodes\",\"type\":\"string[]\"}],\"name\":\"addIpfsNodes\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nodes\",\"type\":\"string[]\"}],\"name\":\"addPublicIpfsNodes\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nodes\",\"type\":\"string[]\"}],\"name\":\"addSignerNodes\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"idx\",\"type\":\"uint32\"}],\"name\":\"deleteEthNodes\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"idx\",\"type\":\"uint32\"}],\"name\":\"deleteIpfsNodes\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"idx\",\"type\":\"uint32\"}],\"name\":\"deletePublicIpfsNodes\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"idx\",\"type\":\"uint32\"}],\"name\":\"deleteSignerNodes\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"enode\",\"type\":\"string\"},{\"name\":\"idx\",\"type\":\"uint32\"}],\"name\":\"replaceNode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"resetEthNode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"resetIpfsNode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"resetPublicIpfsNode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"resetSignerNode\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"getEthNodes\",\"outputs\":[{\"name\":\"\",\"type\":\"string[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getEthNodesString\",\"outputs\":[{\"name\":\"\",\"type\":\"string[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getIpfsNodes\",\"outputs\":[{\"name\":\"\",\"type\":\"string[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPublicIpfsNodes\",\"outputs\":[{\"name\":\"\",\"type\":\"string[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getSignerNodes\",\"outputs\":[{\"name\":\"\",\"type\":\"string[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// AccelerateNode is an auto generated Go binding around an Ethereum contract.
type AccelerateNode struct {
	AccelerateNodeCaller     // Read-only binding to the contract
	AccelerateNodeTransactor // Write-only binding to the contract
	AccelerateNodeFilterer   // Log filterer for contract events
}

// AccelerateNodeCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccelerateNodeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccelerateNodeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccelerateNodeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccelerateNodeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccelerateNodeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccelerateNodeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccelerateNodeSession struct {
	Contract     *AccelerateNode   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccelerateNodeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccelerateNodeCallerSession struct {
	Contract *AccelerateNodeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// AccelerateNodeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccelerateNodeTransactorSession struct {
	Contract     *AccelerateNodeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// AccelerateNodeRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccelerateNodeRaw struct {
	Contract *AccelerateNode // Generic contract binding to access the raw methods on
}

// AccelerateNodeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccelerateNodeCallerRaw struct {
	Contract *AccelerateNodeCaller // Generic read-only contract binding to access the raw methods on
}

// AccelerateNodeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccelerateNodeTransactorRaw struct {
	Contract *AccelerateNodeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccelerateNode creates a new instance of AccelerateNode, bound to a specific deployed contract.
func NewAccelerateNode(address common.Address, backend bind.ContractBackend) (*AccelerateNode, error) {
	contract, err := bindAccelerateNode(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccelerateNode{AccelerateNodeCaller: AccelerateNodeCaller{contract: contract}, AccelerateNodeTransactor: AccelerateNodeTransactor{contract: contract}, AccelerateNodeFilterer: AccelerateNodeFilterer{contract: contract}}, nil
}

// NewAccelerateNodeCaller creates a new read-only instance of AccelerateNode, bound to a specific deployed contract.
func NewAccelerateNodeCaller(address common.Address, caller bind.ContractCaller) (*AccelerateNodeCaller, error) {
	contract, err := bindAccelerateNode(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccelerateNodeCaller{contract: contract}, nil
}

// NewAccelerateNodeTransactor creates a new write-only instance of AccelerateNode, bound to a specific deployed contract.
func NewAccelerateNodeTransactor(address common.Address, transactor bind.ContractTransactor) (*AccelerateNodeTransactor, error) {
	contract, err := bindAccelerateNode(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccelerateNodeTransactor{contract: contract}, nil
}

// NewAccelerateNodeFilterer creates a new log filterer instance of AccelerateNode, bound to a specific deployed contract.
func NewAccelerateNodeFilterer(address common.Address, filterer bind.ContractFilterer) (*AccelerateNodeFilterer, error) {
	contract, err := bindAccelerateNode(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccelerateNodeFilterer{contract: contract}, nil
}

// bindAccelerateNode binds a generic wrapper to an already deployed contract.
func bindAccelerateNode(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccelerateNodeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccelerateNode *AccelerateNodeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccelerateNode.Contract.AccelerateNodeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccelerateNode *AccelerateNodeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccelerateNode.Contract.AccelerateNodeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccelerateNode *AccelerateNodeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccelerateNode.Contract.AccelerateNodeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccelerateNode *AccelerateNodeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AccelerateNode.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccelerateNode *AccelerateNodeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccelerateNode.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccelerateNode *AccelerateNodeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccelerateNode.Contract.contract.Transact(opts, method, params...)
}

// GetEthNodes is a free data retrieval call binding the contract method 0x4dc52f71.
//
// Solidity: function getEthNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeCaller) GetEthNodes(opts *bind.CallOpts) ([]string, error) {
	var (
		ret0 = new([]string)
	)
	out := ret0
	err := _AccelerateNode.contract.Call(opts, out, "getEthNodes")
	return *ret0, err
}

// GetEthNodes is a free data retrieval call binding the contract method 0x4dc52f71.
//
// Solidity: function getEthNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeSession) GetEthNodes() ([]string, error) {
	return _AccelerateNode.Contract.GetEthNodes(&_AccelerateNode.CallOpts)
}

// GetEthNodes is a free data retrieval call binding the contract method 0x4dc52f71.
//
// Solidity: function getEthNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeCallerSession) GetEthNodes() ([]string, error) {
	return _AccelerateNode.Contract.GetEthNodes(&_AccelerateNode.CallOpts)
}

// GetEthNodesString is a free data retrieval call binding the contract method 0x51ae9d17.
//
// Solidity: function getEthNodesString() constant returns(string[])
func (_AccelerateNode *AccelerateNodeCaller) GetEthNodesString(opts *bind.CallOpts) ([]string, error) {
	var (
		ret0 = new([]string)
	)
	out := ret0
	err := _AccelerateNode.contract.Call(opts, out, "getEthNodesString")
	return *ret0, err
}

// GetEthNodesString is a free data retrieval call binding the contract method 0x51ae9d17.
//
// Solidity: function getEthNodesString() constant returns(string[])
func (_AccelerateNode *AccelerateNodeSession) GetEthNodesString() ([]string, error) {
	return _AccelerateNode.Contract.GetEthNodesString(&_AccelerateNode.CallOpts)
}

// GetEthNodesString is a free data retrieval call binding the contract method 0x51ae9d17.
//
// Solidity: function getEthNodesString() constant returns(string[])
func (_AccelerateNode *AccelerateNodeCallerSession) GetEthNodesString() ([]string, error) {
	return _AccelerateNode.Contract.GetEthNodesString(&_AccelerateNode.CallOpts)
}

// GetIpfsNodes is a free data retrieval call binding the contract method 0x5bf09a37.
//
// Solidity: function getIpfsNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeCaller) GetIpfsNodes(opts *bind.CallOpts) ([]string, error) {
	var (
		ret0 = new([]string)
	)
	out := ret0
	err := _AccelerateNode.contract.Call(opts, out, "getIpfsNodes")
	return *ret0, err
}

// GetIpfsNodes is a free data retrieval call binding the contract method 0x5bf09a37.
//
// Solidity: function getIpfsNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeSession) GetIpfsNodes() ([]string, error) {
	return _AccelerateNode.Contract.GetIpfsNodes(&_AccelerateNode.CallOpts)
}

// GetIpfsNodes is a free data retrieval call binding the contract method 0x5bf09a37.
//
// Solidity: function getIpfsNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeCallerSession) GetIpfsNodes() ([]string, error) {
	return _AccelerateNode.Contract.GetIpfsNodes(&_AccelerateNode.CallOpts)
}

// GetPublicIpfsNodes is a free data retrieval call binding the contract method 0x2a303e65.
//
// Solidity: function getPublicIpfsNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeCaller) GetPublicIpfsNodes(opts *bind.CallOpts) ([]string, error) {
	var (
		ret0 = new([]string)
	)
	out := ret0
	err := _AccelerateNode.contract.Call(opts, out, "getPublicIpfsNodes")
	return *ret0, err
}

// GetPublicIpfsNodes is a free data retrieval call binding the contract method 0x2a303e65.
//
// Solidity: function getPublicIpfsNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeSession) GetPublicIpfsNodes() ([]string, error) {
	return _AccelerateNode.Contract.GetPublicIpfsNodes(&_AccelerateNode.CallOpts)
}

// GetPublicIpfsNodes is a free data retrieval call binding the contract method 0x2a303e65.
//
// Solidity: function getPublicIpfsNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeCallerSession) GetPublicIpfsNodes() ([]string, error) {
	return _AccelerateNode.Contract.GetPublicIpfsNodes(&_AccelerateNode.CallOpts)
}

// GetSignerNodes is a free data retrieval call binding the contract method 0xf468f76f.
//
// Solidity: function getSignerNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeCaller) GetSignerNodes(opts *bind.CallOpts) ([]string, error) {
	var (
		ret0 = new([]string)
	)
	out := ret0
	err := _AccelerateNode.contract.Call(opts, out, "getSignerNodes")
	return *ret0, err
}

// GetSignerNodes is a free data retrieval call binding the contract method 0xf468f76f.
//
// Solidity: function getSignerNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeSession) GetSignerNodes() ([]string, error) {
	return _AccelerateNode.Contract.GetSignerNodes(&_AccelerateNode.CallOpts)
}

// GetSignerNodes is a free data retrieval call binding the contract method 0xf468f76f.
//
// Solidity: function getSignerNodes() constant returns(string[])
func (_AccelerateNode *AccelerateNodeCallerSession) GetSignerNodes() ([]string, error) {
	return _AccelerateNode.Contract.GetSignerNodes(&_AccelerateNode.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_AccelerateNode *AccelerateNodeCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _AccelerateNode.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_AccelerateNode *AccelerateNodeSession) IsOwner() (bool, error) {
	return _AccelerateNode.Contract.IsOwner(&_AccelerateNode.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_AccelerateNode *AccelerateNodeCallerSession) IsOwner() (bool, error) {
	return _AccelerateNode.Contract.IsOwner(&_AccelerateNode.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AccelerateNode *AccelerateNodeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AccelerateNode.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AccelerateNode *AccelerateNodeSession) Owner() (common.Address, error) {
	return _AccelerateNode.Contract.Owner(&_AccelerateNode.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AccelerateNode *AccelerateNodeCallerSession) Owner() (common.Address, error) {
	return _AccelerateNode.Contract.Owner(&_AccelerateNode.CallOpts)
}

// AddEthNodes is a paid mutator transaction binding the contract method 0x2d436cec.
//
// Solidity: function addEthNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeTransactor) AddEthNodes(opts *bind.TransactOpts, nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "addEthNodes", nodes)
}

// AddEthNodes is a paid mutator transaction binding the contract method 0x2d436cec.
//
// Solidity: function addEthNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeSession) AddEthNodes(nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.Contract.AddEthNodes(&_AccelerateNode.TransactOpts, nodes)
}

// AddEthNodes is a paid mutator transaction binding the contract method 0x2d436cec.
//
// Solidity: function addEthNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) AddEthNodes(nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.Contract.AddEthNodes(&_AccelerateNode.TransactOpts, nodes)
}

// AddIpfsNodes is a paid mutator transaction binding the contract method 0xfea8d842.
//
// Solidity: function addIpfsNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeTransactor) AddIpfsNodes(opts *bind.TransactOpts, nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "addIpfsNodes", nodes)
}

// AddIpfsNodes is a paid mutator transaction binding the contract method 0xfea8d842.
//
// Solidity: function addIpfsNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeSession) AddIpfsNodes(nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.Contract.AddIpfsNodes(&_AccelerateNode.TransactOpts, nodes)
}

// AddIpfsNodes is a paid mutator transaction binding the contract method 0xfea8d842.
//
// Solidity: function addIpfsNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) AddIpfsNodes(nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.Contract.AddIpfsNodes(&_AccelerateNode.TransactOpts, nodes)
}

// AddPublicIpfsNodes is a paid mutator transaction binding the contract method 0x786c666e.
//
// Solidity: function addPublicIpfsNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeTransactor) AddPublicIpfsNodes(opts *bind.TransactOpts, nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "addPublicIpfsNodes", nodes)
}

// AddPublicIpfsNodes is a paid mutator transaction binding the contract method 0x786c666e.
//
// Solidity: function addPublicIpfsNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeSession) AddPublicIpfsNodes(nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.Contract.AddPublicIpfsNodes(&_AccelerateNode.TransactOpts, nodes)
}

// AddPublicIpfsNodes is a paid mutator transaction binding the contract method 0x786c666e.
//
// Solidity: function addPublicIpfsNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) AddPublicIpfsNodes(nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.Contract.AddPublicIpfsNodes(&_AccelerateNode.TransactOpts, nodes)
}

// AddSignerNodes is a paid mutator transaction binding the contract method 0xc8af5a47.
//
// Solidity: function addSignerNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeTransactor) AddSignerNodes(opts *bind.TransactOpts, nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "addSignerNodes", nodes)
}

// AddSignerNodes is a paid mutator transaction binding the contract method 0xc8af5a47.
//
// Solidity: function addSignerNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeSession) AddSignerNodes(nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.Contract.AddSignerNodes(&_AccelerateNode.TransactOpts, nodes)
}

// AddSignerNodes is a paid mutator transaction binding the contract method 0xc8af5a47.
//
// Solidity: function addSignerNodes(string[] nodes) returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) AddSignerNodes(nodes []string) (*types.Transaction, error) {
	return _AccelerateNode.Contract.AddSignerNodes(&_AccelerateNode.TransactOpts, nodes)
}

// DeleteEthNodes is a paid mutator transaction binding the contract method 0x2cc267ea.
//
// Solidity: function deleteEthNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeTransactor) DeleteEthNodes(opts *bind.TransactOpts, idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "deleteEthNodes", idx)
}

// DeleteEthNodes is a paid mutator transaction binding the contract method 0x2cc267ea.
//
// Solidity: function deleteEthNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeSession) DeleteEthNodes(idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.Contract.DeleteEthNodes(&_AccelerateNode.TransactOpts, idx)
}

// DeleteEthNodes is a paid mutator transaction binding the contract method 0x2cc267ea.
//
// Solidity: function deleteEthNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) DeleteEthNodes(idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.Contract.DeleteEthNodes(&_AccelerateNode.TransactOpts, idx)
}

// DeleteIpfsNodes is a paid mutator transaction binding the contract method 0x2a5333e4.
//
// Solidity: function deleteIpfsNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeTransactor) DeleteIpfsNodes(opts *bind.TransactOpts, idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "deleteIpfsNodes", idx)
}

// DeleteIpfsNodes is a paid mutator transaction binding the contract method 0x2a5333e4.
//
// Solidity: function deleteIpfsNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeSession) DeleteIpfsNodes(idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.Contract.DeleteIpfsNodes(&_AccelerateNode.TransactOpts, idx)
}

// DeleteIpfsNodes is a paid mutator transaction binding the contract method 0x2a5333e4.
//
// Solidity: function deleteIpfsNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) DeleteIpfsNodes(idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.Contract.DeleteIpfsNodes(&_AccelerateNode.TransactOpts, idx)
}

// DeletePublicIpfsNodes is a paid mutator transaction binding the contract method 0xd458001f.
//
// Solidity: function deletePublicIpfsNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeTransactor) DeletePublicIpfsNodes(opts *bind.TransactOpts, idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "deletePublicIpfsNodes", idx)
}

// DeletePublicIpfsNodes is a paid mutator transaction binding the contract method 0xd458001f.
//
// Solidity: function deletePublicIpfsNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeSession) DeletePublicIpfsNodes(idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.Contract.DeletePublicIpfsNodes(&_AccelerateNode.TransactOpts, idx)
}

// DeletePublicIpfsNodes is a paid mutator transaction binding the contract method 0xd458001f.
//
// Solidity: function deletePublicIpfsNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) DeletePublicIpfsNodes(idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.Contract.DeletePublicIpfsNodes(&_AccelerateNode.TransactOpts, idx)
}

// DeleteSignerNodes is a paid mutator transaction binding the contract method 0x2ef4d35b.
//
// Solidity: function deleteSignerNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeTransactor) DeleteSignerNodes(opts *bind.TransactOpts, idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "deleteSignerNodes", idx)
}

// DeleteSignerNodes is a paid mutator transaction binding the contract method 0x2ef4d35b.
//
// Solidity: function deleteSignerNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeSession) DeleteSignerNodes(idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.Contract.DeleteSignerNodes(&_AccelerateNode.TransactOpts, idx)
}

// DeleteSignerNodes is a paid mutator transaction binding the contract method 0x2ef4d35b.
//
// Solidity: function deleteSignerNodes(uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) DeleteSignerNodes(idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.Contract.DeleteSignerNodes(&_AccelerateNode.TransactOpts, idx)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AccelerateNode *AccelerateNodeTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AccelerateNode *AccelerateNodeSession) RenounceOwnership() (*types.Transaction, error) {
	return _AccelerateNode.Contract.RenounceOwnership(&_AccelerateNode.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _AccelerateNode.Contract.RenounceOwnership(&_AccelerateNode.TransactOpts)
}

// ReplaceNode is a paid mutator transaction binding the contract method 0xb51d5308.
//
// Solidity: function replaceNode(string enode, uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeTransactor) ReplaceNode(opts *bind.TransactOpts, enode string, idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "replaceNode", enode, idx)
}

// ReplaceNode is a paid mutator transaction binding the contract method 0xb51d5308.
//
// Solidity: function replaceNode(string enode, uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeSession) ReplaceNode(enode string, idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.Contract.ReplaceNode(&_AccelerateNode.TransactOpts, enode, idx)
}

// ReplaceNode is a paid mutator transaction binding the contract method 0xb51d5308.
//
// Solidity: function replaceNode(string enode, uint32 idx) returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) ReplaceNode(enode string, idx uint32) (*types.Transaction, error) {
	return _AccelerateNode.Contract.ReplaceNode(&_AccelerateNode.TransactOpts, enode, idx)
}

// ResetEthNode is a paid mutator transaction binding the contract method 0x446f7a4a.
//
// Solidity: function resetEthNode() returns()
func (_AccelerateNode *AccelerateNodeTransactor) ResetEthNode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "resetEthNode")
}

// ResetEthNode is a paid mutator transaction binding the contract method 0x446f7a4a.
//
// Solidity: function resetEthNode() returns()
func (_AccelerateNode *AccelerateNodeSession) ResetEthNode() (*types.Transaction, error) {
	return _AccelerateNode.Contract.ResetEthNode(&_AccelerateNode.TransactOpts)
}

// ResetEthNode is a paid mutator transaction binding the contract method 0x446f7a4a.
//
// Solidity: function resetEthNode() returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) ResetEthNode() (*types.Transaction, error) {
	return _AccelerateNode.Contract.ResetEthNode(&_AccelerateNode.TransactOpts)
}

// ResetIpfsNode is a paid mutator transaction binding the contract method 0xf08502a5.
//
// Solidity: function resetIpfsNode() returns()
func (_AccelerateNode *AccelerateNodeTransactor) ResetIpfsNode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "resetIpfsNode")
}

// ResetIpfsNode is a paid mutator transaction binding the contract method 0xf08502a5.
//
// Solidity: function resetIpfsNode() returns()
func (_AccelerateNode *AccelerateNodeSession) ResetIpfsNode() (*types.Transaction, error) {
	return _AccelerateNode.Contract.ResetIpfsNode(&_AccelerateNode.TransactOpts)
}

// ResetIpfsNode is a paid mutator transaction binding the contract method 0xf08502a5.
//
// Solidity: function resetIpfsNode() returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) ResetIpfsNode() (*types.Transaction, error) {
	return _AccelerateNode.Contract.ResetIpfsNode(&_AccelerateNode.TransactOpts)
}

// ResetPublicIpfsNode is a paid mutator transaction binding the contract method 0x9d738796.
//
// Solidity: function resetPublicIpfsNode() returns()
func (_AccelerateNode *AccelerateNodeTransactor) ResetPublicIpfsNode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "resetPublicIpfsNode")
}

// ResetPublicIpfsNode is a paid mutator transaction binding the contract method 0x9d738796.
//
// Solidity: function resetPublicIpfsNode() returns()
func (_AccelerateNode *AccelerateNodeSession) ResetPublicIpfsNode() (*types.Transaction, error) {
	return _AccelerateNode.Contract.ResetPublicIpfsNode(&_AccelerateNode.TransactOpts)
}

// ResetPublicIpfsNode is a paid mutator transaction binding the contract method 0x9d738796.
//
// Solidity: function resetPublicIpfsNode() returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) ResetPublicIpfsNode() (*types.Transaction, error) {
	return _AccelerateNode.Contract.ResetPublicIpfsNode(&_AccelerateNode.TransactOpts)
}

// ResetSignerNode is a paid mutator transaction binding the contract method 0x85351535.
//
// Solidity: function resetSignerNode() returns()
func (_AccelerateNode *AccelerateNodeTransactor) ResetSignerNode(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "resetSignerNode")
}

// ResetSignerNode is a paid mutator transaction binding the contract method 0x85351535.
//
// Solidity: function resetSignerNode() returns()
func (_AccelerateNode *AccelerateNodeSession) ResetSignerNode() (*types.Transaction, error) {
	return _AccelerateNode.Contract.ResetSignerNode(&_AccelerateNode.TransactOpts)
}

// ResetSignerNode is a paid mutator transaction binding the contract method 0x85351535.
//
// Solidity: function resetSignerNode() returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) ResetSignerNode() (*types.Transaction, error) {
	return _AccelerateNode.Contract.ResetSignerNode(&_AccelerateNode.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AccelerateNode *AccelerateNodeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AccelerateNode.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AccelerateNode *AccelerateNodeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AccelerateNode.Contract.TransferOwnership(&_AccelerateNode.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AccelerateNode *AccelerateNodeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AccelerateNode.Contract.TransferOwnership(&_AccelerateNode.TransactOpts, newOwner)
}

// AccelerateNodeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the AccelerateNode contract.
type AccelerateNodeOwnershipTransferredIterator struct {
	Event *AccelerateNodeOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AccelerateNodeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AccelerateNodeOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AccelerateNodeOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AccelerateNodeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AccelerateNodeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AccelerateNodeOwnershipTransferred represents a OwnershipTransferred event raised by the AccelerateNode contract.
type AccelerateNodeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AccelerateNode *AccelerateNodeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*AccelerateNodeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AccelerateNode.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &AccelerateNodeOwnershipTransferredIterator{contract: _AccelerateNode.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AccelerateNode *AccelerateNodeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *AccelerateNodeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _AccelerateNode.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AccelerateNodeOwnershipTransferred)
				if err := _AccelerateNode.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_AccelerateNode *AccelerateNodeFilterer) ParseOwnershipTransferred(log types.Log) (*AccelerateNodeOwnershipTransferred, error) {
	event := new(AccelerateNodeOwnershipTransferred)
	if err := _AccelerateNode.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
