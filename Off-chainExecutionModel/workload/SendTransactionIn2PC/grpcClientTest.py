# # coding=utf-8
# # from __future__ import print_function
# import os
#
# import grpc
# import json
#
# import coordinator_pb2
# import coordinator_pb2_grpc
# from web3 import eth
# from web3 import Web3
# import web3
# from web3 import contract
# from privateChainTestTestSmartContract import compiled_sol_kvtest
#
# # 空的web3对象
# w3 = Web3()
#
# MainContractAddrList = [
#     "0x684c903c66d69777377f0945052160c9f778d689",
#     "0x2ffea410b53de84d04c3b55e52a0d53fcfd0146e",
# ]  # 主合约
#
# AccessList = [
#     "0xbd770416a3345f91e4b34576cb804a576fa48eb1",
#     "0x5a443704dd4b594b382c22a083e2bd3090a6fef3",
#     "0xbd770416a3345f91e4b34576cb804a576fa48eb1",
#     "0x5a443704dd4b594b382c22a083e2bd3090a6fef3",
#     "0xbd770416a3345f91e4b34576cb804a576fa48eb1",
#     "0x5a443704dd4b594b382c22a083e2bd3090a6fef3",
#     "0x5a443704dd4b594b382c22a083e2bd3090a6fef3",
#     "0xbd770416a3345f91e4b34576cb804a576fa48eb1",
#     "0x5a443704dd4b594b382c22a083e2bd3090a6fef3",
#     "0xbd770416a3345f91e4b34576cb804a576fa48eb1",
# ]
#
# AccessListInBytesForm = [
#     bytes.fromhex(AccessList[0][2:]),
#     bytes.fromhex(AccessList[1][2:]),
#     bytes.fromhex(AccessList[2][2:]),
#     bytes.fromhex(AccessList[3][2:]),
#     bytes.fromhex(AccessList[4][2:]),
#     bytes.fromhex(AccessList[5][2:]),
#     bytes.fromhex(AccessList[6][2:]),
#     bytes.fromhex(AccessList[7][2:]),
#     bytes.fromhex(AccessList[8][2:]),
#     bytes.fromhex(AccessList[9][2:]),
# ]
#
#
# # 本程序为基于指定的合约调用形式，生成标准的智能合约输入数据
#
# def run():
#     # grpc连接
#     channel = grpc.insecure_channel('10.20.36.229:3000')
#     stub = coordinator_pb2_grpc.CoordinatorStub(channel)
#
#     # if os.environ.get('https_proxy'):
#     #     del os.environ['https_proxy']
#
#     # res = hexStringTobytes('bd770416a3345f91e4b34576cb804a576fa48eb1')  # 能够解析出来
#     # print("%s" % res)
#     # contract_instance = w3.eth.contract(Web3.toChecksumAddress("0x684c903c66d69777377f0945052160c9f778d689"), abi=abi)
#
#     # 构造事务，对于kvtest和kvstore
#     abi = json.loads(compiled_sol_kvtest['contracts']['kvtest.sol']['Kvs']['metadata'])['output']['abi']
#     # 获取函数选择器hash
#     fn_selector = web3.contract.encode_hex(web3.contract.function_abi_to_4byte_selector(abi[1]))
#     res = contract.encode_abi(
#         w3, abi[1],
#         [
#             (Web3.toChecksumAddress(AccessList[0]), str(1),
#              Web3.toChecksumAddress(AccessList[1]), str(2)),
#             (Web3.toChecksumAddress(AccessList[2]), str(3),
#              Web3.toChecksumAddress(AccessList[3]), str(4)),
#             (Web3.toChecksumAddress(AccessList[4]), str(5),
#              Web3.toChecksumAddress(AccessList[5]), str(6)),
#             (Web3.toChecksumAddress(AccessList[6]), str(7),
#              Web3.toChecksumAddress(AccessList[7]), str(8)),
#             (Web3.toChecksumAddress(AccessList[8]), str(9),
#              Web3.toChecksumAddress(AccessList[9]), str(10))
#         ],
#         fn_selector
#     )
#
#     res = hexStringTobytes(res[2:])  # 这个语句花的时间比较长
#
#     # grpc远程调用协调者的SendTransaction接口，发送交易
#     for i in range(4000):
#
#         accessAddrList = bytes()
#         accessAddrList = accessAddrList + hexStringTobytes(MainContractAddrList[i % 2][2:])
#
#         for addr in AccessListInBytesForm:
#             # print(addr[2:])
#             # addr = hexStringTobytes(addr[2:])  # 去除0x前缀
#             accessAddrList = accessAddrList + addr  # 连接访问列表，由协调者解析。约定accessAddrList包含了11个长度为20字节的以太坊合约地址。
#
#         # SendTransaction 方法没有返回，不等待接收返回结果
#         stub.SendTransaction(
#             coordinator_pb2.Transaction(id=i, numberOfAddress=11, AccessAddressList=accessAddrList, Input=res,
#                                         positionOfCallChain=0))
#         print(i)
#
#
# # 十六进制字符串转bytes
# def hexStringTobytes(str):
#     return bytes.fromhex(str)
#
#
# if __name__ == '__main__':
#     run()
