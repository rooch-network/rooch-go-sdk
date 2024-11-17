package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rooch-network/rooch-go-sdk/crypto"
	"github.com/rooch-network/rooch-go-sdk/utils"

	//"github.com/pkg/errors"

	"github.com/rooch-network/rooch-go-sdk/bcs"
)

const recvPrefix = "0100000000000000"
const sendPrefix = "0000000000000000"

//type BlockHeaderWithDifficultyInfo struct {
//	BlockHeader          BlockHeader `json:"header"`
//	BlockTimeTarget      uint64      `json:"block_time_target,omitempty"`       //: 5260,
//	BlockDifficutyWindow uint        `json:"block_difficulty_window,omitempty"` //: 24,
//	BlockInfo            BlockInfo   `json:"block_info"`
//}

//type BlockHeaderAndBlockInfo struct {
//	BlockHeader BlockHeader `json:"header"`
//	BlockInfo   BlockInfo   `json:"block_info"`
//}

//type TransactionInfo struct {
//	BlockHash              string          `json:"block_hash"`
//	BlockNumber            string          `json:"block_number"`
//	TransactionHash        string          `json:"transaction_hash"`
//	TransactionIndex       int             `json:"transaction_index"`
//	TransactionGlobalIndex string          `json:"transaction_global_index"`
//	StateRootHash          string          `json:"state_root_hash"`
//	EventRootHash          string          `json:"event_root_hash"`
//	GasUsed                string          `json:"gas_used"`
//	Status                 json.RawMessage `json:"status"`
//}

//	pub struct TransactionSequenceInfo {
//		/// The tx order
//		pub tx_order: u64,
//		/// The tx order signature, it is the signature of the sequencer to commit the tx order.
//		pub tx_order_signature: Vec<u8>,
//		/// The tx accumulator root after the tx is append to the accumulator.
//		pub tx_accumulator_root: H256,
//		/// The tx accumulator info after the tx is append to the accumulator.
//		// pub tx_accumulator_info: Option<AccumulatorInfo>,
//		/// The timestamp of the sequencer when the tx is sequenced, in millisecond.
//		pub tx_timestamp: u64,
//
//		/// Frozen subtree roots of the accumulator.
//		pub tx_accumulator_frozen_subtree_roots: Vec<H256>,
//		/// The total number of leaves in the accumulator.
//		pub tx_accumulator_num_leaves: u64,
//		/// The total number of nodes in the accumulator.
//		pub tx_accumulator_num_nodes: u64,
//	}
//type Authenticator = Authenticator

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

//// CallFunction initializes the transaction with function call data
//func (t *Transaction) CallFunction(input CallFunctionInput) {
//	t.info = &input.Info
//	t.data = NewTransactionData(MoveAction.NewCallFunction(input.CallFunctionArgs))
//}

//MarshalBCS(ser *bcs.Serializer)
////UnmarshalBCS(des *bcs.Deserializer)
////GetInfo() *string
////HashData() ([]byte, error)
////getData() *types.TransactionData

//// SetSender sets the sender address for the transaction
//func (t *Transaction) SetSender(input types.Address) {
//	t.getData().Sender = input
//}
//
//// SetAuth sets the authenticator for the transaction
//func (t *Transaction) SetAuth(input *crypto.Authenticator) {
//	t.auth = input
//}
//
//// SetChainId sets the chain ID for the transaction
//func (t *Transaction) SetChainId(input types.U64) {
//	t.getData().ChainId = input
//}
//
//// SetSeqNumber sets the sequence number for the transaction
//func (t *Transaction) SetSeqNumber(input types.U64) {
//	t.getData().SequenceNumber = input
//}

//// HashData returns the hash of the transaction data
//func (t *Transaction) HashData() ([]byte, error) {
//	return t.Data.Hash()
//}
//
//// getData returns the transaction data after validation
//func (t *Transaction) getData() crypto.TransactionData {
//	//t.isValid()
//	return &t.Data
//}

//// isValid checks if the transaction data is initialized
//func (t *Transaction) isValid() error {
//	if t.Data == TransactionData{} {
//		return errors.New("Transaction data is not initialized. Call action first")
//	}
//	return nil
//}

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

//pub struct Authenticator {
//pub auth_validator_id: u64,
//pub payload: Vec<u8>,
//}

//type Authenticator struct {
//	AuthValidatorId uint64 `json:"auth_validator_id"`
//	Payload         []byte `json:"payload"`
//}
//
//func (au *Authenticator) MarshalBCS(ser *bcs.Serializer) {
//	ser.U64(au.AuthValidatorId)
//	ser.WriteBytes(au.Payload)
//}
//func (au *Authenticator) UnmarshalBCS(des *bcs.Deserializer) {
//	au.AuthValidatorId = des.U64()
//	au.Payload = des.ReadBytes()
//}

//pub struct ModuleId {
//address: AccountAddress,
//name: Identifier,
//}

//pub struct FunctionId {
//pub module_id: ModuleId,
//pub function_name: Identifier,
//}

//type FunctionId struct {
//	ModuleId     ModuleId   `json:"module_id"`
//	FunctionName Identifier `json:"function_name"`
//}
//
//func (mod *ModuleId) MarshalBCS(ser *bcs.Serializer) {
//	mod.Address.MarshalBCS(ser)
//	ser.WriteString(mod.Name)
//}
//
//func (mod *ModuleId) UnmarshalBCS(des *bcs.Deserializer) {
//	mod.Address.UnmarshalBCS(des)
//	mod.Name = des.ReadString()
//}

