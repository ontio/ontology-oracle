using Ont.SmartContract.Framework.Services.Ont;
using Ont.SmartContract.Framework;
using System;
using System.ComponentModel;
using Ont.SmartContract.Framework.Services.System;
using Helper = Ont.SmartContract.Framework.Helper;

namespace Ont.SmartContract
{
    public class Oracle : Framework.SmartContract
    {
        public static readonly byte[] admin = "AMAx993nE6NEqZjwBssUfopxnnvTdob9ij".ToScriptHash();

        public class Result
        {
            public byte[] data;
            public string status;
            public string errMessage;
        }

        public static object Main(string operation, params object[] args)
        {
            switch (operation)
            {
                case "CreateOracleRequest":
                    return CreateOracleRequest((string)args[0], (byte[])args[1]);
                case "SetOracleOutcome":
                    return SetOracleOutcome((byte[])args[0], (byte[])args[1], (string)args[2], (string)args[3]);
                case "GetOracleOutcome":
                    return GetOracleOutcome((byte[])args[0]);
                case "Migrate":
                    return Migrate((byte[])args[0], (bool)args[1], (string)args[2], (string)args[3], (string)args[4], (string)args[5], (string)args[6]);
                default:
                    return false;
            }
        }

        public static bool CreateOracleRequest(string request, byte[] address)
        {
            //TODO: check request format

            //TODO: transfer ong

            Transaction tx = (Transaction)ExecutionEngine.ScriptContainer;
            byte[] txHash = tx.Hash;
            //Runtime.Notify("txHash is :", txHash);

            Map<byte[], string> undoRequest = new Map<byte[], string>();
            byte[] v = Storage.Get(Storage.CurrentContext, "UndoRequest");
            if (v.Length != 0) {
                undoRequest = (Map<byte[], string>)Helper.Deserialize(v);
            }
            undoRequest[txHash] = request;

            byte[] b = Helper.Serialize(undoRequest);
            Storage.Put(Storage.CurrentContext, "UndoRequest", b);
            Runtime.Notify("CreateOracleRequest Done");
            return true;
        }

        public static bool SetOracleOutcome(byte[] txHash, byte[] data, string status, string errMessage)
        {
            //check witness
            if (!Runtime.CheckWitness(admin))
            {
                Runtime.Notify("Checkwitness failed.");
                return false;
            }

            Result result = new Result();
            result.data = data;
            result.status = status;
            result.errMessage= errMessage;
            byte[] r = Helper.Serialize(result);
            Storage.Put(Storage.CurrentContext, txHash, r);

            //remove txHash from undoRequest map
            Map<byte[], string> undoRequest = new Map<byte[], string>();
            byte[] v = Storage.Get(Storage.CurrentContext, "UndoRequest");
            if (v.Length != 0) {
                undoRequest = (Map<byte[], string>)Helper.Deserialize(v);
            }
            undoRequest.Remove(txHash);

            byte[] b = Helper.Serialize(undoRequest);
            Storage.Put(Storage.CurrentContext, "UndoRequest", b);
            Runtime.Notify("SetOracleOutcome Done");
            return true;
        }

        public static byte[] GetOracleOutcome(byte[] txHash)
        {
            byte[] v = Storage.Get(Storage.CurrentContext, txHash);

            //TODO: remove txHash from results
            Runtime.Notify("Get oracle outcome:", v);
            return v;
        }

        public static bool Migrate(byte[] code, bool need_storage, string name, string version, string author, string email, string description)
        {
            //check if owner of the contract
            if (!Runtime.CheckWitness(admin))
            {
                Runtime.Notify("Checkwitness failed.");
                return false;
            }

            Ont.SmartContract.Framework.Services.Ont.Contract.Migrate(code, need_storage, name, version, author, email, description);
            return true;
        }

    }
}