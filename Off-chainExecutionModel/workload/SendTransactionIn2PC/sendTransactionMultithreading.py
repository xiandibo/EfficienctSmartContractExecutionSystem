# coding=utf-8

# 为加快发送交易的速度，以多线程的形式发送交易
# 发送固定的合约地址,不是随机生成的

import time

import grpc
import json

import coordinator_pb2
import coordinator_pb2_grpc
from web3 import eth
from web3 import Web3
import web3
from web3 import contract
# from privateChainTestTestSmartContract import compiled_sol_kvtest
import contractAddress
import threading
import constrctInput

# 空的web3对象
w3 = Web3()

# grpc连接
channel = grpc.insecure_channel('10.20.36.229:3000')  # 协调者地址
stub = coordinator_pb2_grpc.CoordinatorStub(channel)


# 从 threading.Thread 继承创建一个新的子类，并实例化后调用 start() 方法启动新线程，
# 即它调用了线程的 run() 方法
class SendTransactionThread(threading.Thread):
    def __init__(self, threadID, startIndex, endIndex):
        threading.Thread.__init__(self)
        self.startIndex = startIndex
        self.endIndex = endIndex
        self.threadID = threadID

    def run(self):
        print("开始线程：", self.threadID)
        send2PCTransaction(self.startIndex, self.endIndex, self.threadID)
        print("退出线程：", self.threadID)


def send2PCTransaction(start, end, threadID):
    for i in range(start, end):
        # 构造事务的访问列表，main contract + 10 * called contract
        accessAddrList = bytes()
        accessAddrList = accessAddrList + hexStringTobytes(contractAddress.MainContractAddrList[i % 2][2:])

        for addr in contractAddress.AccessListInBytesForm:
            accessAddrList = accessAddrList + addr  # 连接访问列表，由协调者解析。约定accessAddrList包含了11个长度为20字节的以太坊合约地址。

        # SendTransaction 方法没有返回，不等待接收返回结果
        stub.SendTransaction(
            coordinator_pb2.Transaction(id=i,
                                        numberOfAddress=11,
                                        AccessAddressList=accessAddrList,
                                        Input=constrctInput.txInput,
                                        positionOfCallChain=0))
        # print(threadID, i)


# 本程序为基于指定的合约调用形式，生成标准的智能合约输入数据
def run():
    time_start = time.time()

    threadList = []
    threadCount = 10
    TxPerThread = 1000

    # 初始化
    for i in range(threadCount):
        thread = SendTransactionThread(i, i * TxPerThread, (i + 1) * TxPerThread)
        threadList.append(thread)

    # 启动线程
    for elem in threadList:
        elem.start()

    # 等待线程执行结束
    for elem in threadList:
        elem.join()

    print("退出主线程")

    time_end = time.time()
    print('totally cost', time_end - time_start)


# 十六进制字符串转bytes
def hexStringTobytes(s):
    return bytes.fromhex(s)


if __name__ == '__main__':
    run()
