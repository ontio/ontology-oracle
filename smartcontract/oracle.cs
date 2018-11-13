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

        public static object Main(string operation, params object[] args)
        {
            switch (operation)
            {
                case "CreateOracleRequest":
                    return CreateOracleRequest((string)args[0], (byte[])args[1]);
                case "SetOracleOutcome":
                    return SetOracleOutcome((byte[])args[0], (byte[])args[1]);
                case "GetOracleOutcome":
                    return GetOracleOutcome((byte[])args[0]);
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

        public static bool SetOracleOutcome(byte[] txHash, byte[] result)
        {
            //TODO: check witness
            if (!Runtime.CheckWitness(admin))
            {
                Runtime.Notify("Checkwitness failed.");
                return false;
            }

            Storage.Put(Storage.CurrentContext, txHash, result);

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
            Runtime.Notify(v);
            return v;
        }

    }
}