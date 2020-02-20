English | [中文](README_CN.md)

# Ontology Oracle
## ontology oracle
Blockchain Oracle is designed to solve the problem that smartcontract can't interact with outside world. Now more and more dapps depend on outside world trigger: for example:

One flight will arrive at 10:00 am, if the flight delays, insurance smartcontract will be triggered, all applicants should get 100 token for compensation.

Now dapps need more and more outside world data, like flight information, stoke price, weather information, match result...

Ontology oracle is exactly the role of data transporter, makes it possible for smartcontract to get outside world data.

## ontology oracle framework

![framework](/resources/framework.png)

### off-chain（oracle operator & data source）
Outside the blockchain, ontology oracle consist of oracle operator and data source.

Operator listen the oracle request of oracle contract in ontology network, process the request and fetch data from data source, serialize result and put it in oracle request.

Oracle request consist of 2 tasks, data fetch and data parse. Data fetch get response of target api, data parse parse the response and serialize the result according to the data structure defined by users.

### on-chain（oracle contract)
On block chain, oracle contract receive oracle request from users and result from oracle operator. Any application contract and call oracle contract to request outside world data and get result.

### Oracle work flow
![workflow](/resources/workflow.png)

## Document
More document see: [Docs](docs)