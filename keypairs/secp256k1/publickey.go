package secp256k1

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/crypto"
	"github.com/rooch-network/rooch-go-sdk/utils"
)

const (
	SchnorrPublicKeySize = 32
)

// Secp256k1PublicKey represents a Secp256k1 public key
type Secp256k1PublicKey struct {
	data []byte
}

// NewSecp256k1PublicKey creates a new Secp256k1PublicKey object
func NewSecp256k1PublicKey(value interface{}) (*Secp256k1PublicKey, error) {
	var data []byte

	switch v := value.(type) {
	case string:
		// Assume base64 encoded string
		decoded, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, err
		}
		data = decoded
	case []byte:
		data = v
	default:
		return nil, errors.New("unsupported public key input type")
	}

	if len(data) != SchnorrPublicKeySize && len(data) != 33 {
		return nil, errors.New("invalid public key input size")
	}

	return &Secp256k1PublicKey{
		data: data,
	}, nil
}

// Equals checks if two Secp256k1 public keys are equal
func (pk *Secp256k1PublicKey) Equals(other crypto.PublicKey[address.AddressView]) bool {
	//return crypto.BytesEqual(pk.data, other.data)
	if pk == nil || other == nil {
		return false
	}
	return utils.BytesEqual(pk.data, other.ToBytes())
}

// ToBytes returns the byte array representation of the Secp256k1 public key
func (pk *Secp256k1PublicKey) ToBytes() []byte {
	return pk.data
}

// String returns the hex string representation of the public key
func (pk *Secp256k1PublicKey) String() string {
	return hex.EncodeToString(pk.data)
}

func (pk *Secp256k1PublicKey) ToBase64() string {
	return utils.ToB64(pk.ToBytes())
}

// ToAddress returns the Bitcoin address associated with this Secp256k1 public key
// func (pk *Secp256k1PublicKey) ToAddress() (*address.AddressView, error) {
func (pk *Secp256k1PublicKey) ToAddress() (*address.AddressView, error) {
	return address.NewAddressView(pk.data)
}

// ToAddressWith returns the Bitcoin address with specified network type
func (pk *Secp256k1PublicKey) ToAddressWith(network address.BitcoinNetworkType) (*address.AddressView, error) {
	return address.NewAddressViewWithNetwork(pk.data, network)
}

// Flag returns the signature scheme flag for Secp256k1
func (pk *Secp256k1PublicKey) Flag() uint8 {
	return uint8(crypto.Secp256k1Flag)
}

// Verify verifies that the signature is valid for the provided message
func (pk *Secp256k1PublicKey) Verify(message []byte, signature []byte) (bool, error) {
	// Create hash of the message
	messageHash := utils.Sha256(message)
	verifyResult := secp256k1.VerifySignature(pk.ToBytes(), messageHash, signature)
	return verifyResult, nil
}
