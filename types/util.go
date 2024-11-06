package types

import (
	"github.com/rooch-network/rooch-go-sdk/address"
	"math/big"
)

// -- Note these are copied from internal/util/util.go to prevent package loops, but still allow devs to use it

// ParseHex Convenience function to deal with 0x at the beginning of hex strings
func ParseHex(hexStr string) ([]byte, error) {
	// This had to be redefined separately to get around a package loop
	return address.ParseHex(hexStr)
}

//// Sha3256Hash takes a hash of the given sets of bytes
//func Sha3256Hash(bytes [][]byte) (output []byte) {
//	return address.Sha3256Hash(bytes)
//}

//// Sha3256 hashes the input bytes using SHA3-256
//func Sha3256(input []byte) (output []byte) {
//	hasher := sha3.New256()
//	hasher.Write(input)
//	//return hasher.Sum(nil)
//	return hasher.Sum([]byte{})
//}
//
//// Blake2b256 hashes the input bytes using Blake2b256-256
//func Blake2b256(input []byte) (output []byte) {
//	hasher, _ := blake2b.New256(nil)
//	hasher.Write(input)
//	//hasher := w.Sum(nil)
//	return hasher.Sum(nil)
//}

// BytesToHex converts bytes to a 0x prefixed hex string
func BytesToHex(bytes []byte) string {
	return address.BytesToHex(bytes)
}

// StrToUint64 converts a string to a uint64
func StrToUint64(s string) (uint64, error) {
	return address.StrToUint64(s)
}

// StrToBigInt converts a string to a big.Int
func StrToBigInt(val string) (num *big.Int, err error) {
	return address.StrToBigInt(val)
}
