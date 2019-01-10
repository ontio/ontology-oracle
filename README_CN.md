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

## 创建oracle request
用户可以在自己的合约中调用部署在ontology网络上的oracle contract来获取外部数据，目前支持httpGet， httpPost， random.Org获取随机数：

目前，该ontology oracle案例已经部署，合约地址为：

测试网：e0d635c7eb2c5eaa7d2207756a4c03a89790934a

主网：a6ee997b142b002d49670ab73803403b09a23fa0
### httpGet
```text
operation = "CreateOracleRequest"
request = """{
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
				"params":
				{
					"data":
					[
						{
							"type": "String",
							"path": ["timestamp"]
						},
						{
							"type": "String",
							"path": ["last"]
						},
						{
							"type": "Float",
							"decimal": 100,
							"path": ["open"]
						}
					 ]
				}
			}
		]
	}"""
args = [request, address]
```
http response:
```text
{
	"high": "5610.00000000",
	"last": "5518.70",
	"timestamp": "1542359479",
	"bid": "5518.12",
	"vwap": "5436.78",
	"volume": "16423.18407040",
	"low": "5199.80000000",
	"ask": "5518.69",
	"open": 5571.12
}
```
params:

url: http get请求的地址url

### JsonParse
该类型的task会以path参数为key逐层遍历，返回得到的value结果。

params:

data: 所要获取的数据结构，以结构体的形式返回。

type: 数据类型，支持int，float（乘以精度按照int处理），string ，array，map，struct。

sub_type: array, map的子类型。

decimal: 浮点数所要乘以的精度。

path: 一个数组，每个元素作为下次获取数据的索引, 如果数据为json, 则是该json数据的key, 如果数据为array, 则是该数据的index。

下面给出一个更复杂的JsonParse案例：

```text
var request = """{
		"scheduler":{
			"type": "runAfter",
			"params": "2018-06-15 08:37:18"
		},
		"tasks":[
			{
			  "type": "httpGet",
			  "params": {
				"url": "http://data.nba.net/prod/v2/20181129/scoreboard.json"
			  }
			},
			{
				"type": "jsonParse",
				"params":
				{
					"data":
					[
						{
							"type": "Array",
							"path": ["games"],
							"sub_type":
							[
								{
									"type": "Struct",
									"sub_type":
									[
										{
											"type": "String",
											"path": ["gameId"]
										},
										{
											"type": "String",
											"path": ["vTeam", "teamId"]
										},
										{
											"type": "String",
											"path": ["vTeam", "score"]
										},
										{
											"type": "String",
											"path": ["hTeam", "teamId"]
										},
										{
											"type": "String",
											"path": ["hTeam", "score"]
										}
									]
								}
							]
						}
					]
				}
			}
		]
	}"""
```

部分原始http response返回格式为：
```text
{
	"numGames": 3,
	"games": [{
			"gameId": "0021800316",
			"vTeam": {
				"teamId": "1610612744",
				"score": "128"
			},
			"hTeam": {
				"teamId": "1610612761",
				"score": "131"
			}
		},
		{
			"gameId": "0021800317",
			"vTeam": {
				"teamId": "2610612744",
				"score": "96"
			},
			"hTeam": {
				"teamId": "2610612761",
				"score": "131"
			}
		},
		{
			"gameId": "0021800318",
			"vTeam": {
				"teamId": "3610612744",
				"score": "128"
			},
			"hTeam": {
				"teamId": "3610612761",
				"score": "131"
			}
		}
	]
}
```
该案例解析原始http response，并根据“data”的结构创建结构体，并序列化。

应用合约通过跨合约调用oracle合约获取数据结果，反序列化后的返回如下：
```text
[
    [
        ["0021800316", "1610612744", "128", "1610612761", "131"],
        ["0021800317", "2610612744", "128", "2610612761", "131"],
        ["0021800318", "3610612744", "128", "3610612761", "131"]
    ]
]
```
### scheduler定义
```text
{
    "type": "",
    "params": "",
}
```
type目前支持RunAfter，在某时刻后执行。 此项不填为立即执行。
#### RunAfter
该类型的job会在指定的时间之后执行，如赛事结果oracle可以将其设置为赛事结束之后。

params: 代表时间的字符串，格式如"2018-06-15 08:37:18"

