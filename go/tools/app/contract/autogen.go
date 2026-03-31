// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// InstructionSenderMetaData contains all meta data concerning the InstructionSender contract.
var InstructionSenderMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_teeExtensionRegistry\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_teeMachineRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"_extensionId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setExtensionId\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"sign\",\"inputs\":[{\"name\":\"_message\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"teeExtensionRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractITeeExtensionRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"teeMachineRegistry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractITeeMachineRegistry\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"updateKey\",\"inputs\":[{\"name\":\"_encryptedKey\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"}]",
	Bin: "0x60c060405234801561000f575f5ffd5b50604051611334380380611334833981810160405281019061003191906100fe565b8173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff1681525050505061013c565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6100cd826100a4565b9050919050565b6100dd816100c3565b81146100e7575f5ffd5b50565b5f815190506100f8816100d4565b92915050565b5f5f60408385031215610114576101136100a0565b5b5f610121858286016100ea565b9250506020610132858286016100ea565b9150509250929050565b60805160a0516111ad6101875f395f818161014e015281816101b7015261061c01525f81816102ff015281816103a60152818161040e015281816104c1015261076401526111ad5ff3fe608060405260043610610054575f3560e01c8063524967d71461005857806376cd7cbc14610082578063a435d58a146100b2578063aa5032c6146100dc578063d473e270146100f2578063e6eb68671461011c575b5f5ffd5b348015610063575f5ffd5b5061006c61014c565b60405161007991906108d5565b60405180910390f35b61009c60048036038101906100979190610960565b610170565b6040516100a991906109c3565b60405180910390f35b3480156100bd575f5ffd5b506100c66103a4565b6040516100d391906109fc565b60405180910390f35b3480156100e7575f5ffd5b506100f06103c8565b005b3480156100fd575f5ffd5b506101066105d0565b6040516101139190610a2d565b60405180910390f35b61013660048036038101906101319190610960565b6105d5565b60405161014391906109c3565b60405180910390f35b7f000000000000000000000000000000000000000000000000000000000000000081565b5f5f5f54036101b4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101ab90610aa0565b60405180910390fd5b5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feeabcbf5f5460016040518363ffffffff1660e01b8152600401610212929190610af7565b5f60405180830381865afa15801561022c573d5f5f3e3d5ffd5b505050506040513d5f823e3d601f19601f820116820180604052508101906102549190610ca1565b905061025e610809565b7f4b45590000000000000000000000000000000000000000000000000000000000815f0181815250507f5349474e0000000000000000000000000000000000000000000000000000000081602001818152505084848080601f0160208091040260200160405190810160405280939291908181526020018383808284375f81840152601f19601f8201169050808301925050505050505081604001819052507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663f731df533484846040518463ffffffff1660e01b8152600401610359929190610f29565b60206040518083038185885af1158015610375573d5f5f3e3d5ffd5b50505050506040513d601f19601f8201168201806040525081019061039a9190610f88565b9250505092915050565b7f000000000000000000000000000000000000000000000000000000000000000081565b5f5f541461040b576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161040290610ffd565b60405180910390fd5b5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663fad5902b6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610475573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906104999190611045565b90505f600190505b818111610592573073ffffffffffffffffffffffffffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632c177358836040518263ffffffff1660e01b81526004016105189190610a2d565b602060405180830381865afa158015610533573d5f5f3e3d5ffd5b505050506040513d601f19601f820116820180604052508101906105579190611070565b73ffffffffffffffffffffffffffffffffffffffff160361057f57805f8190555050506105ce565b808061058a906110c8565b9150506104a1565b506040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105c590611159565b60405180910390fd5b565b5f5481565b5f5f5f5403610619576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161061090610aa0565b60405180910390fd5b5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663feeabcbf5f5460016040518363ffffffff1660e01b8152600401610677929190610af7565b5f60405180830381865afa158015610691573d5f5f3e3d5ffd5b505050506040513d5f823e3d601f19601f820116820180604052508101906106b99190610ca1565b90506106c3610809565b7f4b45590000000000000000000000000000000000000000000000000000000000815f0181815250507f555044415445000000000000000000000000000000000000000000000000000081602001818152505084848080601f0160208091040260200160405190810160405280939291908181526020018383808284375f81840152601f19601f8201169050808301925050505050505081604001819052507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663f731df533484846040518463ffffffff1660e01b81526004016107be929190610f29565b60206040518083038185885af11580156107da573d5f5f3e3d5ffd5b50505050506040513d601f19601f820116820180604052508101906107ff9190610f88565b9250505092915050565b6040518060c001604052805f81526020015f815260200160608152602001606081526020015f67ffffffffffffffff1681526020015f73ffffffffffffffffffffffffffffffffffffffff1681525090565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f819050919050565b5f61089d6108986108938461085b565b61087a565b61085b565b9050919050565b5f6108ae82610883565b9050919050565b5f6108bf826108a4565b9050919050565b6108cf816108b5565b82525050565b5f6020820190506108e85f8301846108c6565b92915050565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f8401126109205761091f6108ff565b5b8235905067ffffffffffffffff81111561093d5761093c610903565b5b60208301915083600182028301111561095957610958610907565b5b9250929050565b5f5f60208385031215610976576109756108f7565b5b5f83013567ffffffffffffffff811115610993576109926108fb565b5b61099f8582860161090b565b92509250509250929050565b5f819050919050565b6109bd816109ab565b82525050565b5f6020820190506109d65f8301846109b4565b92915050565b5f6109e6826108a4565b9050919050565b6109f6816109dc565b82525050565b5f602082019050610a0f5f8301846109ed565b92915050565b5f819050919050565b610a2781610a15565b82525050565b5f602082019050610a405f830184610a1e565b92915050565b5f82825260208201905092915050565b7f657874656e73696f6e204944206e6f74207365740000000000000000000000005f82015250565b5f610a8a601483610a46565b9150610a9582610a56565b602082019050919050565b5f6020820190508181035f830152610ab781610a7e565b9050919050565b5f819050919050565b5f610ae1610adc610ad784610abe565b61087a565b610a15565b9050919050565b610af181610ac7565b82525050565b5f604082019050610b0a5f830185610a1e565b610b176020830184610ae8565b9392505050565b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b610b6482610b1e565b810181811067ffffffffffffffff82111715610b8357610b82610b2e565b5b80604052505050565b5f610b956108ee565b9050610ba18282610b5b565b919050565b5f67ffffffffffffffff821115610bc057610bbf610b2e565b5b602082029050602081019050919050565b5f610bdb8261085b565b9050919050565b610beb81610bd1565b8114610bf5575f5ffd5b50565b5f81519050610c0681610be2565b92915050565b5f610c1e610c1984610ba6565b610b8c565b90508083825260208201905060208402830185811115610c4157610c40610907565b5b835b81811015610c6a5780610c568882610bf8565b845260208401935050602081019050610c43565b5050509392505050565b5f82601f830112610c8857610c876108ff565b5b8151610c98848260208601610c0c565b91505092915050565b5f60208284031215610cb657610cb56108f7565b5b5f82015167ffffffffffffffff811115610cd357610cd26108fb565b5b610cdf84828501610c74565b91505092915050565b5f81519050919050565b5f82825260208201905092915050565b5f819050602082019050919050565b610d1a81610bd1565b82525050565b5f610d2b8383610d11565b60208301905092915050565b5f602082019050919050565b5f610d4d82610ce8565b610d578185610cf2565b9350610d6283610d02565b805f5b83811015610d92578151610d798882610d20565b9750610d8483610d37565b925050600181019050610d65565b5085935050505092915050565b610da8816109ab565b82525050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f610de082610dae565b610dea8185610db8565b9350610dfa818560208601610dc8565b610e0381610b1e565b840191505092915050565b5f82825260208201905092915050565b5f610e2882610ce8565b610e328185610e0e565b9350610e3d83610d02565b805f5b83811015610e6d578151610e548882610d20565b9750610e5f83610d37565b925050600181019050610e40565b5085935050505092915050565b5f67ffffffffffffffff82169050919050565b610e9681610e7a565b82525050565b5f60c083015f830151610eb15f860182610d9f565b506020830151610ec46020860182610d9f565b5060408301518482036040860152610edc8282610dd6565b91505060608301518482036060860152610ef68282610e1e565b9150506080830151610f0b6080860182610e8d565b5060a0830151610f1e60a0860182610d11565b508091505092915050565b5f6040820190508181035f830152610f418185610d43565b90508181036020830152610f558184610e9c565b90509392505050565b610f67816109ab565b8114610f71575f5ffd5b50565b5f81519050610f8281610f5e565b92915050565b5f60208284031215610f9d57610f9c6108f7565b5b5f610faa84828501610f74565b91505092915050565b7f657874656e73696f6e20494420616c72656164792073657400000000000000005f82015250565b5f610fe7601883610a46565b9150610ff282610fb3565b602082019050919050565b5f6020820190508181035f83015261101481610fdb565b9050919050565b61102481610a15565b811461102e575f5ffd5b50565b5f8151905061103f8161101b565b92915050565b5f6020828403121561105a576110596108f7565b5b5f61106784828501611031565b91505092915050565b5f60208284031215611085576110846108f7565b5b5f61109284828501610bf8565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6110d282610a15565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036111045761110361109b565b5b600182019050919050565b7f657874656e73696f6e204944206e6f7420666f756e64000000000000000000005f82015250565b5f611143601683610a46565b915061114e8261110f565b602082019050919050565b5f6020820190508181035f83015261117081611137565b905091905056fea26469706673582212203ae6831d3c4963ce0b42e34445a4d0f7196a73b623a50841f09b6c16627d495564736f6c634300081f0033",
}

