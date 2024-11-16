package ed25519

import (
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/types"
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

type Transaction interface {
	HashData() ([]byte, error)
	//SetAuthenticator(*Authenticator)
}

// Signer a generic interface for any kind of signing
type SignerTest interface {
	// PubKey Retrieve the [PublicKey] for [Signature] verification
	GetPublicKey() PublicKeyTest[address.RoochAddress]

	SignTransaction(tx *Transaction) error
}

type Ed25519KeypairTest struct {
}

//func (k *Ed25519KeypairTest) SignTransaction(tx *Transaction) error {
//	//TODO implement me
//	panic("implement me")
//}

//	func (k *Ed25519KeypairTest) GetPublicKey() PublicKeyTest[address.RoochAddress] {
//		//return k.keypair.PublicKey
//		return &Ed25519PublicKeyTest{}
//	}
func (k *Ed25519KeypairTest) GetPublicKey() PublicKeyTest[address.RoochAddress] {
	return &Ed25519PublicKeyTest{}
}

func (k *Ed25519KeypairTest) SignTransaction(tx *Transaction) error {
	// If you need to use types.Transaction internally:
	if _typedTx, ok := tx.(*types.Transaction); ok {
		// Use typedTx here
		return nil
	}
	return nil
}

func RoochAuthValidatorTest(input []byte, signer SignerTest) error {
	return nil
}

func (k *Ed25519KeypairTest) SignTransactionTest() error {
	return RoochAuthValidatorTest([]byte{}, k)
}
