package crypto

import (
	"bytes"
	"encoding/hex"
	"errors"
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/bcs"
	"github.com/rooch-network/rooch-go-sdk/utils"
	"strings"
)

const (
	BitcoinMessagePrefix = "\u0018Bitcoin Signed Message:\n"
	MessageInfoPrefix    = "Rooch Transaction:\n"
)

type BitcoinSignMessage struct {
	messagePrefix string
	messageInfo   string
	txHash        []byte
}

func NewBitcoinSignMessage(txData []byte, messageInfo string) *BitcoinSignMessage {
	msg := messageInfo
	if !strings.HasPrefix(msg, MessageInfoPrefix) {
		msg = MessageInfoPrefix + messageInfo
	}

	if !strings.HasSuffix(msg, "\n") {
		msg = msg + "\n"
	}

	return &BitcoinSignMessage{
		messagePrefix: BitcoinMessagePrefix,
		messageInfo:   msg,
		txHash:        txData,
	}
}

func (b *BitcoinSignMessage) Raw() string {
	return b.messageInfo + hex.EncodeToString(b.txHash)
}

func (b *BitcoinSignMessage) Encode() []byte {
	msgHex := []byte(hex.EncodeToString(b.txHash))
	infoBytes := []byte(b.messageInfo)
	prefixBytes := bytes.Join([][]byte{
		[]byte(b.messagePrefix),
		utils.VarintByteNum(uint64(len(infoBytes) + len(msgHex))),
	}, []byte{})

	return bytes.Join([][]byte{prefixBytes, infoBytes, msgHex}, []byte{})
}

func (b *BitcoinSignMessage) Hash() []byte {
	return utils.Sha256(b.Encode())
}

//// AccountAuthenticatorImpl an implementation of an authenticator to provide generic verification across multiple types.
////
//// Types:
////   - [Ed25519Authenticator]
////   - [MultiEd25519Authenticator]
////   - [SingleKeyAuthenticator]
////   - [MultiKeyAuthenticator]
//type AccountAuthenticatorImpl interface {
//	bcs.Struct
//
//	// PublicKey is the public key that can be used to verify the signature.  It must be a valid on-chain representation
//	// and cannot be something like [Secp256k1PublicKey] on its own.
//	PublicKey() PublicKey
//
//	// Signature is a typed signature that can be verified by the public key. It must be a valid on-chain representation
//	// and cannot be something like [Secp256k1Signature] on its own.
//	Signature() Signature
//
//	// Verify Return true if the [AccountAuthenticator] can be cryptographically verified
//	Verify(data []byte) bool
//}

//export enum BuiltinAuthValidator {
//ROOCH = 0x00,
//BITCOIN = 0x01,
//// ETHEREUM= 0x02
//}

type AuthValidatorType uint64

const (
	AuthValidatorTypeRooch   AuthValidatorType = 0x00
	AuthValidatorTypeBitcoin AuthValidatorType = 0x01
	//AuthValidatorTypeEthereum AuthValidatorType = 0x02
)

//const (
//	ROOCH BuiltinAuthValidator = iota
//	BITCOIN
//	// ETHEREUM
//)

//pub struct Authenticator {
//pub auth_validator_id: u64,
//pub payload: Vec<u8>,
//}

type Authenticator struct {
	AuthValidatorId uint64 `json:"auth_validator_id"`
	Payload         []byte `json:"payload"`
}

func (au *Authenticator) MarshalBCS(ser *bcs.Serializer) {
	ser.U64(au.AuthValidatorId)
	ser.WriteBytes(au.Payload)
}
func (au *Authenticator) UnmarshalBCS(des *bcs.Deserializer) {
	au.AuthValidatorId = des.U64()
	au.Payload = des.ReadBytes()
}

//static async rooch(input: Bytes, signer: Signer) {
//const signature = await signer.sign(input)
//const pubKeyBytes = signer.getPublicKey().toBytes()
//const serializedSignature = new Uint8Array(1 + signature.length + pubKeyBytes.length)
//serializedSignature.set([SIGNATURE_SCHEME_TO_FLAG[signer.getKeyScheme()]])
//serializedSignature.set(signature, 1)
//serializedSignature.set(signer.getPublicKey().toBytes(), 1 + signature.length)
//
//return new Authenticator(BuiltinAuthValidator.ROOCH, serializedSignature)
//}

func RoochAuthValidator(input []byte, signer Signer[address.RoochAddress]) (*Authenticator, error) {
	signature, err := signer.Sign(input)
	if err != nil {
		return nil, err
	}
	//pubKeyBytes := signer.getPublicKey().toBytes()
	pubKeyBytes := signer.GetPublicKey().ToBytes()
	//serializedSignatureLen := 1 + len(signature.Bytes()) + len(pubKeyBytes)
	//serializedSignature := [serializedSignatureLen]byte{}
	//serializedSignature := make([]byte, 1 + len(signature.Bytes()) + len(pubKeyBytes))
	var serializedSignature []byte
	serializedSignature = append(serializedSignature, byte(SignatureSchemeToFlag[signer.GetKeyScheme()]))
	serializedSignature = append(serializedSignature, signature[:]...)
	serializedSignature = append(serializedSignature, pubKeyBytes[:]...)
	//serializedSignature.set([SIGNATURE_SCHEME_TO_FLAG[signer.getKeyScheme()]])
	//serializedSignature.set(signature, 1)
	//serializedSignature.set(signer.getPublicKey().toBytes(), 1 + signature.length)

	return &Authenticator{
		uint64(AuthValidatorTypeRooch), serializedSignature}, nil
}

func BitcoinAuthValidator(input *BitcoinSignMessage, signer Signer[address.AddressView], signWith string) (*Authenticator, error) {
	if !strings.HasPrefix(input.messageInfo, MessageInfoPrefix) {
		return nil, errors.New("invalid message info")
	}

	messageLength := len([]byte(input.messageInfo)) + len(hex.EncodeToString(input.txHash))

	var signData []byte
	if signWith == "hash" {
		signData = input.Hash()
	} else {
		signData = []byte(input.Raw())
	}

	signature, err := signer.Sign(signData)
	if err != nil {
		return nil, err
	}
	bitcoin_address, err := signer.GetBitcoinAddress()
	if err != nil {
		return nil, err
	}

	bitcoin_payload := BitcoinAuthPayload{
		Signature: signature,
		MessagePrefix: bytes.Join([][]byte{
			[]byte(input.messagePrefix),
			utils.VarintByteNum(uint64(messageLength)),
		}, []byte{}),
		MessageInfo: []byte(input.messageInfo),
		PublicKey:   signer.GetPublicKey().ToBytes(),
		FromAddress: bitcoin_address.ToBytes(),
	}
	payload, err := bcs.Serialize(&bitcoin_payload)
	if err != nil {
		return nil, err
	}

	//return NewAuthenticator(int(BITCOIN), payload), nil
	return &Authenticator{
		uint64(AuthValidatorTypeBitcoin), payload}, nil
}
