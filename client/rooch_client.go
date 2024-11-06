package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/rooch-network/rooch-go-sdk/crypto"
	"github.com/rooch-network/rooch-go-sdk/transactions"
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
	url string
}

func NewRoochClient(url string) RoochClient {
	return RoochClient{
		url: url,
	}
}

func (this *RoochClient) Call(context context.Context, serviceMethod string, reply interface{}, args interface{}) error {
	client, err := jsonrpc.NewClient(this.url)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("call method %s err: ", serviceMethod))
	}

	err = client.Call(context, serviceMethod, reply, args)

	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("call method %s err: ", serviceMethod))
	}

	return nil
}

func (this *RoochClient) GetEvents(context context.Context, eventFilter *EventFilter) ([]Event, error) {
	var result []Event
	params := []interface{}{eventFilter}
	err := this.Call(context, "chain.get_events", &result, params)
	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_events error ")
	}
	return result, nil
}

func (this *RoochClient) GetTransactionByHash(context context.Context, transactionHash string) (*Transaction, error) {
	result := &Transaction{}
	params := []string{transactionHash}
	err := this.Call(context, "chain.get_transaction", result, params)

	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_transaction error ")
	}

	return result, nil
}

func (this *RoochClient) GetTransactionInfoByHash(context context.Context, transactionHash string) (*TransactionInfo, error) {
	result := &TransactionInfo{}
	params := []string{transactionHash}
	err := this.Call(context, "chain.get_transaction_info", result, params)

	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_transaction_info error ")
	}

	return result, nil
}

func (this *RoochClient) GetTransactionEventByHash(context context.Context, transactionHash string) ([]Event, error) {
	var result []Event
	params := []string{transactionHash}
	err := this.Call(context, "chain.get_events_by_txn_hash", &result, params)

	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_events_by_txn_hash error ")
	}

	return result, nil
}

func (this *RoochClient) GetBlockByHash(context context.Context, blockHash string) (*Block, error) {
	result := &Block{}
	params := []string{blockHash}
	err := this.Call(context, "chain.get_block_by_hash", result, params)

	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_block_by_hash ")
	}

	return result, nil
}

func (this *RoochClient) GetBlockByNumber(context context.Context, number int) (*Block, error) {
	result := &Block{}
	params := []int{number}
	err := this.Call(context, "chain.get_block_by_number", result, params)

	if err != nil {
		return nil, errors.Wrap(err, "call method chain.get_block_by_number ")
	}

	return result, nil
}

//func (this *RoochClient) GetBalanceOfStc(context context.Context, address string) (*big.Int, error) {
//	ls, err := this.ListResource(context, address)
//	if err != nil {
//		return nil, err
//	}
//	return ls.GetBalanceOfStc()
//}

//func (this *RoochClient) ListResource(context context.Context, address string) (*ListResource, error) {
//	result := &ListResource{
//		Resources: make(map[string]Resource),
//	}
//	params := []string{address}
//	err := this.Call(context, "state.list_resource", result, params)
//
//	if err != nil {
//		return nil, errors.Wrap(err, "call method state.list_resource ")
//	}
//
//	return result, nil
//}

//func (this *RoochClient) GetResource(context context.Context, address string, restype string, opt GetResourceOption, result interface{}) (interface{}, error) {
//	params := []interface{}{address, restype, opt}
//	err := this.Call(context, "state.get_resource", result, params)
//	if err != nil {
//		return nil, errors.Wrap(err, "call method state.get_resource ")
//	}
//	return result, nil
//}

