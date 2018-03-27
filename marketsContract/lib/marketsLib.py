#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Mon Mar 12 07:00:17 2018

@author: root
"""
from boa.blockchain.vm.Neo.App import RegisterAppCall
from boa.blockchain.vm.Neo.Storage import GetContext, Get, Put
from boa.code.builtins import concat
from boa.blockchain.vm.Neo.Blockchain import GetHeight, GetHeader
from boa.blockchain.vm.Neo.Header import GetTimestamp, GetNextConsensus, GetHash
from boa.blockchain.vm.Neo.Runtime import CheckWitness

OracleContract = RegisterAppCall('c43daa0b3a784f479a2db0f9915222f6fb9d8b47', 'operation', 'args')

class Markets(object):
    
    #param ipfsHash Hash identifying off chain event description  
    def createCategoricalEvent(self, ipfsHash: bytes, spreadMultiplier: int, challengePeriod: int, challengeAmount: int, frontRunnerPeriod: int, endTime: int):
        ctx = GetContext()
        #createDecentralizedOracle
        operation = 'createDecentralizedOracle'
        args = [ipfsHash, spreadMultiplier, challengePeriod, challengeAmount, frontRunnerPeriod]
        isok = OracleContract(operation, args)
        if not isok:
            #print('createDecentralizedOracle Failed')
            return False
        
        keyEndTime = concat(b'EndTime', ipfsHash)
        Put(ctx, keyEndTime, endTime)
        return isok
        
    def betOnOutcome(self, ipfsHash: bytes, outcome: int, amount: int, owner: bytes, contractAddress: bytes):
        #check if witness
        isvoter = CheckWitness(owner)
        if not isvoter:
            #print('Must be voter to vote for outcome')
            return False
        
        #check if event ended
        ctx = GetContext()
        keyEndTime = concat(b'EndTime', ipfsHash)
        endTime = Get(ctx, keyEndTime)
        now = self.now()
        isendTime = (now > endTime)
        if isendTime:
            #print('Event has ended')
            return False
        
        #deposit
        isok = self.deposit(owner, contractAddress, amount)
        
        if not isok:
            #print('Deposit Failed')
            return False
        
        keyOutcomeAmountT = concat(owner, ipfsHash)
        keyOutcomeAmount = concat(outcome, keyOutcomeAmountT)
        outcomeAmount = Get(ctx, keyOutcomeAmount)
        outcomeAmount = outcomeAmount + amount
        Put(ctx, keyOutcomeAmount, outcomeAmount)
        
        keyTotalOutcomeAmounts = concat(outcome, ipfsHash)
        totalOutcomeAmounts = Get(ctx, keyTotalOutcomeAmounts)
        totalOutcomeAmounts = totalOutcomeAmounts + amount
        Put(ctx, keyTotalOutcomeAmounts, totalOutcomeAmounts)
        
        keyTotalAmount = concat(b'TotalAmount', ipfsHash)
        totalAmount = Get(ctx, keyTotalAmount)
        totalAmount = totalAmount + amount
        Put(ctx, keyTotalAmount, totalAmount)
        
        return True

    def isFinalOutcomeSet(self, ipfsHash: bytes):
        operation = 'isFinalOutcomeSet'
        args = [ipfsHash]
        isFinalOutcomeSet = OracleContract(operation, args)
        return isFinalOutcomeSet
        
    def redeemWinnings(self, ipfsHash: bytes, owner: bytes):
        isFinalOutcomeSet = self.isFinalOutcomeSet(ipfsHash)
        if not isFinalOutcomeSet:
            #print('Oracle is not ended yet!')
            return False
        ctx = GetContext()
        finalOutcome = self.getFinalOutcome(ipfsHash)
        
        keyTotalAmount = concat(b'TotalAmount', ipfsHash)
        totalAmount = Get(ctx, keyTotalAmount)
        
        keyOutcomeAmountT = concat(owner, ipfsHash)
        keyOutcomeAmount = concat(finalOutcome, keyOutcomeAmountT)
        ownerAmount = Get(ctx, keyOutcomeAmount)     
        
        keyTotalOutcomeAmounts = concat(finalOutcome, ipfsHash)
        frontRunnerAmount = Get(ctx, keyTotalOutcomeAmounts)        
        
        amount = (totalAmount * ownerAmount) / frontRunnerAmount
        
        return amount

    def getFinalOutcome(self, ipfsHash: bytes):
        operation = 'getFinalOutcome'
        args = [ipfsHash]
        finalOutcome = OracleContract(operation, args)
        return finalOutcome
        
    def now(self):
        height = GetHeight()
        header = GetHeader(height)
        now = GetTimestamp(header)
        return now 
        
    def deposit(from_addr: bytes, to_addr: bytes, amount: int):
        
        #appcall ont contract's transfer method
        
        return True        
        