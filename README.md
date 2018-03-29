[中文版](./README_cn.md)

# Ontology Oracle
Ontology’s Oracle and the Prediction market contract.

# Prediction market
Smart Contracts are confined to Blockchain’s isolated sandbox environment. Contract codes are unable to access data outside the chain, therefore relying on external initiators/factors. For example:  

The World Cup 2018 is nearing its kickoff. Football supporters will be able to place bets using smart contracts. Supporter A bets on Brazil to win the championship, while Supporter B bets on Germany. The gamble begins: Brazil wins. Supporter A is awarded with Supporter B’s wager. However, should none of the teams have won, the betting arrangement would have been automatically invalidated. 

Despite the rules on betting being plain and simple, the essential issue regarding its implementation in a decentralized blockchain ecosystem is: how is the result included in the equation? While blockchain is unable to detect or extract any data deriving from the real world, a mechanism is needed to import data from it onto the blockchain. Although, because of decentralized blockchain’s unique features, there is no way of verifying or guaranteeing the truthfulness of the information being imported. If these mechanisms are not designed carefully enough, participants using smart contracts in a bett might benefit from altering the integrity of the information submitted.

# Oracle
The Oracle is a mechanism equipped to upload data onto the blockchain. 

When unable to guarantee data accuracy originating from authoritative channels (in case of external tampering or hackers), or when visiting inaccessible websites, mistrust emerges concerning the source of the data. These errors can cause severe consequences, therefore presenting us with an alternative: the forecast market; using the gains and losses to stimulate/coerce/enforce groups to provide actual information. People will possess binding interest when placing bets, hence ensuring data accuracy. No one wants to bear the consequences caused by error or inaccurate data – thus reassuring the integrity of the data. 

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

# Compilation
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
