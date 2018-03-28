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

from boa.blockchain.vm.Neo.Runtime import GetTrigger, CheckWitness
from boa.blockchain.vm.Neo.TriggerType import Application, Verification
from lib.decentralizedOracle import DecentralizedOracle
from boa.blockchain.vm.Neo.Action import RegisterAction

OWNER = b''

#register event
CreateRequest = RegisterAction('createRequest', 'result')
SetOutcome = RegisterAction('setOutcome', 'result')
IsOutcomeSet = RegisterAction('isOutcomeSet', 'result')
IsChallenged = RegisterAction('isChallenged', 'result')
IsChallengePeriodOver = RegisterAction('isChallengePeriodOver', 'result')
GetChallengeAmount = RegisterAction('getChallengeAmount', 'result')
ChallengeOutcome = RegisterAction('challengeOutcome', 'result')
VoteForOutcome = RegisterAction('voteForOutcome', 'result')
IsFrontRunnerPeriodOver = RegisterAction('isFrontRunnerPeriodOver', 'result')
IsChallengePeriodOver = RegisterAction('isChallengePeriodOver', 'result')
IsFinalOutcomeSet = RegisterAction('isFinalOutcomeSet', 'result')
GetFinalOutcome = RegisterAction('getFinalOutcome', 'result')
GetFrontRunner = RegisterAction('getFrontRunner', 'result')
RedeemWinnings = RegisterAction('redeemWinnings', 'result')

def Main(operation, args):
    """
    :param operation: str The name of the operation to perform
    :param args: list A list of arguments along with the operation
    :return:
        bytearray: The result of the operation
    """
    
    decentralizedOracle = DecentralizedOracle()
    
    trigger = GetTrigger()

    if trigger == Verification():

        # check if the invoker is the owner of this contract
        is_owner = CheckWitness(OWNER)
        
        # If owner, proceed
        #print("owner verify result")
        return is_owner
        
    elif trigger == Application():

        #print("doing application!")

        if operation == 'createRequest':
            if len(args) != 5:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            spreadMultiplier = args[1]
            challengePeriod = args[2]
            challengeAmount = args[3]
            frontRunnerPeriod = args[4]
            
            r = decentralizedOracle.createRequest(ID, spreadMultiplier, challengePeriod, challengeAmount, frontRunnerPeriod)
            CreateRequest(r)
            #print("createRequest done!")
            return r
        elif operation == 'setOutcome':
            if len(args) != 4:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            outcome = args[1]
            owner = args[2]
            contractAddress = args[3]
            r = decentralizedOracle.setOutcome(ID, outcome, owner, contractAddress)
            SetOutcome(r)
            #print('setOutcome done!')
            return r
        elif operation == 'isOutcomeSet':
            if len(args) != 1:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            r = decentralizedOracle.isOutcomeSet(ID)
            IsOutcomeSet(r)
            #print('isOutcomeSet done!')
            return r
#        elif operation == 'getOutcome':
#            if len(args) != 1:
#                #print('Incorrect Arg Length')
#                return False
#            ID = args[0]
#            r = decentralizedOracle.getOutcome(ID)
#            #print('getOutcome done!')
#            return r
        elif operation == 'isChallenged':
            if len(args) != 1:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            r = decentralizedOracle.isChallenged(ID)
            IsChallenged(r)
            #print('isChallenged done!')
            return r
        elif operation == 'isChallengePeriodOver':
            if len(args) != 1:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            r = decentralizedOracle.isChallengePeriodOver(ID)
            IsChallengePeriodOver(r)
            #print('isChallengePeriodOver done!')
            return r
#        elif operation == 'getOwner':
#            if len(args) != 1:
#                #print('Incorrect Arg Length')
#                return False
#            ID = args[0]
#            r = decentralizedOracle.getOwner(ID)
#            #print('getOwner done!')
#            return r
        elif operation == 'getChallengeAmount':
            if len(args) != 1:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            r = decentralizedOracle.getChallengeAmount(ID)
            GetChallengeAmount(r)
            #print('getChallengeAmount done!')
            return r
        elif operation == 'challengeOutcome':
            if len(args) != 4:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            newOutcome = args[1]
            challenger = args[2]
            contractAddress = args[3]
            r = decentralizedOracle.challengeOutcome(ID, newOutcome, challenger, contractAddress)
            ChallengeOutcome(r)
            #print('challengeOutcome done!')
            return r
        elif operation == 'voteForOutcome':
            if len(args) != 5:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            outcome = args[1]
            amount = args[2]
            voter = args[3]
            contractAddress = args[4]
            r = decentralizedOracle.voteForOutcome(ID, outcome, amount, voter, contractAddress)
            VoteForOutcome(r)
            #print('voteForOutcome done!')
            return r
        elif operation == 'isFrontRunnerPeriodOver':
            if len(args) != 1:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            r = decentralizedOracle.isFrontRunnerPeriodOver(ID)
            IsFrontRunnerPeriodOver(r)
            #print('isFrontRunnerPeriodOver done!')
            return r
        elif operation == 'isFinalOutcomeSet':
            if len(args) != 1:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            r = decentralizedOracle.isFinalOutcomeSet(ID)
            IsFinalOutcomeSet(r)
            #print('isFinalOutcomeSet done!')
            return r
        elif operation == 'getFinalOutcome':
            if len(args) != 1:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            r = decentralizedOracle.getFinalOutcome(ID)
            GetFinalOutcome(r)
            #print('getFinalOutcome done!')
            return r
        elif operation == 'getFrontRunner':
            if len(args) != 1:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            r = decentralizedOracle.getFrontRunner(ID)
            GetFrontRunner(r)
            #print('getFrontRunner done!')
            return r
        elif operation == 'redeemWinnings':
            if len(args) != 2:
                #print('Incorrect Arg Length')
                return False
            ID = args[0]
            owner = args[1]
            r = decentralizedOracle.redeemWinnings(ID, owner)
            RedeemWinnings(r)
            #print('redeemWinnings done!')
            return r
        return False
    return False
    
    
    
    
    
    
    
    
    
    
    