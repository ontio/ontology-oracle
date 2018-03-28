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

from boa.blockchain.vm.Neo.Blockchain import GetHeight, GetHeader
from boa.blockchain.vm.Neo.Header import GetTimestamp, GetNextConsensus, GetHash
from boa.code.builtins import concat
from boa.blockchain.vm.Neo.Runtime import CheckWitness
from boa.blockchain.vm.Neo.Storage import GetContext, Get, Put

class DecentralizedOracle(object):
    
    def createRequest(self, ID: bytes, spreadMultiplier: int, challengePeriod: int, challengeAmount: int, frontRunnerPeriod: int):
        
        outcome = 0
        outcomeSetTimestamp = 0
        frontRunnerSetTimestamp = 0
        
        #gen key index
        keyOutcome = concat(b'Outcome',ID)
        keyOutcomeSetTimestamp = concat(b'OutcomeSetTimestamp',ID)
        keyChallengePeriod = concat(b'ChallengePeriod',ID)
        keyChallengeAmount = concat(b'ChallengeAmount',ID)
        keySpreadMultiplier = concat(b'SpreadMultiplier',ID)
        keyFrontRunnerPeriod = concat(b'FrontRunnerPeriod',ID)
        keyFrontRunnerSetTimestamp = concat(b'FrontRunnerSetTimestamp',ID)
        
        #storage input
        ctx = GetContext()
        Put(ctx, keyOutcome, outcome)
        Put(ctx, keyOutcomeSetTimestamp, outcomeSetTimestamp)
        Put(ctx, keyChallengePeriod, challengePeriod)
        Put(ctx, keyChallengeAmount, challengeAmount)
        Put(ctx, keySpreadMultiplier, spreadMultiplier)
        Put(ctx, keyFrontRunnerPeriod, frontRunnerPeriod)
        Put(ctx, keyFrontRunnerSetTimestamp, frontRunnerSetTimestamp)
        
        return True
     
    #Sets event outcome
    #param outcome Event outcome
    def setOutcome(self, ID: bytes, outcome: int, owner: bytes, contractAddress: bytes):
        ctx = GetContext()
        
        #check if witness
        isowner = CheckWitness(owner)
        if not isowner:
            #print('Must be voter to vote for outcome')
            return False
        
        #check if outcome isset already
        keyOutcomeSetTimestamp = concat(b'OutcomeSetTimestamp',ID)
        outcomeSetTimestamp = Get(ctx, keyOutcomeSetTimestamp)
        if outcomeSetTimestamp != 0:
            #print('Outcome is already set')
            return False
        
        challengeAmount = self.getChallengeAmount(ID)
        
        #transfer challengeAmount to this contract account address
        #deposit
        isok = self.deposit(owner, contractAddress, challengeAmount)
        
        if not isok:
            #print('Deposit Failed')
            return False
        
        keyOutcomeAmountT = concat(owner,ID)
        keyOutcomeAmount = concat(outcome, keyOutcomeAmountT)
        Put(ctx, keyOutcomeAmount, challengeAmount)
        keyTotalOutcomeAmounts = concat(outcome, ID)
        Put(ctx, keyTotalOutcomeAmounts, challengeAmount)
        keyTotalAmount = concat(b'TotalAmount',ID)
        Put(ctx, keyTotalAmount, challengeAmount)
            
        #change outcomeSetTimestamp and storage input
        now = self.now()
        Put(ctx, keyOutcomeSetTimestamp, now)
        
        #change outcome and storage input
        outcome = outcome
        keyOutcome = concat(b'Outcome',ID)
        Put(ctx, keyOutcome, outcome)
        
        return True
        
    #Returns if outcome is set
    #return Is outcome set?
    def isOutcomeSet(self, ID: bytes):
        ctx = GetContext()
        
        keyOutcomeSetTimestamp = concat(b'OutcomeSetTimestamp',ID)
        outcomeSetTimestamp = Get(ctx, keyOutcomeSetTimestamp)
        isSet = (outcomeSetTimestamp != 0)
        return isSet
        
#    #Returns outcome
#    #return Outcome
#    def getOutcome(self, ID: bytes):
#        ctx = GetContext()
#        
#        keyOutcome = concat(b'Outcome',ID)
#        outcome = Get(ctx, keyOutcome)
#        return outcome
    
    #Checks if outcome was challenged
    #return Is challenged?
    def isChallenged(self, ID: bytes):
        ctx = GetContext()
        
        keyFrontRunnerSetTimestamp = concat(b'FrontRunnerSetTimestamp',ID)
        frontRunnerSetTimestamp = Get(ctx, keyFrontRunnerSetTimestamp)
        
        isChallenged = (frontRunnerSetTimestamp != 0)
        return isChallenged
       
    #check if challengePeriod is over
    def isChallengePeriodOver(self, ID: bytes):
        ctx = GetContext()
        
        keyOutcomeSetTimestamp = concat(b'OutcomeSetTimestamp',ID)
        outcomeSetTimestamp = Get(ctx, keyOutcomeSetTimestamp)
        keyChallengePeriod = concat(b'ChallengePeriod',ID)
        challengePeriod = Get(ctx, keyChallengePeriod)
        
        now = self.now()
        isChallengePeriodOver1 = (now - outcomeSetTimestamp > challengePeriod)
        isChallengePeriodOver2 = (outcomeSetTimestamp != 0)
        if not isChallengePeriodOver1:
            return False
        if not isChallengePeriodOver2:
            return False
        else:
            return True
    
    #get contract owner
