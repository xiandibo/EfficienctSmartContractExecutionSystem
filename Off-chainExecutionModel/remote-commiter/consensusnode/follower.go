package consensusnode

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/vadiminshakov/committer/scexecution"
)

//follower互相通信的逻辑，不包含组内共识


//硬编码合约地址

const (
	KVSTORE = "kvstore"
	KVTEST  = "kvtest"
)

var (
	contractAddr = map[string][]common.Address{
		KVSTORE: {
			common.BytesToAddress([]byte("kvstore1")),
			common.BytesToAddress([]byte("kvstore2")),
			common.BytesToAddress([]byte("kvstore3")),
			common.BytesToAddress([]byte("kvstore4")),
			common.BytesToAddress([]byte("kvstore5")),
			common.BytesToAddress([]byte("kvstore6")),
			common.BytesToAddress([]byte("kvstore7")),
			common.BytesToAddress([]byte("kvstore8")),
			common.BytesToAddress([]byte("kvstore9")),
			common.BytesToAddress([]byte("kvstore10")),
		},
		KVTEST: {
			common.BytesToAddress([]byte("kvtest1")),
			common.BytesToAddress([]byte("kvtest2")),
			common.BytesToAddress([]byte("kvtest3")),
			common.BytesToAddress([]byte("kvtest4")),
			common.BytesToAddress([]byte("kvtest5")),
			common.BytesToAddress([]byte("kvtest6")),
			common.BytesToAddress([]byte("kvtest7")),
			common.BytesToAddress([]byte("kvtest8")),
			common.BytesToAddress([]byte("kvtest9")),
			common.BytesToAddress([]byte("kvtest10")),
		},
	}
)

