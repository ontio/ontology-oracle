#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Mon Mar 12 09:24:54 2018

@author: root
"""

from boa.blockchain.vm.Neo.Runtime import GetTrigger, CheckWitness
from boa.blockchain.vm.Neo.TriggerType import Application, Verification
from lib.marketsLib import Markets
from boa.blockchain.vm.Neo.Action import RegisterAction

OWNER = b''

#register event
CreateCategoricalEvent = RegisterAction('createCategoricalEvent', 'result')
BetOnOutcome = RegisterAction('betOnOutcome', 'result')
IsFinalOutcomeSet = RegisterAction('isFinalOutcomeSet', 'result')
GetFinalOutcome = RegisterAction('getFinalOutcome', 'result')
RedeemWinnings = RegisterAction('redeemWinnings', 'result')

def Main(operation, args):
    """
    :param operation: str The name of the operation to perform
    :param args: list A list of arguments along with the operation
    :return:
        bytearray: The result of the operation
    """
    
    markets = Markets()
    
    trigger = GetTrigger()

    if trigger == Verification():

        # check if the invoker is the owner of this contract
        is_owner = CheckWitness(OWNER)
        
        # If owner, proceed
        #print("owner verify result")
        return is_owner
        
    elif trigger == Application():
        
        if operation == 'createCategoricalEvent':
            if len(args) != 6:
                #print('Incorrect Arg Length')
                return False
            ipfsHash = args[0]
            spreadMultiplier = args[1]
            challengePeriod = args[2]
            challengeAmount = args[3]
            frontRunnerPeriod = args[4]
            endTime = args[5]
            
            r = markets.createCategoricalEvent(ipfsHash, spreadMultiplier, challengePeriod, challengeAmount, frontRunnerPeriod, endTime)
            CreateCategoricalEvent(r)
            #print("createCategoricalEvent done!")
            return r
        elif operation == 'betOnOutcome':
            if len(args) != 5:
                #print('Incorrect Arg Length')
                return False
            ipfsHash = args[0]
            outcome = args[1]
            amount = args[2]
            owner = args[3]
            contractAddress = args[4]
            r = markets.betOnOutcome(ipfsHash, outcome, amount, owner, contractAddress)
            BetOnOutcome(r)
            #print('betOnOutcome done!')
            return r
        elif operation == 'isFinalOutcomeSet':
            if len(args) != 1:
                #print('Incorrect Arg Length')
                return False
            ipfsHash = args[0]
            r = markets.isFinalOutcomeSet(ipfsHash)
            IsFinalOutcomeSet(r)
            #print('isFinalOutcomeSet done!')
            return r
        elif operation == 'getFinalOutcome':
            if len(args) != 1:
                #print('Incorrect Arg Length')
                return False
            ipfsHash = args[0]
            r = markets.getFinalOutcome(ipfsHash)
            GetFinalOutcome(r)
            #print('getFinalOutcome done!')
            return r
        elif operation == 'redeemWinnings':
            if len(args) != 2:
                #print('Incorrect Arg Length')
                return False
            ipfsHash = args[0]
            owner = args[1]
            r = markets.redeemWinnings(ipfsHash, owner)
            RedeemWinnings(r)
            #print('redeemWinnings done!')
            return r
        return False
    return False
        
        
        
        
        
        
        
        