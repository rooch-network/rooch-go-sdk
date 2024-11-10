package types

import (
	"encoding/hex"
	"fmt"
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/bcs"
)

//const ROOCH_ADDRESS_LENGTH = 32

type ObjectID struct {
	Address []RoochAddress
}

func NewObjectID(address []RoochAddress) ObjectID {
	return ObjectID{
		Address: address,
	}
}

// ConvertObjectID handles the input
func ConvertObjectID(input interface{}) (ObjectID, error) {
	switch input.(type) {
	case string, []byte:
		address, err := address.ConvertToRoochAddress(input)
		if err != nil {
			return ObjectID{}, fmt.Errorf("failed to convert address: %w", err)
		}

		return NewObjectID([]RoochAddress{*address}), nil
	case [][]byte:
		var addresses []RoochAddress
		bytes, _ := (input).([][]byte)
		for _, value := range bytes {
			address, err := address.ConvertToRoochAddress(value)
			if err != nil {
				return ObjectID{}, fmt.Errorf("failed to convert address: %w", err)
			}
			addresses = append(addresses, *address)
		}

		return NewObjectID(addresses), nil

	default:
		return ObjectID{}, fmt.Errorf("unsupported input type for ObjectID")
	}
}

// Join handles the output transformation (equivalent to transform.output in JS)
func (o *ObjectID) String() string {
	totalLen := len(o.Address) * address.RoochAddressLength
	bytes := make([]byte, totalLen)

	offset := 0
	for _, addr := range o.Address {
		copy(bytes[offset:], addr.Bytes())
		offset += address.RoochAddressLength
	}

	return "0x" + hex.EncodeToString(bytes)
}

// Helper function for hex conversion
func fromHEX(hexStr string) ([]byte, error) {
	// Remove 0x prefix if present
	if len(hexStr) >= 2 && hexStr[0:2] == "0x" {
		hexStr = hexStr[2:]
	}

	return hex.DecodeString(hexStr)
}

func (o *ObjectID) MarshalBCS(ser *bcs.Serializer) {
	bcs.SerializeSequence(o.Address, ser)
}
func (o *ObjectID) UnmarshalBCS(des *bcs.Deserializer) {
	o.Address = bcs.DeserializeSequence[RoochAddress](des)
}
