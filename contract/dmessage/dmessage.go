// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dmessage

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

// ContextABI is the input ABI used to generate the binding from.
const ContextABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// Context is an auto generated Go binding around an Ethereum contract.
type Context struct {
	ContextCaller     // Read-only binding to the contract
	ContextTransactor // Write-only binding to the contract
	ContextFilterer   // Log filterer for contract events
}

// ContextCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContextCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContextTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContextFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContextSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContextSession struct {
	Contract     *Context          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContextCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContextCallerSession struct {
	Contract *ContextCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// ContextTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContextTransactorSession struct {
	Contract     *ContextTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// ContextRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContextRaw struct {
	Contract *Context // Generic contract binding to access the raw methods on
}

// ContextCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContextCallerRaw struct {
	Contract *ContextCaller // Generic read-only contract binding to access the raw methods on
}

// ContextTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContextTransactorRaw struct {
	Contract *ContextTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContext creates a new instance of Context, bound to a specific deployed contract.
func NewContext(address common.Address, backend bind.ContractBackend) (*Context, error) {
	contract, err := bindContext(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Context{ContextCaller: ContextCaller{contract: contract}, ContextTransactor: ContextTransactor{contract: contract}, ContextFilterer: ContextFilterer{contract: contract}}, nil
}

// NewContextCaller creates a new read-only instance of Context, bound to a specific deployed contract.
func NewContextCaller(address common.Address, caller bind.ContractCaller) (*ContextCaller, error) {
	contract, err := bindContext(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContextCaller{contract: contract}, nil
}

// NewContextTransactor creates a new write-only instance of Context, bound to a specific deployed contract.
func NewContextTransactor(address common.Address, transactor bind.ContractTransactor) (*ContextTransactor, error) {
	contract, err := bindContext(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContextTransactor{contract: contract}, nil
}

// NewContextFilterer creates a new log filterer instance of Context, bound to a specific deployed contract.
func NewContextFilterer(address common.Address, filterer bind.ContractFilterer) (*ContextFilterer, error) {
	contract, err := bindContext(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContextFilterer{contract: contract}, nil
}

// bindContext binds a generic wrapper to an already deployed contract.
func bindContext(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContextABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Context.Contract.ContextCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.ContextTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Context *ContextCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Context.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Context *ContextTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Context.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Context *ContextTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Context.Contract.contract.Transact(opts, method, params...)
}

// DMessageABI is the input ABI used to generate the binding from.
const DMessageABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldWriter\",\"type\":\"address\"}],\"name\":\"WritershipDecreased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newWriter\",\"type\":\"address\"}],\"name\":\"WritershipIncreased\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"_id\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_message\",\"type\":\"string\"}],\"name\":\"addMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_ids\",\"type\":\"string[]\"}],\"name\":\"checkAllMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"string\",\"name\":\"_id\",\"type\":\"string\"}],\"name\":\"checkMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_ids\",\"type\":\"string[]\"}],\"name\":\"checkMessages\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"oldWriter\",\"type\":\"address\"}],\"name\":\"decreaseWritership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"_id\",\"type\":\"string\"}],\"name\":\"delMessage\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"},{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_ids\",\"type\":\"string[]\"}],\"name\":\"getIdsMessages\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"_value\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"string\",\"name\":\"_ids\",\"type\":\"string\"}],\"name\":\"getMessage\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"}],\"name\":\"getMessages\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"_value\",\"type\":\"string[]\"},{\"internalType\":\"uint256\",\"name\":\"_size\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getMsgId\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"_value\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMsgIds\",\"outputs\":[{\"internalType\":\"string[]\",\"name\":\"\",\"type\":\"string[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMsgLength\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newWriter\",\"type\":\"address\"}],\"name\":\"increasedWritership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isWriter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceWritership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"_id\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_message\",\"type\":\"string\"}],\"name\":\"updateMessage\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"writer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"writers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_value\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// DMessageFuncSigs maps the 4-byte function signature to its string representation.
var DMessageFuncSigs = map[string]string{
	"631e0c69": "addMessage(string,string)",
	"16313e92": "checkAllMessage(string[])",
	"1873ec68": "checkMessage(string)",
	"834daf0c": "checkMessages(string[])",
	"9f217bdf": "decreaseWritership(address)",
	"2309bc96": "delMessage(string)",
	"6c4ef994": "getIdsMessages(string[])",
	"0cc4e8d8": "getMessage(string)",
	"4be8362d": "getMessages(uint256,uint256)",
	"95cb2638": "getMsgId(uint256)",
	"68bfc43c": "getMsgIds()",
	"085c00e1": "getMsgLength()",
	"14b4d411": "increasedWritership(address)",
	"8f32d59b": "isOwner()",
	"4d6ee9fd": "isWriter()",
	"8da5cb5b": "owner()",
	"715018a6": "renounceOwnership()",
	"4075e83e": "renounceWritership()",
	"f2fde38b": "transferOwnership(address)",
	"75a0f3a4": "updateMessage(string,string)",
	"453a2abc": "writer()",
	"8b2f5369": "writers()",
}

// DMessageBin is the compiled bytecode used for deploying new contracts.
var DMessageBin = "0x608060405234801561001057600080fd5b5060006100246001600160e01b036100e216565b600080546001600160a01b0319166001600160a01b0383169081178255604051929350917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a35060006100816001600160e01b036100e216565b6001600160a01b0381166000818152600160208190526040808320805460ff191690921790915560038290555192935090917f8f2adc2e2f2723d375f8bb1359cb1027f713945395ca36c42b83dd94c5bdc7f29190a25060006007556100e6565b3390565b6118f5806100f56000396000f3fe608060405234801561001057600080fd5b50600436106101425760003560e01c806368bfc43c116100b85780638b2f53691161007c5780638b2f5369146102ab5780638da5cb5b146102c05780638f32d59b146102c857806395cb2638146102d05780639f217bdf146102e3578063f2fde38b146102f657610142565b806368bfc43c146102475780636c4ef9941461025c578063715018a61461026f57806375a0f3a414610277578063834daf0c1461028a57610142565b80632309bc961161010a5780632309bc96146101cd5780634075e83e146101ee578063453a2abc146101f65780634be8362d1461020b5780634d6ee9fd1461022c578063631e0c691461023457610142565b8063085c00e1146101475780630cc4e8d81461016557806314b4d4111461018557806316313e921461019a5780631873ec68146101ba575b600080fd5b61014f610309565b60405161015c9190611787565b60405180910390f35b6101786101733660046111a3565b610310565b60405161015c91906116c6565b610198610193366004611148565b6103be565b005b6101ad6101a836600461116e565b6103f7565b60405161015c91906116b8565b6101ad6101c83660046111a3565b61043f565b6101e06101db3660046111a3565b61046b565b60405161015c9291906116d7565b6101986105d7565b6101fe610670565b60405161015c9190611668565b61021e61021936600461125f565b610699565b60405161015c929190611698565b6101ad610728565b6101ad6102423660046111d8565b610758565b61024f610832565b60405161015c9190611687565b61021e61026a36600461116e565b610938565b6101986109c3565b6101ad6102853660046111d8565b610a31565b61029d61029836600461116e565b610ad7565b60405161015c929190611795565b6102b3610b18565b60405161015c9190611676565b6101fe610ba3565b6101ad610bb2565b6101786102de366004611241565b610bd6565b6101986102f1366004611148565b610c40565b610198610304366004611148565b610c6d565b6007545b90565b6060600582604051610322919061165c565b9081526040805160209281900383018120805460026001821615610100026000190190911604601f810185900485028301850190935282825290929091908301828280156103b15780601f10610386576101008083540402835291602001916103b1565b820191906000526020600020905b81548152906001019060200180831161039457829003601f168201915b505050505090505b919050565b6103c6610bb2565b6103eb5760405162461bcd60e51b81526004016103e290611767565b60405180910390fd5b6103f481610c9a565b50565b6000805b82518110156104365761042083828151811061041357fe5b602002602001015161043f565b61042e5760009150506103b9565b6001016103fb565b50600192915050565b6000600482604051610451919061165c565b9081526040519081900360200190205460ff169050919050565b60606000610477610bb2565b6104935760405162461bcd60e51b81526004016103e290611767565b6004836040516104a3919061165c565b9081526040519081900360200190205460ff1615156001146104d75760405162461bcd60e51b81526004016103e290611757565b60006004846040516104e9919061165c565b908152604051908190036020019020805491151560ff19909216919091179055610511610d6b565b600583604051610521919061165c565b90815260200160405180910390206001818054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156105c65780601f1061059b576101008083540402835291602001916105c6565b820191906000526020600020905b8154815290600101906020018083116105a957829003601f168201915b50505050509150915091505b915091565b6105df610bb2565b6105fb5760405162461bcd60e51b81526004016103e290611767565b600060016000610609610d9e565b6001600160a01b031681526020810191909152604001600020805460ff1916911515919091179055610639610d9e565b6001600160a01b03167fb24c36b70cc67df9ea22f053c15c4a41e7170d7be0637d2637a3445a266ac6af60405160405180910390a2565b600061067a610728565b1515600114156106935761068c610d9e565b905061030d565b50600090565b60606000600754838501106106b057836007540392505b826040519080825280602002602001820160405280156106e457816020015b60608152602001906001900390816106cf5790505b50915060005b8381101561071f57610700610173828701610bd6565b83828151811061070c57fe5b60209081029190910101526001016106ea565b50909391925050565b600060016000610736610d9e565b6001600160a01b0316815260208101919091526040016000205460ff16905090565b6000610762610728565b61077e5760405162461bcd60e51b81526004016103e290611747565b60048360405161078e919061165c565b9081526040519081900360200190205460ff16156107be5760405162461bcd60e51b81526004016103e290611707565b816005846040516107cf919061165c565b908152602001604051809103902090805190602001906107f0929190610f6e565b506001600484604051610803919061165c565b908152604051908190036020019020805491151560ff1990921691909117905561043683610da2565b92915050565b60608060075460405190808252806020026020018201604052801561086b57816020015b60608152602001906001900390816108565790505b50905060005b6007548110156109325760008181526006602090815260409182902080548351601f60026000196101006001861615020190931692909204918201849004840281018401909452808452909183018282801561090e5780601f106108e35761010080835404028352916020019161090e565b820191906000526020600020905b8154815290600101906020018083116108f157829003601f168201915b505050505082828151811061091f57fe5b6020908102919091010152600101610871565b50905090565b60606000825190508060405190808252806020026020018201604052801561097457816020015b606081526020019060019003908161095f5790505b50915060005b818110156109bd5761099e84828151811061099157fe5b6020026020010151610310565b8382815181106109aa57fe5b602090810291909101015260010161097a565b50915091565b6109cb610bb2565b6109e75760405162461bcd60e51b81526004016103e290611767565b600080546040516001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600080546001600160a01b0319169055565b6000610a3b610728565b610a575760405162461bcd60e51b81526004016103e290611747565b600483604051610a67919061165c565b9081526040519081900360200190205460ff161515600114610a9b5760405162461bcd60e51b81526004016103e290611777565b81600584604051610aac919061165c565b90815260200160405180910390209080519060200190610acd929190610f6e565b5060019392505050565b600080805b8351811015610b0b57610af484828151811061041357fe5b610b03579150600090506105d2565b600101610adc565b5060009360019350915050565b6060600354604051908082528060200260200182016040528015610b46578160200160208202803883390190505b50905060005b600354811015610b9f5760008181526002602052604090205482516001600160a01b0390911690839083908110610b7f57fe5b6001600160a01b0390921660209283029190910190910152600101610b4c565b5090565b6000546001600160a01b031690565b600080546001600160a01b0316610bc7610d9e565b6001600160a01b031614905090565b60008181526006602090815260409182902080548351601f60026000196101006001861615020190931692909204918201849004840281018401909452808452606093928301828280156103b15780601f10610386576101008083540402835291602001916103b1565b610c48610bb2565b610c645760405162461bcd60e51b81526004016103e290611767565b6103f481610dd1565b610c75610bb2565b610c915760405162461bcd60e51b81526004016103e290611767565b6103f481610e94565b6001600160a01b038116610cc05760405162461bcd60e51b81526004016103e290611737565b6001600160a01b03811660009081526001602052604090205460ff1615610cf95760405162461bcd60e51b81526004016103e290611727565b6001600160a01b0381166000818152600160208181526040808420805460ff191684179055600380548552600290925280842080546001600160a01b0319168617905581549092019055517f8f2adc2e2f2723d375f8bb1359cb1027f713945395ca36c42b83dd94c5bdc7f29190a250565b60005b610d76610309565b8110156103f457610d896101c882610bd6565b610d9657610d9681610f15565b600101610d6e565b3390565b60075460009081526006602090815260409091208251610dc492840190610f6e565b5050600780546001019055565b6001600160a01b03811660009081526001602081905260409091205460ff16151514610e0f5760405162461bcd60e51b81526004016103e2906116f7565b6001600160a01b038116610e355760405162461bcd60e51b81526004016103e290611737565b6001600160a01b0381166000908152600160205260409020805460ff19169055610e5d610d6b565b6040516001600160a01b038216907fb24c36b70cc67df9ea22f053c15c4a41e7170d7be0637d2637a3445a266ac6af90600090a250565b6001600160a01b038116610eba5760405162461bcd60e51b81526004016103e290611717565b600080546040516001600160a01b03808516939216917f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e091a3600080546001600160a01b0319166001600160a01b0392909216919091179055565b6007548110610f23576103f4565b60075460001990810160009081526006602052604080822084835291208154610f6093919291600261010060018416150290910190911604610fe8565b506007805460001901905550565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610faf57805160ff1916838001178555610fdc565b82800160010185558215610fdc579182015b82811115610fdc578251825591602001919060010190610fc1565b50610b9f92915061105d565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106110215780548555610fdc565b82800160010185558215610fdc57600052602060002091601f016020900482015b82811115610fdc578254825591600101919060010190611042565b61030d91905b80821115610b9f5760008155600101611063565b803561082c81611895565b600082601f83011261109357600080fd5b81356110a66110a1826117d7565b6117b0565b81815260209384019390925082018360005b838110156110e457813586016110ce88826110ee565b84525060209283019291909101906001016110b8565b5050505092915050565b600082601f8301126110ff57600080fd5b813561110d6110a1826117f8565b9150808252602083016020830185838301111561112957600080fd5b61113483828461184f565b50505092915050565b803561082c816118a9565b60006020828403121561115a57600080fd5b60006111668484611077565b949350505050565b60006020828403121561118057600080fd5b813567ffffffffffffffff81111561119757600080fd5b61116684828501611082565b6000602082840312156111b557600080fd5b813567ffffffffffffffff8111156111cc57600080fd5b611166848285016110ee565b600080604083850312156111eb57600080fd5b823567ffffffffffffffff81111561120257600080fd5b61120e858286016110ee565b925050602083013567ffffffffffffffff81111561122b57600080fd5b611237858286016110ee565b9150509250929050565b60006020828403121561125357600080fd5b6000611166848461113d565b6000806040838503121561127257600080fd5b600061127e858561113d565b92505060206112378582860161113d565b600061129b83836112b6565b505060200190565b60006112af8383611395565b9392505050565b6112bf81611833565b82525050565b60006112d082611826565b6112da818561182a565b93506112e583611820565b8060005b838110156113135781516112fd888261128f565b975061130883611820565b9250506001016112e9565b509495945050505050565b600061132982611826565b611333818561182a565b93508360208202850161134585611820565b8060005b8581101561137f578484038952815161136285826112a3565b945061136d83611820565b60209a909a0199925050600101611349565b5091979650505050505050565b6112bf8161183e565b60006113a082611826565b6113aa818561182a565b93506113ba81856020860161185b565b6113c38161188b565b9093019392505050565b60006113d882611826565b6113e281856103b9565b93506113f281856020860161185b565b9290920192915050565b600061140960258361182a565b7f57726974657261626c653a206f6c642077726974657220776173206e6f7420778152643934ba32b960d91b602082015260400192915050565b6000611450601d8361182a565b7f4d6573736167653a20616464206d657373616765206973206578697374000000815260200192915050565b600061148960268361182a565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206181526564647265737360d01b602082015260400192915050565b60006114d1601f8361182a565b7f57726974657261626c653a206e65772077726974657220697320657869737400815260200192915050565b600061150a602a8361182a565b7f57726974657261626c653a206e65772077726974657220697320746865207a65815269726f206164647265737360b01b602082015260400192915050565b600061155660238361182a565b7f577269746561626c653a2063616c6c6572206973206e6f7420746865207772698152623a32b960e91b602082015260400192915050565b600061159b60248361182a565b7f4d6573736167653a2064656c657465206d657373616765206973206e6f7420658152631e1a5cdd60e21b602082015260400192915050565b60006115e160208361182a565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572815260200192915050565b600061161a60248361182a565b7f4d6573736167653a20757064617465206d657373616765206973206e6f7420658152631e1a5cdd60e21b602082015260400192915050565b6112bf8161030d565b60006112af82846113cd565b6020810161082c82846112b6565b602080825281016112af81846112c5565b602080825281016112af818461131e565b604080825281016116a9818561131e565b90506112af6020830184611653565b6020810161082c828461138c565b602080825281016112af8184611395565b604080825281016116e88185611395565b90506112af602083018461138c565b6020808252810161082c816113fc565b6020808252810161082c81611443565b6020808252810161082c8161147c565b6020808252810161082c816114c4565b6020808252810161082c816114fd565b6020808252810161082c81611549565b6020808252810161082c8161158e565b6020808252810161082c816115d4565b6020808252810161082c8161160d565b6020810161082c8284611653565b604081016117a38285611653565b6112af602083018461138c565b60405181810167ffffffffffffffff811182821017156117cf57600080fd5b604052919050565b600067ffffffffffffffff8211156117ee57600080fd5b5060209081020190565b600067ffffffffffffffff82111561180f57600080fd5b506020601f91909101601f19160190565b60200190565b5190565b90815260200190565b600061082c82611843565b151590565b6001600160a01b031690565b82818337506000910152565b60005b8381101561187657818101518382015260200161185e565b83811115611885576000848401525b50505050565b601f01601f191690565b61189e81611833565b81146103f457600080fd5b61189e8161030d56fea365627a7a723158200160f12b57a18e062da7574a2f8614899e01db640f29147603ede50bbd08a0a66c6578706572696d656e74616cf564736f6c634300050c0040"

// DeployDMessage deploys a new Ethereum contract, binding an instance of DMessage to it.
func DeployDMessage(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DMessage, error) {
	parsed, err := abi.JSON(strings.NewReader(DMessageABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DMessageBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DMessage{DMessageCaller: DMessageCaller{contract: contract}, DMessageTransactor: DMessageTransactor{contract: contract}, DMessageFilterer: DMessageFilterer{contract: contract}}, nil
}

// DMessage is an auto generated Go binding around an Ethereum contract.
type DMessage struct {
	DMessageCaller     // Read-only binding to the contract
	DMessageTransactor // Write-only binding to the contract
	DMessageFilterer   // Log filterer for contract events
}

// DMessageCaller is an auto generated read-only Go binding around an Ethereum contract.
type DMessageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DMessageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DMessageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DMessageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DMessageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DMessageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DMessageSession struct {
	Contract     *DMessage         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DMessageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DMessageCallerSession struct {
	Contract *DMessageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// DMessageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DMessageTransactorSession struct {
	Contract     *DMessageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// DMessageRaw is an auto generated low-level Go binding around an Ethereum contract.
type DMessageRaw struct {
	Contract *DMessage // Generic contract binding to access the raw methods on
}

// DMessageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DMessageCallerRaw struct {
	Contract *DMessageCaller // Generic read-only contract binding to access the raw methods on
}

// DMessageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DMessageTransactorRaw struct {
	Contract *DMessageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDMessage creates a new instance of DMessage, bound to a specific deployed contract.
func NewDMessage(address common.Address, backend bind.ContractBackend) (*DMessage, error) {
	contract, err := bindDMessage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DMessage{DMessageCaller: DMessageCaller{contract: contract}, DMessageTransactor: DMessageTransactor{contract: contract}, DMessageFilterer: DMessageFilterer{contract: contract}}, nil
}

// NewDMessageCaller creates a new read-only instance of DMessage, bound to a specific deployed contract.
func NewDMessageCaller(address common.Address, caller bind.ContractCaller) (*DMessageCaller, error) {
	contract, err := bindDMessage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DMessageCaller{contract: contract}, nil
}

// NewDMessageTransactor creates a new write-only instance of DMessage, bound to a specific deployed contract.
func NewDMessageTransactor(address common.Address, transactor bind.ContractTransactor) (*DMessageTransactor, error) {
	contract, err := bindDMessage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DMessageTransactor{contract: contract}, nil
}

// NewDMessageFilterer creates a new log filterer instance of DMessage, bound to a specific deployed contract.
func NewDMessageFilterer(address common.Address, filterer bind.ContractFilterer) (*DMessageFilterer, error) {
	contract, err := bindDMessage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DMessageFilterer{contract: contract}, nil
}

// bindDMessage binds a generic wrapper to an already deployed contract.
func bindDMessage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DMessageABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DMessage *DMessageRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DMessage.Contract.DMessageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DMessage *DMessageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DMessage.Contract.DMessageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DMessage *DMessageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DMessage.Contract.DMessageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DMessage *DMessageCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DMessage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DMessage *DMessageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DMessage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DMessage *DMessageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DMessage.Contract.contract.Transact(opts, method, params...)
}

// CheckAllMessage is a free data retrieval call binding the contract method 0x16313e92.
//
// Solidity: function checkAllMessage(string[] _ids) constant returns(bool)
func (_DMessage *DMessageCaller) CheckAllMessage(opts *bind.CallOpts, _ids []string) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DMessage.contract.Call(opts, out, "checkAllMessage", _ids)
	return *ret0, err
}

// CheckAllMessage is a free data retrieval call binding the contract method 0x16313e92.
//
// Solidity: function checkAllMessage(string[] _ids) constant returns(bool)
func (_DMessage *DMessageSession) CheckAllMessage(_ids []string) (bool, error) {
	return _DMessage.Contract.CheckAllMessage(&_DMessage.CallOpts, _ids)
}

// CheckAllMessage is a free data retrieval call binding the contract method 0x16313e92.
//
// Solidity: function checkAllMessage(string[] _ids) constant returns(bool)
func (_DMessage *DMessageCallerSession) CheckAllMessage(_ids []string) (bool, error) {
	return _DMessage.Contract.CheckAllMessage(&_DMessage.CallOpts, _ids)
}

// CheckMessage is a free data retrieval call binding the contract method 0x1873ec68.
//
// Solidity: function checkMessage(string _id) constant returns(bool)
func (_DMessage *DMessageCaller) CheckMessage(opts *bind.CallOpts, _id string) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DMessage.contract.Call(opts, out, "checkMessage", _id)
	return *ret0, err
}

// CheckMessage is a free data retrieval call binding the contract method 0x1873ec68.
//
// Solidity: function checkMessage(string _id) constant returns(bool)
func (_DMessage *DMessageSession) CheckMessage(_id string) (bool, error) {
	return _DMessage.Contract.CheckMessage(&_DMessage.CallOpts, _id)
}

// CheckMessage is a free data retrieval call binding the contract method 0x1873ec68.
//
// Solidity: function checkMessage(string _id) constant returns(bool)
func (_DMessage *DMessageCallerSession) CheckMessage(_id string) (bool, error) {
	return _DMessage.Contract.CheckMessage(&_DMessage.CallOpts, _id)
}

// CheckMessages is a free data retrieval call binding the contract method 0x834daf0c.
//
// Solidity: function checkMessages(string[] _ids) constant returns(uint256, bool)
func (_DMessage *DMessageCaller) CheckMessages(opts *bind.CallOpts, _ids []string) (*big.Int, bool, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(bool)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _DMessage.contract.Call(opts, out, "checkMessages", _ids)
	return *ret0, *ret1, err
}

// CheckMessages is a free data retrieval call binding the contract method 0x834daf0c.
//
// Solidity: function checkMessages(string[] _ids) constant returns(uint256, bool)
func (_DMessage *DMessageSession) CheckMessages(_ids []string) (*big.Int, bool, error) {
	return _DMessage.Contract.CheckMessages(&_DMessage.CallOpts, _ids)
}

// CheckMessages is a free data retrieval call binding the contract method 0x834daf0c.
//
// Solidity: function checkMessages(string[] _ids) constant returns(uint256, bool)
func (_DMessage *DMessageCallerSession) CheckMessages(_ids []string) (*big.Int, bool, error) {
	return _DMessage.Contract.CheckMessages(&_DMessage.CallOpts, _ids)
}

// GetIdsMessages is a free data retrieval call binding the contract method 0x6c4ef994.
//
// Solidity: function getIdsMessages(string[] _ids) constant returns(string[] _value, uint256 _size)
func (_DMessage *DMessageCaller) GetIdsMessages(opts *bind.CallOpts, _ids []string) (struct {
	Value []string
	Size  *big.Int
}, error) {
	ret := new(struct {
		Value []string
		Size  *big.Int
	})
	out := ret
	err := _DMessage.contract.Call(opts, out, "getIdsMessages", _ids)
	return *ret, err
}

// GetIdsMessages is a free data retrieval call binding the contract method 0x6c4ef994.
//
// Solidity: function getIdsMessages(string[] _ids) constant returns(string[] _value, uint256 _size)
func (_DMessage *DMessageSession) GetIdsMessages(_ids []string) (struct {
	Value []string
	Size  *big.Int
}, error) {
	return _DMessage.Contract.GetIdsMessages(&_DMessage.CallOpts, _ids)
}

// GetIdsMessages is a free data retrieval call binding the contract method 0x6c4ef994.
//
// Solidity: function getIdsMessages(string[] _ids) constant returns(string[] _value, uint256 _size)
func (_DMessage *DMessageCallerSession) GetIdsMessages(_ids []string) (struct {
	Value []string
	Size  *big.Int
}, error) {
	return _DMessage.Contract.GetIdsMessages(&_DMessage.CallOpts, _ids)
}

// GetMessage is a free data retrieval call binding the contract method 0x0cc4e8d8.
//
// Solidity: function getMessage(string _ids) constant returns(string)
func (_DMessage *DMessageCaller) GetMessage(opts *bind.CallOpts, _ids string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _DMessage.contract.Call(opts, out, "getMessage", _ids)
	return *ret0, err
}

// GetMessage is a free data retrieval call binding the contract method 0x0cc4e8d8.
//
// Solidity: function getMessage(string _ids) constant returns(string)
func (_DMessage *DMessageSession) GetMessage(_ids string) (string, error) {
	return _DMessage.Contract.GetMessage(&_DMessage.CallOpts, _ids)
}

// GetMessage is a free data retrieval call binding the contract method 0x0cc4e8d8.
//
// Solidity: function getMessage(string _ids) constant returns(string)
func (_DMessage *DMessageCallerSession) GetMessage(_ids string) (string, error) {
	return _DMessage.Contract.GetMessage(&_DMessage.CallOpts, _ids)
}

// GetMessages is a free data retrieval call binding the contract method 0x4be8362d.
//
// Solidity: function getMessages(uint256 start, uint256 limit) constant returns(string[] _value, uint256 _size)
func (_DMessage *DMessageCaller) GetMessages(opts *bind.CallOpts, start *big.Int, limit *big.Int) (struct {
	Value []string
	Size  *big.Int
}, error) {
	ret := new(struct {
		Value []string
		Size  *big.Int
	})
	out := ret
	err := _DMessage.contract.Call(opts, out, "getMessages", start, limit)
	return *ret, err
}

// GetMessages is a free data retrieval call binding the contract method 0x4be8362d.
//
// Solidity: function getMessages(uint256 start, uint256 limit) constant returns(string[] _value, uint256 _size)
func (_DMessage *DMessageSession) GetMessages(start *big.Int, limit *big.Int) (struct {
	Value []string
	Size  *big.Int
}, error) {
	return _DMessage.Contract.GetMessages(&_DMessage.CallOpts, start, limit)
}

// GetMessages is a free data retrieval call binding the contract method 0x4be8362d.
//
// Solidity: function getMessages(uint256 start, uint256 limit) constant returns(string[] _value, uint256 _size)
func (_DMessage *DMessageCallerSession) GetMessages(start *big.Int, limit *big.Int) (struct {
	Value []string
	Size  *big.Int
}, error) {
	return _DMessage.Contract.GetMessages(&_DMessage.CallOpts, start, limit)
}

// GetMsgId is a free data retrieval call binding the contract method 0x95cb2638.
//
// Solidity: function getMsgId(uint256 _index) constant returns(string _value)
func (_DMessage *DMessageCaller) GetMsgId(opts *bind.CallOpts, _index *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _DMessage.contract.Call(opts, out, "getMsgId", _index)
	return *ret0, err
}

// GetMsgId is a free data retrieval call binding the contract method 0x95cb2638.
//
// Solidity: function getMsgId(uint256 _index) constant returns(string _value)
func (_DMessage *DMessageSession) GetMsgId(_index *big.Int) (string, error) {
	return _DMessage.Contract.GetMsgId(&_DMessage.CallOpts, _index)
}

// GetMsgId is a free data retrieval call binding the contract method 0x95cb2638.
//
// Solidity: function getMsgId(uint256 _index) constant returns(string _value)
func (_DMessage *DMessageCallerSession) GetMsgId(_index *big.Int) (string, error) {
	return _DMessage.Contract.GetMsgId(&_DMessage.CallOpts, _index)
}

// GetMsgIds is a free data retrieval call binding the contract method 0x68bfc43c.
//
// Solidity: function getMsgIds() constant returns(string[])
func (_DMessage *DMessageCaller) GetMsgIds(opts *bind.CallOpts) ([]string, error) {
	var (
		ret0 = new([]string)
	)
	out := ret0
	err := _DMessage.contract.Call(opts, out, "getMsgIds")
	return *ret0, err
}

// GetMsgIds is a free data retrieval call binding the contract method 0x68bfc43c.
//
// Solidity: function getMsgIds() constant returns(string[])
func (_DMessage *DMessageSession) GetMsgIds() ([]string, error) {
	return _DMessage.Contract.GetMsgIds(&_DMessage.CallOpts)
}

// GetMsgIds is a free data retrieval call binding the contract method 0x68bfc43c.
//
// Solidity: function getMsgIds() constant returns(string[])
func (_DMessage *DMessageCallerSession) GetMsgIds() ([]string, error) {
	return _DMessage.Contract.GetMsgIds(&_DMessage.CallOpts)
}

// GetMsgLength is a free data retrieval call binding the contract method 0x085c00e1.
//
// Solidity: function getMsgLength() constant returns(uint256)
func (_DMessage *DMessageCaller) GetMsgLength(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DMessage.contract.Call(opts, out, "getMsgLength")
	return *ret0, err
}

// GetMsgLength is a free data retrieval call binding the contract method 0x085c00e1.
//
// Solidity: function getMsgLength() constant returns(uint256)
func (_DMessage *DMessageSession) GetMsgLength() (*big.Int, error) {
	return _DMessage.Contract.GetMsgLength(&_DMessage.CallOpts)
}

// GetMsgLength is a free data retrieval call binding the contract method 0x085c00e1.
//
// Solidity: function getMsgLength() constant returns(uint256)
func (_DMessage *DMessageCallerSession) GetMsgLength() (*big.Int, error) {
	return _DMessage.Contract.GetMsgLength(&_DMessage.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_DMessage *DMessageCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DMessage.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_DMessage *DMessageSession) IsOwner() (bool, error) {
	return _DMessage.Contract.IsOwner(&_DMessage.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_DMessage *DMessageCallerSession) IsOwner() (bool, error) {
	return _DMessage.Contract.IsOwner(&_DMessage.CallOpts)
}

// IsWriter is a free data retrieval call binding the contract method 0x4d6ee9fd.
//
// Solidity: function isWriter() constant returns(bool)
func (_DMessage *DMessageCaller) IsWriter(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DMessage.contract.Call(opts, out, "isWriter")
	return *ret0, err
}

// IsWriter is a free data retrieval call binding the contract method 0x4d6ee9fd.
//
// Solidity: function isWriter() constant returns(bool)
func (_DMessage *DMessageSession) IsWriter() (bool, error) {
	return _DMessage.Contract.IsWriter(&_DMessage.CallOpts)
}

// IsWriter is a free data retrieval call binding the contract method 0x4d6ee9fd.
//
// Solidity: function isWriter() constant returns(bool)
func (_DMessage *DMessageCallerSession) IsWriter() (bool, error) {
	return _DMessage.Contract.IsWriter(&_DMessage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DMessage *DMessageCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DMessage.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DMessage *DMessageSession) Owner() (common.Address, error) {
	return _DMessage.Contract.Owner(&_DMessage.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DMessage *DMessageCallerSession) Owner() (common.Address, error) {
	return _DMessage.Contract.Owner(&_DMessage.CallOpts)
}

// Writer is a free data retrieval call binding the contract method 0x453a2abc.
//
// Solidity: function writer() constant returns(address)
func (_DMessage *DMessageCaller) Writer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DMessage.contract.Call(opts, out, "writer")
	return *ret0, err
}

// Writer is a free data retrieval call binding the contract method 0x453a2abc.
//
// Solidity: function writer() constant returns(address)
func (_DMessage *DMessageSession) Writer() (common.Address, error) {
	return _DMessage.Contract.Writer(&_DMessage.CallOpts)
}

// Writer is a free data retrieval call binding the contract method 0x453a2abc.
//
// Solidity: function writer() constant returns(address)
func (_DMessage *DMessageCallerSession) Writer() (common.Address, error) {
	return _DMessage.Contract.Writer(&_DMessage.CallOpts)
}

// Writers is a free data retrieval call binding the contract method 0x8b2f5369.
//
// Solidity: function writers() constant returns(address[] _value)
func (_DMessage *DMessageCaller) Writers(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _DMessage.contract.Call(opts, out, "writers")
	return *ret0, err
}

// Writers is a free data retrieval call binding the contract method 0x8b2f5369.
//
// Solidity: function writers() constant returns(address[] _value)
func (_DMessage *DMessageSession) Writers() ([]common.Address, error) {
	return _DMessage.Contract.Writers(&_DMessage.CallOpts)
}

// Writers is a free data retrieval call binding the contract method 0x8b2f5369.
//
// Solidity: function writers() constant returns(address[] _value)
func (_DMessage *DMessageCallerSession) Writers() ([]common.Address, error) {
	return _DMessage.Contract.Writers(&_DMessage.CallOpts)
}

// AddMessage is a paid mutator transaction binding the contract method 0x631e0c69.
//
// Solidity: function addMessage(string _id, string _message) returns(bool)
func (_DMessage *DMessageTransactor) AddMessage(opts *bind.TransactOpts, _id string, _message string) (*types.Transaction, error) {
	return _DMessage.contract.Transact(opts, "addMessage", _id, _message)
}

// AddMessage is a paid mutator transaction binding the contract method 0x631e0c69.
//
// Solidity: function addMessage(string _id, string _message) returns(bool)
func (_DMessage *DMessageSession) AddMessage(_id string, _message string) (*types.Transaction, error) {
	return _DMessage.Contract.AddMessage(&_DMessage.TransactOpts, _id, _message)
}

// AddMessage is a paid mutator transaction binding the contract method 0x631e0c69.
//
// Solidity: function addMessage(string _id, string _message) returns(bool)
func (_DMessage *DMessageTransactorSession) AddMessage(_id string, _message string) (*types.Transaction, error) {
	return _DMessage.Contract.AddMessage(&_DMessage.TransactOpts, _id, _message)
}

// DecreaseWritership is a paid mutator transaction binding the contract method 0x9f217bdf.
//
// Solidity: function decreaseWritership(address oldWriter) returns()
func (_DMessage *DMessageTransactor) DecreaseWritership(opts *bind.TransactOpts, oldWriter common.Address) (*types.Transaction, error) {
	return _DMessage.contract.Transact(opts, "decreaseWritership", oldWriter)
}

// DecreaseWritership is a paid mutator transaction binding the contract method 0x9f217bdf.
//
// Solidity: function decreaseWritership(address oldWriter) returns()
func (_DMessage *DMessageSession) DecreaseWritership(oldWriter common.Address) (*types.Transaction, error) {
	return _DMessage.Contract.DecreaseWritership(&_DMessage.TransactOpts, oldWriter)
}

// DecreaseWritership is a paid mutator transaction binding the contract method 0x9f217bdf.
//
// Solidity: function decreaseWritership(address oldWriter) returns()
func (_DMessage *DMessageTransactorSession) DecreaseWritership(oldWriter common.Address) (*types.Transaction, error) {
	return _DMessage.Contract.DecreaseWritership(&_DMessage.TransactOpts, oldWriter)
}

// DelMessage is a paid mutator transaction binding the contract method 0x2309bc96.
//
// Solidity: function delMessage(string _id) returns(string, bool)
func (_DMessage *DMessageTransactor) DelMessage(opts *bind.TransactOpts, _id string) (*types.Transaction, error) {
	return _DMessage.contract.Transact(opts, "delMessage", _id)
}

// DelMessage is a paid mutator transaction binding the contract method 0x2309bc96.
//
// Solidity: function delMessage(string _id) returns(string, bool)
func (_DMessage *DMessageSession) DelMessage(_id string) (*types.Transaction, error) {
	return _DMessage.Contract.DelMessage(&_DMessage.TransactOpts, _id)
}

// DelMessage is a paid mutator transaction binding the contract method 0x2309bc96.
//
// Solidity: function delMessage(string _id) returns(string, bool)
func (_DMessage *DMessageTransactorSession) DelMessage(_id string) (*types.Transaction, error) {
	return _DMessage.Contract.DelMessage(&_DMessage.TransactOpts, _id)
}

// IncreasedWritership is a paid mutator transaction binding the contract method 0x14b4d411.
//
// Solidity: function increasedWritership(address newWriter) returns()
func (_DMessage *DMessageTransactor) IncreasedWritership(opts *bind.TransactOpts, newWriter common.Address) (*types.Transaction, error) {
	return _DMessage.contract.Transact(opts, "increasedWritership", newWriter)
}

// IncreasedWritership is a paid mutator transaction binding the contract method 0x14b4d411.
//
// Solidity: function increasedWritership(address newWriter) returns()
func (_DMessage *DMessageSession) IncreasedWritership(newWriter common.Address) (*types.Transaction, error) {
	return _DMessage.Contract.IncreasedWritership(&_DMessage.TransactOpts, newWriter)
}

// IncreasedWritership is a paid mutator transaction binding the contract method 0x14b4d411.
//
// Solidity: function increasedWritership(address newWriter) returns()
func (_DMessage *DMessageTransactorSession) IncreasedWritership(newWriter common.Address) (*types.Transaction, error) {
	return _DMessage.Contract.IncreasedWritership(&_DMessage.TransactOpts, newWriter)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DMessage *DMessageTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DMessage.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DMessage *DMessageSession) RenounceOwnership() (*types.Transaction, error) {
	return _DMessage.Contract.RenounceOwnership(&_DMessage.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DMessage *DMessageTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DMessage.Contract.RenounceOwnership(&_DMessage.TransactOpts)
}

// RenounceWritership is a paid mutator transaction binding the contract method 0x4075e83e.
//
// Solidity: function renounceWritership() returns()
func (_DMessage *DMessageTransactor) RenounceWritership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DMessage.contract.Transact(opts, "renounceWritership")
}

// RenounceWritership is a paid mutator transaction binding the contract method 0x4075e83e.
//
// Solidity: function renounceWritership() returns()
func (_DMessage *DMessageSession) RenounceWritership() (*types.Transaction, error) {
	return _DMessage.Contract.RenounceWritership(&_DMessage.TransactOpts)
}

// RenounceWritership is a paid mutator transaction binding the contract method 0x4075e83e.
//
// Solidity: function renounceWritership() returns()
func (_DMessage *DMessageTransactorSession) RenounceWritership() (*types.Transaction, error) {
	return _DMessage.Contract.RenounceWritership(&_DMessage.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DMessage *DMessageTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DMessage.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DMessage *DMessageSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DMessage.Contract.TransferOwnership(&_DMessage.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DMessage *DMessageTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DMessage.Contract.TransferOwnership(&_DMessage.TransactOpts, newOwner)
}

// UpdateMessage is a paid mutator transaction binding the contract method 0x75a0f3a4.
//
// Solidity: function updateMessage(string _id, string _message) returns(bool)
func (_DMessage *DMessageTransactor) UpdateMessage(opts *bind.TransactOpts, _id string, _message string) (*types.Transaction, error) {
	return _DMessage.contract.Transact(opts, "updateMessage", _id, _message)
}

// UpdateMessage is a paid mutator transaction binding the contract method 0x75a0f3a4.
//
// Solidity: function updateMessage(string _id, string _message) returns(bool)
func (_DMessage *DMessageSession) UpdateMessage(_id string, _message string) (*types.Transaction, error) {
	return _DMessage.Contract.UpdateMessage(&_DMessage.TransactOpts, _id, _message)
}

// UpdateMessage is a paid mutator transaction binding the contract method 0x75a0f3a4.
//
// Solidity: function updateMessage(string _id, string _message) returns(bool)
func (_DMessage *DMessageTransactorSession) UpdateMessage(_id string, _message string) (*types.Transaction, error) {
	return _DMessage.Contract.UpdateMessage(&_DMessage.TransactOpts, _id, _message)
}

// DMessageOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DMessage contract.
type DMessageOwnershipTransferredIterator struct {
	Event *DMessageOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DMessageOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DMessageOwnershipTransferred)
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
		it.Event = new(DMessageOwnershipTransferred)
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
func (it *DMessageOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DMessageOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DMessageOwnershipTransferred represents a OwnershipTransferred event raised by the DMessage contract.
type DMessageOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DMessage *DMessageFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DMessageOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DMessage.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DMessageOwnershipTransferredIterator{contract: _DMessage.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DMessage *DMessageFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DMessageOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DMessage.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DMessageOwnershipTransferred)
				if err := _DMessage.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_DMessage *DMessageFilterer) ParseOwnershipTransferred(log types.Log) (*DMessageOwnershipTransferred, error) {
	event := new(DMessageOwnershipTransferred)
	if err := _DMessage.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DMessageWritershipDecreasedIterator is returned from FilterWritershipDecreased and is used to iterate over the raw logs and unpacked data for WritershipDecreased events raised by the DMessage contract.
type DMessageWritershipDecreasedIterator struct {
	Event *DMessageWritershipDecreased // Event containing the contract specifics and raw log

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
func (it *DMessageWritershipDecreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DMessageWritershipDecreased)
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
		it.Event = new(DMessageWritershipDecreased)
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
func (it *DMessageWritershipDecreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DMessageWritershipDecreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DMessageWritershipDecreased represents a WritershipDecreased event raised by the DMessage contract.
type DMessageWritershipDecreased struct {
	OldWriter common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWritershipDecreased is a free log retrieval operation binding the contract event 0xb24c36b70cc67df9ea22f053c15c4a41e7170d7be0637d2637a3445a266ac6af.
//
// Solidity: event WritershipDecreased(address indexed oldWriter)
func (_DMessage *DMessageFilterer) FilterWritershipDecreased(opts *bind.FilterOpts, oldWriter []common.Address) (*DMessageWritershipDecreasedIterator, error) {

	var oldWriterRule []interface{}
	for _, oldWriterItem := range oldWriter {
		oldWriterRule = append(oldWriterRule, oldWriterItem)
	}

	logs, sub, err := _DMessage.contract.FilterLogs(opts, "WritershipDecreased", oldWriterRule)
	if err != nil {
		return nil, err
	}
	return &DMessageWritershipDecreasedIterator{contract: _DMessage.contract, event: "WritershipDecreased", logs: logs, sub: sub}, nil
}

// WatchWritershipDecreased is a free log subscription operation binding the contract event 0xb24c36b70cc67df9ea22f053c15c4a41e7170d7be0637d2637a3445a266ac6af.
//
// Solidity: event WritershipDecreased(address indexed oldWriter)
func (_DMessage *DMessageFilterer) WatchWritershipDecreased(opts *bind.WatchOpts, sink chan<- *DMessageWritershipDecreased, oldWriter []common.Address) (event.Subscription, error) {

	var oldWriterRule []interface{}
	for _, oldWriterItem := range oldWriter {
		oldWriterRule = append(oldWriterRule, oldWriterItem)
	}

	logs, sub, err := _DMessage.contract.WatchLogs(opts, "WritershipDecreased", oldWriterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DMessageWritershipDecreased)
				if err := _DMessage.contract.UnpackLog(event, "WritershipDecreased", log); err != nil {
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

// ParseWritershipDecreased is a log parse operation binding the contract event 0xb24c36b70cc67df9ea22f053c15c4a41e7170d7be0637d2637a3445a266ac6af.
//
// Solidity: event WritershipDecreased(address indexed oldWriter)
func (_DMessage *DMessageFilterer) ParseWritershipDecreased(log types.Log) (*DMessageWritershipDecreased, error) {
	event := new(DMessageWritershipDecreased)
	if err := _DMessage.contract.UnpackLog(event, "WritershipDecreased", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DMessageWritershipIncreasedIterator is returned from FilterWritershipIncreased and is used to iterate over the raw logs and unpacked data for WritershipIncreased events raised by the DMessage contract.
type DMessageWritershipIncreasedIterator struct {
	Event *DMessageWritershipIncreased // Event containing the contract specifics and raw log

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
func (it *DMessageWritershipIncreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DMessageWritershipIncreased)
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
		it.Event = new(DMessageWritershipIncreased)
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
func (it *DMessageWritershipIncreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DMessageWritershipIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DMessageWritershipIncreased represents a WritershipIncreased event raised by the DMessage contract.
type DMessageWritershipIncreased struct {
	NewWriter common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWritershipIncreased is a free log retrieval operation binding the contract event 0x8f2adc2e2f2723d375f8bb1359cb1027f713945395ca36c42b83dd94c5bdc7f2.
//
// Solidity: event WritershipIncreased(address indexed newWriter)
func (_DMessage *DMessageFilterer) FilterWritershipIncreased(opts *bind.FilterOpts, newWriter []common.Address) (*DMessageWritershipIncreasedIterator, error) {

	var newWriterRule []interface{}
	for _, newWriterItem := range newWriter {
		newWriterRule = append(newWriterRule, newWriterItem)
	}

	logs, sub, err := _DMessage.contract.FilterLogs(opts, "WritershipIncreased", newWriterRule)
	if err != nil {
		return nil, err
	}
	return &DMessageWritershipIncreasedIterator{contract: _DMessage.contract, event: "WritershipIncreased", logs: logs, sub: sub}, nil
}

// WatchWritershipIncreased is a free log subscription operation binding the contract event 0x8f2adc2e2f2723d375f8bb1359cb1027f713945395ca36c42b83dd94c5bdc7f2.
//
// Solidity: event WritershipIncreased(address indexed newWriter)
func (_DMessage *DMessageFilterer) WatchWritershipIncreased(opts *bind.WatchOpts, sink chan<- *DMessageWritershipIncreased, newWriter []common.Address) (event.Subscription, error) {

	var newWriterRule []interface{}
	for _, newWriterItem := range newWriter {
		newWriterRule = append(newWriterRule, newWriterItem)
	}

	logs, sub, err := _DMessage.contract.WatchLogs(opts, "WritershipIncreased", newWriterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DMessageWritershipIncreased)
				if err := _DMessage.contract.UnpackLog(event, "WritershipIncreased", log); err != nil {
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

// ParseWritershipIncreased is a log parse operation binding the contract event 0x8f2adc2e2f2723d375f8bb1359cb1027f713945395ca36c42b83dd94c5bdc7f2.
//
// Solidity: event WritershipIncreased(address indexed newWriter)
func (_DMessage *DMessageFilterer) ParseWritershipIncreased(log types.Log) (*DMessageWritershipIncreased, error) {
	event := new(DMessageWritershipIncreased)
	if err := _DMessage.contract.UnpackLog(event, "WritershipIncreased", log); err != nil {
		return nil, err
	}
	return event, nil
}

// OwnableABI is the input ABI used to generate the binding from.
const OwnableABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// OwnableFuncSigs maps the 4-byte function signature to its string representation.
var OwnableFuncSigs = map[string]string{
	"8f32d59b": "isOwner()",
	"8da5cb5b": "owner()",
	"715018a6": "renounceOwnership()",
	"f2fde38b": "transferOwnership(address)",
}

// Ownable is an auto generated Go binding around an Ethereum contract.
type Ownable struct {
	OwnableCaller     // Read-only binding to the contract
	OwnableTransactor // Write-only binding to the contract
	OwnableFilterer   // Log filterer for contract events
}

// OwnableCaller is an auto generated read-only Go binding around an Ethereum contract.
type OwnableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OwnableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OwnableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OwnableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OwnableSession struct {
	Contract     *Ownable          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OwnableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OwnableCallerSession struct {
	Contract *OwnableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// OwnableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OwnableTransactorSession struct {
	Contract     *OwnableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// OwnableRaw is an auto generated low-level Go binding around an Ethereum contract.
type OwnableRaw struct {
	Contract *Ownable // Generic contract binding to access the raw methods on
}

// OwnableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OwnableCallerRaw struct {
	Contract *OwnableCaller // Generic read-only contract binding to access the raw methods on
}

// OwnableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OwnableTransactorRaw struct {
	Contract *OwnableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOwnable creates a new instance of Ownable, bound to a specific deployed contract.
func NewOwnable(address common.Address, backend bind.ContractBackend) (*Ownable, error) {
	contract, err := bindOwnable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ownable{OwnableCaller: OwnableCaller{contract: contract}, OwnableTransactor: OwnableTransactor{contract: contract}, OwnableFilterer: OwnableFilterer{contract: contract}}, nil
}

// NewOwnableCaller creates a new read-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableCaller(address common.Address, caller bind.ContractCaller) (*OwnableCaller, error) {
	contract, err := bindOwnable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableCaller{contract: contract}, nil
}

// NewOwnableTransactor creates a new write-only instance of Ownable, bound to a specific deployed contract.
func NewOwnableTransactor(address common.Address, transactor bind.ContractTransactor) (*OwnableTransactor, error) {
	contract, err := bindOwnable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OwnableTransactor{contract: contract}, nil
}

// NewOwnableFilterer creates a new log filterer instance of Ownable, bound to a specific deployed contract.
func NewOwnableFilterer(address common.Address, filterer bind.ContractFilterer) (*OwnableFilterer, error) {
	contract, err := bindOwnable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OwnableFilterer{contract: contract}, nil
}

// bindOwnable binds a generic wrapper to an already deployed contract.
func bindOwnable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OwnableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.OwnableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.OwnableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ownable *OwnableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Ownable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ownable *OwnableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ownable *OwnableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ownable.Contract.contract.Transact(opts, method, params...)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Ownable *OwnableCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Ownable.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Ownable *OwnableSession) IsOwner() (bool, error) {
	return _Ownable.Contract.IsOwner(&_Ownable.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Ownable *OwnableCallerSession) IsOwner() (bool, error) {
	return _Ownable.Contract.IsOwner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Ownable.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Ownable *OwnableCallerSession) Owner() (common.Address, error) {
	return _Ownable.Contract.Owner(&_Ownable.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Ownable *OwnableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Ownable.Contract.RenounceOwnership(&_Ownable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ownable *OwnableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ownable.Contract.TransferOwnership(&_Ownable.TransactOpts, newOwner)
}

// OwnableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Ownable contract.
type OwnableOwnershipTransferredIterator struct {
	Event *OwnableOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *OwnableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OwnableOwnershipTransferred)
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
		it.Event = new(OwnableOwnershipTransferred)
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
func (it *OwnableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OwnableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OwnableOwnershipTransferred represents a OwnershipTransferred event raised by the Ownable contract.
type OwnableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*OwnableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &OwnableOwnershipTransferredIterator{contract: _Ownable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Ownable *OwnableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *OwnableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Ownable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OwnableOwnershipTransferred)
				if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Ownable *OwnableFilterer) ParseOwnershipTransferred(log types.Log) (*OwnableOwnershipTransferred, error) {
	event := new(OwnableOwnershipTransferred)
	if err := _Ownable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WriteableABI is the input ABI used to generate the binding from.
const WriteableABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldWriter\",\"type\":\"address\"}],\"name\":\"WritershipDecreased\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newWriter\",\"type\":\"address\"}],\"name\":\"WritershipIncreased\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"oldWriter\",\"type\":\"address\"}],\"name\":\"decreaseWritership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newWriter\",\"type\":\"address\"}],\"name\":\"increasedWritership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isWriter\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceWritership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"writer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"writers\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"_value\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// WriteableFuncSigs maps the 4-byte function signature to its string representation.
var WriteableFuncSigs = map[string]string{
	"9f217bdf": "decreaseWritership(address)",
	"14b4d411": "increasedWritership(address)",
	"8f32d59b": "isOwner()",
	"4d6ee9fd": "isWriter()",
	"8da5cb5b": "owner()",
	"715018a6": "renounceOwnership()",
	"4075e83e": "renounceWritership()",
	"f2fde38b": "transferOwnership(address)",
	"453a2abc": "writer()",
	"8b2f5369": "writers()",
}

// Writeable is an auto generated Go binding around an Ethereum contract.
type Writeable struct {
	WriteableCaller     // Read-only binding to the contract
	WriteableTransactor // Write-only binding to the contract
	WriteableFilterer   // Log filterer for contract events
}

// WriteableCaller is an auto generated read-only Go binding around an Ethereum contract.
type WriteableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WriteableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type WriteableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WriteableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type WriteableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// WriteableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type WriteableSession struct {
	Contract     *Writeable        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// WriteableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type WriteableCallerSession struct {
	Contract *WriteableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// WriteableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type WriteableTransactorSession struct {
	Contract     *WriteableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// WriteableRaw is an auto generated low-level Go binding around an Ethereum contract.
type WriteableRaw struct {
	Contract *Writeable // Generic contract binding to access the raw methods on
}

// WriteableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type WriteableCallerRaw struct {
	Contract *WriteableCaller // Generic read-only contract binding to access the raw methods on
}

// WriteableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type WriteableTransactorRaw struct {
	Contract *WriteableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewWriteable creates a new instance of Writeable, bound to a specific deployed contract.
func NewWriteable(address common.Address, backend bind.ContractBackend) (*Writeable, error) {
	contract, err := bindWriteable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Writeable{WriteableCaller: WriteableCaller{contract: contract}, WriteableTransactor: WriteableTransactor{contract: contract}, WriteableFilterer: WriteableFilterer{contract: contract}}, nil
}

// NewWriteableCaller creates a new read-only instance of Writeable, bound to a specific deployed contract.
func NewWriteableCaller(address common.Address, caller bind.ContractCaller) (*WriteableCaller, error) {
	contract, err := bindWriteable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &WriteableCaller{contract: contract}, nil
}

// NewWriteableTransactor creates a new write-only instance of Writeable, bound to a specific deployed contract.
func NewWriteableTransactor(address common.Address, transactor bind.ContractTransactor) (*WriteableTransactor, error) {
	contract, err := bindWriteable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &WriteableTransactor{contract: contract}, nil
}

// NewWriteableFilterer creates a new log filterer instance of Writeable, bound to a specific deployed contract.
func NewWriteableFilterer(address common.Address, filterer bind.ContractFilterer) (*WriteableFilterer, error) {
	contract, err := bindWriteable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &WriteableFilterer{contract: contract}, nil
}

// bindWriteable binds a generic wrapper to an already deployed contract.
func bindWriteable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(WriteableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Writeable *WriteableRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Writeable.Contract.WriteableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Writeable *WriteableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Writeable.Contract.WriteableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Writeable *WriteableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Writeable.Contract.WriteableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Writeable *WriteableCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Writeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Writeable *WriteableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Writeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Writeable *WriteableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Writeable.Contract.contract.Transact(opts, method, params...)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Writeable *WriteableCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Writeable.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Writeable *WriteableSession) IsOwner() (bool, error) {
	return _Writeable.Contract.IsOwner(&_Writeable.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Writeable *WriteableCallerSession) IsOwner() (bool, error) {
	return _Writeable.Contract.IsOwner(&_Writeable.CallOpts)
}

// IsWriter is a free data retrieval call binding the contract method 0x4d6ee9fd.
//
// Solidity: function isWriter() constant returns(bool)
func (_Writeable *WriteableCaller) IsWriter(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Writeable.contract.Call(opts, out, "isWriter")
	return *ret0, err
}

// IsWriter is a free data retrieval call binding the contract method 0x4d6ee9fd.
//
// Solidity: function isWriter() constant returns(bool)
func (_Writeable *WriteableSession) IsWriter() (bool, error) {
	return _Writeable.Contract.IsWriter(&_Writeable.CallOpts)
}

// IsWriter is a free data retrieval call binding the contract method 0x4d6ee9fd.
//
// Solidity: function isWriter() constant returns(bool)
func (_Writeable *WriteableCallerSession) IsWriter() (bool, error) {
	return _Writeable.Contract.IsWriter(&_Writeable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Writeable *WriteableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Writeable.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Writeable *WriteableSession) Owner() (common.Address, error) {
	return _Writeable.Contract.Owner(&_Writeable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Writeable *WriteableCallerSession) Owner() (common.Address, error) {
	return _Writeable.Contract.Owner(&_Writeable.CallOpts)
}

// Writer is a free data retrieval call binding the contract method 0x453a2abc.
//
// Solidity: function writer() constant returns(address)
func (_Writeable *WriteableCaller) Writer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Writeable.contract.Call(opts, out, "writer")
	return *ret0, err
}

// Writer is a free data retrieval call binding the contract method 0x453a2abc.
//
// Solidity: function writer() constant returns(address)
func (_Writeable *WriteableSession) Writer() (common.Address, error) {
	return _Writeable.Contract.Writer(&_Writeable.CallOpts)
}

// Writer is a free data retrieval call binding the contract method 0x453a2abc.
//
// Solidity: function writer() constant returns(address)
func (_Writeable *WriteableCallerSession) Writer() (common.Address, error) {
	return _Writeable.Contract.Writer(&_Writeable.CallOpts)
}

// Writers is a free data retrieval call binding the contract method 0x8b2f5369.
//
// Solidity: function writers() constant returns(address[] _value)
func (_Writeable *WriteableCaller) Writers(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Writeable.contract.Call(opts, out, "writers")
	return *ret0, err
}

// Writers is a free data retrieval call binding the contract method 0x8b2f5369.
//
// Solidity: function writers() constant returns(address[] _value)
func (_Writeable *WriteableSession) Writers() ([]common.Address, error) {
	return _Writeable.Contract.Writers(&_Writeable.CallOpts)
}

// Writers is a free data retrieval call binding the contract method 0x8b2f5369.
//
// Solidity: function writers() constant returns(address[] _value)
func (_Writeable *WriteableCallerSession) Writers() ([]common.Address, error) {
	return _Writeable.Contract.Writers(&_Writeable.CallOpts)
}

// DecreaseWritership is a paid mutator transaction binding the contract method 0x9f217bdf.
//
// Solidity: function decreaseWritership(address oldWriter) returns()
func (_Writeable *WriteableTransactor) DecreaseWritership(opts *bind.TransactOpts, oldWriter common.Address) (*types.Transaction, error) {
	return _Writeable.contract.Transact(opts, "decreaseWritership", oldWriter)
}

// DecreaseWritership is a paid mutator transaction binding the contract method 0x9f217bdf.
//
// Solidity: function decreaseWritership(address oldWriter) returns()
func (_Writeable *WriteableSession) DecreaseWritership(oldWriter common.Address) (*types.Transaction, error) {
	return _Writeable.Contract.DecreaseWritership(&_Writeable.TransactOpts, oldWriter)
}

// DecreaseWritership is a paid mutator transaction binding the contract method 0x9f217bdf.
//
// Solidity: function decreaseWritership(address oldWriter) returns()
func (_Writeable *WriteableTransactorSession) DecreaseWritership(oldWriter common.Address) (*types.Transaction, error) {
	return _Writeable.Contract.DecreaseWritership(&_Writeable.TransactOpts, oldWriter)
}

// IncreasedWritership is a paid mutator transaction binding the contract method 0x14b4d411.
//
// Solidity: function increasedWritership(address newWriter) returns()
func (_Writeable *WriteableTransactor) IncreasedWritership(opts *bind.TransactOpts, newWriter common.Address) (*types.Transaction, error) {
	return _Writeable.contract.Transact(opts, "increasedWritership", newWriter)
}

// IncreasedWritership is a paid mutator transaction binding the contract method 0x14b4d411.
//
// Solidity: function increasedWritership(address newWriter) returns()
func (_Writeable *WriteableSession) IncreasedWritership(newWriter common.Address) (*types.Transaction, error) {
	return _Writeable.Contract.IncreasedWritership(&_Writeable.TransactOpts, newWriter)
}

// IncreasedWritership is a paid mutator transaction binding the contract method 0x14b4d411.
//
// Solidity: function increasedWritership(address newWriter) returns()
func (_Writeable *WriteableTransactorSession) IncreasedWritership(newWriter common.Address) (*types.Transaction, error) {
	return _Writeable.Contract.IncreasedWritership(&_Writeable.TransactOpts, newWriter)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Writeable *WriteableTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Writeable.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Writeable *WriteableSession) RenounceOwnership() (*types.Transaction, error) {
	return _Writeable.Contract.RenounceOwnership(&_Writeable.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Writeable *WriteableTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Writeable.Contract.RenounceOwnership(&_Writeable.TransactOpts)
}

// RenounceWritership is a paid mutator transaction binding the contract method 0x4075e83e.
//
// Solidity: function renounceWritership() returns()
func (_Writeable *WriteableTransactor) RenounceWritership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Writeable.contract.Transact(opts, "renounceWritership")
}

// RenounceWritership is a paid mutator transaction binding the contract method 0x4075e83e.
//
// Solidity: function renounceWritership() returns()
func (_Writeable *WriteableSession) RenounceWritership() (*types.Transaction, error) {
	return _Writeable.Contract.RenounceWritership(&_Writeable.TransactOpts)
}

// RenounceWritership is a paid mutator transaction binding the contract method 0x4075e83e.
//
// Solidity: function renounceWritership() returns()
func (_Writeable *WriteableTransactorSession) RenounceWritership() (*types.Transaction, error) {
	return _Writeable.Contract.RenounceWritership(&_Writeable.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Writeable *WriteableTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Writeable.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Writeable *WriteableSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Writeable.Contract.TransferOwnership(&_Writeable.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Writeable *WriteableTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Writeable.Contract.TransferOwnership(&_Writeable.TransactOpts, newOwner)
}

// WriteableOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Writeable contract.
type WriteableOwnershipTransferredIterator struct {
	Event *WriteableOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *WriteableOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WriteableOwnershipTransferred)
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
		it.Event = new(WriteableOwnershipTransferred)
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
func (it *WriteableOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WriteableOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WriteableOwnershipTransferred represents a OwnershipTransferred event raised by the Writeable contract.
type WriteableOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Writeable *WriteableFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*WriteableOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Writeable.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &WriteableOwnershipTransferredIterator{contract: _Writeable.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Writeable *WriteableFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *WriteableOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Writeable.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WriteableOwnershipTransferred)
				if err := _Writeable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Writeable *WriteableFilterer) ParseOwnershipTransferred(log types.Log) (*WriteableOwnershipTransferred, error) {
	event := new(WriteableOwnershipTransferred)
	if err := _Writeable.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WriteableWritershipDecreasedIterator is returned from FilterWritershipDecreased and is used to iterate over the raw logs and unpacked data for WritershipDecreased events raised by the Writeable contract.
type WriteableWritershipDecreasedIterator struct {
	Event *WriteableWritershipDecreased // Event containing the contract specifics and raw log

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
func (it *WriteableWritershipDecreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WriteableWritershipDecreased)
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
		it.Event = new(WriteableWritershipDecreased)
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
func (it *WriteableWritershipDecreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WriteableWritershipDecreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WriteableWritershipDecreased represents a WritershipDecreased event raised by the Writeable contract.
type WriteableWritershipDecreased struct {
	OldWriter common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWritershipDecreased is a free log retrieval operation binding the contract event 0xb24c36b70cc67df9ea22f053c15c4a41e7170d7be0637d2637a3445a266ac6af.
//
// Solidity: event WritershipDecreased(address indexed oldWriter)
func (_Writeable *WriteableFilterer) FilterWritershipDecreased(opts *bind.FilterOpts, oldWriter []common.Address) (*WriteableWritershipDecreasedIterator, error) {

	var oldWriterRule []interface{}
	for _, oldWriterItem := range oldWriter {
		oldWriterRule = append(oldWriterRule, oldWriterItem)
	}

	logs, sub, err := _Writeable.contract.FilterLogs(opts, "WritershipDecreased", oldWriterRule)
	if err != nil {
		return nil, err
	}
	return &WriteableWritershipDecreasedIterator{contract: _Writeable.contract, event: "WritershipDecreased", logs: logs, sub: sub}, nil
}

// WatchWritershipDecreased is a free log subscription operation binding the contract event 0xb24c36b70cc67df9ea22f053c15c4a41e7170d7be0637d2637a3445a266ac6af.
//
// Solidity: event WritershipDecreased(address indexed oldWriter)
func (_Writeable *WriteableFilterer) WatchWritershipDecreased(opts *bind.WatchOpts, sink chan<- *WriteableWritershipDecreased, oldWriter []common.Address) (event.Subscription, error) {

	var oldWriterRule []interface{}
	for _, oldWriterItem := range oldWriter {
		oldWriterRule = append(oldWriterRule, oldWriterItem)
	}

	logs, sub, err := _Writeable.contract.WatchLogs(opts, "WritershipDecreased", oldWriterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WriteableWritershipDecreased)
				if err := _Writeable.contract.UnpackLog(event, "WritershipDecreased", log); err != nil {
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

// ParseWritershipDecreased is a log parse operation binding the contract event 0xb24c36b70cc67df9ea22f053c15c4a41e7170d7be0637d2637a3445a266ac6af.
//
// Solidity: event WritershipDecreased(address indexed oldWriter)
func (_Writeable *WriteableFilterer) ParseWritershipDecreased(log types.Log) (*WriteableWritershipDecreased, error) {
	event := new(WriteableWritershipDecreased)
	if err := _Writeable.contract.UnpackLog(event, "WritershipDecreased", log); err != nil {
		return nil, err
	}
	return event, nil
}

// WriteableWritershipIncreasedIterator is returned from FilterWritershipIncreased and is used to iterate over the raw logs and unpacked data for WritershipIncreased events raised by the Writeable contract.
type WriteableWritershipIncreasedIterator struct {
	Event *WriteableWritershipIncreased // Event containing the contract specifics and raw log

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
func (it *WriteableWritershipIncreasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(WriteableWritershipIncreased)
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
		it.Event = new(WriteableWritershipIncreased)
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
func (it *WriteableWritershipIncreasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *WriteableWritershipIncreasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// WriteableWritershipIncreased represents a WritershipIncreased event raised by the Writeable contract.
type WriteableWritershipIncreased struct {
	NewWriter common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWritershipIncreased is a free log retrieval operation binding the contract event 0x8f2adc2e2f2723d375f8bb1359cb1027f713945395ca36c42b83dd94c5bdc7f2.
//
// Solidity: event WritershipIncreased(address indexed newWriter)
func (_Writeable *WriteableFilterer) FilterWritershipIncreased(opts *bind.FilterOpts, newWriter []common.Address) (*WriteableWritershipIncreasedIterator, error) {

	var newWriterRule []interface{}
	for _, newWriterItem := range newWriter {
		newWriterRule = append(newWriterRule, newWriterItem)
	}

	logs, sub, err := _Writeable.contract.FilterLogs(opts, "WritershipIncreased", newWriterRule)
	if err != nil {
		return nil, err
	}
	return &WriteableWritershipIncreasedIterator{contract: _Writeable.contract, event: "WritershipIncreased", logs: logs, sub: sub}, nil
}

// WatchWritershipIncreased is a free log subscription operation binding the contract event 0x8f2adc2e2f2723d375f8bb1359cb1027f713945395ca36c42b83dd94c5bdc7f2.
//
// Solidity: event WritershipIncreased(address indexed newWriter)
func (_Writeable *WriteableFilterer) WatchWritershipIncreased(opts *bind.WatchOpts, sink chan<- *WriteableWritershipIncreased, newWriter []common.Address) (event.Subscription, error) {

	var newWriterRule []interface{}
	for _, newWriterItem := range newWriter {
		newWriterRule = append(newWriterRule, newWriterItem)
	}

	logs, sub, err := _Writeable.contract.WatchLogs(opts, "WritershipIncreased", newWriterRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(WriteableWritershipIncreased)
				if err := _Writeable.contract.UnpackLog(event, "WritershipIncreased", log); err != nil {
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

// ParseWritershipIncreased is a log parse operation binding the contract event 0x8f2adc2e2f2723d375f8bb1359cb1027f713945395ca36c42b83dd94c5bdc7f2.
//
// Solidity: event WritershipIncreased(address indexed newWriter)
func (_Writeable *WriteableFilterer) ParseWritershipIncreased(log types.Log) (*WriteableWritershipIncreased, error) {
	event := new(WriteableWritershipIncreased)
	if err := _Writeable.contract.UnpackLog(event, "WritershipIncreased", log); err != nil {
		return nil, err
	}
	return event, nil
}