//pub enum MoveAction {
////Execute a Move script
//Script(ScriptCall),
////Execute a Move function
//Function(FunctionCall),
////Publish Move modules
//ModuleBundle(Vec<Vec<u8>>),
//}

type MoveActionVariant uint32

const (
	MoveActionVariantScript   MoveActionVariant = iota
	MoveActionVariantFunction                   // Deprecated
	MoveActionVariantModuleBundle
	//MoveActionVariantUnknown       MoveActionVariant = 10
)

type MoveActionImpl interface {
	bcs.Struct
	MoveActionType() MoveActionVariant // This is specifically to ensure that wrong types don't end up here
}

// TransactionPayload the actual instructions of which functions to call on chain
type MoveAction struct {
	//Payload TransactionPayloadImpl
	//Variant MoveActionVariant  `json:"variant"`
	//Action  MoveActionImpl  `json:"action"`
	Action MoveActionImpl
}

//func (ltd *MoveAction) MarshalBCS(ser *bcs.Serializer) {
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

func (ma *MoveAction) MarshalBCS(ser *bcs.Serializer) {
	if ma == nil || ma.Action == nil {
		ser.SetError(fmt.Errorf("Move action is nil"))
		return
	}
	ser.Uleb128(uint32(ma.Action.MoveActionType()))
	ma.Action.MarshalBCS(ser)
}
func (ma *MoveAction) UnmarshalBCS(des *bcs.Deserializer) {
	moveActionType := MoveActionVariant(des.Uleb128())
	switch moveActionType {
	case MoveActionVariantScript:
		ma.Action = &ScriptCall{}
	case MoveActionVariantFunction:
		ma.Action = &FunctionCall{}
	case MoveActionVariantModuleBundle:
		ma.Action = &ModuleBundle{}
	default:
		des.SetError(fmt.Errorf("Invalid move action type, %d", moveActionType))
		return
	}

	ma.Action.UnmarshalBCS(des)
}

//// UnmarshalJSON unmarshals the [TransactionPayload] from JSON handling conversion between types
//func (ma *MoveAction) UnmarshalJSON(b []byte) error {
//	type inner struct {
//		Variant uint32 `json:"variant"`
//	}
//	data := &inner{}
//	err := json.Unmarshal(b, &data)
//	if err != nil {
//		return err
//	}
//	ma.MoveActionType = MoveActionVariant(data.Variant)
//	switch ma.Variant {
//	case MoveActionVariantScript:
//		ma.Action = &ScriptCall{}
//	case MoveActionVariantFunction:
//		ma.Action = &FunctionCall{}
//	case MoveActionVariantModuleBundle:
//		ma.Action = &TransactionPayloadModuleBundle{}
//	default:
//		// Make sure it doesn't crash with new variants
//		ma.Action = &MoveActionUnknown{Variant: uint32(ma.Variant)}
//		ma.Variant = MoveActionVariantUnknown
//
//		return json.Unmarshal(b, &ma.Action.(*MoveActionUnknown).Action)
//	}
//	return json.Unmarshal(b, ma.Action)
//}

//
//// BoolTag represents the bool type in Move
//type BoolTag struct{}
//
////region BoolTag TypeTagImpl
//
//func (xt *BoolTag) String() string {
//	return "bool"
//}
//
//func (xt *BoolTag) GetType() TypeTagVariant {
//	return TypeTagBool
//}
//
////endregion
//
////region BoolTag bcs.struct
//
//func (xt *BoolTag) MarshalBCS(_ *bcs.Serializer)     {}
//func (xt *BoolTag) UnmarshalBCS(_ *bcs.Deserializer) {}

//pub enum MoveAction {
////Execute a Move script
//Script(ScriptCall),
////Execute a Move function
//Function(FunctionCall),
////Publish Move modules
//ModuleBundle(Vec<Vec<u8>>),
//}

// This is a fallback type for unknown move action.
type MoveActionUnknown struct {
	Variant uint32         `json:"variant"` // Variant is the actual type field from the JSON.
	Action  map[string]any `json:"action"`  // Action is the raw JSON payload.
}

// ModuleBundle is long deprecated and no longer used, but exist as an enum position in TransactionPayload
type ModuleBundle struct {
	Value [][]byte `json:"value"`
}

func (mb *ModuleBundle) MoveActionType() MoveActionVariant {
	return MoveActionVariantModuleBundle
}

func (mb *ModuleBundle) MarshalBCS(ser *bcs.Serializer) {
	//ser.SetError(errors.New("ModuleBundle unimplemented"))

	//ser.WriteBytes(sc.Code)
	//bcs.SerializeSequence(sc.TyArgs, ser)
	ser.Uleb128(uint32(len(mb.Value)))
	for _, a := range mb.Value {
		ser.WriteBytes(a)
	}
}
func (mb *ModuleBundle) UnmarshalBCS(des *bcs.Deserializer) {
	//des.SetError(errors.New("ModuleBundle unimplemented"))

	//sc.Code = des.ReadBytes()
	//sc.TyArgs = bcs.DeserializeSequence[TypeTag](des)
	alen := des.Uleb128()
	mb.Value = make([][]byte, alen)
	for i := range alen {
		mb.Value[i] = des.ReadBytes()
	}
}

//pub struct ScriptCall {
//#[serde(with = "serde_bytes")]
//pub code: Vec<u8>,
//pub ty_args: Vec<TypeTag>,
////TOOD custom serialize
//pub args: Vec<Vec<u8>>,
//}

type ScriptCall struct {
	Code   []byte    `json:"code"`
	TyArgs []TypeTag `json:"ty_args"`
	Args   [][]byte  `json:"args"`
}

