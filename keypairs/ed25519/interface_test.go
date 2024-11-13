package ed25519

import (
	"github.com/rooch-network/rooch-go-sdk/address"
)

type PublicKeyTest[T any] interface {
	// ToAddress converts the public key to its corresponding address
	ToAddress() (*T, error)
}

// Ed25519PublicKey represents an Ed25519 public key
type Ed25519PublicKeyTest struct {
	//data []byte
}

// ToAddress returns the Rooch address associated with this Ed25519 public key
func (pk *Ed25519PublicKeyTest) ToAddress() (*address.RoochAddress, error) {
	return &address.RoochAddress{}, nil
}

// Signer a generic interface for any kind of signing
type SignerTest interface {
	// PubKey Retrieve the [PublicKey] for [Signature] verification
	GetPublicKey() PublicKeyTest[address.RoochAddress]
}

type Ed25519KeypairTest struct {
}

//	func (k *Ed25519KeypairTest) GetPublicKey() PublicKeyTest[address.RoochAddress] {
//		//return k.keypair.PublicKey
//		return &Ed25519PublicKeyTest{}
//	}
func (k *Ed25519KeypairTest) GetPublicKey() PublicKeyTest[address.RoochAddress] {
	return &Ed25519PublicKeyTest{}
}

func RoochAuthValidatorTest(input []byte, signer SignerTest) error {
	return nil
}

func (k *Ed25519KeypairTest) SignTransactionTest() error {
	return RoochAuthValidatorTest([]byte{}, k)
}
