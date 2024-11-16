package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/rooch-network/rooch-go-sdk/crypto"
	"github.com/rooch-network/rooch-go-sdk/types/transactions"
	"math/big"
	"strings"

	"github.com/rooch-network/rooch-go-sdk/types"

	"github.com/pkg/errors"
	"github.com/rooch-network/rooch-go-sdk/client/jsonrpc"
)

const DefaultMaxGasAmount = 10000000
const GasTokenCode = "0x3::gas_coin::RGas"
const DefaultTransactionExpirationSeconds = 2 * 60 * 60
const emptyString = ""

type Signer crypto.Signer

type RoochClient struct {
	client  *jsonrpc.Client // jsonrpc client to use for requests
	url     string
	chainID uint64 // Chain ID of the network e.g. 2 for Testnet
	//baseUrl *url.URL
}

func NewRoochClient(url string) (RoochClient, error) {
	client, err := jsonrpc.NewClient(url)
	if err != nil {
		return RoochClient{}, err
	}
	return RoochClient{
		client: client,
		url:    url,
	}, nil
}

func (c *RoochClient) Call(context context.Context, serviceMethod string, reply interface{}, args interface{}) error {
	client, err := jsonrpc.NewClient(c.url)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("call method %s err: ", serviceMethod))
	}

	err = client.Call(context, serviceMethod, reply, args)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("call method %s err: ", serviceMethod))
	}

	return nil
}

func (c *RoochClient) GetEvents(context context.Context, eventFilter *EventFilter) ([]Event, error) {
	var result []Event
	params := []interface{}{eventFilter}
	err := c.Call(context, "chain.get_events", &result, params)
	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_events error ")
	}
	return result, nil
}

func (c *RoochClient) GetTransactionByHash(context context.Context, transactionHash string) (*Transaction, error) {
	result := &Transaction{}
	params := []string{transactionHash}
	err := c.Call(context, "chain.get_transaction", result, params)

	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_transaction error ")
	}

	return result, nil
}

func (c *RoochClient) GetTransactionInfoByHash(context context.Context, transactionHash string) (*TransactionInfo, error) {
	result := &TransactionInfo{}
	params := []string{transactionHash}
	err := c.Call(context, "chain.get_transaction_info", result, params)

	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_transaction_info error ")
	}

	return result, nil
}

func (c *RoochClient) GetTransactionEventByHash(context context.Context, transactionHash string) ([]Event, error) {
	var result []Event
	params := []string{transactionHash}
	err := c.Call(context, "chain.get_events_by_txn_hash", &result, params)

	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_events_by_txn_hash error ")
	}

	return result, nil
}

func (c *RoochClient) GetBlockByHash(context context.Context, blockHash string) (*Block, error) {
	result := &Block{}
	params := []string{blockHash}
	err := c.Call(context, "chain.get_block_by_hash", result, params)

	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_block_by_hash ")
	}

	return result, nil
}

func (c *RoochClient) GetBlockByNumber(context context.Context, number int) (*Block, error) {
	result := &Block{}
	params := []int{number}
	err := c.Call(context, "chain.get_block_by_number", result, params)

	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_block_by_number ")
	}

	return result, nil
}

//func (c *RoochClient) GetBalanceOfStc(context context.Context, address string) (*big.Int, error) {
//	ls, err := c.ListResource(context, address)
//	if err != nil {
//		return nil, err
//	}
//	return ls.GetBalanceOfStc()
//}

//func (c *RoochClient) ListResource(context context.Context, address string) (*ListResource, error) {
//	result := &ListResource{
//		Resources: make(map[string]Resource),
//	}
//	params := []string{address}
//	err := c.Call(context, "state.list_resource", result, params)
//
//	if err != nil {
//		return nil, errors.Wrap(err, "call method state.list_resource ")
//	}
//
//	return result, nil
//}

