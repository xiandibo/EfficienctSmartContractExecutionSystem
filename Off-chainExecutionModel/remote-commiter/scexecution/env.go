package scexecution

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/vadiminshakov/committer/core/state"
	"github.com/vadiminshakov/committer/core/vm"
	"math/big"
)

// Config is a basic type specifying certain configuration flags for running
// the EVM.
type Config struct {
	ChainConfig *params.ChainConfig
	Difficulty  *big.Int
	Origin      common.Address
	Coinbase    common.Address
	BlockNumber *big.Int
	Time        *big.Int
	GasLimit    uint64
	GasPrice    *big.Int
	Value       *big.Int
	Debug       bool
	EVMConfig   vm.Config

	State     *state.StateDB
	StateCopy *state.StateDB //用于防止并发错误，这个StateCopy只负责装载一个合约
	GetHashFn func(n uint64) common.Hash
}

func NewEnv(cfg *Config) *vm.EVM {
	context := vm.Context{
		CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
		Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
		GetHash:     cfg.GetHashFn,
		Origin:      cfg.Origin,
		Coinbase:    cfg.Coinbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		Difficulty:  cfg.Difficulty,
		GasLimit:    cfg.GasLimit,
		GasPrice:    cfg.GasPrice,
	}

	//设定我们想要的jumptable
	//temp := vm.NewEVM(context, cfg.State, cfg.ChainConfig, cfg.EVMConfig)
	//temp.ChangeChainRule()

	//此处不使用cfg.state, 使用其copy
	return vm.NewEVM(context, cfg.State, cfg.ChainConfig, cfg.EVMConfig)
	//return temp
}

// 输入多了一个statecopy，它只装载了一个智能合约
func NewEnv1(cfg *Config, statecopy *state.StateDB) *vm.EVM {
	context := vm.Context{
		CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
		Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
		GetHash:     cfg.GetHashFn,
		Origin:      cfg.Origin,
		Coinbase:    cfg.Coinbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		Difficulty:  cfg.Difficulty,
		GasLimit:    cfg.GasLimit,
		GasPrice:    cfg.GasPrice,
	}

	//设定我们想要的jumptable
	//temp := vm.NewEVM(context, cfg.State, cfg.ChainConfig, cfg.EVMConfig)
	//temp.ChangeChainRule()

	//此处不使用cfg.state, 使用其copy
	return vm.NewEVM(context, statecopy, cfg.ChainConfig, cfg.EVMConfig)
	//return temp
}