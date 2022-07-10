package scexecution

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	"strings"
)


func ExecuteSmartContract() []byte {
	//var definition = `[{"constant":true,"inputs":[],"name":"seller","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[],"name":"abort","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"value","outputs":[{"name":"","type":"uint256"}],"type":"function"},{"constant":false,"inputs":[],"name":"refund","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"buyer","outputs":[{"name":"","type":"address"}],"type":"function"},{"constant":false,"inputs":[],"name":"confirmReceived","outputs":[],"type":"function"},{"constant":true,"inputs":[],"name":"state","outputs":[{"name":"","type":"uint8"}],"type":"function"},{"constant":false,"inputs":[],"name":"confirmPurchase","outputs":[],"type":"function"},{"inputs":[],"type":"constructor"},{"anonymous":false,"inputs":[],"name":"Aborted","type":"event"},{"anonymous":false,"inputs":[],"name":"PurchaseConfirmed","type":"event"},{"anonymous":false,"inputs":[],"name":"ItemReceived","type":"event"},{"anonymous":false,"inputs":[],"name":"Refunded","type":"event"}]`

	var definition = `[
	{
		"constant": false,
		"inputs": [],
		"name": "testwithoutinput",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "key",
				"type": "string"
			}
		],
		"name": "get",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "key",
				"type": "string"
			},
			{
				"name": "value",
				"type": "string"
			}
		],
		"name": "set",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	}
]`

	var code = common.Hex2Bytes("608060405260043610610051576000357c010000000000000000000000000000000000000000000000000000000090048063566261f914610056578063693ec85e146100e6578063e942b51614610227575b600080fd5b34801561006257600080fd5b5061006b610386565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100ab578082015181840152602081019050610090565b50505050905090810190601f1680156100d85780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156100f257600080fd5b506101ac6004803603602081101561010957600080fd5b810190808035906020019064010000000081111561012657600080fd5b82018360208201111561013857600080fd5b8035906020019184600183028401116401000000008311171561015a57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506103c8565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101ec5780820151818401526020810190506101d1565b50505050905090810190601f1680156102195780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561023357600080fd5b506103846004803603604081101561024a57600080fd5b810190808035906020019064010000000081111561026757600080fd5b82018360208201111561027957600080fd5b8035906020019184600183028401116401000000008311171561029b57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290803590602001906401000000008111156102fe57600080fd5b82018360208201111561031057600080fd5b8035906020019184600183028401116401000000008311171561033257600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506104d5565b005b6060806040805190810160405280600781526020017f616263646566670000000000000000000000000000000000000000000000000081525090508091505090565b60606000826040518082805190602001908083835b60208310151561040257805182526020820191506020810190506020830392506103dd565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104c95780601f1061049e576101008083540402835291602001916104c9565b820191906000526020600020905b8154815290600101906020018083116104ac57829003601f168201915b50505050509050919050565b806000836040518082805190602001908083835b60208310151561050e57805182526020820191506020810190506020830392506104e9565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390209080519060200190610554929190610559565b505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061059a57805160ff19168380011785556105c8565b828001600101855582156105c8579182015b828111156105c75782518255916020019190600101906105ac565b5b5090506105d591906105d9565b5090565b6105fb91905b808211156105f75760008160009055506001016105df565b5090565b9056fea165627a7a7230582077bcca77f75f59c126b2500f24b623fd031c4b546dacbeb63443e3ca68fba9f30029")

	abi, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return nil
	}

	test, err := abi.Pack("testwithoutinput")
	if err != nil {
		return nil
	}


	res, _, _ := Execute(code, test, nil)
	return res
}

/*func ExecuteTx(tx *types.Transaction) {

	var (
		signer     = types.HomesteadSigner{}
		//testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		db         = rawdb.NewMemoryDatabase() //新建测试用的内存数据库
		gspec      = &core.Genesis{
			Config: params.TestChainConfig,
		}
		//genesis       = gspec.MustCommit(db)
		//新建区块链
		blockchain, _ = core.NewBlockChain(db, nil, gspec.Config, ethash.NewFaker(), vm.Config{}, nil, nil)
	)
	defer blockchain.Stop()

	//暂时创建一条随机链，可以往里面插入块
	//先处理交易再生产块
	//我们不一定需要生产块

	msg, _ := tx.AsMessage(signer)


	//types.NewMessage(*tx.To(), tx.To(), tx.Nonce(), tx.Value(), tx.Gas(), new(big.Int).Set(tx.GasPrice()), tx.Data(), tx.AccessList(), true)

	//初始化db
	rawdb.NewMemoryDatabase() //用于模拟仿真


	parent := blockchain.CurrentBlock()
	timestamp := time.Now().Unix()
	if parent.Time() >= uint64(timestamp) {
		timestamp = int64(parent.Time() + 1)
	}
	num := parent.Number()

	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		GasLimit:   9999999999,
		//Extra:      w.extra,
		Time:       uint64(timestamp),
	}

	statedb, _ := blockchain.StateAt(parent.Root())

	blockContext := core.NewEVMBlockContext(header, blockchain, nil)
	vmenv := vm.NewEVM(blockContext, vm.TxContext{}, statedb, gspec.Config, *blockchain.GetVMConfig())
	gasPool := new(core.GasPool).AddGas(header.GasLimit)
	//usedGas := new(uint64)
	//core.applyTransaction(msg, gspec.Config, blockchain, nil, gasPool, statedb, header, tx, usedGas, vmenv)

	txContext := core.NewEVMTxContext(msg)
	vmenv.Reset(txContext, statedb)

	// Apply the transaction to the current state (included in the env).
	result, _ := core.ApplyMessage(vmenv, msg, gasPool)
	println(result)
	//暂时不持久化
	//statedb.Finalise(true)

}*/