//func (c *RoochClient) GetResource(context context.Context, address string, restype string, opt GetResourceOption, result interface{}) (interface{}, error) {
//	params := []interface{}{address, restype, opt}
//	err := c.Call(context, "state.get_resource", result, params)
//	if err != nil {
//		return nil, errors.Wrap(err, "call method state.get_resource ")
//	}
//	return result, nil
//}

func (c *RoochClient) GetAccountSequenceNumber(context context.Context, address string) (uint64, error) {
	//state, err := c.GetState(context, address)
	restype := "0x00000000000000000000000000000001::Account::Account"
	result := &RawResource{}
	opt := GetResourceOption{
		Decode: false,
	}
	r, err := c.GetResource(context, address, restype, opt, result)
	if r == nil {
		return 0, errors.New("get \"0x1::Account::Account\" resource return nil")
	}
	if err != nil {
		return 0, errors.Wrap(err, "call method GetResource ")
	}
	if err != nil {
		return 0, err
	}
	bs, err := HexStringToBytes(r.(*RawResource).Raw)
	if err != nil {
		return 0, err
	}
	accountResource, err := types.BcsDeserializeAccountResource(bs)
	if err != nil {
		return 0, errors.Wrap(err, "Bcs Deserialize AccountResource failed")
	}
	return accountResource.SequenceNumber, nil
}

//func (c *RoochClient) GetState(context context.Context, address string) (*types.AccountResource, error) {
//	var result []byte
//	params := []string{address + "/1/0x00000000000000000000000000000001::Account::Account"}
//	err := c.Call(context, "state.get", &result, params)
//
//	if err != nil {
//		return nil, errors.Wrap(err, "call method state.get ")
//	}
//
//	accountResource, err := types.BcsDeserializeAccountResource(result)
//	if err != nil {
//		return nil, errors.Wrap(err, "Bcs Deserialize AccountResource failed")
//	}
//
//	return &accountResource, nil
//}

func (c *RoochClient) SubmitTransaction(context context.Context, signer Signer,
	transaction transactions.Transaction) (string, error) {

	authenticator, err := signer.SignTransaction(transaction)
	if err != nil {
		return emptyString, errors.Wrap(err, "gen SignedUserTransaction failed")
	}

	signedUserTransactionBytes, err := signedUserTransaction.BcsSerialize()
	if err != nil {
		return emptyString, errors.Wrap(err, "Bcs Serialize  SignedUserTransaction failed")
	}

	var result string
	params := []string{hex.EncodeToString(signedUserTransactionBytes)}
	err = c.Call(context, "txpool.submit_hex_transaction", &result, params)

	if err != nil {
		return emptyString, errors.Wrap(err, "call txpool.submit_hex_transaction ")
	}

	return result, nil
}

func (c *RoochClient) SubmitSignedTransaction(context context.Context,
	userTxn *types.SignedUserTransaction) (string, error) {
	signedUserTransactionBytes, err := userTxn.BcsSerialize()
	if err != nil {
		return emptyString, errors.Wrap(err, "Bcs Serialize  SignedUserTransaction failed")
	}

	var result string
	params := []string{hex.EncodeToString(signedUserTransactionBytes)}
	err = c.Call(context, "txpool.submit_hex_transaction", &result, params)

	if err != nil {
		return emptyString, errors.Wrap(err, "call txpool.submit_hex_transaction ")
	}

	return result, nil
}

func (c *RoochClient) SubmitSignedTransactionBytes(context context.Context,
	userTxn []byte) (string, error) {
	var result string
	params := []string{hex.EncodeToString(userTxn)}
	err := c.Call(context, "txpool.submit_hex_transaction", &result, params)

	if err != nil {
		return emptyString, errors.Wrap(err, "call txpool.submit_hex_transaction ")
	}

	return result, nil
}

