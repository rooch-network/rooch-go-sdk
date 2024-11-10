package types

import (
	"encoding/hex"
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/bcs"
	"math/big"
)

// ArgType represents the type of argument
type ArgType string

const (
	ArgTypeU8       ArgType = "u8"
	ArgTypeU16      ArgType = "u16"
	ArgTypeU32      ArgType = "u32"
	ArgTypeU64      ArgType = "u64"
	ArgTypeU128     ArgType = "u128"
	ArgTypeU256     ArgType = "u256"
	ArgTypeBool     ArgType = "bool"
	ArgTypeString   ArgType = "string"
	ArgTypeObject   ArgType = "object"
	ArgTypeObjectID ArgType = "ObjectID"
	ArgTypeAddress  ArgType = "address"
)

// Args represents the BCS serialized value
type Args struct {
	value []byte
}

// NewArgs constructs a new Args instance
func NewArgs(input []byte) *Args {
	return &Args{value: input}
}

// Consistent with JS encodeWithHex()
func (a *Args) EncodeWithHex() string {
	return "0x" + hex.EncodeToString(a.value)
}

// Consistent with JS encode()
func (a *Args) Encode() []byte {
	return a.value
}

//func (xt *StructTag) MarshalBCS(ser *bcs.Serializer) {
//	xt.Address.MarshalBCS(ser)
//	ser.WriteString(xt.Module)
//	ser.WriteString(xt.Name)
//	bcs.SerializeSequence(xt.TypeParams, ser)
//}
//func (xt *StructTag) UnmarshalBCS(des *bcs.Deserializer) {
//	xt.Address.UnmarshalBCS(des)
//	xt.Module = des.ReadString()
//	xt.Name = des.ReadString()
//	xt.TypeParams = bcs.DeserializeSequence[TypeTag](des)
//}

// Consistent with JS static u8()
func ArgU8(input uint8) (*Args, error) {
	//serialized := bcs.NewSerializer().SerializeU8(input)
	bytes, err := bcs.SerializeU8(input)
	if err != nil {
		return nil, err
	}
	return NewArgs(bytes), nil
}

// Consistent with JS static u16()
func ArgU16(input uint16) (*Args, error) {
	//serialized := bcs.NewSerializer().SerializeU16(input)
	bytes, err := bcs.SerializeU16(input)
	if err != nil {
		return nil, err
	}
	return NewArgs(bytes), nil
}

// Consistent with JS static u32()
func ArgU32(input uint32) (*Args, error) {
	//serialized := bcs.NewSerializer().SerializeU32(input)
	//return NewArgs(serialized.Bytes())

	bytes, err := bcs.SerializeU32(input)
	if err != nil {
		return nil, err
	}
	return NewArgs(bytes), nil
}

// Consistent with JS static u64()
func ArgU64(input uint64) (*Args, error) {
	//serialized := bcs.NewSerializer().SerializeU64(input)
	//return NewArgs(serialized.Bytes())

	bytes, err := bcs.SerializeU64(input)
	if err != nil {
		return nil, err
	}
	return NewArgs(bytes), nil
}

// Consistent with JS static u128()
func ArgU128(input big.Int) (*Args, error) {
	//serialized := bcs.NewSerializer().SerializeU128(input)
	//return NewArgs(serialized.Bytes())

	bytes, err := bcs.SerializeU128(input)
	if err != nil {
		return nil, err
	}
	return NewArgs(bytes), nil
}

// Consistent with JS static u256()
func ArgU256(input big.Int) (*Args, error) {
	//serialized := bcs.NewSerializer().SerializeU256(input)
	//return NewArgs(serialized.Bytes())

	bytes, err := bcs.SerializeU256(input)
	if err != nil {
		return nil, err
	}
	return NewArgs(bytes), nil
}

// Consistent with JS static bool()
func ArgBool(input bool) (*Args, error) {
	//serialized := bcs.NewSerializer().SerializeBool(input)
	//return NewArgs(serialized.Bytes())

	bytes, err := bcs.SerializeBool(input)
	if err != nil {
		return nil, err
	}
	return NewArgs(bytes), nil
}

// Consistent with JS static string()
func ArgString(input string) (*Args, error) {
	//serialized := bcs.NewSerializer().SerializeString(input)
	//return NewArgs(serialized.Bytes())

	bytes, err := bcs.SerializeString(input)
	if err != nil {
		return nil, err
	}
	return NewArgs(bytes), nil
}

// Consistent with JS static address()
func ArgAddress(input interface{}) (*Args, error) {
	//serialized := Address{}.Serialize(input)
	//return NewArgs(serialized)

	roochAddress, err := address.ConvertToRoochAddress(input)
	if err != nil {
		return nil, err
	}
	serializer := &bcs.Serializer{}
	roochAddress.MarshalBCS(serializer)

	bytes := serializer.ToBytes()
	return NewArgs(bytes), nil
}

// Consistent with JS static object()
func ArgObject(input StructTag) (*Args, error) {
	objectID, err := StructTagToObjectID(&input)
	if err != nil {
		return nil, err
	}
	serializer := &bcs.Serializer{}
	objectID.MarshalBCS(serializer)

	bytes := serializer.ToBytes()
	return NewArgs(bytes), nil
}

