package types

import (
	"encoding/json"
	"fmt"
	"github.com/rooch-network/rooch-go-sdk/crypto"

	"github.com/rooch-network/rooch-go-sdk/bcs"
)

type Transaction struct {
	Data TransactionData `json:"data"`
	//Data          crypto.TransactionData `json:"data"`
	Authenticator crypto.Authenticator `json:"authenticator"`
	Info          string               `json:"info"`
}

func (t *Transaction) MarshalBCS(ser *bcs.Serializer) {
	t.Data.MarshalBCS(ser)
	t.Authenticator.MarshalBCS(ser)
	ser.WriteString(t.Info)
}
func (t *Transaction) UnmarshalBCS(des *bcs.Deserializer) {
	t.Data.UnmarshalBCS(des)
	t.Authenticator.UnmarshalBCS(des)
	t.Info = des.ReadString()
}

// GetInfo returns the transaction info
func (t *Transaction) GetInfo() *string {
	return &t.Info
}

// HashData returns the hash of the transaction data
func (t *Transaction) HashData() ([]byte, error) {
	return t.Data.Hash()
}

// getData returns the transaction data after validation
func (t *Transaction) GetData() crypto.TransactionData {
	//t.isValid()
	return &t.Data
}

type TransactionSequenceInfo struct {
	TxOrder                         uint64   `json:"tx_order"`
	TxOrderSignature                []byte   `json:"tx_order_signature"`
	TxAccumulatorRoot               string   `json:"tx_accumulator_root"`
	TxTimestamp                     uint64   `json:"tx_timestamp"`
	TxAccumulatorFrozenSubtreeRoots []string `json:"tx_accumulator_frozen_subtree_roots"`
	TxAccumulatorNumLeaves          uint64   `json:"tx_accumulator_num_leaves"`
	TxAccumulatorNumNodes           uint64   `json:"tx_accumulator_num_nodes"`
}

//pub struct TransactionExecutionInfo {
///// The hash of this transaction.
//pub tx_hash: H256,
//
///// The root hash of Sparse Merkle Tree describing the world state at the end of this
///// transaction.
//pub state_root: H256,
//
///// The root Object count of Sparse Merkle Tree describing the world state at the end of this transaction.
//pub size: u64,
//
///// The root hash of Merkle Accumulator storing all events emitted during this transaction.
//pub event_root: H256,
//
///// The amount of gas used.
//pub gas_used: u64,
//
///// The vm status. If it is not `Executed`, this will provide the general error class. Execution
///// failures and Move abort's receive more detailed information. But other errors are generally
///// categorized with no status code or other information
//pub status: KeptVMStatus,
//}

type TransactionExecutionInfo struct {
	TxHash    string `json:"tx_hash"`
	StateRoot string `json:"state_root"`
	Size      uint64 `json:"size"`
	EventRoot string `json:"event_root"`
	GasUsed   uint64 `json:"gas_used"`
	/// The vm status. If it is not `Executed`, this will provide the general error class. Execution
	/// failures and Move abort's receive more detailed information. But other errors are generally
	/// categorized with no status code or other information
	//pub status: KeptVMStatus,
	//Status string `json:"status"`
	Status json.RawMessage `json:"status"`
}

//pub struct L1Block {
//pub chain_id: MultiChainID,
//pub block_height: u64,
//pub block_hash: Vec<u8>,
//}

type L1Block struct {
	ChainId     MultiChainIDVariant `json:"chain_id"`
	BlockHeight uint64              `json:"block_height"`
	BlockHash   []byte              `json:"block_hash"`
	//GasUsed      uint64 `json:"gas_used"`
	/// The vm status. If it is not `Executed`, this will provide the general error class. Execution
	/// failures and Move abort's receive more detailed information. But other errors are generally
	/// categorized with no status code or other information
	//pub status: KeptVMStatus,
	//Status string `json:"status"`
	//Status json.RawMessage `json:"status"`
}

func (lb *L1Block) LedgerTxDataType() LedgerTxDataVariant {
	return LedgerTxDataVariantL1Block
}

