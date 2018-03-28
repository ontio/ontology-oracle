# Ontology Oracle
本体网络上的oracle和预测市场合约。

# 预测市场
智能合约是在区块链提供的沙盒环境中运行，沙盒是个封闭环境，使合约代码不能读取链外数据，但很多时候智能合约又必须依赖外部触发条件。比如：

2018 世界杯即将开赛，球迷可用智能合约来实现对赌，比如有球迷 A 预测巴西队会夺冠，另有球迷 B 预测德国队夺冠。这就可以开设一个赌局：巴西队夺冠，则 B 的赌注判给 A，反之亦然，如果两队都没有夺冠，则赌约自动解除。

赌约的规则很简单，但在区块链去中心化体系下存在一个关键问题是，如何将比赛结果放进这个赌约当中去？区块链系统中却无从判断外面现实世界发生的事件，这就需要我们引入一种机制将现实社会的事件输入区块链之中。不过，因为区块链“去中心化”的特点，没有一个节点可以对输入信息的真伪做出裁决，如果这种机制设计得不够周密，那么参与智能合约赌局的一方就很有可能为了利益而否认事实。

# 预言机
这种机制就是Oracle预言机，Oracle的核心功能是提供数据上链服务。

我们没办法保证官方渠道的准确性（被黑客篡改），和网站无法访问的情况，也就无法信任这个数据来源。而这类错误又将造成严重后果，在这种情况下我们还有另一个选择，那就是预测市场，用利益得失来强迫群体说真话。人们下注，博弈，以利益绑定的形式来保证数据的准确性。没人愿意为错误的比赛结果买单，所以这个数据更可信。

# Requirements
Usage requires Python 3.5.

# Installation
Clone the repository and navigate into the project directory. Make a Python 3 virtual environment and activate it via

```
python3 -m venv venv
source venv/bin/activate
```

or to explicitly install Python 3.5 via

```
virtualenv -p /usr/local/bin/python3.5 venv
source venv/bin/activate
```
Then install the requirements via

```
pip install -r requirements.txt
```

#Compilation
The template may be compiled as follows

```python
from boa.compiler import Compiler

Compiler.load_and_save('oracleContract/oracle.py')
```
This will compile your conrtact to oracle.avm

Same as marketsContract


# Contributing

Contributions are welcome!

Feel free to open issues for discussion, and create pull requests to post your updates.


# License

This project is under LGPL v3.0 license.