// Consistent with JS static ObjectID()
func ArgObjectID(input string) (*Args, error) {
	objectID, err := ConvertObjectID(input)
	if err != nil {
		return nil, err
	}

	serializer := &bcs.Serializer{}
	objectID.MarshalBCS(serializer)

	bytes := serializer.ToBytes()
	return NewArgs(bytes), nil
}

// Consistent with JS static struct()
func ArgStruct(input interface{}) (*Args, error) {
	var bytes []byte
	switch v := input.(type) {
	case []byte:
		bytes = v
	case StructTag:
		ser := &bcs.Serializer{}
		struct_tag := input.(StructTag)
		struct_tag.MarshalBCS(ser)
		bytes = ser.ToBytes()
	}
	return NewArgs(bytes), nil
}

// Consistent with JS static vec()
func ArgVec(argType ArgType, input interface{}) (*Args, error) {
	ser := &bcs.Serializer{}

	switch argType {
	case ArgTypeU8:
		bcs.SerializeSequence(input.([]uint8), ser)
		//bcs.SerializeSequence(input.([]uint8))
	case ArgTypeU16:
		bcs.SerializeSequence(input.([]uint16), ser)
	case ArgTypeU32:
		bcs.SerializeSequence(input.([]uint32), ser)
	case ArgTypeU64:
		bcs.SerializeSequence(input.([]uint64), ser)
	case ArgTypeU128:
		bcs.SerializeSequence(input.([][]byte), ser)
	case ArgTypeU256:
		bcs.SerializeSequence(input.([][]byte), ser)
	case ArgTypeBool:
		bcs.SerializeSequence(input.([]bool), ser)
	case ArgTypeString:
		bcs.SerializeSequence(input.([]string), ser)
	case ArgTypeObject:
		structTags := input.([]StructTag)
		objectIDs := make([]ObjectID, len(structTags))
		for i, tag := range structTags {
			obj_id, err := StructTagToObjectID(&tag)
			if err != nil {
				return nil, err
			}
			objectIDs[i] = obj_id
		}
		bcs.SerializeSequence(objectIDs, ser)
	case ArgTypeObjectID:
		ids := input.([]interface{})
		objectIDs := make([]ObjectID, len(ids))
		for i, id := range ids {
			obj_id, err := ConvertObjectID(id)
			if err != nil {
				return nil, err
			}
			objectIDs[i] = obj_id
		}
		bcs.SerializeSequence(objectIDs, ser)
	case ArgTypeAddress:
		//bcs.SerializeSequence(Address{}.SerializeVec(input.([]string)), ser)
		ids := input.([]interface{})
		addresses := make([]RoochAddress, len(ids))
		for i, id := range ids {
			address, err := address.ConvertToRoochAddress(id)
			if err != nil {
				return nil, err
			}
			addresses[i] = *address
		}
		bcs.SerializeSequence(addresses, ser)
	}

	return NewArgs(ser.ToBytes()), nil
}

//// Helper interfaces
//type BcsSerializable interface {
//	Serialize() []byte
//}
//
//// bcs package interface (implement according to your BCS package)
//type Serializer struct {
//	data []byte
//}
//
//func (s *Serializer) SerializeU8(v uint8) *Serializer {
//	s.data = append(s.data, v)
//	return s
//}
//
//func (s *Serializer) SerializeU16(v uint16) *Serializer {
//	// Implement U16 serialization
//	return s
//}
//
//func (s *Serializer) SerializeU32(v uint32) *Serializer {
//	// Implement U32 serialization
//	return s
//}
//
//func (s *Serializer) SerializeU64(v uint64) *Serializer {
//	// Implement U64 serialization
//	return s
//}
//
//func (s *Serializer) SerializeU128(v []byte) *Serializer {
//	// Implement U128 serialization
//	return s
//}
//
//func (s *Serializer) SerializeU256(v []byte) *Serializer {
//	// Implement U256 serialization
//	return s
//}
//
//func (s *Serializer) SerializeBool(v bool) *Serializer {
//	if v {
//		s.data = append(s.data, 1)
//	} else {
//		s.data = append(s.data, 0)
//	}
//	return s
//}
//
//func (s *Serializer) SerializeString(v string) *Serializer {
//	// Implement string serialization
//	return s
//}
//
//func (s *Serializer) SerializeVec(v interface{}) *Serializer {
//	// Implement vector serialization based on type
//	return s
//}
//
//func (s *Serializer) Bytes() []byte {
//	return s.data
//}
//
//func NewSerializer() *Serializer {
//	return &Serializer{
//		data: make([]byte, 0),
//	}
//}
//```
//
//Example usage:
//
//```go
//func Example() {
//	// Create u8 argument
//	u8Arg := U8(123)
//
//	// Create string argument
//	strArg := String("hello")
//
//	// Create vector argument
//	vecArg := Vec(ArgTypeU8, []uint8{1, 2, 3})
//
//	// Get hex encoding
//	hexStr := u8Arg.EncodeWithHex()
//
//	// Get raw bytes
//	rawBytes := strArg.Encode()
//
//	// Create address argument
//	addrArg := Address("0x123...")
//
//	// Create vector of addresses
//	addrVec := Vec(ArgTypeAddress, []string{"0x123...", "0x456..."})
//}
//```
