# coding=utf-8
# 使用多进程发送随机交易
# 2021/8/17 把随机生成交易改成服从zipf分布

import time
from multiprocessing import Process

import numpy as np
from web3 import Web3
import random
import coordinator_pb2
import coordinator_pb2_grpc
import grpc
import constrctInput
import initContractAddr
from zipfianGenerator import randZipf



def processFunction(start, end, _stub):
    # 空的web3对象
    w3 = Web3()

    # kvtest
    # kvtestContractList = []
    # kvstore
    # kvstoreContractList = []

    kvtestContractList = initContractAddr.InitkvtestContract()
    kvstoreContractList = initContractAddr.InitkvstoreContract()

    totalContractCount = 11
    for i in range(start, end):
        # 每100条输出一次，估计发送速度
        if i % 100 == 0:
            print(i)
        # 构造事务的访问列表，main contract + 10 * called contract
        # 生成11个CONTRACTCOUNT以内的随机数
        numberOfAddr = 100000
        conflictRate = 0
        # randomNumberList = random_list(0, numberOfAddr - 1, 11)
        # randomNumberList = randZipf(numberOfAddr - 1, conflictRate, totalContractCount)
        randomNumberList = np.random.randint(0, numberOfAddr, size=totalContractCount)
        # 根据生成的随机数列表，构造accessAddrList   main contract + 10 * called contract
        accessAddrList_hex = [kvtestContractList[randomNumberList[0]]]
        for j in range(1, totalContractCount):
            accessAddrList_hex.append(kvstoreContractList[randomNumberList[j]])

        accessAddrList = bytes()
        # accessAddrList = accessAddrList + hexStringTobytes(accessAddrList_hex[0][2:])

        for addr_hex in accessAddrList_hex:
            addr = bytes.fromhex(addr_hex[2:])
            accessAddrList = accessAddrList + addr  # 连接访问列表，由协调者解析。约定accessAddrList包含了11个长度为20字节的以太坊合约地址。

        inputData = constrctInput.ContructTransaction(accessAddrList_hex[1:])

        # SendTransaction 方法没有返回，不等待接收返回结果
        _stub.SendTransaction(
            coordinator_pb2.Transaction(id=i,
                                        numberOfAddress=4,
                                        AccessAddressList=accessAddrList[:4*20],
                                        Input=inputData,
                                        positionOfCallChain=0))
        # print(threadID, i)


# 十六进制字符串转bytes
def hexStringTobytes(s):
    return bytes.fromhex(s)


def random_list(start, stop, length):
    if length >= 0:
        length = int(length)
    start, stop = (int(start), int(stop)) if start <= stop else (int(stop), int(start))
    random_list = []
    for i in range(length):
        random_list.append(random.randint(start, stop))
    return random_list


if __name__ == '__main__':
    # sam = randZipf(100000, 0, 11)
    # print(sam)
    time_start = time.time()

    process_list = []
    processCount = 20
    TxPerProcess = 1000

    # grpc连接
    channel = grpc.insecure_channel('10.20.82.234:3000')  # 协调者地址
    stub = coordinator_pb2_grpc.CoordinatorStub(channel)

    for i in range(processCount):  # 开启5个子进程执行fun1函数
        p = Process(target=processFunction, args=(i * TxPerProcess, (i + 1) * TxPerProcess, stub,))  # 实例化进程对象
        p.start()
        process_list.append(p)

    for i in process_list:
        p.join()

    print('结束测试')
    time_end = time.time()
    print('totally cost', time_end - time_start)