func (lb *L1Block) MarshalBCS(ser *bcs.Serializer) {
	ser.U64(uint64(lb.ChainId))
	ser.U64(lb.BlockHeight)
	ser.WriteBytes(lb.BlockHash)
}
func (lb *L1Block) UnmarshalBCS(des *bcs.Deserializer) {
	lb.ChainId = MultiChainIDVariant(des.U64())
	lb.BlockHeight = des.U64()
	lb.BlockHash = des.ReadBytes()
}

//pub struct L1Transaction {
//pub chain_id: MultiChainID,
//pub block_hash: Vec<u8>,
///// The original L1 transaction id, usually the hash of the transaction
//pub txid: Vec<u8>,
//}

type L1Transaction struct {
	ChainId   MultiChainIDVariant `json:"chain_id"`
	BlockHash []byte              `json:"block_hash"`
	TxID      []byte              `json:"txid"`
	//GasUsed      uint64 `json:"gas_used"`
	/// The vm status. If it is not `Executed`, this will provide the general error class. Execution
	/// failures and Move abort's receive more detailed information. But other errors are generally
	/// categorized with no status code or other information
	//pub status: KeptVMStatus,
	//Status string `json:"status"`
	//Status json.RawMessage `json:"status"`
}

func (lt *L1Transaction) LedgerTxDataType() LedgerTxDataVariant {
	return LedgerTxDataVariantL1Tx
}

func (lt *L1Transaction) MarshalBCS(ser *bcs.Serializer) {
	ser.U64(uint64(lt.ChainId))
	ser.WriteBytes(lt.BlockHash)
	ser.WriteBytes(lt.TxID)
}
func (lt *L1Transaction) UnmarshalBCS(des *bcs.Deserializer) {
	lt.ChainId = MultiChainIDVariant(des.U64())
	lt.BlockHash = des.ReadBytes()
	lt.TxID = des.ReadBytes()
}

//pub struct RoochTransaction {
//pub data: TransactionData,
//pub authenticator: Authenticator,
//
//#[serde(skip_serializing, skip_deserializing)]
//data_hash: Option<H256>,
//}

type RoochTransaction struct {
	Data          TransactionData      `json:"data"`
	Authenticator crypto.Authenticator `json:"authenticator"`
	DataHash      string               `json:"data_hash,omitempty"`
	//GasUsed      uint64 `json:"gas_used"`
	/// The vm status. If it is not `Executed`, this will provide the general error class. Execution
	/// failures and Move abort's receive more detailed information. But other errors are generally
	/// categorized with no status code or other information
	//pub status: KeptVMStatus,
	//Status string `json:"status"`
	//Status json.RawMessage `json:"status"`
}

func (rt *RoochTransaction) LedgerTxDataType() LedgerTxDataVariant {
	return LedgerTxDataVariantL2Tx
}

func (rt *RoochTransaction) MarshalBCS(ser *bcs.Serializer) {
	//ser.U64(uint64(rt.Data))
	rt.Data.MarshalBCS(ser)
	rt.Authenticator.MarshalBCS(ser)
	ser.WriteString(rt.DataHash)
}
func (rt *RoochTransaction) UnmarshalBCS(des *bcs.Deserializer) {
	rt.Data.UnmarshalBCS(des)
	rt.Authenticator.UnmarshalBCS(des)
	rt.DataHash = des.ReadString()
}

//pub enum LedgerTxData {
//L1Block(L1Block),
//L1Tx(L1Transaction),
//L2Tx(RoochTransaction),
//}

//type LedgerTxData struct {
//	data          LedgerTxData            `json:"data"`
//	sequence_info TransactionSequenceInfo `json:"sequence_info"`
//}

//type TransactionPayloadVariant uint32
//
//const (
//	TransactionPayloadVariantScript        TransactionPayloadVariant = 0
//	TransactionPayloadVariantModuleBundle  TransactionPayloadVariant = 1 // Deprecated
//	TransactionPayloadVariantEntryFunction TransactionPayloadVariant = 2
//	TransactionPayloadVariantMultisig      TransactionPayloadVariant = 3
//)

