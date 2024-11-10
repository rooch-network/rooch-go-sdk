package address

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/rooch-network/rooch-go-sdk/bcs"
	"github.com/rooch-network/rooch-go-sdk/crypto"
	"strings"
)

var (
	ErrInvalidAddress     = errors.New("invalid address")
	ErrInvalidAddressLen  = errors.New("invalid address length")
	ErrInvalidHexAddress  = errors.New("invalid hex address")
	ErrInvalidBech32      = errors.New("invalid bech32 address")
	ErrUnsupportedAddress = errors.New("unsupported address type")
)

// RoochAddress a 32-byte representation of an on-chain address
//
// Implements:
//   - [bcs.Marshaler]
//   - [bcs.Unmarshaler]
//   - [json.Marshaler]
//   - [json.Unmarshaler]
//
// type RoochAddress H256
type RoochAddress [32]byte

const RoochBech32Prefix = "rooch"

const RoochAddressLength = 32

//const DEFAULT_MAX_GAS_AMOUNT = 10000000
//const GAS_TOKEN_CODE = "0x3::gas_coin::RGas"

// AddressZero is [RoochAddress] 0x0
var AddressZero = RoochAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// AddressOne is [RoochAddress] 0x1
var AddressOne = RoochAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

// AddressTwo is [RoochAddress] 0x2
var AddressTwo = RoochAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}

// AddressThree is [RoochAddress] 0x3
var AddressThree = RoochAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3}

// AddressFour is [RoochAddress] 0x4
var AddressFour = RoochAddress{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4}

// IsSpecial Returns whether the address is a "special" address. Addresses are considered
// special if the first 63 characters of the hex string are zero. In other words,
// an address is special if the first 31 bytes are zero and the last byte is
// smaller than `0b10000` (16). In other words, special is defined as an address
// that matches the following regex: `^0x0{63}[0-9a-f]$`. In short form this means
// the addresses in the range from `0x0` to `0xf` (inclusive) are special.
// For more details see the v1 address standard defined as part of AIP-40:
func (ra *RoochAddress) IsSpecial() bool {
	for _, b := range ra[:31] {
		if b != 0 {
			return false
		}
	}
	return ra[31] < 0x10
}

// String Returns the canonical string representation of the [RoochAddress]
//
// Please use [RoochAddress.StringLong] for all indexer queries.
func (ra *RoochAddress) String() string {
	if ra.IsSpecial() {
		return fmt.Sprintf("0x%x", ra[31])
	} else {
		//	return "0x" + hex.EncodeToString(a.Bytes())
		return BytesToHex(ra[:])
	}
}

// FromAuthKey converts [crypto.AuthenticationKey] to [RoochAddress]
func (ra *RoochAddress) FromAuthKey(authKey *crypto.AuthenticationKey) {
	copy(ra[:], authKey[:])
}

// AuthKey converts [RoochAddress] to [crypto.AuthenticationKey]
func (ra *RoochAddress) AuthKey() *crypto.AuthenticationKey {
	authKey := &crypto.AuthenticationKey{}
	copy(authKey[:], ra[:])
	return authKey
}

// StringLong Returns the long string representation of the RoochAddress
//
// This is most commonly used for all indexer queries.
func (ra *RoochAddress) StringLong() string {
	return BytesToHex(ra[:])
}

// MarshalBCS Converts the RoochAddress to BCS encoded bytes
func (ra *RoochAddress) MarshalBCS(ser *bcs.Serializer) {
	//ser.FixedBytes(ra[:])
	ser.WriteBytes(ra[:])
}

// UnmarshalBCS Converts the RoochAddress from BCS encoded bytes
func (ra *RoochAddress) UnmarshalBCS(des *bcs.Deserializer) {
	//des.ReadFixedBytesInto((*ra)[:])
	//des.ReadFixedBytesInto((*ra)[:])
	des.ReadBytes()
}

