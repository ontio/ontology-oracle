package utils

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ontio/ontology/common"
	"github.com/ontio/ontology/core/types"
)

const (
	RPC_GET_STORAGE              = "getstorage"
	RPC_GET_SMART_CONTRACT_EVENT = "getsmartcodeevent"
	RPC_SEND_TRANSACTION         = "sendrawtransaction"
)

//JsonRpc version
const JSON_RPC_VERSION = "2.0"

//Http rpcClient for ontology rpc api
type RpcClient struct {
	qid        uint64
	addr       string
	httpClient *http.Client
}

//JsonRpcRequest object in rpc
type JsonRpcRequest struct {
	Version string        `json:"jsonrpc"`
	Id      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

//JsonRpcResponse object response for JsonRpcRequest
type JsonRpcResponse struct {
	Error  int64           `json:"error"`
	Desc   string          `json:"desc"`
	Result json.RawMessage `json:"result"`
}

//NewRpcClient return RpcClient instance
func NewRpcClient() *RpcClient {
	return &RpcClient{
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   5,
				DisableKeepAlives:     false, //enable keepalive
				IdleConnTimeout:       time.Second * 300,
				ResponseHeaderTimeout: time.Second * 300,
			},
			Timeout: time.Second * 300, //timeout for http response
		},
	}
}

//GetStorage return smart contract storage item.
//addr is smart contact address
//key is the key of value in smart contract
func (rpc *RpcClient) GetStorage(smartContractAddress common.Address, key []byte) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := smartContractAddress.Serialize(buf)
	if err != nil {
		return nil, fmt.Errorf("Address Serialize error:%s", err)
	}
	hexString := hex.EncodeToString(buf.Bytes())
	data, err := rpc.sendRpcRequest(RPC_GET_STORAGE, []interface{}{hexString, hex.EncodeToString(key)})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}
	hexData := ""
	err = json.Unmarshal(data, &hexData)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal error:%s", err)
	}
	value, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, fmt.Errorf("hex.DecodeString error:%s", err)
	}
	return value, nil
}

func (rpc *RpcClient) getQid() string {
	return fmt.Sprintf("%d", atomic.AddUint64(&rpc.qid, 1))
}

//sendRpcRequest send Rpc request to ontology
func (rpc *RpcClient) sendRpcRequest(method string, params []interface{}) ([]byte, error) {
	rpcReq := &JsonRpcRequest{
		Version: JSON_RPC_VERSION,
		Id:      rpc.getQid(),
		Method:  method,
		Params:  params,
	}
	data, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, fmt.Errorf("JsonRpcRequest json.Marsha error:%s", err)
	}
	resp, err := rpc.httpClient.Post(rpc.addr, "application/json", strings.NewReader(string(data)))
	if err != nil {
		return nil, fmt.Errorf("http post request:%s error:%s", data, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read rpc response body error:%s", err)
	}

	rpcRsp := &JsonRpcResponse{}
	err = json.Unmarshal(body, rpcRsp)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal JsonRpcResponse:%s error:%s", body, err)
	}
	if rpcRsp.Error != 0 {
		return nil, fmt.Errorf("sendRpcRequest error code:%d desc:%s", rpcRsp.Error, rpcRsp.Desc)
	}
	return rpcRsp.Result, nil
}

//SmartContactEvent object for event of transaction
type SmartContactEvent struct {
	Address interface{} `json:"CodeHash"`
	States  []interface{}
	TxHash  interface{}
}

//GetSmartContractEvent return smart contract event execute by invoke transaction.
func (rpc *RpcClient) GetSmartContractEvent(txHash common.Uint256) ([]*SmartContactEvent, error) {
	return rpc.GetSmartContractEventWithHexString(hex.EncodeToString(txHash.ToArray()))
}

//GetSmartContractEvent return smart contract event execute by invoke transaction by hex string code
func (rpc *RpcClient) GetSmartContractEventWithHexString(txHash string) ([]*SmartContactEvent, error) {
	data, err := rpc.sendRpcRequest(RPC_GET_SMART_CONTRACT_EVENT, []interface{}{txHash})
	if err != nil {
		return nil, fmt.Errorf("sendRpcRequest error:%s", err)
	}
	events := make([]*SmartContactEvent, 0)
	err = json.Unmarshal(data, &events)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal SmartContactEvent:%s error:%s", data, err)
	}
	return events, nil
}

//SetAddress set rpc server address. Simple http://localhost:20336
func (rpc *RpcClient) SetAddress(addr string) *RpcClient {
	rpc.addr = addr
	return rpc
}

//SendRawTransaction send a transaction to ontology network, and return hash of the transaction
func (this *RpcClient) SendRawTransaction(tx *types.Transaction) error {
	var buffer bytes.Buffer
	err := tx.Serialize(&buffer)
	if err != nil {
		return fmt.Errorf("serialize error:%s", err)
	}
	txData := hex.EncodeToString(buffer.Bytes())
	_, err = this.sendRpcRequest(RPC_SEND_TRANSACTION, []interface{}{txData})
	if err != nil {
		return fmt.Errorf("sendRpcRequest error:%s", err)
	}
	return nil
}