func (this *RoochClient) GetAccountSequenceNumber(context context.Context, address string) (uint64, error) {
	//state, err := this.GetState(context, address)
	restype := "0x00000000000000000000000000000001::Account::Account"
	result := &RawResource{}
	opt := GetResourceOption{
		Decode: false,
	}
	r, err := this.GetResource(context, address, restype, opt, result)
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

//func (this *RoochClient) GetState(context context.Context, address string) (*types.AccountResource, error) {
//	var result []byte
//	params := []string{address + "/1/0x00000000000000000000000000000001::Account::Account"}
//	err := this.Call(context, "state.get", &result, params)
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

func (this *RoochClient) SubmitTransaction(context context.Context, signer Signer,
	rawUserTransaction *transactions.Transaction) (string, error) {
	signedUserTransaction, err := signTxn(privateKey, rawUserTransaction)
	if err != nil {
		return emptyString, errors.Wrap(err, "gen SignedUserTransaction failed")
	}

	signedUserTransactionBytes, err := signedUserTransaction.BcsSerialize()
	if err != nil {
		return emptyString, errors.Wrap(err, "Bcs Serialize  SignedUserTransaction failed")
	}

	var result string
	params := []string{hex.EncodeToString(signedUserTransactionBytes)}
	err = this.Call(context, "txpool.submit_hex_transaction", &result, params)

	if err != nil {
		return emptyString, errors.Wrap(err, "call txpool.submit_hex_transaction ")
	}

	return result, nil
}

func (this *RoochClient) SubmitSignedTransaction(context context.Context,
	userTxn *types.SignedUserTransaction) (string, error) {
	signedUserTransactionBytes, err := userTxn.BcsSerialize()
	if err != nil {
		return emptyString, errors.Wrap(err, "Bcs Serialize  SignedUserTransaction failed")
	}

	var result string
	params := []string{hex.EncodeToString(signedUserTransactionBytes)}
	err = this.Call(context, "txpool.submit_hex_transaction", &result, params)

	if err != nil {
		return emptyString, errors.Wrap(err, "call txpool.submit_hex_transaction ")
	}

	return result, nil
}

func (this *RoochClient) SubmitSignedTransactionBytes(context context.Context,
	userTxn []byte) (string, error) {
	var result string
	params := []string{hex.EncodeToString(userTxn)}
	err := this.Call(context, "txpool.submit_hex_transaction", &result, params)

	if err != nil {
		return emptyString, errors.Wrap(err, "call txpool.submit_hex_transaction ")
	}

	return result, nil
}

func (this *RoochClient) BuildRawUserTransaction(context context.Context, sender types.AccountAddress, payload transactions.TransactionPayload,
	gasPrice int, gasLimit uint64, seq uint64) (*types.RawUserTransaction, error) {
	nodeInfo, err := this.GetNodeInfo(context)
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

//func (this *RoochClient) GetGasUnitPrice(context context.Context) (int, error) {
//	var result string
//	err := this.Call(context, "txpool.gas_price", &result, nil)
//
//	if err != nil {
//		return 1, errors.Wrap(err, "call method txpool.gas_price ")
//	}
//
//	return strconv.Atoi(result)
//}

func (this *RoochClient) CallContract(context context.Context, call ContractCall) (interface{}, error) {
	var result []interface{}
	err := this.Call(context, "contract.call_v2", &result, []interface{}{call})

	if err != nil {
		return 1, errors.Wrap(err, "call method contract.call_v2 ")
	}

	return result, nil
}

//func (this *RoochClient) EstimateGasByDryRunRaw(context context.Context, txn types.RawUserTransaction, publicKey types.Ed25519PublicKey) (*big.Int, error) {
//	result, err := this.DryRunRaw(context, txn, publicKey)
//	if err != nil {
//		return nil, errors.Wrap(err, "call method DryRunRaw ")
//	}
//	return extractGasUsed(result)
//}

func (this *RoochClient) DryRunRaw(context context.Context, txn types.RawUserTransaction, publicKey types.Ed25519PublicKey) (*DryRunResult, error) {
	var result DryRunResult
	data, err := txn.BcsSerialize()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = this.Call(context, "contract.dry_run_raw", &result, []interface{}{BytesToHexString(data), BytesToHexString(publicKey)})
	if err != nil {
		return nil, errors.Wrap(err, "call method contract.dry_run_raw ")
	}
	return &result, nil
}

func (this *RoochClient) EstimateGas(context context.Context, chainId int, gasUnitPrice int, maxGasAmount uint64,
	senderAddress string, publicKey types.Ed25519PublicKey, accountSeqNumber *uint64,
	code string, typeArgs []string, args []string) (*big.Int, error) {
	result, err := this.DryRun(context, chainId, gasUnitPrice, maxGasAmount, senderAddress, publicKey, accountSeqNumber, code, typeArgs, args)
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

func (this *RoochClient) DryRun(context context.Context, chainId int, gasUnitPrice int, maxGasAmount uint64,
	senderAddress string, publicKey types.Ed25519PublicKey, accountSeqNumber *uint64,
	code string, typeArgs []string, args []string) (*DryRunResult, error) {
	var result DryRunResult
	dryRunParam := DryRunParam{
		ChainId:         chainId & 0xFF,
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
	err := this.Call(context, "contract.dry_run", &result, []interface{}{dryRunParam})
	if err != nil {
		return nil, errors.Wrap(err, "call method contract.dry_run ")
	}
	return &result, nil
}

//func (this *RoochClient) DeployContract(context context.Context, sender types.AccountAddress, privateKey types.Ed25519PrivateKey,
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
//	price, err := this.GetGasUnitPrice(context)
//	if err != nil {
//		return "", errors.Wrap(err, "get gas unit price failed ")
//	}
//
//	state, err := this.GetState(context, "0x"+hex.EncodeToString(sender[:]))
//
//	if err != nil {
//		return "", errors.Wrap(err, "call txpool.submit_hex_transaction ")
//	}
//
//	rawTransactoin, err := this.BuildRawUserTransaction(context, sender, &packagePayload, price, DefaultMaxGasAmount, state.SequenceNumber)
//	if err != nil {
//		return emptyString, errors.Wrap(err, "build raw user txn failed")
//	}
//
//	return this.SubmitTransaction(context, privateKey, rawTransactoin)
//}
