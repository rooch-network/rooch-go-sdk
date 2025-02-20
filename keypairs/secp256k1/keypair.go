package secp256k1

import (
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/crypto"
	"github.com/rooch-network/rooch-go-sdk/utils"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	//ethCrypto "github.com/ethereum/go-ethereum/crypto"
)

const DefaultSecp256k1DerivationPath = "m/54'/784'/0'/0/0"

// Secp256k1KeypairData represents the keypair data structure
type Secp256k1KeypairData struct {
	PublicKey []byte
	SecretKey []byte
}

// Secp256k1Keypair represents a Secp256k1 keypair
type Secp256k1Keypair struct {
	keypair Secp256k1KeypairData
}

// NewSecp256k1Keypair creates a new keypair instance
func NewSecp256k1Keypair(keypair *Secp256k1KeypairData) (*Secp256k1Keypair, error) {
	if keypair != nil {
		return &Secp256k1Keypair{keypair: *keypair}, nil
	}

	// Generate random keypair
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return nil, err
	}

	//privateKey, err := ethCrypto.GenerateKey()
	//if err != nil {
	//	return nil, err
	//}

	//return &Secp256k1PrivateKey{priv}, nil

	return &Secp256k1Keypair{
		keypair: Secp256k1KeypairData{
			PublicKey: privateKey.PubKey().SerializeCompressed(),
			SecretKey: privateKey.Serialize(),
		},
	}, nil
}

// Generate generates a new random keypair
func GenerateSecp256k1Keypair() (*Secp256k1Keypair, error) {
	return NewSecp256k1Keypair(nil)
}

// FromSecretKey creates a keypair from a raw secret key byte array
// func FromSecp256k1SecretKey(secretKey []byte, skipValidation bool) (*Secp256k1Keypair, error) {
func FromSecp256k1SecretKey(secretKey interface{}, skipValidation bool) (*Secp256k1Keypair, error) {
	//const decodeSecretKey =
	//	typeof secretKey === 'string'
	//? (() => {
	//	const decoded = decodeRoochSercetKey(secretKey)
	//	if (decoded.schema !== 'Secp256k1') {
	//		throw new Error('provided secretKey is invalid')
	//	}
	//	return decoded.secretKey
	//})()
	//: secretKey
	//DecodeRoochSecretKey(value string) (*ParsedKeypair, error)
	decodedSecretKey := []byte{}
	switch v := secretKey.(type) {
	case string:
		parsedKeypair, err := crypto.DecodeRoochSecretKey(v)
		if err != nil {
			return nil, err
		}
		if parsedKeypair.Schema != crypto.Secp256k1Scheme {
			return nil, errors.New("provided secretKey is invalid")
		}
		decodedSecretKey = parsedKeypair.SecretKey
	case []byte:
		decodedSecretKey = v
	default:
		return nil, errors.New("invalid secretKey")
	}

	//privateKey, publicKey := btcec.PrivKeyFromBytes(secretKey)
	//publicKey: Uint8Array = secp256k1.getPublicKey(decodeSecretKey, true)
	privKey, _ := btcec.PrivKeyFromBytes(decodedSecretKey)
	pubKey := privKey.PubKey().SerializeCompressed()
	if !skipValidation {
		// Perform validation similar to TypeScript implementation
		msg := []byte("rooch validation")
		//hash := sha256.Sum256(msg)
		msgHash := utils.Blake2b256(msg)

		signature, err := secp256k1.Sign(msgHash[:], decodedSecretKey)
		if err != nil {
			return nil, err
		}

		if !secp256k1.VerifySignature(pubKey, msgHash, signature) {
			return nil, errors.New("provided secretKey is invalid")
		}
	}

	return &Secp256k1Keypair{
		keypair: Secp256k1KeypairData{
			PublicKey: pubKey,
			SecretKey: privKey.Serialize(),
		},
	}, nil
}

// FromSeed generates a keypair from a 32 byte seed
func FromSecp256k1Seed(seed []byte) (*Secp256k1Keypair, error) {
	privateKey, _ := btcec.PrivKeyFromBytes(seed)

	return &Secp256k1Keypair{
		keypair: Secp256k1KeypairData{
			PublicKey: privateKey.PubKey().SerializeCompressed(),
			SecretKey: privateKey.Serialize(),
		},
	}, nil
}

// GetPublicKey returns the public key
//
//	func (kp *Secp256k1Keypair) GetPublicKey() []byte {
//		return kp.keypair.PublicKey
//	}
func (kp *Secp256k1Keypair) GetPublicKey() crypto.PublicKey[address.AddressView] {
	return &Secp256k1PublicKey{kp.keypair.PublicKey}
}

// GetSchnorrPublicKey returns the Schnorr public key
func (kp *Secp256k1Keypair) GetSchnorrPublicKey() crypto.PublicKey[address.AddressView] {
	privateKey, _ := btcec.PrivKeyFromBytes(kp.keypair.SecretKey)
	schnorrPubKey := schnorr.SerializePubKey(privateKey.PubKey())
	return &Secp256k1PublicKey{schnorrPubKey}
}

// GetSecretKey returns the secret key
func (kp *Secp256k1Keypair) GetSecretKey() []byte {
	return kp.keypair.SecretKey
}

// Sign signs the provided data
func (kp *Secp256k1Keypair) Sign(input []byte) ([]byte, error) {
	hash := utils.Sha256(input)
	//privateKey, _ := btcec.PrivKeyFromBytes(kp.keypair.SecretKey)

	signature, err := secp256k1.Sign(hash[:], kp.GetSecretKey())
	if err != nil {
		return nil, err
	}

	return signature, nil
}

// DeriveKeypair derives a keypair from mnemonics and path
func DeriveSecp256k1Keypair(mnemonics string, path string) (*Secp256k1Keypair, error) {
	if path == "" {
		path = DefaultSecp256k1DerivationPath
	}

	seed := bip39.NewSeed(mnemonics, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	// Derive the key using the path
	derivedKey := masterKey
	// Implementation of path derivation would go here
	// You'll need to parse the path and derive each level

	return &Secp256k1Keypair{
		keypair: Secp256k1KeypairData{
			PublicKey: derivedKey.PublicKey().Key,
			SecretKey: derivedKey.Key,
		},
	}, nil
}
