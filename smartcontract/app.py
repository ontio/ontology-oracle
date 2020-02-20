from ontology.interop.System.App import RegisterAppCall
from ontology.interop.System.ExecutionEngine import GetExecutingScriptHash
from ontology.interop.System.Runtime import Notify, Serialize, Deserialize
from ontology.builtins import *
from ontology.interop.Ontology.Runtime import Base58ToAddress

oracleContract = RegisterAppCall('e0d635c7eb2c5eaa7d2207756a4c03a89790934a', 'operation', 'args')

def main(operation, args):
    if operation == 'genRandom':
        return genRandom()
    if operation == 'getRandom':
        if len(args) == 1:
            return getRandom(args[0])
    return False


def genRandom():

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

    res = oracleContract('CreateOracleRequest',[req,Base58ToAddress('AKQj23VWXpdRou1FCUX3g8XBcSxfBLiezp')])

    return True

def getRandom(txHash):
    res = oracleContract('GetOracleOutcome', [txHash])
    if not res:
        return ''
    a = Deserialize(res)
    b = Deserialize(a[0])
    """
    the structure of b is:
    [
        [
            ["0021800316", "1610612744", "128", "1610612761", "131"],
            ["0021800317", "2610612744", "96", "2610612761", "131"],
            ["0021800318", "3610612744", "128", "3610612761", "131"]
        ]
    ]
    """
    Notify(b)
    Notify(b[0])
    Notify(b[0][0])
    Notify(b[0][0][0])
    return True