// MarshalJSON converts the RoochAddress to JSON
func (ra *RoochAddress) MarshalJSON() ([]byte, error) {
	return json.Marshal(ra.String())
}

//// UnmarshalJSON converts the RoochAddress from JSON
//func (ra *RoochAddress) UnmarshalJSON(b []byte) error {
//	var str string
//	err := json.Unmarshal(b, &str)
//	if err != nil {
//		return fmt.Errorf("failed to convert input to AccountAdddress: %w", err)
//	}
//	err = ra.ParseStringRelaxed(str)
//	if err != nil {
//		return fmt.Errorf("failed to convert input to AccountAdddress: %w", err)
//	}
//	return nil
//}

//// NamedObjectAddress derives a named object address based on the input address as the creator
//func (ra *RoochAddress) NamedObjectAddress(seed []byte) (accountAddress RoochAddress) {
//	return ra.DerivedAddress(seed, crypto.NamedObjectScheme)
//}
//
//// ObjectAddressFromObject derives an object address based on the input address as the creator object
//func (ra *RoochAddress) ObjectAddressFromObject(objectAddress *RoochAddress) (accountAddress RoochAddress) {
//	return ra.DerivedAddress(objectAddress[:], crypto.DeriveObjectScheme)
//}
//
//// ResourceAccount derives an object address based on the input address as the creator
//func (ra *RoochAddress) ResourceAccount(seed []byte) (accountAddress RoochAddress) {
//	return ra.DerivedAddress(seed, crypto.ResourceAccountScheme)
//}

// DerivedAddress addresses are derived by the address, the seed, then the type byte
func (ra *RoochAddress) DerivedAddress(seed []byte, typeByte uint8) (accountAddress RoochAddress) {
	authKey := ra.AuthKey()
	authKey.FromBytesAndScheme(append(authKey[:], seed[:]...), typeByte)
	copy(accountAddress[:], authKey[:])
	return
}

//var (
//	ErrInvalidAddress     = errors.New("invalid address")
//	ErrInvalidAddressLen  = errors.New("invalid address length")
//	ErrInvalidHexAddress  = errors.New("invalid hex address")
//	ErrInvalidBech32      = errors.New("invalid bech32 address")
//	ErrUnsupportedAddress = errors.New("unsupported address type")
//)
//
//// Address represents a Rooch address which can be in different formats
//type Address interface {
//	String() string
//	Bytes() []byte
//}

//// RoochAddress implements the Address interface
//type RoochAddress struct {
//	bytes []byte
//}

//// BitcoinAddress implements the Address interface
//type BitcoinAddress struct {
//	address string
//}

// NewRoochAddress creates a new RoochAddress from bytes
func NewRoochAddress(bytes []byte) (*RoochAddress, error) {
	if len(bytes) != RoochAddressLength {
		return nil, ErrInvalidAddressLen
	}
	return (*RoochAddress)(bytes), nil
}

// NewRoochAddressFromString creates a RoochAddress from a string
func NewRoochAddressFromString(address string) (*RoochAddress, error) {
	// Try to decode as hex
	if strings.HasPrefix(strings.ToLower(address), "0x") {
		bytes, err := hex.DecodeString(address[2:])
		if err != nil {
			return nil, ErrInvalidHexAddress
		}
		return NewRoochAddress(bytes)
	}

	// Try to decode as bech32
	if strings.HasPrefix(address, RoochBech32Prefix) {
		_, data, err := bech32.Decode(address)
		if err != nil {
			return nil, ErrInvalidBech32
		}
		converted, err := bech32.ConvertBits(data, 5, 8, true)
		if err != nil {
			return nil, err
		}
		return NewRoochAddress(converted)
	}

	//// Try to convert from Bitcoin address
	//btcAddr, err := NewBitcoinAddress(address)
	//if err == nil {
	//	return btcAddr.ToRoochAddress()
	//}

	return nil, ErrInvalidAddress
}

//// String returns the hex representation of the address
//func (a *RoochAddress) String() string {
//	return "0x" + hex.EncodeToString(a.Bytes())
//}

