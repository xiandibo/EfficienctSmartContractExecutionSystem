# 编译kvtest合约，获取合约abi
import json

import numpy as np
import web3
from web3 import contract
from web3 import eth
from web3 import Web3
from privateChainTestTestSmartContract import compiled_sol_kvtest
import contractAddress


# 十六进制字符串转bytes
def hexStringTobytes(str):
    return bytes.fromhex(str)


w3 = Web3()

abi = json.loads(compiled_sol_kvtest['contracts']['kvtest.sol']['Kvs']['metadata'])['output']['abi']
# 获取函数选择器hash，此处abi[1] 表示获取合约的第二个函数test,也就是我们调用的函数
fn_selector = web3.contract.encode_hex(web3.contract.function_abi_to_4byte_selector(abi[1]))
# 构造事务的input,也即事务的data字段
txInput = contract.encode_abi(
    w3, abi[1],
    [
        (Web3.toChecksumAddress(contractAddress.AccessList[0]), str(1),
         Web3.toChecksumAddress(contractAddress.AccessList[1]), str(2)),
        (Web3.toChecksumAddress(contractAddress.AccessList[2]), str(3),
         Web3.toChecksumAddress(contractAddress.AccessList[3]), str(4)),
        (Web3.toChecksumAddress(contractAddress.AccessList[4]), str(5),
         Web3.toChecksumAddress(contractAddress.AccessList[5]), str(6)),
        (Web3.toChecksumAddress(contractAddress.AccessList[6]), str(7),
         Web3.toChecksumAddress(contractAddress.AccessList[7]), str(8)),
        (Web3.toChecksumAddress(contractAddress.AccessList[8]), str(9),
         Web3.toChecksumAddress(contractAddress.AccessList[9]), str(10))
    ],
    fn_selector
)

# 把input编码成bytes，这里gprc通信传输的数据类型
txInput = hexStringTobytes(txInput[2:])  # 这个语句花的时间比较长


# 生成随机合约事务的版本
def ContructTransaction(accessList):
    keyList = np.random.randint(0, 1000, size=10)
    input = contract.encode_abi(
        w3, abi[1],
        [
            (Web3.toChecksumAddress(accessList[0]), str(keyList[0]),
             Web3.toChecksumAddress(accessList[1]), str(keyList[1])),
            (Web3.toChecksumAddress(accessList[2]), str(keyList[2]),
             Web3.toChecksumAddress(accessList[3]), str(keyList[3])),
            (Web3.toChecksumAddress(accessList[4]), str(keyList[4]),
             Web3.toChecksumAddress(accessList[5]), str(keyList[5])),
            (Web3.toChecksumAddress(accessList[6]), str(keyList[6]),
             Web3.toChecksumAddress(accessList[7]), str(keyList[7])),
            (Web3.toChecksumAddress(accessList[8]), str(keyList[8]),
             Web3.toChecksumAddress(accessList[9]), str(keyList[9]))
        ],
        fn_selector
    )
    input = hexStringTobytes(input[2:])  # 这个语句花的时间比较长

    return input
