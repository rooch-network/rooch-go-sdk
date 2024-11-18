package types

import (
	"errors"
	"fmt"
	"github.com/rooch-network/rooch-go-sdk/utils"

	//"github.com/pkg/errors"

	"github.com/rooch-network/rooch-go-sdk/bcs"
)

const recvPrefix = "0100000000000000"
const sendPrefix = "0000000000000000"

//type Transaction struct {
//	Data TransactionData `json:"data"`
//	//Data          crypto.TransactionData `json:"data"`
//	Authenticator crypto.Authenticator `json:"authenticator"`
//	Info          string               `json:"info"`
//}
//
//func (t *Transaction) MarshalBCS(ser *bcs.Serializer) {
//	t.Data.MarshalBCS(ser)
//	t.Authenticator.MarshalBCS(ser)
//	ser.WriteString(t.Info)
//}
//func (t *Transaction) UnmarshalBCS(des *bcs.Deserializer) {
//	t.Data.UnmarshalBCS(des)
//	t.Authenticator.UnmarshalBCS(des)
//	t.Info = des.ReadString()
//}
//
//// GetInfo returns the transaction info
//func (t *Transaction) GetInfo() *string {
//	return &t.Info
//}
//
//// HashData returns the hash of the transaction data
//func (t *Transaction) HashData() ([]byte, error) {
//	return t.Data.Hash()
//}
//
//// getData returns the transaction data after validation
//func (t *Transaction) GetData() crypto.TransactionData {
//	//t.isValid()
//	return &t.Data
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
	TypeArgs   []TypeTag  `json:"ty_args"`
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
	bcs.SerializeSequence(fc.TypeArgs, ser)
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
	fc.TypeArgs = bcs.DeserializeSequence[TypeTag](des)
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
