package address

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"golang.org/x/crypto/blake2b"
)

// BitcoinNetworkType represents different Bitcoin networks
type BitcoinNetworkType int

const (
	BitcoinNetworkBitcoin BitcoinNetworkType = iota
	BitcoinNetworkTestnet
	BitcoinNetworkSignet
	BitcoinNetworkRegtest
)

// BitcoinAddressType represents different Bitcoin address types
type BitcoinAddressType int

const (
	PKH     BitcoinAddressType = 0
	SH      BitcoinAddressType = 1
	WITNESS BitcoinAddressType = 2
)

const (
	PubkeyAddressPrefixMain = 0x00
	PubkeyAddressPrefixTest = 0x6F
	ScriptAddressPrefixMain = 0x05
	ScriptAddressPrefixTest = 0xC4
)

// BitcoinNetwork represents a Bitcoin network configuration
type BitcoinNetwork struct {
	network BitcoinNetworkType
}

// NewBitcoinNetwork creates a new BitcoinNetwork instance
func NewBitcoinNetwork(network BitcoinNetworkType) *BitcoinNetwork {
	return &BitcoinNetwork{network: network}
}

// FromBech32Prefix creates a BitcoinNetwork from a bech32 prefix
func FromBech32Prefix(prefix string) *BitcoinNetwork {
	switch prefix {
	case "bc":
		return NewBitcoinNetwork(BitcoinNetworkBitcoin)
	case "tb":
		return NewBitcoinNetwork(BitcoinNetworkTestnet)
	case "bcrt":
		return NewBitcoinNetwork(BitcoinNetworkRegtest)
	default:
		return nil
	}
}

// Bech32HRP returns the human-readable prefix for the network
func (bn *BitcoinNetwork) Bech32HRP() string {
	switch bn.network {
	case BitcoinNetworkBitcoin:
		return "bc"
	case BitcoinNetworkTestnet:
		return "tb"
	case BitcoinNetworkSignet:
		return "tb"
	case BitcoinNetworkRegtest:
		return "bcrt"
	default:
		return ""
	}
}

// AddressInfo represents decoded address information
type BitcoinAddressInfo struct {
	Bytes   []byte
	Type    BitcoinAddressType
	Version uint8
}

//// BitcoinAddress represents a Bitcoin address
//type BitcoinAddress struct {
//	rawAddress string
//	bytes      []byte
//	network    BitcoinNetworkType
//}

// BitcoinAddress represents a Bitcoin address
type BitcoinAddress struct {
	bytes        []byte
	rawAddress   string
	roochAddress *RoochAddress
}

// NewBitcoinAddress creates a new BitcoinAddress instance
func NewBitcoinAddress(input string, network BitcoinNetworkType) (*BitcoinAddress, error) {
	ba := &BitcoinAddress{rawAddress: input}

	if isHex(input) {
		// Handle hex input
		hexStr := input
		if len(input) >= 2 && input[:2] == "0x" {
			hexStr = input[2:]
		}
		decoded, err := hex.DecodeString(hexStr)
		if err != nil {
			return nil, err
		}
		ba.bytes = decoded

		var prefixed []byte
		version := ba.bytes[1]

		switch BitcoinAddressType(ba.bytes[0]) {
		case PKH:
			prefixed = make([]byte, 22)
			prefixed[0] = version
			prefixed[1] = ba.GetPubkeyAddressPrefix(network)
			copy(prefixed[2:], ba.bytes[2:])
			ba.rawAddress = base58.CheckEncode(prefixed[1:], prefixed[0])

		case SH:
			prefixed = make([]byte, 22)
			prefixed[0] = version
			prefixed[1] = ba.GetScriptAddressPrefix(network)
			copy(prefixed[2:], ba.bytes[2:])
			ba.rawAddress = base58.CheckEncode(prefixed[1:], prefixed[0])

		case WITNESS:
			hrp := NewBitcoinNetwork(network).Bech32HRP()
			program := ba.bytes[2:]
			conv, err := bech32.ConvertBits(program, 8, 5, true)
			if err != nil {
				return nil, err
			}
			encoded, err := bech32.EncodeM(hrp, conv)
			if err != nil {
				return nil, err
			}
			ba.rawAddress = encoded
		}
	} else {
		// Handle non-hex input
		info, err := ba.Decode()
		if err != nil {
			return nil, err
		}
		ba.bytes = ba.WrapAddress(info.Type, info.Bytes, info.Version)
	}

	return ba, nil
}

func BitcoinAddressOnlyFromPublicKey(publicKey []byte) (*BitcoinAddress, error) {
	return BitcoinAddressFromPublicKey(publicKey, BitcoinNetworkSignet)
}

// FromPublicKey creates a BitcoinAddress from a public key
func BitcoinAddressFromPublicKey(publicKey []byte, network BitcoinNetworkType) (*BitcoinAddress, error) {
	// Implement taproot public key to address conversion
	pubKey, err := schnorr.ParsePubKey(publicKey)
	if err != nil {
		return nil, err
	}

	// Convert to taproot output key
	taprootKey := schnorr.SerializePubKey(pubKey)

	// Convert to 5-bit Bech32m words
	program, err := bech32.ConvertBits(taprootKey, 8, 5, true)
	if err != nil {
		return nil, err
	}

	// Add witness version
	version := []byte{0x01}
	program = append(version, program...)

	// Encode with Bech32m
	hrp := NewBitcoinNetwork(network).Bech32HRP()
	address, err := bech32.EncodeM(hrp, program)
	if err != nil {
		return nil, err
	}

	return NewBitcoinAddress(address, network)
}

