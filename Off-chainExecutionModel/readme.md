# 智能合约链下执行模块

## 文件描述

| 文件路径      | 描述                                       |
| :------------- | :----------------------------------------------- |
| ./remote-commiter | 源代码        |
| ./workload | 测试代码             |


## 启动模块

在执行以下步骤前，切换当前目录到 `./Off-chainExecutionModel/remote-commiter`

链下执行模块由一个协调者和多个共识组组成，每个共识组由一个代理节点和多个随从节点组成。该模块的编译命令和节点部署命令存在 `./Makefile` 文件中

1. 编译源代码

    一份代码表示一个节点的文件，若想要部署多个节点，需要重复部署多个代码。对于每一个节点，在启动节点之前需要编译代码：
    ```
    $ make prepare
    ```

2. 启动节点

    (1) 启动协调者节点
    
    在命令行输入以下指令启动协调者：
    ```
    $ make run-example-coordinator
    ```

    (2) 启动共识组

    启动一个共识组包括启动代理节点和随从节点。
    
    启动代理节点：
    ```
    make run-example-follower1
    ```

    启动多个随从节点：
    ```
    make run-example-voter-4111
    ```

    ```
    make run-example-voter-4112
    ```

    其他类似的节点的启动命令与此类似，被定义在`Makefile`中
    
    启动节点包括智能合约的自动部署和节点连接。

3. 测试代码
   
    使用`python3`来实现实验代码。 代码目录为: `./ Off-chainExecutionModel/test/SendTransactionIn2PC`

    运行`sendTransactionMultiPro.py`文件来发送事务，随后在协调者节点的控制台输入`run`命令启动事务执行即可开启测试 
