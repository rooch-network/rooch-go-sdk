package types

import (
	"crypto/sha256"
	//"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/bcs"
	//"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
	//"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/bcs"
	//"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"
	//"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/bcs"
	//"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"
	//"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/bcs"
	//"github.com/novifinancial/serde-reflection/serde-generate/runtime/golang/serde"
)

type H256 [32]byte

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

// Hash160 performs RIPEMD160(SHA256(data))
func Hash160(data []byte) []byte {
	sha := sha256.Sum256(data)
	ripemd := ripemd160.New()
	ripemd.Write(sha[:])
	return ripemd.Sum(nil)
}

//
//func (obj *HashValue) Serialize(serializer serde.Serializer) error {
//	if err := serializer.IncreaseContainerDepth(); err != nil {
//		return err
//	}
//	if err := serializer.SerializeBytes(([]byte)(*obj)); err != nil {
//		return err
//	}
//	serializer.DecreaseContainerDepth()
//	return nil
//}
//
//func (obj *HashValue) BcsSerialize() ([]byte, error) {
//	if obj == nil {
//		return nil, fmt.Errorf("Cannot serialize null object")
//	}
//	serializer := bcs.NewSerializer()
//	if err := obj.Serialize(serializer); err != nil {
//		return nil, err
//	}
//	return serializer.GetBytes(), nil
//}
//
//func DeserializeHashValue(deserializer serde.Deserializer) (HashValue, error) {
//	var obj []byte
//	if err := deserializer.IncreaseContainerDepth(); err != nil {
//		return (HashValue)(obj), err
//	}
//	if val, err := deserializer.DeserializeBytes(); err == nil {
//		obj = val
//	} else {
//		return (HashValue)(obj), err
//	}
//	deserializer.DecreaseContainerDepth()
//	return (HashValue)(obj), nil
//}
//
//func BcsDeserializeHashValue(input []byte) (HashValue, error) {
//	if input == nil {
//		var obj HashValue
//		return obj, fmt.Errorf("Cannot deserialize null array")
//	}
//	deserializer := bcs.NewDeserializer(input)
//	obj, err := DeserializeHashValue(deserializer)
//	if err == nil && deserializer.GetBufferOffset() < uint64(len(input)) {
//		return obj, fmt.Errorf("Some input bytes were not read")
//	}
//	return obj, err
//}
