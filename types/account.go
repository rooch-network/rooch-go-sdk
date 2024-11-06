package types

import (
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/crypto"
)

// Re-export types so that way the user experience doesn't change

type RoochAddress = address.RoochAddress

// AccountAddress is a 32 byte address on the Rooch blockchain
// It can represent an Object, an Account, and much more.
type AccountAddress = address.AccountAddress

// Account is a wrapper for a signer, handling the AccountAddress and signing
type Account = address.Account

// AccountZero represents the 0x0 address
var AccountZero = address.AccountZero

// AccountOne represents the 0x1 address
var AccountOne = address.AccountOne

// AccountTwo represents the 0x2 address
var AccountTwo = address.AccountTwo

// AccountThree represents the 0x3 address
var AccountThree = address.AccountThree

// AccountFour represents the 0x4 address
var AccountFour = address.AccountFour

// NewAccountFromSigner creates an account from a Signer, which is most commonly a private key
func NewAccountFromSigner(signer crypto.Signer, authKey ...crypto.AuthenticationKey) (*Account, error) {
	return address.NewAccountFromSigner(signer, authKey...)
}

// NewEd25519Account creates a legacy Ed25519 account, this is most commonly used in wallets
func NewEd25519Account() (*Account, error) {
	return address.NewEd25519Account()
}

// NewEd25519SingleSenderAccount creates a single signer Ed25519 account
func NewEd25519SingleSenderAccount() (*Account, error) {
	return address.NewEd25519SingleSignerAccount()
}

// NewSecp256k1Account creates a Secp256k1 account
func NewSecp256k1Account() (*Account, error) {
	return address.NewSecp256k1Account()
}
