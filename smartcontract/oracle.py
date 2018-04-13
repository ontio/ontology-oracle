#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Copyright (C) 2018 The ontology Authors
This file is part of The ontology library.

The ontology is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

The ontology is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with The ontology.  If not, see <http://www.gnu.org/licenses/>.
"""

from boa.blockchain.vm.Neo.TriggerType import Application, Verification
from lib.decentralizedOracle import DecentralizedOracle
from boa.blockchain.vm.Neo.Action import RegisterAction

OWNER = b''

#register event
CreateOracleRequest = RegisterAction('createOracleRequest', 'result')
SetOracleOutcome = RegisterAction('setOracleOutcome', 'result', 'txHash')
SetOracleCronOutcome = RegisterAction('setOracleCronOutcome', 'result', 'txHash')
IsOracleOutcomeSet = RegisterAction('isOracleOutcomeSet', 'result')
GetOracleOutcome = RegisterAction('getOracleOutcome', 'result')

def Main(operation, args):
    """
    :param operation: str The name of the operation to perform
    :param args: list A list of arguments along with the operation
    :return:
        bytearray: The result of the operation
    """
    
    decentralizedOracle = DecentralizedOracle()

    if operation == 'createOracleRequest':
        if len(args) != 2:
            #print('Incorrect Arg Length')
            return False
        request = args[0]
        owner = args[1]
        r = decentralizedOracle.createOracleRequest(request, owner)
        CreateOracleRequest(request)
    elif operation == 'setOracleOutcome':
        if len(args) != 2:
            #print('Incorrect Arg Length')
            return False
        txHash = args[0]
        outcome = args[1]
        r = decentralizedOracle.setOracleOutcome(txHash, outcome)
        SetOracleOutcome(r, txHash)
    elif operation == 'setOracleCronOutcome':
        if len(args) != 2:
            #print('Incorrect Arg Length')
            return False
        txHash = args[0]
        outcome = args[1]
        r = decentralizedOracle.setOracleCronOutcome(txHash, outcome)
        SetOracleCronOutcome(r, txHash)
    elif operation == 'isOracleOutcomeSet':
        if len(args) != 1:
            #print('Incorrect Arg Length')
            return False
        txHash = args[0]
        r = decentralizedOracle.isOracleOutcomeSet(txHash)
        IsOracleOutcomeSet(r)
    elif operation == 'getOracleOutcome':
        if len(args) != 2:
            #print('Incorrect Arg Length')
            return False
        txHash = args[0]
        addr = args[1]
        r = decentralizedOracle.getOracleOutcome(txHash, addr)
        GetOracleOutcome(r)
        
    return False
    
    
    
    
    
    
    
    
    
    
    