//#[derive(Clone, Debug, Hash, Eq, PartialEq, Serialize, Deserialize)]
//pub enum LedgerTxData {
//L1Block(L1Block),
//L1Tx(L1Transaction),
//L2Tx(RoochTransaction),
//}

type LedgerTxDataVariant uint32

const (
	LedgerTxDataVariantL1Block LedgerTxDataVariant = 0
	LedgerTxDataVariantL1Tx    LedgerTxDataVariant = 1 // Deprecated
	LedgerTxDataVariantL2Tx    LedgerTxDataVariant = 2
)

type LedgerTxDataImpl interface {
	bcs.Struct
	LedgerTxDataType() LedgerTxDataVariant // This is specifically to ensure that wrong types don't end up here
}

// TransactionPayload the actual instructions of which functions to call on chain
type LedgerTxData struct {
	//Payload TransactionPayloadImpl
	//Variant LedgerTxDataVariant `json:"variant"`
	TxData LedgerTxDataImpl `json:"tx_data"`
}

//func (ltd *LedgerTxData) MarshalBCS(ser *bcs.Serializer) {
//	ser.Uleb128(uint32(ltd.Variant))
//	ser.Struct(ltd.TxData)
//}
//func (ltd *LedgerTxData) UnmarshalBCS(des *bcs.Deserializer) {
//	variant := LedgerTxDataVariant(des.Uleb128())
//	switch variant {
//	case LedgerTxDataVariantL1Block:
//		ltd.TxData = &L1Block{}
//	case LedgerTxDataVariantL1Tx:
//		ltd.TxData = &L1Transaction{}
//	case LedgerTxDataVariantL2Tx:
//		ltd.TxData = &RoochTransaction{}
//	default:
//		des.SetError(fmt.Errorf("bad variant %d for MultisigTransactionPayload", variant))
//		return
//	}
//	des.Struct(ltd.TxData)
//}

func (ltd *LedgerTxData) MarshalBCS(ser *bcs.Serializer) {
	if ltd == nil || ltd.TxData == nil {
		ser.SetError(fmt.Errorf("Ledger tx data is nil"))
		return
	}
	ser.Uleb128(uint32(ltd.TxData.LedgerTxDataType()))
	ltd.TxData.MarshalBCS(ser)
}
func (ltd *LedgerTxData) UnmarshalBCS(des *bcs.Deserializer) {
	txDataType := LedgerTxDataVariant(des.Uleb128())
	switch txDataType {
	case LedgerTxDataVariantL1Block:
		ltd.TxData = &L1Block{}
	case LedgerTxDataVariantL1Tx:
		ltd.TxData = &L1Transaction{}
	case LedgerTxDataVariantL2Tx:
		ltd.TxData = &RoochTransaction{}
	default:
		des.SetError(fmt.Errorf("Invalid ledger tx data type, %d", txDataType))
		return
	}

	ltd.TxData.UnmarshalBCS(des)
}

//type LedgerTxData struct {
//	data          LedgerTxData            `json:"data"`
//	sequence_info TransactionSequenceInfo `json:"sequence_info"`
//}

//endregion
//endregion

//#[derive(Clone, Debug, Eq, PartialEq, Serialize, Deserialize)]
//pub struct LedgerTransaction {
//pub data: LedgerTxData,
//pub sequence_info: TransactionSequenceInfo,
//}

type LedgerTransaction struct {
	data          LedgerTxData            `json:"data"`
	sequence_info TransactionSequenceInfo `json:"sequence_info"`
}

//#[derive(Debug, Clone)]
//pub struct TransactionWithInfo {
//pub transaction: LedgerTransaction,
//pub execution_info: Option<TransactionExecutionInfo>,
//}

type TransactionWithInfo struct {
	Transaction    LedgerTransaction        `json:"transaction"`
	execution_info TransactionExecutionInfo `json:"execution_info,omitempty"`
}

//type AccumulatorInfo struct {
//	AccumulatorRoot    string   `json:"accumulator_root"`
//	FrozenSubtreeRoots []string `json:"frozen_subtree_roots"`
//	NumLeaves          string   `json:"num_leaves"`
//	NumNodes           string   `json:"num_nodes"`
//}
