package address

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"

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

// BitcoinAddress represents a Bitcoin address
type BitcoinAddress struct {
	rawAddress string
	bytes      []byte
	network    BitcoinNetworkType
}

// NewBitcoinAddress creates a new BitcoinAddress instance
func NewBitcoinAddress(input string, network BitcoinNetworkType) (*BitcoinAddress, error) {
	addr := &BitcoinAddress{
		rawAddress: input,
		network:    network,
	}

	if isHex(input) {
		// Handle hex input
		hexBytes, err := hex.DecodeString(stripHexPrefix(input))
		if err != nil {
			return nil, err
		}
		addr.bytes = hexBytes

		var prefixed []byte
		version := hexBytes[1]

		switch hexBytes[0] {
		case byte(PKH):
			prefixed = make([]byte, 22)
			prefixed[0] = version
			prefixed[1] = addr.getPubkeyAddressPrefix()
			copy(prefixed[2:], hexBytes[2:])
			addr.rawAddress = base58.CheckEncode(prefixed[1:], prefixed[0])

		case byte(SH):
			prefixed = make([]byte, 22)
			prefixed[0] = version
			prefixed[1] = addr.getScriptAddressPrefix()
			copy(prefixed[2:], hexBytes[2:])
			addr.rawAddress = base58.CheckEncode(prefixed[1:], prefixed[0])

		case byte(WITNESS):
			hrp := NewBitcoinNetwork(network).Bech32HRP()
			conv, err := bech32.ConvertBits(hexBytes[2:], 8, 5, true)
			if err != nil {
				return nil, err
			}
			finalConv := append([]byte{version}, conv...)

			var encodedAddr string
			if version == 0 {
				encodedAddr, err = bech32.EncodeM(hrp, finalConv)
			} else {
				encodedAddr, err = bech32.Encode(hrp, finalConv)
			}
			if err != nil {
				return nil, err
			}
			addr.rawAddress = encodedAddr
		}
	} else {
		// Handle address string input
		info, err := addr.decode()
		if err != nil {
			return nil, err
		}
		addr.bytes = addr.wrapAddress(info)
	}

	return addr, nil
}

// FromPublicKey creates a BitcoinAddress from a public key
func FromPublicKey(publicKey []byte, network BitcoinNetworkType) (*BitcoinAddress, error) {
	// Note: This is a simplified implementation
	// You'll need to implement the actual taproot logic here
	hrp := NewBitcoinNetwork(network).Bech32HRP()

	// Convert public key to program
	program := sha256.Sum256(publicKey)

	// Convert to 5-bit words
	conv, err := bech32.ConvertBits(program[:], 8, 5, true)
	if err != nil {
		return nil, err
	}

	// Add version 1 for taproot
	words := append([]byte{1}, conv...)

	// Encode with bech32m
	addr, err := bech32.EncodeM(hrp, words)
	if err != nil {
		return nil, err
	}

	return NewBitcoinAddress(addr, network)
}

// ToBytes returns the address bytes
func (ba *BitcoinAddress) ToBytes() []byte {
	return []byte(ba.rawAddress)
}

// GenMultiChainAddress generates a multi-chain address
func (ba *BitcoinAddress) GenMultiChainAddress() []byte {
	// Implementation depends on your MultiChainAddress structure
	// This is a placeholder implementation
	return append([]byte{0x01}, ba.bytes...)
}

// GenRoochAddress generates a Rooch address
func (ba *BitcoinAddress) GenRoochAddress() ([]byte, error) {
	hash, err := blake2b.New(RoochAddressLength, nil)
	if err != nil {
		return nil, err
	}
	hash.Write(ba.bytes)
	return hash.Sum(nil), nil
}

// decode decodes the raw address
func (ba *BitcoinAddress) decode() (*BitcoinAddressInfo, error) {
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

// Helper functions

func (ba *BitcoinAddress) getPubkeyAddressPrefix() byte {
	if ba.network == BitcoinNetworkBitcoin {
		return PubkeyAddressPrefixMain
	}
	return PubkeyAddressPrefixTest
}

func (ba *BitcoinAddress) getScriptAddressPrefix() byte {
	if ba.network == BitcoinNetworkBitcoin {
		return ScriptAddressPrefixMain
	}
	return ScriptAddressPrefixTest
}

func (ba *BitcoinAddress) wrapAddress(info *BitcoinAddressInfo) []byte {
	result := make([]byte, len(info.Bytes)+2)
	result[0] = byte(info.Type)
	result[1] = info.Version
	copy(result[2:], info.Bytes)
	return result
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

//// Hash160 performs RIPEMD160(SHA256(data))
//func Hash160(data []byte) []byte {
//	sha := sha256.Sum256(data)
//	ripemd := ripemd160.New()
//	ripemd.Write(sha[:])
//	return ripemd.Sum(nil)
//}