func (c *RoochClient) BuildRawUserTransaction(context context.Context, sender types.AccountAddress, payload transactions.TransactionPayload,
	gasPrice int, gasLimit uint64, seq uint64) (*types.RawUserTransaction, error) {
	nodeInfo, err := c.GetNodeInfo(context)
	if err != nil {
		return nil, errors.Wrap(err, "get node info failed ")
	}
	return &types.RawUserTransaction{
		Sender:                  sender,
		SequenceNumber:          seq,
		Payload:                 payload,
		MaxGasAmount:            gasLimit,
		GasUnitPrice:            uint64(gasPrice),
		GasTokenCode:            GasTokenCode,
		ExpirationTimestampSecs: uint64(nodeInfo.NowSeconds + DefaultTransactionExpirationSeconds),
		ChainId:                 types.ChainId{Id: uint8(nodeInfo.PeerInfo.ChainInfo.ChainID)},
	}, nil
}

//func (c *RoochClient) GetGasUnitPrice(context context.Context) (int, error) {
//	var result string
//	err := c.Call(context, "txpool.gas_price", &result, nil)
//
//	if err != nil {
//		return 1, errors.Wrap(err, "call method txpool.gas_price ")
//	}
//
//	return strconv.Atoi(result)
//}

//func (c *RoochClient) CallContract(context context.Context, call ContractCall) (interface{}, error) {
//	var result []interface{}
//	err := c.Call(context, "contract.call_v2", &result, []interface{}{call})
//
//	if err != nil {
//		return 1, errors.Wrap(err, "call method contract.call_v2 ")
//	}
//
//	return result, nil
//}

//func (c *RoochClient) EstimateGasByDryRunRaw(context context.Context, txn types.RawUserTransaction, publicKey types.Ed25519PublicKey) (*big.Int, error) {
//	result, err := c.DryRunRaw(context, txn, publicKey)
//	if err != nil {
//		return nil, errors.Wrap(err, "call method DryRunRaw ")
//	}
//	return extractGasUsed(result)
//}

func (c *RoochClient) DryRunRaw(context context.Context, txn types.RawUserTransaction, publicKey types.Ed25519PublicKey) (*DryRunResult, error) {
	var result DryRunResult
	data, err := txn.BcsSerialize()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = c.Call(context, "contract.dry_run_raw", &result, []interface{}{BytesToHexString(data), BytesToHexString(publicKey)})
	if err != nil {
		return nil, errors.Wrap(err, "call method contract.dry_run_raw ")
	}
	return &result, nil
}

func (c *RoochClient) EstimateGas(context context.Context, chainID int, gasUnitPrice int, maxGasAmount uint64,
	senderAddress string, publicKey types.Ed25519PublicKey, accountSeqNumber *uint64,
	code string, typeArgs []string, args []string) (*big.Int, error) {
	result, err := c.DryRun(context, chainID, gasUnitPrice, maxGasAmount, senderAddress, publicKey, accountSeqNumber, code, typeArgs, args)
	if err != nil {
		return nil, errors.Wrap(err, "call method DryRun ")
	}
	return extractGasUsed(result)
}

func extractGasUsed(result *DryRunResult) (*big.Int, error) {
	if !strings.EqualFold("Executed", result.ExplainedStatus) {
		return nil, fmt.Errorf("DryRun result ExplainedStatus is not 'Executed' ")
	}
	i := new(big.Int)
	i.SetString(result.GasUsed, 10)
	return i, nil
}

func (c *RoochClient) DryRun(context context.Context, chainID int, gasUnitPrice int, maxGasAmount uint64,
	senderAddress string, publicKey types.Ed25519PublicKey, accountSeqNumber *uint64,
	code string, typeArgs []string, args []string) (*DryRunResult, error) {
	var result DryRunResult
	dryRunParam := DryRunParam{
		ChainId:         chainID & 0xFF,
		GasUnitPrice:    gasUnitPrice,
		Sender:          senderAddress,
		SenderPublicKey: BytesToHexString(publicKey),
		SequenceNumber:  accountSeqNumber,
		MaxGasAmount:    maxGasAmount,
		Script: DryRunParamScript{
			Code:     code,
			TypeArgs: typeArgs,
			Args:     args,
		},
	}
	err := c.Call(context, "contract.dry_run", &result, []interface{}{dryRunParam})
	if err != nil {
		return nil, errors.Wrap(err, "call method contract.dry_run ")
	}
	return &result, nil
}

