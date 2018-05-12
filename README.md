# Ontology Oracle
## ontology oracle
区块链预言机（Blockchain Oracles）"概念的提出主要是为了解决区块链协议自身存在的局限性：区块链这种去中心化的网络（包括在其上构建的应用/ 智能合约）不能与外部内容进行交互。但很多时候智能合约又必须依赖外部触发条件，比如合约：

当明天某地的天气在14:00至15:00的平均温度如果高于38摄氏度，则向通讯录中的所有人发放一个50代币的红包。

又或者...

某次航班将于10:00am到达目的地，如果航班延误将触发智能合约，投保人将获得100代币的补偿。

随着区块链应用变得越来越复杂，迫切需要解决“围墙花园”的限制。如航班信息，股票价格，贵金属价格，等等大量的数据需要上链服务。对于这样的数据，智能合约是无法从外部网站获取的。所以就需要预言机来搬运，将外部数据写入到区块链中，使合约得以完成。

Ontology oracle正是这样的一个数据搬运工的角色，它使得在智能合约中获取链外数据成为可能，未来还将会针去中心化oracle的数据正确性问题进行一定的博弈。

## ontology oracle基础架构

![framework](/resources/framework.png)

### 链下部分（oracle node & data source）
在链外，ontology oracle由许多oracle node组成，这些node和ontology网络相连，并独立的处理数据请求，未来可以支持更多的区块链网络。
节点的工作由一系列的子任务组成，每个子任务定义了相应的操作，如从数据源获取数据，解析数据，上链等。每个子任务的结果都是下一个子任务的输入。不同节点通过这一系列任务的方式将获取得到的外部数据发送到部署在ontology网络上的oracle contract，由该合约来进行结果聚集。

### 链上部分（oracle contract)
oracle contract主要通过某种共识机制对node发送的数据进行聚集，最终达成一个结果。

### Oracle运作流程
![workflow](/resources/workflow.png)

## 使用说明
### 申请注册Oracle Node
调用Oracle contract的RegisterOracleNode方法，参数为：
type RegisterOracleNodeParam struct {
	Address  string `json:"address"`
	Guaranty uint64 `json:"guaranty"`
}
Address要申请成为Oracle Node的钱包地址，Guaranty需要抵押至少1000个ONT最为押金，该接口调用需要额外消耗500ONG。
申请完成之后，会有管理员审核该申请，同意之后即可成为Oracle Node。

### 配置
```text
{
  "LogLevel": 0,
  "Port": "6688",
  "MaxLogSize": 50,
  "ONTRPCAdress": "http://127.0.0.1:20336",
  "ScannerInterval": 3,
}
```
LogLevel配置日志级别， port配置node的http端口，MaxLogSize配置单个日志文件的大小，不指定默认20M，ONTRPCAdress配置ontology网络的rpc地址，ScannerInterval配置node扫描ontology网络中oracle请求的时间间隔。

### 启动oracle node
生成钱包文件wallet.dat，-p参数为钱包密码
```text
go run main.go node -p passwordtest
```

### 创建oracle request
调用部署在ontology网络上的oracle contract，示例参数：
```go
var operation = "createOracleRequest"
var request = `{
"scheduler":{
    "type": "runAfter",
    "params": "2018-06-15 08:37:18"
},
"tasks":[
  {
    "type": "httpGet",
    "params": {
      "url": "https://bitstamp.net/api/ticker/"
    }
  },
  {
    "type": "jsonParse",
    "params": {
      "path": ["last"]
    }
  }
]
}`
var args = []interface{}{request, address}
```
### job定义
```text
{
    "scheduler":{}，
    "tasks": [{}, {}，...]
}
```
job由scheduler和一系列的task组成，schedule的类型目前支持runAfter，task的类型目前支持HttpGet, HttpPost, JsonParse.

### scheduler定义
```text
{
    "type": "",
    "params": "",
}
```
#### RunAfter
该类型的job会在指定的时间之后执行，如赛事结果oracle可以将其设置为赛事结束之后。

param: 代表时间的字符串，格式如"2018-06-15 08:37:18"

### task定义
```text
{
    "type": "",
    "params": "{json}",
}
```
#### HttpGet
该类型的task会向指定url发送Get请求。

param:

url: get请求的发送地址

#### HttpPost
该类型的task会向指定url发送Post请求。

param:

url: post请求的发送地址

data: post请求要发送的数据内容

#### JsonParse
该类型的task会以path参数为key逐层遍历，返回得到的value结果。

param:

path: 一个string数组，每个string作为下次获取的json数据的key。

### 获取oracle request的结果
```go
var operation = "getOracleOutcome"
var args = []interface{}{txhash, address}
```
txhash为oracle request在链上的交易hash。


