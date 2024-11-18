package api

import (
	"encoding/json"
	"github.com/rooch-network/rooch-go-sdk/utils"
	"regexp"
	"strings"
)

// U64 is a type for handling JSON string representations of the uint64
type U64 uint64

// UnmarshalJSON deserializes a JSON data blob into a [U64]
func (u *U64) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	uv, err := utils.StrToUint64(str)
	if err != nil {
		return err
	}
	*u = U64(uv)
	return nil
}

// ToUint64 converts a [U64] to an uint64
//
// We can guarantee that it's safe to convert a [U64] to an uint64 because we've already validated the input on JSON parsing.
func (u *U64) ToUint64() uint64 {
	return uint64(*u)
}

// HexBytes is a type for handling Bytes encoded as hex in JSON
type HexBytes []byte

// UnmarshalJSON deserializes a JSON data blob into a [HexBytes]
//
// Example:
//
//	"0x123456" -> []byte{0x12, 0x34, 0x56}
func (u *HexBytes) UnmarshalJSON(b []byte) error {
	var str string
	err := json.Unmarshal(b, &str)
	if err != nil {
		return err
	}
	bytes, err := utils.ParseHex(str)
	if err != nil {
		return err
	}
	*u = bytes
	return nil
}

// Hash is a representation of a hash as Hex in JSON
//
// # This is always represented as a 32-byte hash in hexadecimal format
//
// Example:
//
//	0xf4d07fdb8b5151971886a910e516d418a790dd5f6e068b0588066518a395a600
type Hash = string // TODO: do we make this a 32 byte array? or byte array?

var (
	VectorRegex = regexp.MustCompile(`^vector<(.+)>$`)
	StructRegex = regexp.MustCompile(`^([^:]+)::([^:]+)::([^<]+)(<(.+)>)?`)
)

func SplitGenericParameters(str string) []string {
	var result []string
	var current string
	var depth int

	for _, char := range str {
		switch char {
		case '<':
			depth++
			current += string(char)
		case '>':
			depth--
			current += string(char)
		case ',':
			if depth == 0 {
				result = append(result, strings.TrimSpace(current))
				current = ""
			} else {
				current += string(char)
			}
		default:
			current += string(char)
		}
	}

	if current != "" {
		result = append(result, strings.TrimSpace(current))
	}

	return result
}

//func Str2Uint64(str string) (uint64, error) {
//	i, err := strconv.ParseInt(str, 10, 64)
//	return uint64(i), err
//}