#    def getOwner(self, ID: bytes):
#        ctx = GetContext()
#        
#        keyOwner = concat(b'Owner',ID)
#        owner = Get(ctx, keyOwner)
#        return owner
        
    #get challengAmount
    def getChallengeAmount(self, ID: bytes):
        ctx = GetContext()
        
        keyChallengeAmount = concat(b'ChallengeAmount',ID)    
        challengeAmount = Get(ctx, keyChallengeAmount)
        return challengeAmount
        
    #challenge current outcome with a new outcome
    def challengeOutcome(self, ID: bytes, newOutcome: int, challenger: bytes, contractAddress: bytes): 
        ctx = GetContext()
        
        #check witness
        ischallenger = CheckWitness(challenger)
        if not ischallenger:
            #print('Must be challenger to challenge')
            return False
        
        #check if outcome is set
        isSet = self.isOutcomeSet(ID)
        if not isSet:
            #print('Outcome is not set yet')
            return False
        
        #check if is already challenged
        isChallenged = self.isChallenged(ID)
        if isChallenged:
            #print('Oracle is already be challenged')
            return False
            
        #check if challengePeriod over
        isChallengePeriodOver = self.isChallengePeriodOver(ID)
        if isChallengePeriodOver:
            #print('Challenge period is over')
            return False
           
        challengeAmount = self.getChallengeAmount(ID)
        
        #transfer challengeAmount to this contract account address
        #deposit
        isok = self.deposit(challenger, contractAddress, challengeAmount)
        
        if not isok:
            #print('Deposit Failed')
            return False
        
        keyOutcomeAmountT = concat(challenger, ID)
        keyOutcomeAmount = concat(newOutcome, keyOutcomeAmountT)
        Put(ctx, keyOutcomeAmount, challengeAmount)
        
        keyTotalOutcomeAmounts = concat(newOutcome, ID)
        Put(ctx, keyTotalOutcomeAmounts, challengeAmount)
        
        keyTotalAmount = concat(b'TotalAmount', ID)
        totalAmount = Get(ctx, keyTotalAmount)
        totalAmount = totalAmount + challengeAmount
        Put(ctx, keyTotalAmount, totalAmount)
        
        keyFrontRunner = concat(b'FrontRunner', ID)
        Put(ctx, keyFrontRunner, newOutcome)
        
        now = self.now()
        keyFrontRunnerSetTimestamp = concat(b'FrontRunnerSetTimestamp',ID)
        Put(ctx, keyFrontRunnerSetTimestamp, now)
            
        return True
        
    #vote for outcome that user want to bet    
    def voteForOutcome(self, ID: bytes, outcome: int, amount: int, voter: bytes, contractAddress: bytes):
        ctx = GetContext()
        
        #check if witness
        isvoter = CheckWitness(voter)
        if not isvoter:
            #print('Must be voter to vote for outcome')
            return False 
        
        #check if Oracle is challenged
        isChallenged = self.isChallenged(ID)
        if not isChallenged:
            #print('Oracle is not be challenged yet')
            return False
        
        #check if front runner period is not over yet
        isFrontRunnerPeriodOver = self.isFrontRunnerPeriodOver(ID)
        if isFrontRunnerPeriodOver:
            #print('Front runner period is over, oracle is done')
            return False
            
        #read param from storage
        keySpreadMultiplier = concat(b'SpreadMultiplier',ID)
        spreadMultiplier = Get(ctx, keySpreadMultiplier)
        keyTotalAmount = concat(b'TotalAmount', ID)
        totalAmount = Get(ctx, keyTotalAmount)
        keyTotalOutcomeAmounts = concat(outcome, ID)
        totalOutcomeAmounts = Get(ctx, keyTotalOutcomeAmounts)
        
        maxAmount = (totalAmount - totalOutcomeAmounts) * spreadMultiplier
        if maxAmount > totalOutcomeAmounts:
            maxAmount = maxAmount - totalOutcomeAmounts
        else:
            maxAmount = 0
        if amount > maxAmount:
            amount = maxAmount
        
        #deposit
        isok = self.deposit(voter, contractAddress, amount)
        
        if not isok:
            #print('Deposit Failed')
            return False
        
        keyOutcomeAmountT = concat(voter, ID)
        keyOutcomeAmount = concat(outcome, keyOutcomeAmountT)
        outcomeAmount = Get(ctx, keyOutcomeAmount)
        outcomeAmount = outcomeAmount + amount
        Put(ctx, keyOutcomeAmount, outcomeAmount)
        
        keyTotalOutcomeAmounts = concat(outcome, ID)
        totalOutcomeAmounts = Get(ctx, keyTotalOutcomeAmounts)
        totalOutcomeAmounts = totalOutcomeAmounts + amount
        Put(ctx, keyTotalOutcomeAmounts, totalOutcomeAmounts)
        
        keyTotalAmount = concat(b'TotalAmount', ID)
        totalAmount = Get(ctx, keyTotalAmount)
        totalAmount = totalAmount + amount
        Put(ctx, keyTotalAmount, totalAmount)
        
        #check if front runner changes
        keyFrontRunner = concat(b'FrontRunner', ID)
        frontRunner = Get(ctx, keyFrontRunner)
        keyTotalOutcomeAmountsF = concat(frontRunner, ID)
        totalOutcomeAmountsF = Get(ctx, keyTotalOutcomeAmountsF)
        
        isfrontrunnerChanges1 = (outcome != frontRunner)
        isfrontrunnerChanges2 = (totalOutcomeAmounts > totalOutcomeAmountsF)
        if not isfrontrunnerChanges1:
            return False
        if not isfrontrunnerChanges2:
            return False
        else:
            Put(ctx, keyFrontRunner, outcome)
            now = self.now()
            keyFrontRunnerSetTimestamp = concat(b'FrontRunnerSetTimestamp',ID)
            Put(ctx, keyFrontRunnerSetTimestamp, now)
            return True
    
    def getFrontRunner(self, ID: bytes):
        ctx = GetContext()
        
        keyFrontRunner = concat(b'FrontRunner', ID)
        frontRunner = Get(ctx, keyFrontRunner)
        return frontRunner
        
    def isFrontRunnerPeriodOver(self, ID: bytes):
        ctx = GetContext()
        
        keyFrontRunnerSetTimestamp = concat(b'FrontRunnerSetTimestamp',ID)
        frontRunnerSetTimestamp = Get(ctx, keyFrontRunnerSetTimestamp)
        keyFrontRunnerPeriod = concat(b'FrontRunnerPeriod',ID)
        frontRunnerPeriod = Get(ctx, keyFrontRunnerPeriod)
        now = self.now()
        
        isFrontRunnerPeriodOver1 = (frontRunnerSetTimestamp != 0)
        isFrontRunnerPeriodOver2 = (now - frontRunnerSetTimestamp > frontRunnerPeriod)
        if not isFrontRunnerPeriodOver1:
            return False
        if not isFrontRunnerPeriodOver2:
            return False
        else:
            return True
        
    #Returns if winning outcome is set
    def isFinalOutcomeSet(self, ID: bytes):
        isChallengePeriodOver = self.isChallengePeriodOver(ID)
        isChallenged = self.isChallenged(ID)
        isFrontRunnerPeriodOver = self.isFrontRunnerPeriodOver(ID)
        
        #isFinalOutcomeSet = isChallengePeriodOver & (not isChallenged) | isFrontRunnerPeriodOver

        if isFrontRunnerPeriodOver:
            return True
        elif not isChallengePeriodOver:
            return False
        elif isChallenged:
            return False
        return True
        
    #Returns winning outcom
    def getFinalOutcome(self, ID: bytes):
        isFinalOutcomeSet = self.isFinalOutcomeSet(ID)
        if not isFinalOutcomeSet:
            #print('Oracle is not ended yet!')
            return False
        ctx = GetContext()
        
        isFrontRunnerPeriodOver = self.isFrontRunnerPeriodOver(ID)
        keyFrontRunner = concat(b'FrontRunner', ID)
        frontRunner = Get(ctx, keyFrontRunner)
        keyOutcome = concat(b'Outcome',ID)
        outcome = Get(ctx, keyOutcome)
        
        if isFrontRunnerPeriodOver:
            return frontRunner
        return outcome
            
    def now(self):
        height = GetHeight()
        header = GetHeader(height)
        now = GetTimestamp(header)
        return now 

    def deposit(from_addr: bytes, to_addr: bytes, amount: int):
        
        #appcall ont contract's transfer method
        
        return True
    
    #after method must change outcome amount to 0
    def redeemWinnings(self, ID: bytes, owner: bytes):
        isFinalOutcomeSet = self.isFinalOutcomeSet(ID)
        if not isFinalOutcomeSet:
            #print('Oracle is not ended yet!')
            return False
        ctx = GetContext()
        finalOutcome = self.getFinalOutcome(ID)
        
        keyTotalAmount = concat(b'TotalAmount', ID)
        totalAmount = Get(ctx, keyTotalAmount)
        
        keyOutcomeAmountT = concat(owner, ID)
        keyOutcomeAmount = concat(finalOutcome, keyOutcomeAmountT)
        ownerAmount = Get(ctx, keyOutcomeAmount)     
        
        keyTotalOutcomeAmounts = concat(finalOutcome, ID)
        frontRunnerAmount = Get(ctx, keyTotalOutcomeAmounts)        
        
        amount = (totalAmount * ownerAmount) / frontRunnerAmount
        
        return amount
        
        
        
        
        
        
        
        
        
        
        
        
        