// InstructionSenderABI is the input ABI used to generate the binding from.
// Deprecated: Use InstructionSenderMetaData.ABI instead.
var InstructionSenderABI = InstructionSenderMetaData.ABI

// InstructionSenderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use InstructionSenderMetaData.Bin instead.
var InstructionSenderBin = InstructionSenderMetaData.Bin

// DeployInstructionSender deploys a new Ethereum contract, binding an instance of InstructionSender to it.
func DeployInstructionSender(auth *bind.TransactOpts, backend bind.ContractBackend, _teeExtensionRegistry common.Address, _teeMachineRegistry common.Address) (common.Address, *types.Transaction, *InstructionSender, error) {
	parsed, err := InstructionSenderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(InstructionSenderBin), backend, _teeExtensionRegistry, _teeMachineRegistry)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &InstructionSender{InstructionSenderCaller: InstructionSenderCaller{contract: contract}, InstructionSenderTransactor: InstructionSenderTransactor{contract: contract}, InstructionSenderFilterer: InstructionSenderFilterer{contract: contract}}, nil
}

// InstructionSender is an auto generated Go binding around an Ethereum contract.
type InstructionSender struct {
	InstructionSenderCaller     // Read-only binding to the contract
	InstructionSenderTransactor // Write-only binding to the contract
	InstructionSenderFilterer   // Log filterer for contract events
}

