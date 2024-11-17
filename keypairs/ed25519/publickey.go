package ed25519

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/bcs"
	"github.com/rooch-network/rooch-go-sdk/crypto"
	"github.com/rooch-network/rooch-go-sdk/utils"
)

const PublicKeySize = 32

// Ed25519PublicKey represents an Ed25519 public key
type Ed25519PublicKey struct {
	data []byte
}

// NewEd25519PublicKey creates a new Ed25519PublicKey object
// value can be either a base64 encoded string or a byte slice
func NewEd25519PublicKey(value interface{}) (*Ed25519PublicKey, error) {
	var data []byte

	switch v := value.(type) {
	case string:
		var err error
		data, err = base64.StdEncoding.DecodeString(v)
		if err != nil {
			return nil, err
		}
	case []byte:
		data = v
	default:
		return nil, errors.New("invalid public key input type")
	}

	if len(data) != PublicKeySize {
		return nil, errors.New("invalid public key size")
	}

	return &Ed25519PublicKey{
		data: data,
	}, nil
}

// Equals checks if two Ed25519 public keys are equal
func (pk *Ed25519PublicKey) Equals(other crypto.PublicKey[address.RoochAddress]) bool {
	if pk == nil || other == nil {
		return false
	}
	return utils.BytesEqual(pk.data, other.ToBytes())
}

// ToBytes returns the byte array representation of the Ed25519 public key
func (pk *Ed25519PublicKey) ToBytes() []byte {
	return pk.data
}

// String returns the hex string representation of the public key
func (pk *Ed25519PublicKey) String() string {
	return hex.EncodeToString(pk.data)
}

// Flag returns the signature scheme flag for Ed25519
func (pk *Ed25519PublicKey) Flag() uint8 {
	return uint8(crypto.Ed25519Flag)
}

// Verify verifies that the signature is valid for the provided message
func (pk *Ed25519PublicKey) Verify(message, signature []byte) (bool, error) {
	return ed25519.Verify(pk.data, message, signature), nil
}

// ToAddress returns the Rooch address associated with this Ed25519 public key
func (pk *Ed25519PublicKey) ToAddress() (*address.RoochAddress, error) {
	tmp := make([]byte, PublicKeySize+1)
	tmp[0] = byte(crypto.SignatureSchemeSize[crypto.Ed25519Scheme])
	//tmp.set([SIGNATURE_SCHEME_TO_FLAG.ED25519])
	//copy(tmp[1:], pk.data)
	copy(tmp[1:], pk.ToBytes())

	//hash := utils.Blake2b(tmp, 32)
	//addressBytes := hash[:address.ROOCH_ADDRESS_LENGTH*2]
	//return address.NewRoochAddress(addressBytes)

	addressBytes := utils.Blake2b256(tmp)[:address.RoochAddressLength*2]
	return address.NewRoochAddressFromBytes(addressBytes)
}

func (pk *Ed25519PublicKey) ToBase64() string {
	return utils.ToB64(pk.ToBytes())
}

func (pk *Ed25519PublicKey) ToHex() string {
	return utils.BytesToHex(pk.ToBytes())
}

// // FromHex sets the [Ed25519PublicKey] to the bytes represented by the hex string, with or without a leading 0x
// //
// // Errors if the hex string is not valid, or if the bytes length is not [ed25519.PublicKeySize].
// //
// // Implements:
// //   - [CryptoMaterial]
func (pk *Ed25519PublicKey) FromHex(hexStr string) (err error) {
	bytes, err := utils.ParseHex(hexStr)
	if err != nil {
		return err
	}
	return pk.FromBytes(bytes)
}

func (pk *Ed25519PublicKey) FromBytes(bytes []byte) (err error) {
	if len(bytes) != ed25519.PublicKeySize {
		return errors.New("invalid ed25519 public key size")
	}
	pk.data = bytes
	return nil
}

//endregion

//region Ed25519PublicKey bcs.Struct implementation

// MarshalBCS serializes the [Ed25519PublicKey] to BCS bytes
//
// Implements:
//   - [bcs.Marshaler]
func (pk *Ed25519PublicKey) MarshalBCS(ser *bcs.Serializer) {
	ser.WriteBytes(pk.data)
}

// UnmarshalBCS deserializes the [Ed25519PublicKey] from BCS bytes
//
// Sets [bcs.Deserializer.Error] if the bytes length is not [ed25519.PublicKeySize], or if it fails to read the required bytes.
//
// Implements:
//   - [bcs.Unmarshaler]
func (pk *Ed25519PublicKey) UnmarshalBCS(des *bcs.Deserializer) {
	kb := des.ReadBytes()
	if des.Error() != nil {
		return
	}
	err := pk.FromBytes(kb)
	if err != nil {
		des.SetError(err)
		return
	}
}
