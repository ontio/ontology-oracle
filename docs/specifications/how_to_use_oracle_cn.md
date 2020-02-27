[English](how_to_use_oracle_EN.md) | 中文

# 如何在app中使用Oracle服务

用户需要在自己的app智能合约中添加对Oracle合约的调用来使用Oracle服务，如在需要外部数据时，构造数据请求，发送给Oracle合约，并通过Oracle合约获取数据请求的结果。步骤如下：

### 注册对Oracle合约的跨合约调用：

```python
from ontology.interop.System.App import RegisterAppCall

oracleContract = RegisterAppCall('e0d635c7eb2c5eaa7d2207756a4c03a89790934a', 'operation', 'args')
```

目前，已经部署的案例Oracle合约地址为：

测试网：`e0d635c7eb2c5eaa7d2207756a4c03a89790934a`

主网：`a6ee997b142b002d49670ab73803403b09a23fa0`

该Oracle合约提供一些公共数据的获取方法和从random.org获取随机数的方法。

### 构造Oracle Request

```python
req = """{
		"scheduler":{
			"type": "runAfter",
			"params": "2018-06-15 08:37:18"
		},
		"tasks":[
			{
			    "type": "httpGet",
			    "params": {
				    "url": "https://data.nba.net/prod/v2/20181129/scoreboard.json"
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

Oracle Request的类型和构造方法详见[Oracle_Request](oracle_request.md)

### 调用Oracle合约发送数据请求

```python
 oracleContract('CreateOracleRequest', [req, Base58ToAddress('AKQj23VWXpdRou1FCUX3g8XBcSxfBLiezp')])
```

可以看到传入参数为request内容和地址，其中地址为调用者的地址，该地址会为该请求支付一定的手续费，该手续费由Oracle提供方决定。测试网和主网的案例Oracle合约额外收取0.01 ong作为上传数据请求结果的手续费。

### 调用Oracle合约获取数据请求结果

```python
res = oracleContract('GetOracleOutcome', [txHash])

a = Deserialize(res)
b = Deserialize(a[0])
```

返回值是一个数组：

```py
res = state(data, status, errMessage)
```

返回值中第一个元素是数据请求的结果，第二个元素为数据请求的状态，类型为string，"completed"代表数据请求成功，"errored"代表数据请求失败，第三个元素为失败的error message，类型为string。

因此b是最终解析到的数据请求结果，b的具体格式和类型取决于用户在Oracle Request中的定义，详见[Oracle_Request](oracle_request.md)。

app合约样例见：[合约使用Oracle服务样例](https://github.com/ontio/ontology-oracle/blob/master/smartcontract/app.py)