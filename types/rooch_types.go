// Copyright (c) RoochNetwork
// SPDX-License-Identifier: Apache-2.0

package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
)

// Basic numeric types
type (
	U8  uint8
	U16 uint16
	U32 uint32
	U64 uint64
)

// Large numeric types using big.Int for u128 and u256
type (
	U128 *big.Int
	U256 *big.Int
)

// Boolean type
type Bool bool

// Object and address types
type (
	//ObjectID string
	Address interface{} // Can be string, *address.Address, or []byte
)

// Bytes is a type alias for []byte
type Bytes []byte

// EmptyBytes represents an empty byte slice
var EmptyBytes = make([]byte, 0)

func ConvertToBytes(input string, encoding string) ([]byte, error) {
	switch encoding {
	case "hex":
		return hex.DecodeString(input)
	case "utf8":
		return []byte(input), nil
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

func ConvertToVector(input interface{}, encoding string) ([]byte, error) {
	switch v := input.(type) {
	case string:
		return ConvertToBytes(v, encoding)
	case []byte:
		return v, nil
	default:
		return nil, errors.New("invalid input type for Vector")
	}
}