// InstructionSenderCaller is an auto generated read-only Go binding around an Ethereum contract.
type InstructionSenderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InstructionSenderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type InstructionSenderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InstructionSenderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type InstructionSenderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// InstructionSenderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type InstructionSenderSession struct {
	Contract     *InstructionSender // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// InstructionSenderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type InstructionSenderCallerSession struct {
	Contract *InstructionSenderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// InstructionSenderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type InstructionSenderTransactorSession struct {
	Contract     *InstructionSenderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// InstructionSenderRaw is an auto generated low-level Go binding around an Ethereum contract.
type InstructionSenderRaw struct {
	Contract *InstructionSender // Generic contract binding to access the raw methods on
}

// InstructionSenderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type InstructionSenderCallerRaw struct {
	Contract *InstructionSenderCaller // Generic read-only contract binding to access the raw methods on
}

// InstructionSenderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type InstructionSenderTransactorRaw struct {
	Contract *InstructionSenderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewInstructionSender creates a new instance of InstructionSender, bound to a specific deployed contract.
func NewInstructionSender(address common.Address, backend bind.ContractBackend) (*InstructionSender, error) {
	contract, err := bindInstructionSender(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &InstructionSender{InstructionSenderCaller: InstructionSenderCaller{contract: contract}, InstructionSenderTransactor: InstructionSenderTransactor{contract: contract}, InstructionSenderFilterer: InstructionSenderFilterer{contract: contract}}, nil
}

// NewInstructionSenderCaller creates a new read-only instance of InstructionSender, bound to a specific deployed contract.
func NewInstructionSenderCaller(address common.Address, caller bind.ContractCaller) (*InstructionSenderCaller, error) {
	contract, err := bindInstructionSender(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &InstructionSenderCaller{contract: contract}, nil
}

// NewInstructionSenderTransactor creates a new write-only instance of InstructionSender, bound to a specific deployed contract.
func NewInstructionSenderTransactor(address common.Address, transactor bind.ContractTransactor) (*InstructionSenderTransactor, error) {
	contract, err := bindInstructionSender(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &InstructionSenderTransactor{contract: contract}, nil
}

// NewInstructionSenderFilterer creates a new log filterer instance of InstructionSender, bound to a specific deployed contract.
func NewInstructionSenderFilterer(address common.Address, filterer bind.ContractFilterer) (*InstructionSenderFilterer, error) {
	contract, err := bindInstructionSender(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &InstructionSenderFilterer{contract: contract}, nil
}

// bindInstructionSender binds a generic wrapper to an already deployed contract.
func bindInstructionSender(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := InstructionSenderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InstructionSender *InstructionSenderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InstructionSender.Contract.InstructionSenderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InstructionSender *InstructionSenderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InstructionSender.Contract.InstructionSenderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InstructionSender *InstructionSenderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InstructionSender.Contract.InstructionSenderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_InstructionSender *InstructionSenderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _InstructionSender.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_InstructionSender *InstructionSenderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InstructionSender.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_InstructionSender *InstructionSenderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _InstructionSender.Contract.contract.Transact(opts, method, params...)
}

// ExtensionId is a free data retrieval call binding the contract method 0xd473e270.
//
// Solidity: function _extensionId() view returns(uint256)
func (_InstructionSender *InstructionSenderCaller) ExtensionId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _InstructionSender.contract.Call(opts, &out, "_extensionId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExtensionId is a free data retrieval call binding the contract method 0xd473e270.
//
// Solidity: function _extensionId() view returns(uint256)
func (_InstructionSender *InstructionSenderSession) ExtensionId() (*big.Int, error) {
	return _InstructionSender.Contract.ExtensionId(&_InstructionSender.CallOpts)
}

// ExtensionId is a free data retrieval call binding the contract method 0xd473e270.
//
// Solidity: function _extensionId() view returns(uint256)
func (_InstructionSender *InstructionSenderCallerSession) ExtensionId() (*big.Int, error) {
	return _InstructionSender.Contract.ExtensionId(&_InstructionSender.CallOpts)
}

// TeeExtensionRegistry is a free data retrieval call binding the contract method 0xa435d58a.
//
// Solidity: function teeExtensionRegistry() view returns(address)
func (_InstructionSender *InstructionSenderCaller) TeeExtensionRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _InstructionSender.contract.Call(opts, &out, "teeExtensionRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TeeExtensionRegistry is a free data retrieval call binding the contract method 0xa435d58a.
//
// Solidity: function teeExtensionRegistry() view returns(address)
func (_InstructionSender *InstructionSenderSession) TeeExtensionRegistry() (common.Address, error) {
	return _InstructionSender.Contract.TeeExtensionRegistry(&_InstructionSender.CallOpts)
}

// TeeExtensionRegistry is a free data retrieval call binding the contract method 0xa435d58a.
//
// Solidity: function teeExtensionRegistry() view returns(address)
func (_InstructionSender *InstructionSenderCallerSession) TeeExtensionRegistry() (common.Address, error) {
	return _InstructionSender.Contract.TeeExtensionRegistry(&_InstructionSender.CallOpts)
}

// TeeMachineRegistry is a free data retrieval call binding the contract method 0x524967d7.
//
// Solidity: function teeMachineRegistry() view returns(address)
func (_InstructionSender *InstructionSenderCaller) TeeMachineRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _InstructionSender.contract.Call(opts, &out, "teeMachineRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TeeMachineRegistry is a free data retrieval call binding the contract method 0x524967d7.
//
// Solidity: function teeMachineRegistry() view returns(address)
func (_InstructionSender *InstructionSenderSession) TeeMachineRegistry() (common.Address, error) {
	return _InstructionSender.Contract.TeeMachineRegistry(&_InstructionSender.CallOpts)
}

// TeeMachineRegistry is a free data retrieval call binding the contract method 0x524967d7.
//
// Solidity: function teeMachineRegistry() view returns(address)
func (_InstructionSender *InstructionSenderCallerSession) TeeMachineRegistry() (common.Address, error) {
	return _InstructionSender.Contract.TeeMachineRegistry(&_InstructionSender.CallOpts)
}

// SetExtensionId is a paid mutator transaction binding the contract method 0xaa5032c6.
//
// Solidity: function setExtensionId() returns()
func (_InstructionSender *InstructionSenderTransactor) SetExtensionId(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _InstructionSender.contract.Transact(opts, "setExtensionId")
}

// SetExtensionId is a paid mutator transaction binding the contract method 0xaa5032c6.
//
// Solidity: function setExtensionId() returns()
func (_InstructionSender *InstructionSenderSession) SetExtensionId() (*types.Transaction, error) {
	return _InstructionSender.Contract.SetExtensionId(&_InstructionSender.TransactOpts)
}

// SetExtensionId is a paid mutator transaction binding the contract method 0xaa5032c6.
//
// Solidity: function setExtensionId() returns()
func (_InstructionSender *InstructionSenderTransactorSession) SetExtensionId() (*types.Transaction, error) {
	return _InstructionSender.Contract.SetExtensionId(&_InstructionSender.TransactOpts)
}

// Sign is a paid mutator transaction binding the contract method 0x76cd7cbc.
//
// Solidity: function sign(bytes _message) payable returns(bytes32)
func (_InstructionSender *InstructionSenderTransactor) Sign(opts *bind.TransactOpts, _message []byte) (*types.Transaction, error) {
	return _InstructionSender.contract.Transact(opts, "sign", _message)
}

// Sign is a paid mutator transaction binding the contract method 0x76cd7cbc.
//
// Solidity: function sign(bytes _message) payable returns(bytes32)
func (_InstructionSender *InstructionSenderSession) Sign(_message []byte) (*types.Transaction, error) {
	return _InstructionSender.Contract.Sign(&_InstructionSender.TransactOpts, _message)
}

// Sign is a paid mutator transaction binding the contract method 0x76cd7cbc.
//
// Solidity: function sign(bytes _message) payable returns(bytes32)
func (_InstructionSender *InstructionSenderTransactorSession) Sign(_message []byte) (*types.Transaction, error) {
	return _InstructionSender.Contract.Sign(&_InstructionSender.TransactOpts, _message)
}

// UpdateKey is a paid mutator transaction binding the contract method 0xe6eb6867.
//
// Solidity: function updateKey(bytes _encryptedKey) payable returns(bytes32)
func (_InstructionSender *InstructionSenderTransactor) UpdateKey(opts *bind.TransactOpts, _encryptedKey []byte) (*types.Transaction, error) {
	return _InstructionSender.contract.Transact(opts, "updateKey", _encryptedKey)
}

// UpdateKey is a paid mutator transaction binding the contract method 0xe6eb6867.
//
// Solidity: function updateKey(bytes _encryptedKey) payable returns(bytes32)
func (_InstructionSender *InstructionSenderSession) UpdateKey(_encryptedKey []byte) (*types.Transaction, error) {
	return _InstructionSender.Contract.UpdateKey(&_InstructionSender.TransactOpts, _encryptedKey)
}

// UpdateKey is a paid mutator transaction binding the contract method 0xe6eb6867.
//
// Solidity: function updateKey(bytes _encryptedKey) payable returns(bytes32)
func (_InstructionSender *InstructionSenderTransactorSession) UpdateKey(_encryptedKey []byte) (*types.Transaction, error) {
	return _InstructionSender.Contract.UpdateKey(&_InstructionSender.TransactOpts, _encryptedKey)
}
