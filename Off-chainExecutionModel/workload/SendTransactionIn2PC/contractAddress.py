# 存放合约地址

# 主合约 Main contract
# 合约调用合约的发起者
MainContractAddrList = [
    "0x5e5eb12d9f899d4284bea7f0509c10418dcf934c",
    "0xdbab14855cecb76a7aff8f37e8ec4d2b4611ec4a",
]

# kvstore 合约地址
# 此合约负责存放键值对，被主合约调用
AccessList = [
    "0xf901eaabd215d8c47906d9ee82266d46ca8f4f11",
    "0xe8d8062b51fa8e79e865456a64224d34e8e1c30e",
    "0xf901eaabd215d8c47906d9ee82266d46ca8f4f11",
    "0xe8d8062b51fa8e79e865456a64224d34e8e1c30e",
    "0xf901eaabd215d8c47906d9ee82266d46ca8f4f11",
    "0xe8d8062b51fa8e79e865456a64224d34e8e1c30e",
    "0xf901eaabd215d8c47906d9ee82266d46ca8f4f11",
    "0xe8d8062b51fa8e79e865456a64224d34e8e1c30e",
    "0xf901eaabd215d8c47906d9ee82266d46ca8f4f11",
    "0xe8d8062b51fa8e79e865456a64224d34e8e1c30e",
]

# 把AccessList的元素转换成bytes，方便组织成gprc使用的bytes类型
# [2:]为除去0x前缀
AccessListInBytesForm = [
    bytes.fromhex(AccessList[0][2:]),
    bytes.fromhex(AccessList[1][2:]),
    bytes.fromhex(AccessList[2][2:]),
    bytes.fromhex(AccessList[3][2:]),
    bytes.fromhex(AccessList[4][2:]),
    bytes.fromhex(AccessList[5][2:]),
    bytes.fromhex(AccessList[6][2:]),
    bytes.fromhex(AccessList[7][2:]),
    bytes.fromhex(AccessList[8][2:]),
    bytes.fromhex(AccessList[9][2:]),
]


