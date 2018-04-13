#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Created on Tue Apr  3 07:47:58 2018

@author: root
"""

from boa.blockchain.vm.Neo.Transaction import GetHash, GetUnspentCoins
from boa.blockchain.vm.System.ExecutionEngine import GetScriptContainer

def getTxHash():
    sc = GetScriptContainer()
    txHash = GetHash(sc)
    return txHash