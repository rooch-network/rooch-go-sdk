package crypto

import (
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/bcs"
)

type TransactionData interface {
	Hash() ([]byte, error)
	bcs.Struct
}

type Transaction interface {
	GetData() TransactionData
	GetInfo() *string
	HashData() ([]byte, error)
	bcs.Struct
}

// Signer a generic interface for any kind of signing
type Signer[T any] interface {
	Sign(msg []byte) ([]byte, error)

	SignTransaction(tx Transaction) (*Authenticator, error)

	GetBitcoinAddress() (*address.BitcoinAddress, error)

	GetRoochAddress() (*address.RoochAddress, error)

	GetKeyScheme() SignatureScheme

	//GetPublicKey() PublicKey[address.RoochAddress]
	GetPublicKey() PublicKey[T]
}
