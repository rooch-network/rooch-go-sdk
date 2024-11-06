package crypto

import (
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/bcs"
	"github.com/rooch-network/rooch-go-sdk/transactions"
)

type SignatureScheme = string

// Seeds for deriving addresses from addresses
const (
	Ed25519Scheme   SignatureScheme = "ED25519"   // Ed25519Scheme is the default scheme for deriving the AuthenticationKey
	Secp256k1Scheme SignatureScheme = "Secp256k1" // MultiEd25519Scheme is the scheme for deriving the AuthenticationKey for Multi-ed25519 accounts
)

// FromPublicKey for private / public key pairs, the [AuthenticationKey] is derived from the [PublicKey] directly
func GetSignatureFlag(scheme SignatureScheme) uint8 {
	//flag := 0x00
	switch scheme {
	case Ed25519Scheme:
		return 0x00
	case Secp256k1Scheme:
		return 0x01
	default:
		return 0x00
	}

	//return uint8(flag)
}

// Signer a generic interface for any kind of signing
type Signer interface {
	//// Sign signs a transaction and returns an associated [AccountAuthenticator]
	//Sign(msg []byte) (authenticator *AccountAuthenticator, err error)
	//
	//// SignMessage signs a message and returns the raw [Signature] without a [PublicKey] for verification
	//SignMessage(msg []byte) (signature Signature, err error)
	//
	//// SimulationAuthenticator creates a new [AccountAuthenticator] for simulation purposes
	//SimulationAuthenticator() *AccountAuthenticator
	//
	//// AuthKey gives the [AuthenticationKey] associated with the [Signer]
	//AuthKey() *AuthenticationKey
	//
	//// PubKey Retrieve the [PublicKey] for [Signature] verification
	//PubKey() PublicKey

	// Sign signs a transaction and returns an associated [AccountAuthenticator]
	Sign(msg []byte) (sig Signature, err error)

	// SignMessage signs a message and returns the raw [Signature] without a [PublicKey] for verification
	SignTransaction(tx transactions.Transaction) (transactions.Authenticator, error)

	// SimulationAuthenticator creates a new [AccountAuthenticator] for simulation purposes
	//GetBitcoinAddress() *BitcoinAddress

	// AuthKey gives the [AuthenticationKey] associated with the [Signer]
	GetRoochAddress() address.RoochAddress

	// PubKey Retrieve the [PublicKey] for [Signature] verification
	GetKeyScheme() SignatureScheme

	// PubKey Retrieve the [PublicKey] for [Signature] verification
	GetPublicKey() PublicKey

	//abstract sign(input: Bytes): Promise<Bytes>
	//
	//abstract signTransaction(input: Transaction): Promise<Authenticator>
	//abstract getBitcoinAddress(): BitcoinAddress
	//
	//abstract getRoochAddress(): RoochAddress
	//
	///**
	// * Get the key scheme of the keypair: Secp256k1 or ED25519
	// */
	//abstract getKeyScheme(): SignatureScheme
	//
	///**
	// * The public key for this keypair
	// */
	//abstract getPublicKey(): PublicKey<Address>
}

// MessageSigner a generic interface for a signing private key, a private key isn't always a signer, see SingleSender
//
//// This is not BCS serializable, because this doesn't go on-chain.  An example is [Secp256k1PrivateKey]
//type MessageSigner interface {
//	// SignMessage signs a message and returns the raw [Signature] without a [VerifyingKey]
//	SignMessage(msg []byte) (signature Signature, err error)
//
//	// EmptySignature creates an empty signature for use in simulation
//	EmptySignature() Signature
//
//	// VerifyingKey Retrieve the [VerifyingKey] for signature verification.
//	VerifyingKey() VerifyingKey
//}

// PublicKey is an interface for a public key that can be used to verify transactions in a TransactionAuthenticator
type PublicKey interface {
	VerifyingKey

	// AuthKey gives the [AuthenticationKey] associated with the [PublicKey]
	AuthKey() *AuthenticationKey

	// AuthKey gives the [AuthenticationKey] associated with the [Signer]
	ToRoochAddress() address.RoochAddress

	// Scheme The [DeriveScheme] used for address derivation
	Scheme() DeriveScheme
}

// VerifyingKey a generic interface for a public key associated with the private key, but it cannot necessarily stand on
// its own as a [PublicKey] for authentication on Rooch.  An example is [Secp256k1PublicKey].  All [PublicKey]s are also
// VerifyingKeys.
type VerifyingKey interface {
	bcs.Struct
	CryptoMaterial

	// Verify verifies a message with the public key
	Verify(msg []byte, sig Signature) bool
}

// Signature is an identifier for a serializable [Signature] for on-chain representation
type Signature interface {
	bcs.Struct
	CryptoMaterial
}

// CryptoMaterial is a set of functions for serializing and deserializing a key to and from bytes and hex
// This mirrors the trait in Rust
type CryptoMaterial interface {
	// Bytes outputs the raw byte representation of the [CryptoMaterial]
	Bytes() []byte

	// FromBytes loads the [CryptoMaterial] from the raw bytes
	FromBytes([]byte) error

	// ToHex outputs the hex representation of the [CryptoMaterial] with a leading `0x`
	ToHex() string

	// FromHex parses the hex representation of the [CryptoMaterial] with or without a leading `0x`
	FromHex(string) error
}
