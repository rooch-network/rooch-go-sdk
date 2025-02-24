package ed25519

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/types"
	"strings"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"

	"github.com/rooch-network/rooch-go-sdk/crypto"
)

const (
	DefaultEd25519DerivationPath = "m/44'/784'/0'/0'/0'"
	PrivateKeySize               = 32
)

// Ed25519KeypairData represents the keypair data structure
type Ed25519KeypairData struct {
	PublicKey []byte
	SecretKey []byte
}

// Ed25519Keypair represents an Ed25519 keypair used for signing transactions
type Ed25519Keypair struct {
	keypair Ed25519KeypairData
}

// NewEd25519Keypair creates a new Ed25519 keypair instance
func NewEd25519Keypair(keypair *Ed25519KeypairData) (*Ed25519Keypair, error) {
	if keypair == nil {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
		return &Ed25519Keypair{
			keypair: Ed25519KeypairData{
				PublicKey: pub,
				SecretKey: priv,
			},
		}, nil
	}
	return &Ed25519Keypair{keypair: *keypair}, nil
}

// GetKeyScheme returns the key scheme of the keypair
func (k *Ed25519Keypair) GetKeyScheme() crypto.SignatureScheme {
	return crypto.Ed25519Scheme
}

// Generate generates a new random Ed25519 keypair
func GenerateEd25519Keypair() (*Ed25519Keypair, error) {
	return NewEd25519Keypair(nil)
}

// FromSecretKey creates an Ed25519 keypair from a raw secret key byte array
func FromEd25519SecretKey(secretKey []byte, skipValidation bool) (*Ed25519Keypair, error) {
	if len(secretKey) != PrivateKeySize {
		return nil, fmt.Errorf("wrong secretKey size. Expected %d bytes, got %d", PrivateKeySize, len(secretKey))
	}

	privKey := ed25519.NewKeyFromSeed(secretKey)
	pubKey := privKey.Public().(ed25519.PublicKey)

	keypair := &Ed25519KeypairData{
		PublicKey: pubKey,
		SecretKey: privKey,
	}

	if !skipValidation {
		message := []byte("rooch validation")
		signature := ed25519.Sign(privKey, message)
		if !ed25519.Verify(pubKey, message, signature) {
			return nil, errors.New("provided secretKey is invalid")
		}
	}

	return NewEd25519Keypair(keypair)
}

// GetBitcoinAddress returns the Bitcoin address (not implemented for Ed25519)
func (k *Ed25519Keypair) GetBitcoinAddress() (*address.BitcoinAddress, error) {
	return nil, errors.New("method not implemented in Ed25519")
}

// GetRoochAddress returns the Rooch address
func (k *Ed25519Keypair) GetRoochAddress() (*address.RoochAddress, error) {
	return k.GetPublicKey().ToAddress()
}

// GetPublicKey returns the public key for this Ed25519 keypair
func (k *Ed25519Keypair) GetPublicKey() crypto.PublicKey[address.RoochAddress] {
	//return k.keypair.PublicKey
	return &Ed25519PublicKey{k.keypair.PublicKey}
}

// The Bech32 secret key string for this Ed25519 keypair
func (k *Ed25519Keypair) GetSecretKey() (string, error) {
	//return k.keypair.SecretKey[:PrivateKeySize]

	return crypto.EncodeRoochSecretKey(
		k.keypair.SecretKey[:PrivateKeySize],
		k.GetKeyScheme(),
	)
}

// Sign returns the signature for the provided data using Ed25519
func (k *Ed25519Keypair) Sign(input []byte) ([]byte, error) {
	return ed25519.Sign(k.keypair.SecretKey, input), nil
}

// SignTransaction signs a transaction and returns an authenticator
func (k *Ed25519Keypair) SignTransaction(tx crypto.Transaction) (*crypto.Authenticator, error) {
	// If you need to use types.Transaction internally:
	if typesTx, ok := tx.(*types.Transaction); ok {
		// Use typedTx here
		hash, err := typesTx.HashData()
		if err != nil {
			return nil, err
		}
		return crypto.RoochAuthValidator(hash, k)
	} else {
		return nil, errors.New("invalid transaction")
	}
}

// DeriveKeypair derives Ed25519 keypair from mnemonics and path
func DeriveEd25519Keypair(mnemonics string, path string) (*Ed25519Keypair, error) {
	if path == "" {
		path = DefaultEd25519DerivationPath
	}

	if !isValidHardenedPath(path) {
		return nil, errors.New("invalid derivation path")
	}

	seed := bip39.NewSeed(mnemonics, "")
	return deriveEd25519KeypairFromSeed(hex.EncodeToString(seed), path)
}

// DeriveKeypairFromSeed derives Ed25519 keypair from seed and path
func DeriveEd25519KeypairFromSeed(seedHex string, path string) (*Ed25519Keypair, error) {
	if path == "" {
		path = DefaultEd25519DerivationPath
	}

	if !isValidHardenedPath(path) {
		return nil, errors.New("invalid derivation path")
	}

	return deriveEd25519KeypairFromSeed(seedHex, path)
}

// Helper functions

func deriveEd25519KeypairFromSeed(seedHex string, path string) (*Ed25519Keypair, error) {
	seed, err := hex.DecodeString(seedHex)
	if err != nil {
		return nil, err
	}

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}

	// Parse and derive the path
	segments := strings.Split(path[2:], "/")
	currentKey := masterKey

	for _, segment := range segments {
		if segment == "" {
			continue
		}

		hardened := strings.HasSuffix(segment, "'")
		index := uint32(0)
		if hardened {
			segment = segment[:len(segment)-1]
			// Add hardened offset
			index = 0x80000000
		}

		// Parse index
		var segmentIndex uint32
		_, err := fmt.Sscanf(segment, "%d", &segmentIndex)
		if err != nil {
			return nil, err
		}
		index += segmentIndex

		currentKey, err = currentKey.NewChildKey(index)
		if err != nil {
			return nil, err
		}
	}

	return FromEd25519SecretKey(currentKey.Key, false)
}

func isValidHardenedPath(path string) bool {
	if !strings.HasPrefix(path, "m/") {
		return false
	}

	segments := strings.Split(path[2:], "/")
	for _, segment := range segments {
		if segment == "" {
			continue
		}
		if !strings.HasSuffix(segment, "'") {
			return false
		}
	}

	return true
}
