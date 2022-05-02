# 从文件中读取合约地址，组织成列表
def InitkvtestContract():
    with open('./kvtest.txt') as f:
        kvtestList = f.readlines()
        for i in range(len(kvtestList)):
            kvtestList[i] = kvtestList[i][:len(kvtestList[i])-1]

        # print(kvtestList)

        f.close()
        return kvtestList


def InitkvstoreContract():
    with open('./kvstore.txt') as f:
        kvstoreList = f.readlines()
        for i in range(len(kvstoreList)):
            kvstoreList[i] = kvstoreList[i][:len(kvstoreList[i])-1]

        # print(kvstoreList)

        f.close()
        return kvstoreList


def run():
    kvtestList = InitkvtestContract()
    print(kvtestList)
    kvstoreList = InitkvstoreContract()
    print(kvstoreList)


if __name__ == '__main__':
    run()