// ToBytes returns the address bytes
func (ba *BitcoinAddress) ToBytes() []byte {
	return []byte(ba.rawAddress)
}

//type MultiChainAddress struct {
//	MultiChainID int
//	RawAddress   []byte
//}

// GenMultiChainAddress generates a multi-chain address
func (ba *BitcoinAddress) GenMultiChainAddress() []byte {
	// Implement BCS serialization here
	return nil
}

// GenRoochAddress generates a Rooch address
func (ba *BitcoinAddress) GenRoochAddress() (*RoochAddress, error) {
	//hash, err := blake2b.New(RoochAddressLength, nil)
	//if err != nil {
	//	return nil, err
	//}
	//hash.Write(ba.bytes)
	////return hash.Sum(nil), nil
	//return NewRoochAddressFromBytes(hash.Sum(nil))
	if ba.roochAddress == nil {
		hash, err := blake2b.New(RoochAddressLength, nil)
		if err != nil {
			return nil, err
		}
		hash.Write(ba.bytes)
		roochAddress, err := NewRoochAddressFromBytes(hash.Sum(nil))
		if err != nil {
			return nil, err
		}
		ba.roochAddress = roochAddress
	}
	return ba.roochAddress, nil
}

// decode decodes the raw address
func (ba *BitcoinAddress) Decode() (*BitcoinAddressInfo, error) {
	if len(ba.rawAddress) < 14 || len(ba.rawAddress) > 74 {
		return nil, errors.New("invalid address length")
	}

	// Try bech32 first
	hrp, decoded, err := bech32.Decode(ba.rawAddress)
	if err == nil {
		network := FromBech32Prefix(hrp)
		if network != nil {
			version := decoded[0]
			program, err := bech32.ConvertBits(decoded[1:], 5, 8, false)
			if err != nil {
				return nil, err
			}

			if err := validateWitness(version, program); err != nil {
				return nil, err
			}

			return &BitcoinAddressInfo{
				Bytes:   program,
				Type:    WITNESS,
				Version: version,
			}, nil
		}
	}

	// Try base58check
	decoded58 := base58.Decode(ba.rawAddress)
	if len(decoded58) != 21 {
		return nil, errors.New("invalid base58 address")
	}

	switch decoded58[0] {
	case PubkeyAddressPrefixMain:
		return &BitcoinAddressInfo{
			Bytes: decoded58[1:],
			Type:  PKH,
		}, nil
	case ScriptAddressPrefixMain:
		return &BitcoinAddressInfo{
			Bytes: decoded58[1:],
			Type:  SH,
		}, nil
	default:
		return nil, fmt.Errorf("invalid address prefix: %d", decoded58[0])
	}
}

// WrapAddress wraps the address bytes with type and version
func (ba *BitcoinAddress) WrapAddress(addrType BitcoinAddressType, data []byte, version byte) []byte {
	if version != 0 {
		result := make([]byte, len(data)+2)
		result[0] = byte(addrType)
		result[1] = version
		copy(result[2:], data)
		return result
	}

	result := make([]byte, len(data)+1)
	result[0] = byte(addrType)
	copy(result[1:], data)
	return result
}

// Helper functions

// GetPubkeyAddressPrefix returns the prefix for public key addresses
func (ba *BitcoinAddress) GetPubkeyAddressPrefix(network BitcoinNetworkType) byte {
	if network == BitcoinNetworkBitcoin {
		return PubkeyAddressPrefixMain
	}
	return PubkeyAddressPrefixTest
}

// GetScriptAddressPrefix returns the prefix for script addresses
func (ba *BitcoinAddress) GetScriptAddressPrefix(network BitcoinNetworkType) byte {
	if network == BitcoinNetworkBitcoin {
		return ScriptAddressPrefixMain
	}
	return ScriptAddressPrefixTest
}

func validateWitness(version byte, program []byte) error {
	if version == 0 {
		if len(program) != 20 && len(program) != 32 {
			return errors.New("invalid witness program length for version 0")
		}
	} else if version == 1 {
		if len(program) != 32 {
			return errors.New("invalid witness program length for version 1")
		}
	} else {
		return fmt.Errorf("unsupported witness version: %d", version)
	}
	return nil
}

func isHex(s string) bool {
	if len(s) >= 2 && s[0:2] == "0x" {
		s = s[2:]
	}
	_, err := hex.DecodeString(s)
	return err == nil
}

func stripHexPrefix(s string) string {
	if len(s) >= 2 && s[0:2] == "0x" {
		return s[2:]
	}
	return s
}

//func taggedHash(tag string, msg []byte) []byte {
//	tagHash := sha256.Sum256([]byte(tag))
//	h := sha256.New()
//	h.Write(tagHash[:])
//	h.Write(tagHash[:])
//	h.Write(msg)
//	return h.Sum(nil)
//}