// Bytes returns the raw bytes of the address
func (a *RoochAddress) Bytes() []byte {
	//return ([]byte)(a)
	return a[:]
}

// ToBech32 converts the address to bech32 format
func (a *RoochAddress) ToBech32() (string, error) {
	converted, err := bech32.ConvertBits(a.Bytes(), 8, 5, true)
	if err != nil {
		return "", err
	}
	encoded, err := bech32.Encode(RoochBech32Prefix, converted)
	if err != nil {
		return "", err
	}
	return encoded, nil
}

//// NewBitcoinAddress creates a new BitcoinAddress
//func NewBitcoinAddress(address string) (*BitcoinAddress, error) {
//	// Here you would implement Bitcoin address validation
//	// This is a simplified version
//	if !isValidBitcoinAddress(address) {
//		return nil, ErrInvalidAddress
//	}
//	return &BitcoinAddress{address: address}, nil
//}

//// String returns the Bitcoin address string
//func (a *BitcoinAddress) String() string {
//	return a.address
//}
//
//// Bytes returns the raw bytes of the Bitcoin address
//func (a *BitcoinAddress) Bytes() []byte {
//	// This is a placeholder - you would implement actual Bitcoin address parsing
//	return []byte(a.address)
//}
//
//// ToRoochAddress converts a Bitcoin address to a Rooch address
//func (a *BitcoinAddress) ToRoochAddress() (*RoochAddress, error) {
//	// This is a placeholder - you would implement actual conversion logic
//	// For now, we'll create a dummy address of the correct length
//	bytes := make([]byte, RoochAddressLength)
//	// TODO: Implement actual conversion logic
//	return NewRoochAddress(bytes)
//}
//
//// Helper functions
//
//func isValidBitcoinAddress(address string) bool {
//	// This is a placeholder - you would implement actual Bitcoin address validation
//	// Could use btcd/btcutil or other Bitcoin libraries
//	return len(address) > 0 // Dummy validation
//}

//// NormalizeAddress normalizes a Rooch address string
//func NormalizeAddress(address string, forceAdd0x bool) string {
//	address = strings.ToLower(address)
//	if !forceAdd0x && strings.HasPrefix(address, "0x") {
//		address = address[2:]
//	}
//	// Pad with zeros to correct length
//	paddedLen := RoochAddressLength * 2 // Each byte is 2 hex chars
//	if len(address) < paddedLen {
//		address = strings.Repeat("0", paddedLen-len(address)) + address
//	}
//	return "0x" + address
//}

// NormalizeRoochAddress normalizes a Rooch address
func NormalizeRoochAddress(input string, forceAdd0x bool) string {
	addr := strings.ToLower(input)
	if strings.HasPrefix(addr, "0x") {
		addr = addr[2:]
	}
	targetLen := RoochAddressLength * 2
	addr = strings.Repeat("0", targetLen-len(addr)) + addr
	if forceAdd0x {
		addr = "0x" + addr
	}
	return addr
}

// ToCanonicalRoochAddress returns the canonical form of a Rooch address
func ToCanonicalRoochAddress(input string, forceAdd0x bool) string {
	return NormalizeRoochAddress(input, forceAdd0x)
}

// IsValidRoochAddress checks if the given address is a valid Rooch address
func IsValidRoochAddress(address string) bool {
	_, err := NewRoochAddressFromString(address)
	return err == nil
}

// ConvertToRoochAddress converts various input formats to a RoochAddress
func ConvertToRoochAddress(input interface{}) (*RoochAddress, error) {
	switch v := input.(type) {
	case string:
		return NewRoochAddressFromString(v)
	case []byte:
		return NewRoochAddress(v)
	case RoochAddress:
		return NewRoochAddress(v.Bytes())
	default:
		return nil, ErrUnsupportedAddress
	}
}

// Utility functions for testing
func BytesEqual(a, b []byte) bool {
	return bytes.Equal(a, b)
}
