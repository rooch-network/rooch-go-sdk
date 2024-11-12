// Copyright (c) RoochNetwork
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"crypto/sha256"
	"crypto/sha512"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

// Sha256 computes SHA256 hash of the input bytes
func Sha256(msg []byte) []byte {
	hash := sha256.New()
	hash.Write(msg)
	return hash.Sum(nil)
}

// Sha3_256 computes SHA3-256 hash of the input bytes
func Sha3_256(msg []byte) []byte {
	hash := sha3.New256()
	hash.Write(msg)
	return hash.Sum(nil)
}

// Sha512 computes SHA512 hash of the input bytes
func Sha512(msg []byte) []byte {
	hash := sha512.New()
	hash.Write(msg)
	return hash.Sum(nil)
}

// Blake2b computes BLAKE2b hash of the input bytes
func Blake2b(msg []byte) ([]byte, error) {
	hash, err := blake2b.New512(nil)
	if err != nil {
		return nil, err
	}
	hash.Write(msg)
	return hash.Sum(nil), nil
}

// Hash160 computes RIPEMD160(SHA256(msg))
func Hash160(msg []byte) []byte {
	sha256Hash := Sha256(msg)
	ripemd := ripemd160.New()
	ripemd.Write(sha256Hash)
	return ripemd.Sum(nil)
}

// Sha256Double computes SHA256(SHA256(concatenated msgs))
func Sha256Double(msgs ...[]byte) ([]byte, error) {
	concatenated, err := ConcatBytes(msgs...)
	if err != nil {
		return nil, err
	}
	firstHash := Sha256(concatenated)
	return Sha256(firstHash), nil
}
