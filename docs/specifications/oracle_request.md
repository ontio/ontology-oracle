English | [中文](oracle_request_cn.md)

## Oracle request
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
raw http response:
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
Application contract deserialize result as follows:
```text
[
    "1542359479",
    "5518.70",
    5571.12,
]
```
params:

url: url of http get request

### JsonParse
jsonParse will parse the http response, key is the "path" list in params. Then serialize result as "data" structure defined by users. And write result in oracle contract.

params:

data: data structure defined by users, used to serialize and deserialize result.

type: data type, support int, float(* decimal as int), string, array, map, struct.

sub_type: sub_type of array, map and struct.

decimal: decimal of float。

path: iterator list of json parse key, if data is json, write key in list, if data is array, write index as string in list.

Following is a complex JsonParse example

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

parts of raw http response is:
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
This example parse raw http response, then create a struct according to "data" structure. Then serialize the struct.

Application contract call oracle contract to get result, then deserialize the result as follows:
```text
[
    [
        ["0021800316", "1610612744", "128", "1610612761", "131"],
        ["0021800317", "2610612744", "128", "2610612761", "131"],
        ["0021800318", "3610612744", "128", "3610612761", "131"]
    ]
]
```
### scheduler
```text
{
    "type": "",
    "params": "",
}
```
type only support RunAfter, means task run after a given time, empty means run right now.
#### RunAfter
This type of task will run after given time, for example, users can get the result of a game after the game ended.

params: string indicates time, format in "2018-06-15 08:37:18"

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
raw http response:
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

Application contract deserialize result as follows:
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
The http post example above get random number from random.org, this ontology oracle package a more convenient method "randomOrg" to get random number.

#### signed random number
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

n: amount of random numbers

min: min of random numbers

max: max of random numbers

replacement: true(random number can occur several times), false(random number is unique)

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
#### unsigned random number
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

n: amount of random numbers

min: min of random numbers

max: max of random numbers

replacement: true(random number can occur several times), false(random number is unique)

response:
```go
type IntegerData struct {
	Data           []interface{} `json:"data"`
	CompletionTime string        `json:"completionTime"`
}
```
