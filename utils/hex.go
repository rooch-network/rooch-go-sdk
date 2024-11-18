package utils

import (
	"encoding/hex"
	"strings"
)

// IsHex checks if a string is a valid hexadecimal.
func IsHex(input interface{}) bool {
	switch v := input.(type) {
	case string:
		// Remove 0x prefix if present
		cleaned := strings.TrimPrefix(strings.ToLower(v), "0x")
		if len(cleaned)%2 != 0 {
			return false
		}
		_, err := hex.DecodeString(cleaned)
		return err == nil
	case []uint8:
		for _, b := range v {
			if b > 255 {
				return false
			}
		}
		return true
	default:
		return false
	}
}

// GetHexByteLength returns the byte length of a hex string.
func GetHexByteLength(input string) float64 {
	cleaned := strings.TrimPrefix(strings.ToLower(input), "0x")
	if cleaned == "" {
		return 0
	}
	return float64(len(cleaned)) / 2
}

// NormalizeHex removes the "0x" prefix from a hex string if present.
func NormalizeHex(input string) string {
	return strings.TrimPrefix(input, "0x")
}

// FromHEX converts a hex string to a byte array.
func FromHEX(input string) []uint8 {
	if len(input)%2 != 0 {
		input = "0" + input
	}
	decoded, err := hex.DecodeString(input)
	if err != nil {
		// Return array of zeros for invalid input
		return make([]uint8, len(input)/2)
	}
	return decoded
}

// ToHEX converts a byte array to a hex string.
func ToHEX(input []uint8) (string, error) {
	return hex.EncodeToString(input), nil
}

func HexToBytes(h string) ([]byte, error) {
	var bs []byte
	var err error
	if !strings.HasPrefix(h, "0x") {
		bs, err = hex.DecodeString(h)
	} else {
		bs, err = hex.DecodeString(h[2:])
	}
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func HexTo4Uint8(h string) ([4]uint8, error) {
	var us [4]uint8
	bs, err := HexToBytes(h)
	if err != nil {
		return us, err
	}
	copy(us[:], bs[:4])
	return us, nil
}

func HexTo32Uint8(h string) ([32]uint8, error) {
	var us [32]uint8

	bs, err := HexToBytes(h)
	if err != nil {
		return us, err
	}

	if len(bs) > 32 {
		copy(us[:], bs[:32])
	} else {
		copy(us[:], bs[:])
	}
	return us, nil
}
