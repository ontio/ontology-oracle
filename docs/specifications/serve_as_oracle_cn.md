[English](serve_as_oracle_EN.md) | 中文

# 如何提供Oracle服务

任何人都可以成为Oracle服务提供商，其既可以提供一些公共数据的获取服务，也可以提供一些只有自己掌握的专业数据的获取服务，并按照每次请求收取用户一定的手续费。下面将详细描述如何成为一个Oracle服务提供商：

## 部署自己的Oracle合约

首先数据提供方需要部署自己的Oracle合约，协议标准见：[Oracle协议标准](https://github.com/ontio/OEPs/blob/master/OEPS/OEP-34.mediawiki)

Oracle合约样例参考：[Oracle合约样例](https://github.com/ontio/ontology-oracle/blob/master/smartcontract/oracle.py)

## 下载Oracle client

从Oracle仓库克隆源代码：

```shell
$ git clone https://github.com/ontio/ontology-oracle.git
```

## 添加自己的专业数据服务

Oracle仓库中默认有公共数据的获取服务，不需要任何权限就可以获取到数据。以及random.org的专业数据服务，需要申请key来获取数据，样例中使用了免费的基础开发版API key。

数据提供方可以添加任何专业数据服务，也可以只提供样例中存在的公共数据服务。

数据服务在文件夹[runners](https://github.com/ontio/ontology-oracle/tree/master/runners)中：

`runner.go`定义了runner接口以及不同runner的路由。

`http.go`提供了基础的额公共数据获取服务，通过http的方式获取数据。

`json_parse.go`提供了数据解析服务。

`random_org.go`提供了获取随机数的专业数据服务，包括普通随机数和签名随机数。

数据提供方可以参考`random_org.go`来实现自己的专业数据服务。

## 编译Oracle client

```shell
go build main.go
```

## 修改配置文件

修改配置文件`config.json`：

```json
{
  "WalletFile": "./wallet.dat",
  "LogLevel": 0,
  "ONTRPCAddress": "http://127.0.0.1:20336",
  "ScannerInterval": 10,
  "GasPrice": 0,
  "GasLimit": 20000,
  "ContractAddress": "b869eed573863f8efdb3ae39d4963a781e81d4b3",
  "RandomOrgKey": "c7511065-c88d-4f28-af4f-293c91ad20d9"
}
```

`WalletFile`: ontology的钱包路径，此钱包用于支付上传数据请求的结果会消耗gas。

`LogLevel`: 日志级别

```go
const (
	DebugLog = iota
	InfoLog
	WarnLog
	ErrorLog
	FatalLog
	MaxLevelLog
)
```

`ONTRPCAddress`: ontology节点的rpc地址和端口，上传数据请求结果的交易发往该节点。

`ScannerInterval`: 扫描时间间隔，每隔固定时间oracle client从oracle合约中扫描一次数据请求。

`ContractAddress`: oracle合约的合约地址hash

`RandomOrgKey`: random.org的免费基础开发版key，可以添加其他专业数据的获取权限证明。

## 启动Oracle client

```shell
./main node
```

输入钱包密码即可启动。

## 公开Oracle合约地址和提供的数据服务说明

Oracle合约部署完成，Oracle client启动完成后，数据提供商可以公开自己的Oracle合约地址和收费标准，以及提供的数据服务内容和request发送格式，供用户使用。