### httpPost
```text
operation = "CreateOracleRequest"
request = """{
        "scheduler":{
            "type": "runAfter",
            "params": "2018-06-15 08:37:18"
        },
        "tasks":[
            {
              "type": "httpPost",
              "params": {
                "url": "https://api.random.org/json-rpc/1/invoke",
                "contentType": "application/json-rpc",
                "body": "{\\"jsonrpc\\": \\"2.0\\",\\"method\\": \\"generateSignedIntegers\\",\\"params\\": {\\"apiKey\\": \\"c7511065-c88d-4f28-af4f-293c91ad20d9\\",\\"n\\": 6,\\"min\\": 1,\\"max\\": 10,\\"replacement\\": false,\\"base\\": 10},\\"id\\": 1}"
              }
            },
            {
                "type": "jsonParse",
                "params":
                {
                    "data":
                    [
                        {
                            "type": "Array",
                            "path": ["result", "random", "data"],
                            "sub_type":
                                [
                                    {
                                        "type": "Int"
                                    }
                                ]
                        }
                    ]
                }
            }
        ]
    }"""
args = [request, address]
```
http response:
```text
{
    "jsonrpc": "2.0",
    "result": {
        "random": {
            "method": "generateSignedIntegers",
            "hashedApiKey": "oT3AdLMVZKajz0pgW/8Z+t5sGZkqQSOnAi1aB8Li0tXgWf8LolrgdQ1wn9sKx1ehxhUZmhwUIpAtM8QeRbn51Q==",
            "n": 6,
            "min": 1,
            "max": 6,
            "replacement": true,
            "base": 10,
            "data": [
                2,
                4,
                4,
                1,
                5,
                3
            ],
            "completionTime": "2013-09-30 14:58:03Z",
            "serialNumber": 69260
        },
        "signature": "BxHxajeRg7Q+XGjBdFS1c7wkZbJgJlverfZ5TVDyzCKqo2K5A4pD+54EMqmysRYwkL3w2NS2DFLVrsyO1o96bW9BGp5zjjrEegz9mB+04iOTaRwmdQnLJAj/m3WRptA+qzodPCTaqud8YWBifqWCM34q98XwjX+nlahyHVHT9vf5KO0YVkD/yRI1WN5M/qX21chVvSxhWdmIrdCkrovGnysFq8SzCRNhpYx+/1P+YT2IKsH8jth9z82IAz1ANVh918H/UdpuD1dR7TD6nk3ntRgGrIiu2qqVzFi8A7/6viVgRqtffE4KVZY6O9mUJ+sGkF5Ohayms7LHSFy1VC8wMbMgwod+A8nr5yzjAC4SCUkT1bKAyWNF3SdVcLtvWdcf97Ew6RjohzCW4Vs3jUlh6jF/pj3b3++U3lBHCh43IIonw8MQ7afwpqP12yvyDym1isNjhMKYjmzWRerSvnsMyQIH8xFW7IHt2g/0qnzJgABFmUNBRKJPCD9CMgjh60sSwW7EyrGMy7/qisfE0IU74P/F7KCty/g1jIlXX5/O1lQjwY34wnoP0NXL08QteukRZZUfJQnscx1NGE+HX1c9bMBI8LC0ZFYFk+uY6ib/0rCV5OcLLE9PihCdC8WoI1x3bobr8tbtfgnXMTjogxwVXiiSN1TMnTIWlJ+KM5eSWrw=",
        "bitsUsed": 16,
        "bitsLeft": 932400,
        "requestsLeft": 199991,
        "advisoryDelay": 1000
    },
    "id": 1
}
```
该案例返回一个结构体，结构体中有一个数组，数组中有6个元素。

应用合约解析后的返回如下：
```text
[
    [
        2,
        4,
        4,
        1,
        5,
        3
    ]
]
```
### randomOrg
上述httpPost的例子其实是在random.org获取签名随机数，目前的oracle模板专门为random.org的随机数封装了一个更简便的调用方法randomOrg。
#### 签名随机数
```text
operation = "CreateOracleRequest"
request = """{
        "scheduler":{
            "type": "runAfter",
            "params": "2018-06-15 08:37:18"
        },
        "tasks":[
            {
              "type": "randomOrg",
              "params": {
                "method": "GenerateSignedIntegers",
                "n": 10,
                "min": 1,
                "max": 10,
                "replacement": false
              }
            },
            {
                "type": "jsonParse",
                "params":
                {
                    "data":
                    [
                        {
                            "type": "Array",
                            "path": ["data"],
                            "sub_type":
                                [
                                    {
                                        "type": "Int"
                                    }
                                ]
                        },
                        {
                            "type": "String",
                            "path": ["signature"]
                        }
                    ]
                }
            }
        ]
    }"""
args = [request, address]
```
params:

n: 随机数的个数

min: 随机数的最小值

max: 随机数的最大值

replacement: 是否允许重复(true:是 / false:否)

response:
```go
type SignedIntegerData struct {
	Raw          json.RawMessage `json:"raw"`
	HashedApiKey string          `json:"hashedApiKey"`
	SerialNumber int             `json:"serialNumber"`
	Data         []int           `json:"data"`
	Signature    string          `json:"signature"`
}
```
#### 非签名随机数
```text
operation = "{
        "scheduler":{
            "type": "runAfter",
            "params": "2018-06-15 08:37:18"
        },
        "tasks":[
            {
              "type": "randomOrg",
              "params": {
                "method": "GenerateIntegers",
                "n": 10,
                "min": 1,
                "max": 10,
                "replacement": false
              }
            },
            {
                "type": "jsonParse",
                "params":
                {
                    "data":
                    [
                        {
                            "type": "Array",
                            "path": ["data"],
                            "sub_type":
                                [
                                    {
                                        "type": "Int"
                                    }
                                ]
                        },
                        {
                            "type": "String",
                            "path": ["completionTime"]
                        }
                    ]
                }
            }
        ]
    }"""
args = [request, address]
```
params:

n: 随机数的个数

min: 随机数的最小值

max: 随机数的最大值

replacement: 是否允许重复(true:是 / false:否)

response:
```go
type IntegerData struct {
	Data           []interface{} `json:"data"`
	CompletionTime string        `json:"completionTime"`
}
```

### 获取oracle request的结果
```text
operation = "GetOracleOutcome"
args = txhash
```
txhash为oracle request在链上的交易hash。

# 代码示例
Oracle contract模板见smartcontract/oracle.py

使用Oracle的应用合约模板见smartcontract/app.py
