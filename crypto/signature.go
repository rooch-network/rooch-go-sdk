// Copyright (c) RoochNetwork
// SPDX-License-Identifier: Apache-2.0

package crypto

// SignatureScheme represents the type of signature scheme
type SignatureScheme string

// SignatureFlag represents the byte flag for signature schemes
type SignatureFlag uint8

const (
	// Signature schemes
	Ed25519Scheme   SignatureScheme = "ED25519"
	Secp256k1Scheme SignatureScheme = "Secp256k1"
)

// Signature scheme flags
const (
	Ed25519Flag   SignatureFlag = 0x00
	Secp256k1Flag SignatureFlag = 0x01
)

// SignatureSchemeSizes maps signature schemes to their respective sizes
var SignatureSchemeSize = map[SignatureScheme]int{
	Ed25519Scheme:   32,
	Secp256k1Scheme: 33,
}

// SchemeToFlag maps signature schemes to their flags
var SignatureSchemeToFlag = map[SignatureScheme]SignatureFlag{
	Ed25519Scheme:   Ed25519Flag,
	Secp256k1Scheme: Secp256k1Flag,
}

// FlagToScheme maps signature flags to their schemes
var SignatureFlagToScheme = map[SignatureFlag]SignatureScheme{
	Ed25519Flag:   Ed25519Scheme,
	Secp256k1Flag: Secp256k1Scheme,
}
