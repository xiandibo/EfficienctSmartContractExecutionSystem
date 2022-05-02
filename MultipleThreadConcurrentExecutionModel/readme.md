# 运行智能合约多线程并发执行模块

## 文件描述

| 文件路径      | 描述                                       |
| :------------- | :----------------------------------------------- |
| ./code | 源代码        |
| ./workload | 测试代码             |


## 启动模块

智能合约多线程并发执行模块实现在每个区块链节点上。在执行以下步骤前，切换当前目录到`./MultipleThreadConcurrentExecutionModel/code/segment/go-ethereum`

1. 编译源代码
   
    编译源代码需要Golang (1.14版) 和C编译器。运行以下命令来编译源代码：
    ```
    $ make geth
    ```

2. 初始化节点数据

    在`./poadata/`中准备了4个文件夹来存储节点数据：signer1、node1、node2和node3。
	
    以初始化signer1节点为例。首先删除旧的节点数据：
    ```
    $ rm -rf ./poadata/signer1/data/geth/
    ```

    初始化节点数据：
    ```
    $ ./build/bin/geth --datadir poadata/signer1/data init poadata/poatest.json
    ```

3. 运行节点

    如果是运行挖矿节点，执行以下命令：
    ```
    $ ./build/bin/geth --datadir ./poadata/signer1/data --rpc --rpcaddr "0.0.0.0" --rpcport 5409 --rpcapi="eth,net,web3,personal,miner" --networkid 540980442 --port 9002 --allow-insecure-unlock --unlock 0xF886F1916926Dc60F9C75941dBEa1856f997FDe1 --password ./poadata/signer1/data/password  console
	```

    然后输入以下命令来启动挖矿和生成新区块。
    ```
    $ miner.start(1)
	```

    如果是运行验证节点，执行以下命令：
    ```
    $ ./build/bin/geth --datadir ./poadata/node1/data --rpc --rpcaddr "0.0.0.0" --rpcport 5409 --rpcapi="eth,net,web3,personal,miner" --networkid 540980442 --port 9002 --allow-insecure-unlock  console
    ```

    注意：如果是在同一个物理节点上运行多个节点，需要注意使用不同的端口号 `– port port number`

4. 运行一个区块链系统

    要搭建一个区块链系统，需要先运行多个区块链节点。本项目提供了4个节点数据：一个挖矿节点signer1和三个验证节点node1、node2和node3。

    要构建一个区块链系统，需要将三个验证节点连接到挖矿节点。

    在挖矿节点的控制台界面，输入：
    ```
    $ admin.nodeInfo.enode
    >"enode://ac6e7ea07632938ede95a368eb8d4b062472ac17aa24bd0e32e7d442a829847dc64bd6b4a337b973d1c52d0e44c0b64c5735ade383526923907d34b0e4979fa6@127.0.0.1:9002"
	```
    
    在验证节点的控制台界面，输入：
    ```
    $ admin.addPeer("enode://ac6e7ea07632938ede95a368eb8d4b062472ac17aa24bd0e32e7d442a829847dc64bd6b4a337b973d1c52d0e44c0b64c5735ade383526923907d34b0e4979fa6@127.0.0.1:9002")
    >true
    ```

	此时表示验证节点已经连接上挖矿节点。

5. 测试代码

    使用`python3`来实现实验代码。 代码目录为:`./MultipleThreadConcurrentExecutionModel/workload/SmartContractTest`

    运行实验所需的依赖库:`web3 5.13.0`和`py-solc 3.2.0`  
    
    首先运行`deploy .py`来部署智能合约。然后运行`test.py`发送事务 