//func (c *RoochClient) DeployContract(context context.Context, sender types.AccountAddress, privateKey types.Ed25519PrivateKey,
//	function types.ScriptFunction, code []byte) (string, error) {
//	module := types.Module{
//		Code: code,
//	}
//	pk := types.Package{
//		PackageAddress: sender,
//		Modules:        []types.Module{module},
//		InitScript:     &function,
//	}
//	packagePayload := types.TransactionPayload__Package{
//		Value: pk,
//	}
//
//	price, err := c.GetGasUnitPrice(context)
//	if err != nil {
//		return "", errors.Wrap(err, "get gas unit price failed ")
//	}
//
//	state, err := c.GetState(context, "0x"+hex.EncodeToString(sender[:]))
//
//	if err != nil {
//		return "", errors.Wrap(err, "call txpool.submit_hex_transaction ")
//	}
//
//	rawTransactoin, err := c.BuildRawUserTransaction(context, sender, &packagePayload, price, DefaultMaxGasAmount, state.SequenceNumber)
//	if err != nil {
//		return emptyString, errors.Wrap(err, "build raw user txn failed")
//	}
//
//	return c.SubmitTransaction(context, privateKey, rawTransactoin)
//}

// getRpcApiVersion gets RPC API version
func (c *RoochClient) getRpcApiVersion() (string, error) {
	var resp struct {
		Info struct {
			Version string `json:"version"`
		} `json:"info"`
	}
	err := c.transport.request("rpc.discover", nil, &resp)
	return resp.Info.Version, err

	var result string
	params := []string{hex.EncodeToString(userTxn)}
	err := c.Call(context, "txpool.submit_hex_transaction", &result, params)

	if err != nil {
		return emptyString, errors.Wrap(err, "call txpool.submit_hex_transaction ")
	}

	return result, nil

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("call method %s err: ", serviceMethod))
	}
}

// getChainId gets chain ID
func (c *RoochClient) getChainId() (uint64, error) {
	if c.chainID > 0 {
		return c.chainID, nil
	}

	var result string
	params := []string{hex.EncodeToString(userTxn)}
	err := c.Call(context, "txpool.submit_hex_transaction", &result, params)

	if err != nil {
		return emptyString, errors.Wrap(err, "call txpool.submit_hex_transaction ")
	}
	var result string
	err := c.client.request("rooch_getChainID", nil, &result)
	if err != nil {
		return nil, err
	}

	c.chainID = new(big.Int)
	c.chainID.SetString(result, 10)
	return c.chainID, nil
}

// executeViewFunction executes view function
func (c *RoochClient) executeViewFunction(input CallFunctionArgs) (*AnnotatedFunctionResultView, error) {
	callFunction := NewCallFunction(input)

	params := map[string]interface{}{
		"function_id": callFunction.functionId(),
		"args":        callFunction.encodeArgs(),
		"ty_args":     callFunction.typeArgs,
	}

	var result AnnotatedFunctionResultView
	err := c.transport.request("rooch_executeViewFunction", []interface{}{params}, &result)
	return &result, err
}

// signAndExecuteTransaction signs and executes transaction
func (c *RoochClient) signAndExecuteTransaction(params struct {
	Transaction interface{}
	Signer      Signer
	Option      *struct{ WithOutput bool }
}) (*ExecuteTransactionResponseView, error) {
	var transactionHex string

	switch tx := params.Transaction.(type) {
	case []byte:
		transactionHex = str(HEX, tx)
	case *Transaction:
		sender := params.Signer.getRoochAddress().toHexAddress()

		chainID, err := c.getChainId()
		if err != nil {
			return nil, err
		}
		tx.setChainId(chainID)

		seqNum, err := c.getSequenceNumber(sender)
		if err != nil {
			return nil, err
		}
		tx.setSeqNumber(seqNum)
		tx.setSender(sender)

		auth, err := params.Signer.signTransaction(tx)
		if err != nil {
			return nil, err
		}
		tx.setAuth(auth)

		encoded, err := tx.encode()
		if err != nil {
			return nil, err
		}
		transactionHex = "0x" + encoded.toHex()
	default:
		return nil, fmt.Errorf("unsupported transaction type")
	}

	var result ExecuteTransactionResponseView
	err := c.transport.request("rooch_executeRawTransaction",
		[]interface{}{transactionHex, params.Option}, &result)
	return &result, err
}

