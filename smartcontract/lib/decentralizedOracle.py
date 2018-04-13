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

from boa.code.builtins import concat
from boa.blockchain.vm.Neo.Runtime import CheckWitness
from boa.blockchain.vm.Neo.Storage import GetContext, Get, Put
from lib.utils import getTxHash

oracleAddr = b''

class DecentralizedOracle(object):
    
    def createOracleRequest(self, request: str, addr: bytes):
        ctx = GetContext()
        txHash = getTxHash()
        
        #get undo txhash
        numPre = Get(ctx, 'Num')
        undoListPre = Get(ctx, 'Undo')
        undoList = concat(undoListPre, txHash)
        
        #gen key index
        keyOwner = concat('Owner', txHash)
        keyRequest = concat('Request', txHash)
        
        num = numPre + 1
        
        #storage input
        Put(ctx, keyOwner, addr)
        Put(ctx, keyRequest, request)
        Put(ctx, 'Undo', undoList)
        Put(ctx, 'Num', num)
        
        return txHash
     
    def setOracleOutcome(self, txHash: bytes, outcome):
        ctx = GetContext()
        
#        #check witness
#        isowner = CheckWitness(oracleAddr)
#        if not isowner:
#            return False
        
        #check if outcome isset already
        keyOutcome = concat('Outcome', txHash)
        keyIsSet = concat('IsSet', txHash)
        isSet = Get(ctx, keyIsSet)
        if isSet != 0:
            return "fail"
        
        Put(ctx, keyOutcome, outcome)
        Put(ctx, keyIsSet, 1)
        
        #remove from undo list
        undoListPre = Get(ctx, 'Undo')
        numPre = Get(ctx, 'Num')
        i = 0
        undoList = b''
        while i < numPre:
            txHashPre = undoListPre[32*i:32*i+32]
            if txHashPre != txHash:
                undoList = concat(undoList, txHashPre)
            i = i + 1
        
        num = numPre - 1
        Put(ctx, 'Undo', undoList)
        Put(ctx, 'Num', num)

        return "success"
        
    def setOracleCronOutcome(self, txHash: bytes, outcome):
        ctx = GetContext()
        
#        #check witness
#        isowner = CheckWitness(oracleAddr)
#        if not isowner:
#            return False
        
        keyOutcome = concat('Outcome', txHash)    
        Put(ctx, keyOutcome, outcome)

        return "success"
        
    #Returns if outcome is set
    #return Is outcome set?
    def isOracleOutcomeSet(self, txHash: bytes):
        ctx = GetContext()
        
        keyIsSet = concat('IsSet', txHash)
        isSet = Get(ctx, keyIsSet)
        return isSet
        
    #Returns outcome
    #return Outcome
    def getOracleOutcome(self, txHash: bytes, addr: bytes):
        ctx = GetContext()
        
#        #check if is the owner of the request
#        keyOwner = concat('Owner', txHash)
#        owner = Get(ctx, keyOwner)
#        if owner != addr:
#            return "Not Owner"
        
        keyOutcome = concat('Outcome', txHash)
        outcome = Get(ctx, keyOutcome)
        return outcome

        
        
        
        
        
        
        
        
        
        
        
        
        