func (sc *ScriptCall) MoveActionType() MoveActionVariant {
	return MoveActionVariantScript
}

//endregion

//region EntryFunction bcs.Struct

func (sc *ScriptCall) MarshalBCS(ser *bcs.Serializer) {
	ser.WriteBytes(sc.Code)
	bcs.SerializeSequence(sc.TyArgs, ser)
	ser.Uleb128(uint32(len(sc.Args)))
	for _, a := range sc.Args {
		ser.WriteBytes(a)
	}
}
func (sc *ScriptCall) UnmarshalBCS(des *bcs.Deserializer) {
	sc.Code = des.ReadBytes()
	sc.TyArgs = bcs.DeserializeSequence[TypeTag](des)
	alen := des.Uleb128()
	sc.Args = make([][]byte, alen)
	for i := range alen {
		sc.Args[i] = des.ReadBytes()
	}
}

//pub struct FunctionCall {
//pub function_id: FunctionId,
//pub ty_args: Vec<TypeTag>,
//pub args: Vec<Vec<u8>>,
//}

type FunctionCall struct {
	FunctionId FunctionId `json:"function_id"`
	TyArgs     []TypeTag  `json:"ty_args"`
	Args       [][]byte   `json:"args"`
}

func (fc *FunctionCall) MoveActionType() MoveActionVariant {
	return MoveActionVariantFunction
}

//endregion

//region EntryFunction bcs.Struct

func (fc *FunctionCall) MarshalBCS(ser *bcs.Serializer) {
	//sc.Code.MarshalBCS(ser)
	//ser.WriteString(sc.Function)
	//bcs.SerializeSequence(sc.ArgTypes, ser)
	//ser.Uleb128(uint32(len(sc.Args)))
	//for _, a := range sc.Args {
	//	ser.WriteBytes(a)
	//}

	//ser.WriteBytes(fc.Code)
	fc.FunctionId.MarshalBCS(ser)
	bcs.SerializeSequence(fc.TyArgs, ser)
	ser.Uleb128(uint32(len(fc.Args)))
	for _, a := range fc.Args {
		ser.WriteBytes(a)
	}
}
func (fc *FunctionCall) UnmarshalBCS(des *bcs.Deserializer) {
	//fc.FunctionId.UnmarshalBCS(des)
	//fc.Function = des.ReadString()
	//fc.ArgTypes = bcs.DeserializeSequence[TypeTag](des)
	//alen := des.Uleb128()
	//fc.Args = make([][]byte, alen)
	//for i := range alen {
	//	fc.Args[i] = des.ReadBytes()
	//}

	fc.FunctionId.UnmarshalBCS(des)
	//fc.Code = des.ReadBytes()
	fc.TyArgs = bcs.DeserializeSequence[TypeTag](des)
	alen := des.Uleb128()
	fc.Args = make([][]byte, alen)
	for i := range alen {
		fc.Args[i] = des.ReadBytes()
	}
}

//pub enum MoveAction {
////Execute a Move script
//Script(ScriptCall),
////Execute a Move function
//Function(FunctionCall),
////Publish Move modules
//ModuleBundle(Vec<Vec<u8>>),
//}

//pub struct TransactionData {
///// Sender's address.
//pub sender: RoochAddress,
//// Sequence number of this transaction corresponding to sender's account.
//pub sequence_number: u64,
//// The ChainID of the transaction.
//pub chain_id: u64,
//// The max gas to be used.
//pub max_gas_amount: u64,
//// The MoveAction to execute.
//pub action: MoveAction,
//}

type TransactionData struct {
	Sender         RoochAddress `json:"sender"`
	SequenceNumber uint64       `json:"sequence_number"`
	ChainId        uint64       `json:"chain_id"`
	MaxGasAmount   uint64       `json:"max_gas_amount"`
	Action         MoveAction   `json:"action"`
	//GasUsed      uint64 `json:"gas_used"`
	/// The vm status. If it is not `Executed`, this will provide the general error class. Execution
	/// failures and Move abort's receive more detailed information. But other errors are generally
	/// categorized with no status code or other information
	//pub status: KeptVMStatus,
	//Status string `json:"status"`
	//Status json.RawMessage `json:"status"`
}

func (rd *TransactionData) MarshalBCS(ser *bcs.Serializer) {
	rd.Sender.MarshalBCS(ser)
	ser.U64(rd.SequenceNumber)
	ser.U64(rd.ChainId)
	ser.U64(rd.MaxGasAmount)
	rd.Action.MarshalBCS(ser)
}
func (rd *TransactionData) UnmarshalBCS(des *bcs.Deserializer) {
	rd.Sender.UnmarshalBCS(des)
	rd.SequenceNumber = des.U64()
	rd.ChainId = des.U64()
	rd.MaxGasAmount = des.U64()
	rd.Action.UnmarshalBCS(des)
}

//hash(): Bytes {
//return sha3_256(this.encode())
//}

//Hash() ([]byte, error)
//bcs.Struct