// getStates gets states by access path
func (c *RoochClient) getStates(params GetStatesParams) ([]ObjectStateView, error) {
	var result []ObjectStateView
	err := c.transport.request("rooch_getStates",
		[]interface{}{params.AccessPath, params.StateOption}, &result)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 || result[0] == nil {
		return []ObjectStateView{}, nil
	}
	return result, nil
}

// listStates lists states
func (c *RoochClient) listStates(params ListStatesParams) (*PaginatedStateKVViews, error) {
	var result PaginatedStateKVViews
	err := c.transport.request("rooch_listStates",
		[]interface{}{params.AccessPath, params.Cursor, params.Limit, params.StateOption}, &result)
	return &result, err
}

// getModuleAbi gets module ABI
func (c *RoochClient) getModuleAbi(params GetModuleABIParams) (*ModuleABIView, error) {
	var result ModuleABIView
	err := c.transport.request("rooch_getModuleABI",
		[]interface{}{params.ModuleAddr, params.ModuleName}, &result)
	return &result, err
}

// getEvents gets events by event handle
func (c *RoochClient) getEvents(params GetEventsByEventHandleParams) (*PaginatedEventViews, error) {
	var result PaginatedEventViews
	err := c.transport.request("rooch_getEventsByEventHandle",
		[]interface{}{
			params.EventHandleType,
			params.Cursor,
			params.Limit,
			params.DescendingOrder,
			params.EventOptions,
		}, &result)
	return &result, err
}

// getBalance gets balance for coin type
func (c *RoochClient) getBalance(params GetBalanceParams) (*BalanceInfoView, error) {
	if !isValidRoochAddress(params.Owner) {
		return nil, fmt.Errorf("invalid rooch address")
	}

	var result BalanceInfoView
	err := c.transport.request("rooch_getBalance",
		[]interface{}{params.Owner, params.CoinType}, &result)
	return &result, err
}

// transfer transfers coins
func (c *RoochClient) transfer(params struct {
	Signer    Signer
	Recipient string
	Amount    *big.Int
	CoinType  TypeArgs
}) (*ExecuteTransactionResponseView, error) {
	tx := NewTransaction()
	tx.callFunction(CallFunctionArgs{
		Target:   "0x3::transfer::transfer_coin",
		Args:     []Args{Args.Address(params.Recipient), Args.U256(params.Amount)},
		TypeArgs: []string{normalizeTypeArgsToStr(params.CoinType)},
	})

	return c.signAndExecuteTransaction(struct {
		Transaction interface{}
		Signer      Signer
		Option      *struct{ WithOutput bool }
	}{
		Transaction: tx,
		Signer:      params.Signer,
	})
}

// Helper methods

// getSequenceNumber gets sequence number for address
func (c *RoochClient) getSequenceNumber(address string) (*big.Int, error) {
	resp, err := c.executeViewFunction(CallFunctionArgs{
		Target: "0x2::account::sequence_number",
		Args:   []Args{Args.Address(address)},
	})
	if err != nil {
		return nil, err
	}

	if resp != nil && len(resp.ReturnValues) > 0 {
		seqNum := new(big.Int)
		seqNum.SetString(fmt.Sprint(resp.ReturnValues[0].DecodedValue), 10)
		return seqNum, nil
	}

	return big.NewInt(0), nil
}