// kvstore合约代码
var kvstorecode = common.Hex2Bytes("608060405260043610610051576000357c010000000000000000000000000000000000000000000000000000000090048063693ec85e14610056578063e942b51614610197578063fe913865146102f6575b600080fd5b34801561006257600080fd5b5061011c6004803603602081101561007957600080fd5b810190808035906020019064010000000081111561009657600080fd5b8201836020820111156100a857600080fd5b803590602001918460018302840111640100000000831117156100ca57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610331565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561015c578082015181840152602081019050610141565b50505050905090810190601f1680156101895780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156101a357600080fd5b506102f4600480360360408110156101ba57600080fd5b81019080803590602001906401000000008111156101d757600080fd5b8201836020820111156101e957600080fd5b8035906020019184600183028401116401000000008311171561020b57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192908035906020019064010000000081111561026e57600080fd5b82018360208201111561028057600080fd5b803590602001918460018302840111640100000000831117156102a257600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929050505061043e565b005b34801561030257600080fd5b5061032f6004803603602081101561031957600080fd5b81019080803590602001909291905050506104c2565b005b60606000826040518082805190602001908083835b60208310151561036b5780518252602082019150602081019050602083039250610346565b6001836020036101000a03801982511681845116808217855250505050505090500191505090815260200160405180910390208054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156104325780601f1061040757610100808354040283529160200191610432565b820191906000526020600020905b81548152906001019060200180831161041557829003601f168201915b50505050509050919050565b806000836040518082805190602001908083835b6020831015156104775780518252602082019150602081019050602083039250610452565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090805190602001906104bd9291906106a9565b505050565b6060816040519080825280602002602001820160405280156104f35781602001602082028038833980820191505090505b50905060008090505b815181101561053157808303828281518110151561051657fe5b906020019060200201818152505080806001019150506104fc565b506105428160006001845103610546565b5050565b600082905060008290508082141561055f5750506106a4565b600085600286860381151561057057fe5b05860181518110151561057f57fe5b9060200190602002015190505b8183131515610678575b8086848151811015156105a557fe5b9060200190602002015110156105c2578280600101935050610596565b5b85828151811015156105d157fe5b906020019060200201518110156105f0578180600190039250506105c3565b818313151561067357858281518110151561060757fe5b90602001906020020151868481518110151561061f57fe5b90602001906020020151878581518110151561063757fe5b906020019060200201888581518110151561064e57fe5b9060200190602002018281525082815250505082806001019350508180600190039250505b61058c565b8185121561068c5761068b868684610546565b5b838312156106a05761069f868486610546565b5b5050505b505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106106ea57805160ff1916838001178555610718565b82800160010185558215610718579182015b828111156107175782518255916020019190600101906106fc565b5b5090506107259190610729565b5090565b61074b91905b8082111561074757600081600090555060010161072f565b5090565b9056fea165627a7a72305820c848c0abc3818295e01e3eabe828ee894c361ac32efd7b93123791756ff340ba0029")
// kvtest合约代码
var kvtestcode = common.Hex2Bytes("608060405260043610610046576000357c010000000000000000000000000000000000000000000000000000000090048063da778c8a1461004b578063fe91386514610074575b600080fd5b34801561005757600080fd5b50610072600480360361006d919081019061118e565b61009d565b005b34801561008057600080fd5b5061009b6004803603610096919081019061127d565b610dee565b005b6060600080600a9050876000015191508173ffffffffffffffffffffffffffffffffffffffff1663693ec85e89602001516040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040161010691906112dc565b600060405180830381600087803b15801561012057600080fd5b505af1158015610134573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525061015d919081019061114d565b925061016881610dee565b8173ffffffffffffffffffffffffffffffffffffffff1663e942b51689602001518a606001516040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016101c79291906112fe565b600060405180830381600087803b1580156101e157600080fd5b505af11580156101f5573d6000803e3d6000fd5b50505050876040015191508173ffffffffffffffffffffffffffffffffffffffff1663693ec85e89606001516040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040161025991906112dc565b600060405180830381600087803b15801561027357600080fd5b505af1158015610287573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052506102b0919081019061114d565b92506102bb81610dee565b8173ffffffffffffffffffffffffffffffffffffffff1663e942b51689606001518a602001516040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040161031a9291906112fe565b600060405180830381600087803b15801561033457600080fd5b505af1158015610348573d6000803e3d6000fd5b50505050866000015191508173ffffffffffffffffffffffffffffffffffffffff1663693ec85e88602001516040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016103ac91906112dc565b600060405180830381600087803b1580156103c657600080fd5b505af11580156103da573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250610403919081019061114d565b925061040e81610dee565b8173ffffffffffffffffffffffffffffffffffffffff1663e942b516886020015189606001516040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040161046d9291906112fe565b600060405180830381600087803b15801561048757600080fd5b505af115801561049b573d6000803e3d6000fd5b50505050866040015191508173ffffffffffffffffffffffffffffffffffffffff1663693ec85e88606001516040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016104ff91906112dc565b600060405180830381600087803b15801561051957600080fd5b505af115801561052d573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250610556919081019061114d565b925061056181610dee565b8173ffffffffffffffffffffffffffffffffffffffff1663e942b516886060015189602001516040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016105c09291906112fe565b600060405180830381600087803b1580156105da57600080fd5b505af11580156105ee573d6000803e3d6000fd5b50505050856000015191508173ffffffffffffffffffffffffffffffffffffffff1663693ec85e87602001516040518263ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040161065291906112dc565b600060405180830381600087803b15801561066c57600080fd5b505af1158015610680573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052506106a9919081019061114d565b92506106b481610dee565b8173ffffffffffffffffffffffffffffffffffffffff1663e942b516876020015188606001516040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016107139291906112fe565b600060405180830381600087803b15801561072d57600080fd5b505af1158015610741573d6000803e3d6000fd5b50505050856040015191508173ffffffffffffffffffffffffffffffffffffffff1663693ec85e87606001516040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016107a591906112dc565b600060405180830381600087803b1580156107bf57600080fd5b505af11580156107d3573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052506107fc919081019061114d565b925061080781610dee565b8173ffffffffffffffffffffffffffffffffffffffff1663e942b516876060015188602001516040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016108669291906112fe565b600060405180830381600087803b15801561088057600080fd5b505af1158015610894573d6000803e3d6000fd5b50505050846000015191508173ffffffffffffffffffffffffffffffffffffffff1663693ec85e86602001516040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016108f891906112dc565b600060405180830381600087803b15801561091257600080fd5b505af1158015610926573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525061094f919081019061114d565b925061095a81610dee565b8173ffffffffffffffffffffffffffffffffffffffff1663e942b516866020015187606001516040518363ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016109b99291906112fe565b600060405180830381600087803b1580156109d357600080fd5b505af11580156109e7573d6000803e3d6000fd5b50505050846040015191508173ffffffffffffffffffffffffffffffffffffffff1663693ec85e86606001516040518263ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401610a4b91906112dc565b600060405180830381600087803b158015610a6557600080fd5b505af1158015610a79573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250610aa2919081019061114d565b9250610aad81610dee565b8173ffffffffffffffffffffffffffffffffffffffff1663e942b516866060015187602001516040518363ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401610b0c9291906112fe565b600060405180830381600087803b158015610b2657600080fd5b505af1158015610b3a573d6000803e3d6000fd5b50505050836000015191508173ffffffffffffffffffffffffffffffffffffffff1663693ec85e85602001516040518263ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401610b9e91906112dc565b600060405180830381600087803b158015610bb857600080fd5b505af1158015610bcc573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250610bf5919081019061114d565b9250610c0081610dee565b8173ffffffffffffffffffffffffffffffffffffffff1663e942b516856020015186606001516040518363ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401610c5f9291906112fe565b600060405180830381600087803b158015610c7957600080fd5b505af1158015610c8d573d6000803e3d6000fd5b50505050836040015191508173ffffffffffffffffffffffffffffffffffffffff1663693ec85e85606001516040518263ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401610cf191906112dc565b600060405180830381600087803b158015610d0b57600080fd5b505af1158015610d1f573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250610d48919081019061114d565b9250610d5381610dee565b8173ffffffffffffffffffffffffffffffffffffffff1663e942b516856060015186602001516040518363ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401610db29291906112fe565b600060405180830381600087803b158015610dcc57600080fd5b505af1158015610de0573d6000803e3d6000fd5b505050505050505050505050565b606081604051908082528060200260200182016040528015610e1f5781602001602082028038833980820191505090505b50905060008090505b8151811015610e5d578083038282815181101515610e4257fe5b90602001906020020181815250508080600101915050610e28565b50610e6e8160006001845103610e72565b5050565b6000829050600082905080821415610e8b575050610fd0565b6000856002868603811515610e9c57fe5b058601815181101515610eab57fe5b9060200190602002015190505b8183131515610fa4575b808684815181101515610ed157fe5b906020019060200201511015610eee578280600101935050610ec2565b5b8582815181101515610efd57fe5b90602001906020020151811015610f1c57818060019003925050610eef565b8183131515610f9f578582815181101515610f3357fe5b906020019060200201518684815181101515610f4b57fe5b906020019060200201518785815181101515610f6357fe5b9060200190602002018885815181101515610f7a57fe5b9060200190602002018281525082815250505082806001019350508180600190039250505b610eb8565b81851215610fb857610fb7868684610e72565b5b83831215610fcc57610fcb868486610e72565b5b5050505b505050565b6000610fe182356113e5565b905092915050565b600082601f8301121515610ffc57600080fd5b813561100f61100a82611362565b611335565b9150808252602083016020830185838301111561102b57600080fd5b611036838284611401565b50505092915050565b600082601f830112151561105257600080fd5b81516110656110608261138e565b611335565b9150808252602083016020830185838301111561108157600080fd5b61108c838284611410565b50505092915050565b6000608082840312156110a757600080fd5b6110b16080611335565b905060006110c184828501610fd5565b600083015250602082013567ffffffffffffffff8111156110e157600080fd5b6110ed84828501610fe9565b602083015250604061110184828501610fd5565b604083015250606082013567ffffffffffffffff81111561112157600080fd5b61112d84828501610fe9565b60608301525092915050565b600061114582356113f7565b905092915050565b60006020828403121561115f57600080fd5b600082015167ffffffffffffffff81111561117957600080fd5b6111858482850161103f565b91505092915050565b600080600080600060a086880312156111a657600080fd5b600086013567ffffffffffffffff8111156111c057600080fd5b6111cc88828901611095565b955050602086013567ffffffffffffffff8111156111e957600080fd5b6111f588828901611095565b945050604086013567ffffffffffffffff81111561121257600080fd5b61121e88828901611095565b935050606086013567ffffffffffffffff81111561123b57600080fd5b61124788828901611095565b925050608086013567ffffffffffffffff81111561126457600080fd5b61127088828901611095565b9150509295509295909350565b60006020828403121561128f57600080fd5b600061129d84828501611139565b91505092915050565b60006112b1826113ba565b8084526112c5816020860160208601611410565b6112ce81611443565b602085010191505092915050565b600060208201905081810360008301526112f681846112a6565b905092915050565b6000604082019050818103600083015261131881856112a6565b9050818103602083015261132c81846112a6565b90509392505050565b6000604051905081810181811067ffffffffffffffff8211171561135857600080fd5b8060405250919050565b600067ffffffffffffffff82111561137957600080fd5b601f19601f8301169050602081019050919050565b600067ffffffffffffffff8211156113a557600080fd5b601f19601f8301169050602081019050919050565b600081519050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006113f0826113c5565b9050919050565b6000819050919050565b82818337600083830152505050565b60005b8381101561142e578082015181840152602081019050611413565b8381111561143d576000848401525b50505050565b6000601f19601f830116905091905056fea265627a7a72305820d1ce986c9af5aad255f5ad0a6eedc0bc7f1f4470d37adf6f54fe0c7379a477ef6c6578706572696d656e74616cf50037")


// 初始化给定的合约地址和代码
// 每个节点都初始化相同的合约地址和代码，我们再进行事先规定的映射
func Init(cfg *scexecution.Config) {


	for _, kvstoreaddr := range contractAddr[KVSTORE] {
		cfg.State.CreateAccount(kvstoreaddr)
		// set the receiver's (the executing contract) code for execution.
		cfg.State.SetCode(kvstoreaddr, kvstorecode)
	}

	for _, kvtestaddr := range contractAddr[KVTEST] {
		cfg.State.CreateAccount(kvtestaddr)
		// set the receiver's (the executing contract) code for execution.
		cfg.State.SetCode(kvtestaddr, kvtestcode)
	}

}






