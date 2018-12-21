from boa.interop.Ontology.Contract import Migrate
from boa.interop.System.Storage import GetContext, Get, Put, Delete
from boa.interop.System.Runtime import CheckWitness, GetTime, Notify, Serialize, Deserialize
from boa.interop.System.ExecutionEngine import GetExecutingScriptHash, GetScriptContainer
from boa.interop.Ontology.Native import Invoke
from boa.interop.System.Transaction import GetTransactionHash
from boa.builtins import *
from ontology.interop.Ontology.Runtime import Base58ToAddress

######################### Global info ########################
Admin = Base58ToAddress('AMAx993nE6NEqZjwBssUfopxnnvTdob9ij')
ONGAddress = bytearray(b'\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x02')
Fee = "Fee"
UndoRequestKey = "UndoRequest"

"""
https://github.com/ONT-Avocados/python-template/blob/master/libs/Utils.py
"""
def Revert():
    """
    Revert the transaction. The opcodes of this function is `09f7f6f5f4f3f2f1f000f0`,
    but it will be changed to `ffffffffffffffffffffff` since opcode THROW doesn't
    work, so, revert by calling unused opcode.
    """
    raise Exception(0xF1F1F2F2F3F3F4F4)


"""
https://github.com/ONT-Avocados/python-template/blob/master/libs/SafeCheck.py
"""
def Require(condition):
    """
	If condition is not satisfied, return false
	:param condition: required condition
	:return: True or false
	"""
    if not condition:
        Revert()
    return True

def RequireWitness(witness):
    """
	Checks the transaction sender is equal to the witness. If not
	satisfying, revert the transaction.
	:param witness: required transaction sender
	:return: True if transaction sender or revert the transaction.
	"""
    Require(CheckWitness(witness))
    return True

def Main(operation, args):
    if operation == "CreateOracleRequest":
        if len(args) != 2:
            return False
        request = args[0]
        address = args[1]
        CreateOracleRequest(request, address);
    if operation == "SetOracleOutcome":
        if len(args) != 4:
            return False
        txHash = args[0]
        data = args[1]
        status = args[2]
        errMessage = args[3]
        return SetOracleOutcome(txHash, data, status, errMessage)
    if operation == "GetOracleOutcome":
        if len(args) != 1:
            return False
        txHash = args[0]
        return GetOracleOutcome(txHash)
    if operation == "SetFee":
        if len(args) != 1:
            return False
        fee = args[0]
        return SetFee(fee)
    if operation == "MigrateContract":
        if len(args) !=7:
            return False
        code = args[0]
        needStorage = args[1]
        name = args[2]
        version = args[3]
        author = args[4]
        email = args[5]
        description = args[6]
        return MigrateContract(code, needStorage, name, version, author, email, description)

def CreateOracleRequest(request, address):
    #check witness
    RequireWitness(address)

    fee = Get(GetContext(), Fee)
    #transfer ong to oracle admin address
    res = TransferONG(address, Admin, fee)
    if res == False:
        Notify(["transferONG Error"])
        return False

    #get transaction hash
    txHash = GetTransactionHash(GetScriptContainer())

    #update undoRequestMap
    undoRequestMap = GetUndoRequestMap()
    undoRequestMap[txHash] = request
    b = Serialize(undoRequestMap)
    Put(GetContext(), UndoRequestKey, b)
    Notify(["CreateOracleRequest Done", txHash, request])
    return True

def SetOracleOutcome(txHash, data, status, errMessage):
    #check witness
    RequireWitness(Admin)

    #get undoRequest map
    undoRequestMap = GetUndoRequestMap()

    #TODO : check if key exist

    #put result into storage
    result = state(data, status, errMessage)
    r = Serialize(result)
    Put(GetContext(), txHash, r)

    #remove txHash from undoRequest map
    undoRequestMap.remove(txHash)
    b = Serialize(undoRequestMap)
    Put(GetContext(), UndoRequestKey, b)
    Notify(["SetOracleOutcome Done", txHash, status, errMessage])
    return True

def GetOracleOutcome(txHash):
    v = Get(GetContext(), txHash)
    Notify(["Get oracle outcome", v])
    return v

def SetFee(fee):
    RequireWitness(Admin)
    Require(fee >= 0)
    Put(GetContext(), Fee, fee)

def MigrateContract(code, needStorage, name, version, author, email, description):
    RequireWitness(Admin)
    res = Migrate(code, needStorage, name, version, author, email, description)
    Require(res)
    Notify(["Migrate Contract successfully", Admin, GetTime()])
    return True

def GetUndoRequestMap():
    undoRequestMap = {}
    v = Get(GetContext(), UndoRequestKey)
    if len(v) != 0:
        undoRequestMap = Deserialize(v)
    return undoRequestMap

def TransferONG(fromAcct, toAcct, ongAmount):
    """
    transfer ONT
    :param fromacct:
    :param toacct:
    :param amount:
    :return:
    """
    RequireWitness(fromAcct)
    param = state(fromAcct, toAcct, ongAmount)
    res = Invoke(0, ONGAddress, 'transfer', [param])
    if res and res == b'\x01':
        return True
    else:
        return False