package common

import (
	"bytes"
	"encoding/base32"
	"github.com/btcsuite/btcd/btcutil/base58"

	//"encoding/base58"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

type CodecType string

const (
	UTF8      CodecType = "utf8"
	HEX       CodecType = "hex"
	BASE16    CodecType = "base16"
	BASE32    CodecType = "base32"
	BASE64    CodecType = "base64"
	BASE64URL CodecType = "base64url"
	BASE58    CodecType = "base58"
	BASE58XMR CodecType = "base58xmr"
)

var (
	coders = map[CodecType]struct{}{
		UTF8:      {},
		HEX:       {},
		BASE16:    {},
		BASE32:    {},
		BASE64:    {},
		BASE64URL: {},
		BASE58:    {},
		BASE58XMR: {},
	}
	coderTypeError = "Invalid encoding type. Available types: utf8, hex, base16, base32, base64, base64url, base58, base58xmr"
)

// bytesEqual matches JS 'bytesEqual'
func BytesEqual(a, b []byte) bool {
	if &a == &b {
		return true
	}
	return bytes.Equal(a, b)
}

// IsBytes matches JS 'IsBytes'
func IsBytes(a interface{}) bool {
	_, ok := a.([]byte)
	return ok
}

// BytesToString matches JS 'BytesToString'
func BytesToString(typ CodecType, bytes []byte) (string, error) {
	if _, ok := coders[typ]; !ok {
		return "", fmt.Errorf(coderTypeError)
	}
	if !IsBytes(bytes) {
		return "", fmt.Errorf("BytesToString() expects Uint8Array")
	}

	switch typ {
	case UTF8:
		return string(bytes), nil
	case HEX, BASE16:
		return hex.EncodeToString(bytes), nil
	case BASE32:
		return base32.StdEncoding.EncodeToString(bytes), nil
	case BASE64:
		return base64.StdEncoding.EncodeToString(bytes), nil
	case BASE64URL:
		return base64.URLEncoding.EncodeToString(bytes), nil
	case BASE58:
		//return base58.base58Encode(bytes), nil
		return base58.Encode(bytes), nil
	//case BASE58XMR:
	//	return base58.XMREncode(bytes), nil
	default:
		return "", fmt.Errorf(coderTypeError)
	}
}

// str matches JS 'str'
func Str(typ CodecType, bytes []byte) (string, error) {
	return BytesToString(typ, bytes)
}

// StringToBytes matches JS 'StringToBytes'
func StringToBytes(typ CodecType, str string) ([]byte, error) {
	if _, ok := coders[typ]; !ok {
		return nil, fmt.Errorf(coderTypeError)
	}

	switch typ {
	case UTF8:
		return []byte(str), nil
	case HEX, BASE16:
		return hex.DecodeString(str)
	case BASE32:
		return base32.StdEncoding.DecodeString(str)
	case BASE64:
		return base64.StdEncoding.DecodeString(str)
	case BASE64URL:
		return base64.URLEncoding.DecodeString(str)
	case BASE58:
		return base58.Decode(str), nil
	//case BASE58XMR:
	//	return base58XMRDecode(str)
	default:
		return nil, fmt.Errorf(coderTypeError)
	}
}

// ConcatBytes matches JS 'ConcatBytes'
func ConcatBytes(arrays ...[]byte) ([]byte, error) {
	var totalLen int
	for _, arr := range arrays {
		if !IsBytes(arr) {
			return nil, fmt.Errorf("Uint8Array expected")
		}
		totalLen += len(arr)
	}

	result := make([]byte, totalLen)
	var offset int
	for _, arr := range arrays {
		copy(result[offset:], arr)
		offset += len(arr)
	}
	return result, nil
}

// VarintByteNum matches JS 'VarintByteNum'
func VarintByteNum(input uint64) []byte {
	if input < 253 {
		return []byte{byte(input)}
	} else if input < 0x10000 {
		buf := make([]byte, 3)
		buf[0] = 253
		buf[1] = byte(input)
		buf[2] = byte(input >> 8)
		return buf
	} else if input < 0x100000000 {
		buf := make([]byte, 5)
		buf[0] = 254
		buf[1] = byte(input)
		buf[2] = byte(input >> 8)
		buf[3] = byte(input >> 16)
		buf[4] = byte(input >> 24)
		return buf
	} else {
		buf := make([]byte, 9)
		buf[0] = 255
		// Low 32 bits
		buf[1] = byte(input)
		buf[2] = byte(input >> 8)
		buf[3] = byte(input >> 16)
		buf[4] = byte(input >> 24)
		// High 32 bits
		buf[5] = byte(input >> 32)
		buf[6] = byte(input >> 40)
		buf[7] = byte(input >> 48)
		buf[8] = byte(input >> 56)
		return buf
	}
}