func (rd *TransactionData) Hash() ([]byte, error) {
	bytes, err := bcs.Serialize(rd)
	if err != nil {
		return nil, errors.New("unable to serialize Transaction Data")
	}

	return utils.Sha3256(bytes), nil
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

////region TransactionPayload bcs.Struct
//
//func (txn *TransactionPayload) MarshalBCS(ser *bcs.Serializer) {
//	if txn == nil || txn.Payload == nil {
//		ser.SetError(fmt.Errorf("nil transaction payload"))
//		return
//	}
//	ser.Uleb128(uint32(txn.Payload.PayloadType()))
//	txn.Payload.MarshalBCS(ser)
//}
//func (txn *TransactionPayload) UnmarshalBCS(des *bcs.Deserializer) {
//	payloadType := TransactionPayloadVariant(des.Uleb128())
//	switch payloadType {
//	case TransactionPayloadVariantScript:
//		txn.Payload = &Script{}
//	case TransactionPayloadVariantModuleBundle:
//		// Deprecated, should never be in production
//		des.SetError(fmt.Errorf("module bundle is not supported as a transaction payload"))
//		return
//	case TransactionPayloadVariantEntryFunction:
//		txn.Payload = &EntryFunction{}
//	case TransactionPayloadVariantMultisig:
//		txn.Payload = &Multisig{}
//	default:
//		des.SetError(fmt.Errorf("bad txn payload kind, %d", payloadType))
//		return
//	}
//
//	txn.Payload.UnmarshalBCS(des)
//}

//endregion
//endregion

//region ModuleBundle

//// ModuleBundle is long deprecated and no longer used, but exist as an enum position in TransactionPayload
//type ModuleBundle struct{}
//
//func (txn *ModuleBundle) PayloadType() TransactionPayloadVariant {
//	return TransactionPayloadVariantModuleBundle
//}
//
//func (txn *ModuleBundle) MarshalBCS(ser *bcs.Serializer) {
//	ser.SetError(errors.New("ModuleBundle unimplemented"))
//}
//func (txn *ModuleBundle) UnmarshalBCS(des *bcs.Deserializer) {
//	des.SetError(errors.New("ModuleBundle unimplemented"))
//}

//endregion ModuleBundle

//region EntryFunction

//// EntryFunction call a single published entry function arguments are ordered BCS encoded bytes
//type EntryFunction struct {
//	Module   ModuleId
//	Function string
//	ArgTypes []TypeTag
//	Args     [][]byte
//}
//
////region EntryFunction TransactionPayloadImpl
//
//func (sf *EntryFunction) PayloadType() TransactionPayloadVariant {
//	return TransactionPayloadVariantEntryFunction
//}
//
////endregion
//
////region EntryFunction bcs.Struct
//
//func (sf *EntryFunction) MarshalBCS(ser *bcs.Serializer) {
//	sf.Module.MarshalBCS(ser)
//	ser.WriteString(sf.Function)
//	bcs.SerializeSequence(sf.ArgTypes, ser)
//	ser.Uleb128(uint32(len(sf.Args)))
//	for _, a := range sf.Args {
//		ser.WriteBytes(a)
//	}
//}
//func (sf *EntryFunction) UnmarshalBCS(des *bcs.Deserializer) {
//	sf.Module.UnmarshalBCS(des)
//	sf.Function = des.ReadString()
//	sf.ArgTypes = bcs.DeserializeSequence[TypeTag](des)
//	alen := des.Uleb128()
//	sf.Args = make([][]byte, alen)
//	for i := range alen {
//		sf.Args[i] = des.ReadBytes()
//	}
//}
//
////endregion
////endregion
//
////region Multisig
//
//// Multisig is an on-chain multisig transaction, that calls an entry function associated
//type Multisig struct {
//	MultisigAddress AccountAddress
//	Payload         *MultisigTransactionPayload // Optional
//}
//
////region Multisig TransactionPayloadImpl
//
//func (sf *Multisig) PayloadType() TransactionPayloadVariant {
//	return TransactionPayloadVariantMultisig
//}
//
////endregion
//
////region Multisig bcs.Struct
//
//func (sf *Multisig) MarshalBCS(ser *bcs.Serializer) {
//	ser.Struct(&sf.MultisigAddress)
//	if sf.Payload == nil {
//		ser.Bool(false)
//	} else {
//		ser.Bool(true)
//		ser.Struct(sf.Payload)
//	}
//}
//func (sf *Multisig) UnmarshalBCS(des *bcs.Deserializer) {
//	des.Struct(&sf.MultisigAddress)
//	if des.Bool() {
//		sf.Payload = &MultisigTransactionPayload{}
//		des.Struct(sf.Payload)
//	}
//}
//
////endregion
////endregion
//
////region MultisigTransactionPayload
//
//type MultisigTransactionPayloadVariant uint32
//
//const (
//	MultisigTransactionPayloadVariantEntryFunction MultisigTransactionPayloadVariant = 0
//)
//
//type MultisigTransactionImpl interface {
//	bcs.Struct
//}
//
//// MultisigTransactionPayload is an enum allowing for multiple types of transactions to be called via multisig
////
//// Note this does not implement TransactionPayloadImpl
//type MultisigTransactionPayload struct {
//	Variant MultisigTransactionPayloadVariant
//	Payload MultisigTransactionImpl
//}
//
////region MultisigTransactionPayload bcs.Struct
//
//func (sf *MultisigTransactionPayload) MarshalBCS(ser *bcs.Serializer) {
//	ser.Uleb128(uint32(sf.Variant))
//	ser.Struct(sf.Payload)
//}
//func (sf *MultisigTransactionPayload) UnmarshalBCS(des *bcs.Deserializer) {
//	variant := MultisigTransactionPayloadVariant(des.Uleb128())
//	switch variant {
//	case MultisigTransactionPayloadVariantEntryFunction:
//		sf.Payload = &EntryFunction{}
//	default:
//		des.SetError(fmt.Errorf("bad variant %d for MultisigTransactionPayload", variant))
//		return
//	}
//	des.Struct(sf.Payload)
//}

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

//type BlockHeader struct {
//	Timestamp            string  `json:"timestamp"`
//	Author               string  `json:"author"`
//	AuthorAuthKey        *string `json:"author_auth_key"`
//	BlockAccumulatorRoot string  `json:"block_accumulator_root"`
//	BlockHash            string  `json:"block_hash"`
//	BodyHash             string  `json:"body_hash"`
//	ChainId              int     `json:"chain_id"`
//	DifficultyHexStr     string  `json:"difficulty"`
//	Difficulty           uint64  `json:"difficulty_number"`
//	Extra                string  `json:"extra"`
//	GasUsed              string  `json:"gas_used"`
//	Nonce                uint64  `json:"Nonce"`
//	Height               string  `json:"number"`
//	ParentHash           string  `json:"parent_hash"`
//	StateRoot            string  `json:"state_root"`
//	TxnAccumulatorRoot   string  `json:"txn_accumulator_root"`
//}

//func (accumulator *AccumulatorInfo) ToTypesAccumulatorInfo() (*types.AccumulatorInfo, error) {
//	accumulatorRoot, err := hexToBytes(accumulator.AccumulatorRoot)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	subtreeRoots := make([]types.HashValue, 0)
//	for i := 0; i < len(accumulator.FrozenSubtreeRoots); i++ {
//		sub, err := hexToBytes(accumulator.FrozenSubtreeRoots[i])
//		if err != nil {
//			return nil, errors.WithStack(err)
//		}
//		subtreeRoots = append(subtreeRoots, sub)
//	}
//	nl, err := parseUint64(accumulator.NumLeaves)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	nn, err := parseUint64(accumulator.NumNodes)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	return &types.AccumulatorInfo{
//		AccumulatorRoot:    accumulatorRoot,
//		FrozenSubtreeRoots: subtreeRoots,
//		NumLeaves:          nl,
//		NumNodes:           nn,
//	}, nil
//}
//
////func (info *BlockInfo) ToTypesBlockInfo() (*types.BlockInfo, error) {
////	blockHash, err := hexToBytes(info.BlockHash)
////	if err != nil {
////		return nil, errors.WithStack(err)
////	}
////	blockAccumulatorInfo, err := info.BlockAccumulatorInfo.ToTypesAccumulatorInfo()
////	if err != nil {
////		return nil, errors.WithStack(err)
////	}
////	txnAccumulatorInfo, err := info.TxnAccumulatorInfo.ToTypesAccumulatorInfo()
////	if err != nil {
////		return nil, errors.WithStack(err)
////	}
////	diff := types.ToBcsDifficulty(info.TotalDifficulty)
////
////	return &types.BlockInfo{
////		BlockHash:            blockHash,
////		BlockAccumulatorInfo: *blockAccumulatorInfo,
////		TotalDifficulty:      diff,
////		TxnAccumulatorInfo:   *txnAccumulatorInfo,
////	}, nil
////}

//func (info *TransactionInfo) ToTypesTransactionInfo() (*types.TransactionInfo, error) {
//	txnHash, err := hexToBytes(info.TransactionHash)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	stateRootHash, err := hexToBytes(info.StateRootHash)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	eventRootHash, err := hexToBytes(info.EventRootHash)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	gasUsed, err := parseUint64(info.GasUsed)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	status, err := ToTypesKeptVMStatus(string(info.Status))
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	return &types.TransactionInfo{
//		TransactionHash: txnHash,
//		StateRootHash:   stateRootHash,
//		EventRootHash:   eventRootHash,
//		GasUsed:         gasUsed,
//		Status:          status,
//	}, nil
//}

//func ToTypesKeptVMStatus(status string) (types.KeptVMStatus, error) {
//	if len(status) < 1 {
//		return nil, fmt.Errorf("ToTypesKeptVMStatus status is null.")
//	}
//	if strings.EqualFold("\"Executed\"", status) {
//		return &types.KeptVMStatus__Executed{}, nil
//	}
//	if strings.EqualFold("\"OutOfGas\"", status) {
//		return &types.KeptVMStatus__OutOfGas{}, nil
//	}
//	//todo add other error parse
//	return &types.KeptVMStatus__MiscellaneousError{}, nil
//}

//func (header *BlockHeader) ToTypesHeader() (*types.BlockHeader, error) {
//	parentHash, err := hexToBytes(header.ParentHash)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	ts, err := parseUint64(header.Timestamp)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	number, err := parseUint64(header.Height)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	txnRoot, err := hexToBytes(header.TxnAccumulatorRoot)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	blockRoot, err := hexToBytes(header.BlockAccumulatorRoot)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	stateRoot, err := hexToBytes(header.StateRoot)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	gasUsed, err := parseUint64(header.GasUsed)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	bodyHash, err := hexToBytes(header.BodyHash)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	author, err := hexToAccountAddress(header.Author)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	var authorAuthKey *types.AuthenticationKey
//	if header.AuthorAuthKey != nil {
//		var a types.AuthenticationKey
//		a, err := hexToBytes(*header.AuthorAuthKey)
//		if err != nil {
//			return nil, errors.WithStack(err)
//		}
//		authorAuthKey = &a
//	}
//
//	diff := types.ToBcsDifficulty(header.DifficultyHexStr)
//
//	extra, err := hexTo4Uint8(header.Extra)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	return &types.BlockHeader{
//		ParentHash:           parentHash, //HashValue
//		Timestamp:            ts,
//		Number:               number,                                          //uint64
//		Author:               *author,                                         //AccountAddress
//		AuthorAuthKey:        authorAuthKey,                                   //*AuthenticationKey
//		TxnAccumulatorRoot:   txnRoot,                                         //HashValue
//		BlockAccumulatorRoot: blockRoot,                                       //HashValue
//		StateRoot:            stateRoot,                                       //HashValue
//		GasUsed:              gasUsed,                                         //uint64
//		Difficulty:           diff,                                            //[32]uint8
//		BodyHash:             bodyHash,                                        //HashValue
//		ChainId:              types.ChainId{Id: uint8(header.ChainId & 0xFF)}, //ChainId
//		Nonce:                uint32(header.Nonce),                            //uint32
//		Extra:                extra,                                           //type BlockHeaderExtra [4]uint8
//	}, nil
//}
//func (header *BlockHeader) Hash() ([]byte, error) {
//	h, err := header.ToTypesHeader()
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	hash, err := h.GetHash()
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	return *hash, nil
//}

//type BlockMetadata struct {
//	Author        string `json:"author"`
//	ChainID       string `json:"chain_id"`
//	Number        string `json:"number"`
//	ParentGasUsed int    `json:"parent_gas_used"`
//	ParentHash    string `json:"parent_hash"`
//	Timestamp     int64  `json:"timestamp"`
//	Uncles        string `json:"uncles"`
//}

//type Transaction struct {
//	BlockHash        string          `json:"block_hash"`
//	BlockNumber      string          `json:"block_number"`
//	TransactionHash  string          `json:"transaction_hash"`
//	TransactionIndex int             `json:"transaction_index"`
//	BlockMetadata    BlockMetadata   `json:"block_metadata"`
//	UserTransaction  UserTransaction `json:"user_transaction"`
//}
//
//type UserTransaction struct {
//	TransactionHash string         `json:"transaction_hash"`
//	RawTransaction  RawTransaction `json:"raw_txn"`
//	Authenticator   Authenticator  `json:"authenticator"`
//}

//type TransactionData struct {
//	Sender                  string `json:"sender"`
//	SequenceNumber          string `json:"sequence_number"`
//	Payload                 string `json:"payload"`
//	MaxGasAmount            string `json:"max_gas_amount"`
//	GasUnitPrice            string `json:"gas_unit_price"`
//	GasTokenCode            string `json:"gas_token_code"`
//	ExpirationTimestampSecs string `json:"expiration_timestamp_secs"`
//	ChainID                 int    `json:"chain_id"`
//}
//
////type Authenticator struct {
////	Ed25519 Ed25519 `json:"Ed25519"`
////}
////
//////pub struct Authenticator {
//////pub auth_validator_id: u64,
//////pub payload: Vec<u8>,
//////}
////
////func (rt *RoochTransaction) MarshalBCS(ser *bcs.Serializer) {
////	//ser.U64(uint64(rt.Data))
////	rt.Data.MarshalBCS(ser)
////	rt.Authenticator.MarshalBCS(ser)
////	//ser.WriteBytes(rt.Authenticator)
////	ser.WriteString(rt.DataHash)
////}
////func (rt *RoochTransaction) UnmarshalBCS(des *bcs.Deserializer) {
////	//rt.Data = MultiChainIDVariant(des.U64())
////	rt.Data.UnmarshalBCS(des)
////	rt.Authenticator.UnmarshalBCS(des)
////	//rt.Authenticator = des.ReadBytes()
////	//rt.Authenticator = des.ReadBytes()
////	rt.DataHash = des.ReadString()
////}
//
//type Ed25519 struct {
//	PublicKey string `json:"public_key"`
//	Signature string `json:"signature"`
//}
//
//type TransactionInfo struct {
//	BlockHash              string          `json:"block_hash"`
//	BlockNumber            string          `json:"block_number"`
//	TransactionHash        string          `json:"transaction_hash"`
//	TransactionIndex       int             `json:"transaction_index"`
//	TransactionGlobalIndex string          `json:"transaction_global_index"`
//	StateRootHash          string          `json:"state_root_hash"`
//	EventRootHash          string          `json:"event_root_hash"`
//	GasUsed                string          `json:"gas_used"`
//	Status                 json.RawMessage `json:"status"`
//}
//
//type Event struct {
//	BlockHash              string `json:"block_hash"`
//	BlockNumber            string `json:"block_number"`
//	TransactionHash        string `json:"transaction_hash"`
//	TransactionIndex       int    `json:"transaction_index"`
//	Data                   string `json:"data"`
//	TypeTag                string `json:"type_tag"`
//	EventKey               string `json:"event_key"`
//	EventSeqNumber         string `json:"event_seq_number"`
//	TransactionGlobalIndex string `json:"transaction_global_index"`
//	EventIndex             int    `json:"event_index"`
//}
//
//type TransactionProof struct {
//	TransactionInfo TransactionInfo `json:"transaction_info"`
//	Proof           struct {
//		Siblings []string `json:"siblings"`
//	} `json:"proof"`
//	EventProof struct {
//		Event string `json:"event"`
//		Proof struct {
//			Siblings []string `json:"siblings"`
//		} `json:"proof"`
//	} `json:"event_proof"`
//	StateProof json.RawMessage `json:"state_proof"` //todo??
//}
//
//type Block struct {
//	BlockHeader BlockHeader   `json:"header"`
//	BlockBody   BlockBody     `json:"body"`
//	Uncles      []BlockHeader `json:"uncles"`
//}
//
//type DryRunParam struct {
//	ChainId         int               `json:"chain_id"`
//	GasUnitPrice    int               `json:"gas_unit_price"`
//	Sender          string            `json:"sender"`
//	SenderPublicKey string            `json:"sender_public_key"`
//	SequenceNumber  *uint64           `json:"sequence_number,omitempty"`
//	MaxGasAmount    uint64            `json:"max_gas_amount"`
//	Script          DryRunParamScript `json:"script"`
//}
//
//type DryRunParamScript struct {
//	Code     string   `json:"code"`
//	TypeArgs []string `json:"type_args"`
//	Args     []string `json:"args"`
//}
//
//type DryRunResult struct {
//	ExplainedStatus string     `json:"explained_status"`
//	Events          []Event    `json:"events"`
//	GasUsed         string     `json:"gas_used"`
//	Status          string     `json:"status"`
//	WriteSet        []WriteSet `json:"write_set"`
//}
//
//type WriteSet struct {
//	AccessPath string `json:"access_path"`
//	Action     string `json:"action"`
//	Value      struct {
//		Resource struct {
//			Raw  string `json:"raw"`
//			JSON struct {
//				Fee struct {
//					Value int `json:"value"`
//				} `json:"fee"`
//			} `json:"json"`
//		} `json:"Resource"`
//	} `json:"value"`
//}
//
//func (block Block) GetHeader() (*types.BlockHeader, error) {
//	parentHash, err := HexStringToBytes(block.BlockHeader.ParentHash)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	ts, err := strconv.Atoi(block.BlockHeader.Timestamp)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	number, err := strconv.Atoi(block.BlockHeader.Height)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	author, err := types.ToAccountAddress(block.BlockHeader.Author)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	var authKey *types.AuthenticationKey
//	if block.BlockHeader.AuthorAuthKey != nil {
//		var a types.AuthenticationKey
//		a, err = HexStringToBytes(*block.BlockHeader.AuthorAuthKey)
//		if err != nil {
//			return nil, errors.WithStack(err)
//		}
//		authKey = &a
//	}
//	txnAccumulatorRoot, err := HexStringToBytes(block.BlockHeader.TxnAccumulatorRoot)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	blockAccumulatorRoot, err := HexStringToBytes(block.BlockHeader.BlockAccumulatorRoot)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	stateRoot, err := HexStringToBytes(block.BlockHeader.StateRoot)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	gasUsed, err := strconv.Atoi(block.BlockHeader.GasUsed)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	difficulty := types.ToBcsDifficulty(block.BlockHeader.DifficultyHexStr)
//
//	bodyHash, err := HexStringToBytes(block.BlockHeader.BodyHash)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	extra, err := HexStringToBytes(block.BlockHeader.Extra)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	var blockExtra [4]uint8
//	copy(blockExtra[:], extra[:4])
//
//	result := types.BlockHeader{
//		ParentHash:           parentHash,
//		Timestamp:            uint64(ts),
//		Number:               uint64(number),
//		Author:               *author,
//		AuthorAuthKey:        authKey,
//		TxnAccumulatorRoot:   txnAccumulatorRoot,
//		BlockAccumulatorRoot: blockAccumulatorRoot,
//		StateRoot:            stateRoot,
//		GasUsed:              uint64(gasUsed),
//		Difficulty:           difficulty,
//		BodyHash:             bodyHash,
//		ChainId:              types.ChainId{Id: uint8(block.BlockHeader.ChainId)},
//		Nonce:                uint32(block.BlockHeader.Nonce),
//		Extra:                blockExtra,
//	}
//	return &result, nil
//}
//
//func (block Block) GetHeaderHash() (*types.HashValue, error) {
//	header, err := block.GetHeader()
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	headerBytes, err := header.BcsSerialize()
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//
//	var result types.HashValue
//	result = types.Hash(types.PrefixHash("BlockHeader"), headerBytes)
//
//	return &result, nil
//}
//
//type BlockBody struct {
//	UserTransactions []UserTransaction `json:"Full"`
//}
//
//type Resource struct {
//	Raw string `json:"raw"`
//}
//
//type ListResource struct {
//	Resources map[string]Resource `json:"resources"`
//}
//
//type GetResourceOption struct {
//	Decode    bool    `json:"decode"`
//	StateRoot *string `json:"state_root,omitempty"`
//}
//
//type EpochResource struct {
//	Json *EpochJson `json:"json"`
//	Resource
//}
//
//type EpochJson struct {
//	Number                uint64 `json:"number"`                   //: 8385,
//	StartTime             uint64 `json:"start_time"`               //: 1637658309479,
//	StartBlockNumber      uint64 `json:"start_block_number"`       //: 2001384,
//	EndBlockNumber        uint64 `json:"end_block_number"`         //: 2001624,
//	BlockTimeTarget       uint64 `json:"block_time_target"`        //: 5260,
//	RewardPerBlock        uint64 `json:"reward_per_block"`         //: 5260000000,
//	RewardPerUnclePercent uint64 `json:"reward_per_uncle_percent"` //: 10,
//	BlockDifficutyWindow  uint   `json:"block_difficulty_window"`  //: 24,
//	MaxUnclesPerBlock     uint   `json:"max_uncles_per_block"`     //: 2,
//	BlockGasLimit         uint64 `json:"block_gas_limit"`          //: 50000000,
//	Strategy              uint   `json:"strategy"`                 //: 3,
//	NewEpochEvents        struct {
//		Counter uint64 `json:"counter"` //: 8385,
//		Guid    string `json:"guid"`    //: "0x090000000000000000000000000000000000000000000001"
//	} `json:"new_epoch_events"`
//}
//
//func (this ListResource) GetBalances() (map[string]*big.Int, error) {
//	result := make(map[string]*big.Int)
//
//	for k, v := range this.Resources {
//		if strings.Contains(k, "Balance") {
//			data, err := hex.DecodeString(strings.Replace(v.Raw, "0x", "", 1))
//			if err != nil {
//				return nil, errors.Wrap(err, "can't decode balance data")
//			}
//
//			balance, err := decode_u128_argument(data)
//			if err != nil {
//				return nil, errors.Wrap(err, "can't parse data to u128")
//			}
//
//			result[k] = U128ToBigInt(balance)
//		}
//	}
//
//	return result, nil
//}
//
//func (this ListResource) GetBalanceOfStc() (*big.Int, error) {
//	balances, err := this.GetBalances()
//
//	if err != nil {
//		return nil, errors.Wrap(err, "get stc balance error")
//	}
//
//	for k, v := range balances {
//		if k == "0x00000000000000000000000000000001::Account::Balance<0x00000000000000000000000000000001::STC::STC>" {
//			return v, nil
//		}
//	}
//
//	return big.NewInt(0), nil
//}
//
//type PendingTransaction struct {
//	Authenticator   Authenticator  `json:"authenticator"`
//	RawTransaction  RawTransaction `json:"raw_txn"`
//	Timestamp       int64          `json:"timestamp"`
//	TransactionHash string         `json:"transaction_hash"`
//}
//
//type Kind struct {
//	Type     int    `json:"type"`
//	TypeName string `json:"type_name"`
//}
//
//type EventFilter struct {
//	Address   []string `json:"addrs,omitempty"`
//	TypeTags  []string `json:"type_tags"`
//	FromBlock uint64   `json:"from_block"`
//	ToBlock   *uint64  `json:"to_block,omitempty"`
//	EventKeys []string `json:"event_keys"`
//}
//
////type NodeInfo struct {
////	PeerInfo    PeerInfo `json:"peer_info"`
////	SelfAddress string   `json:"self_address"`
////	Net         string   `json:"net"`
////	Consensus   struct {
////		Type string `json:"type"`
////	} `json:"consensus"`
////	NowSeconds int `json:"now_seconds"`
////}
////
////type PeerInfo struct {
////	PeerID         string    `json:"peer_id"`
////	ChainInfo      ChainInfo `json:"chain_info"`
////	NotifProtocols string    `json:"notif_protocols"`
////	RPCProtocols   string    `json:"rpc_protocols"`
////}
//
////type ChainInfo struct {
////	ChainID     int         `json:"chain_id"`
////	GenesisHash string      `json:"genesis_hash"`
////	Header      BlockHeader `json:"head"`
////	BlockInfo   BlockInfo   `json:"block_info"`
////}
//
//type BlockInfo struct {
//	BlockHash            string          `json:"block_hash"`
//	TotalDifficulty      string          `json:"total_difficulty"`
//	TxnAccumulatorInfo   AccumulatorInfo `json:"txn_accumulator_info"`
//	BlockAccumulatorInfo AccumulatorInfo `json:"block_accumulator_info"`
//}
//
//type RawResource struct {
//	Raw string `json:"raw"`
//}
//
//func (info NodeInfo) GetBlockNumber() (uint64, error) {
//	number, err := strconv.Atoi(info.PeerInfo.ChainInfo.Header.Height)
//	if err != nil {
//		return 0, errors.WithStack(err)
//	}
//
//	return uint64(number), nil
//}
//
//type ContractCall struct {
//	FunctionId string   `json:"function_id"`
//	TypeArgs   []string `json:"type_args"`
//	Args       []string `json:"args"`
//}
//
//func NewSendRecvEventFilters(addr string, fromBlock uint64) EventFilter {
//	addr = strings.ReplaceAll(addr, "0x", "")
//	eventKeys := []string{fmt.Sprintf("%s%s", recvPrefix, addr), fmt.Sprintf("%s%s", sendPrefix, addr)}
//	return EventFilter{
//		FromBlock: fromBlock,
//		EventKeys: eventKeys,
//	}
//}

//type ParseError struct{}
//
//func hexTo4Uint8(h string) ([4]uint8, error) {
//	var us [4]uint8
//	bs, err := hexToBytes(h)
//	if err != nil {
//		return us, errors.WithStack(err)
//	}
//	copy(us[:], bs[:4])
//	return us, nil
//}
//
//func hexTo32Uint8(h string) ([32]uint8, error) {
//	var us [32]uint8
//
//	bs, err := hexToBytes(h)
//	if err != nil {
//		return us, errors.WithStack(err)
//	}
//
//	if len(bs) > 32 {
//		copy(us[:], bs[:32])
//	} else {
//		copy(us[:], bs[:])
//	}
//	return us, nil
//}
//
//func hexToAccountAddress(addr string) (*types.AccountAddress, error) {
//	accountBytes, err := hexToBytes(addr)
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	var addressArray types.AccountAddress
//	copy(addressArray[:], accountBytes[:16])
//	return &addressArray, nil
//}
//
//func parseUint64(str string) (uint64, error) {
//	base := 10
//	if strings.HasPrefix(str, "0x") {
//		str = str[2:]
//		base = 16
//	}
//	i, err := strconv.ParseUint(str, base, 64)
//	if err != nil {
//		return 0, errors.WithStack(err)
//	}
//	return i, nil
//}
//
//func bytesToHex(b []byte) string {
//	return "0x" + hex.EncodeToString(b)
//}
//
//func hexToBytes(h string) ([]byte, error) {
//	var bs []byte
//	var err error
//	if !strings.HasPrefix(h, "0x") {
//		bs, err = hex.DecodeString(h)
//	} else {
//		bs, err = hex.DecodeString(h[2:])
//	}
//	if err != nil {
//		return nil, errors.WithStack(err)
//	}
//	return bs, nil
//}
