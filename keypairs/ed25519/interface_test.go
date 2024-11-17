package ed25519

import (
	"errors"
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
}

// Signer a generic interface for any kind of signing
type SignerTest interface {
	GetPublicKey() PublicKeyTest[address.RoochAddress]

	SignTransaction(tx Transaction) (*AuthenticatorTest, error)
}

type AuthenticatorTest struct {
	AuthValidatorId uint64 `json:"auth_validator_id"`
	Payload         []byte `json:"payload"`
}

type Ed25519KeypairTest struct {
}

func (k *Ed25519KeypairTest) GetPublicKey() PublicKeyTest[address.RoochAddress] {
	return &Ed25519PublicKeyTest{}
}

func (k *Ed25519KeypairTest) SignTransaction(tx Transaction) (*AuthenticatorTest, error) {
	// If you need to use types.Transaction internally:
	if typedTx, ok := tx.(*types.Transaction); ok {
		// Use typedTx here
		hash, err := typedTx.HashData()
		if err != nil {
			return nil, err
		}
		return RoochAuthValidatorTest(hash, k)
	} else {
		return nil, errors.New("invalid transaction")
	}
}

func RoochAuthValidatorTest(input []byte, signer SignerTest) (*AuthenticatorTest, error) {
	return &AuthenticatorTest{}, nil
}

//func (k *Ed25519KeypairTest) SignTransactionTest() (*AuthenticatorTest, error) {
//	return RoochAuthValidatorTest([]byte{}, k)
//}
