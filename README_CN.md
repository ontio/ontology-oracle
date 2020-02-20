[English](README.md) | 中文

# Ontology Oracle
## ontology oracle
区块链预言机（Blockchain Oracles）"概念的提出主要是为了解决区块链协议自身存在的局限性：区块链这种去中心化的网络（包括在其上构建的应用/ 智能合约）不能与外部内容进行交互。但很多时候智能合约又必须依赖外部触发条件，比如合约：

某次航班将于10:00am到达目的地，如果航班延误将触发智能合约，投保人将获得100代币的补偿。

随着区块链应用变得越来越复杂，迫切需要解决“围墙花园”的限制。如航班信息，股票价格，贵金属价格，等等大量的数据需要上链服务。对于这样的数据，智能合约是无法从外部网站获取的。所以就需要预言机来搬运，将外部数据写入到区块链中，使合约得以完成。

Ontology oracle正是这样的一个数据搬运工的角色，它使得在智能合约中获取链外数据成为可能，未来还将会针去中心化oracle的数据正确性问题进行一定的博弈。

## ontology oracle基础架构

![framework](/resources/framework.png)

### 链下部分（oracle node & data source）
在链外，ontology oracle由oracle node组成，node和ontology网络相连，并独立的处理数据请求，未来可以支持更多的区块链网络。

节点的工作由两个子任务组成，数据爬取和数据解析，数据爬取从api中获取response，数据解析解析该response，并按照用户定义的数据结构序列化，写入oracle合约。

### 链上部分（oracle contract)
oracle contract主要对node发送的数据进行聚集和存储，供其他合约调用。

### Oracle运作流程
![workflow](/resources/workflow.png)

## 文档
[如何在应用中使用Oracle服务](docs/specifications/how_to_use_oracle_cn.md)

[如何提供一个Oracle服务](docs/specifications/serve_as_oracle_cn.md)