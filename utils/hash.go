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

//// Sha3_256 computes SHA3-256 hash of the input bytes
//func Sha3_256(msg []byte) []byte {
//	hash := sha3.New256()
//	hash.Write(msg)
//	return hash.Sum(nil)
//}

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

// Sha3256 hashes the input bytes using SHA3-256
func Sha3256(data []byte) (output []byte) {
	hasher := sha3.New256()
	hasher.Write(data)
	//return hasher.Sum(nil)
	return hasher.Sum([]byte{})
}

// Blake2b256 hashes the input bytes using Blake2b256-256
func Blake2b256(data []byte) (output []byte) {
	hasher, _ := blake2b.New256(nil)
	hasher.Write(data)
	//hasher := w.Sum(nil)
	return hasher.Sum(nil)
}

// Blake2b computes BLAKE2b hash of the input bytes
func Blake2b512(msg []byte) ([]byte, error) {
	hash, err := blake2b.New512(nil)
	if err != nil {
		return nil, err
	}
	hash.Write(msg)
	return hash.Sum(nil), nil
}

//// Hash160 performs RIPEMD160(SHA256(data))
//func Hash160(data []byte) []byte {
//	sha := sha256.Sum256(data)
//	ripemd := ripemd160.New()
//	ripemd.Write(sha[:])
//	return ripemd.Sum(nil)
//}
