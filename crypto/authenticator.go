package crypto

import (
	"github.com/rooch-network/rooch-go-sdk/bcs"
	"github.com/rooch-network/rooch-go-sdk/transactions"
)

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

func RoochAuthValidator(input []byte, signer Signer) (transactions.Authenticator, error) {
	signature, err := signer.Sign(input)
	if err != nil {
		return transactions.Authenticator{}, err
	}
	//pubKeyBytes := signer.getPublicKey().toBytes()
	pubKeyBytes := signer.GetPublicKey().Bytes()
	//serializedSignatureLen := 1 + len(signature.Bytes()) + len(pubKeyBytes)
	//serializedSignature := [serializedSignatureLen]byte{}
	//serializedSignature := make([]byte, 1 + len(signature.Bytes()) + len(pubKeyBytes))
	var serializedSignature []byte
	serializedSignature = append(serializedSignature, GetSignatureFlag(signer.GetKeyScheme()))
	serializedSignature = append(serializedSignature, signature.Bytes()[:]...)
	serializedSignature = append(serializedSignature, pubKeyBytes[:]...)
	//serializedSignature.set([SIGNATURE_SCHEME_TO_FLAG[signer.getKeyScheme()]])
	//serializedSignature.set(signature, 1)
	//serializedSignature.set(signer.getPublicKey().toBytes(), 1 + signature.length)

	return Authenticator{
		uint64(AuthValidatorTypeRooch), serializedSignature}, nil
}

//
//static async bitcoin(
//input: BitcoinSignMessage,
//signer: Signer,
//signWith: 'hash' | 'raw' = 'hash',
//): Promise<Authenticator> {
//if (!input.messageInfo.startsWith(MessageInfoPrefix)) {
//throw Error('invalid message info')
//}
//
//const messageLength = bytes('utf8', input.messageInfo).length + toHEX(input.txHash).length
//const sign = await signer.sign(signWith === 'hash' ? input.hash() : bytes('utf8', input.raw()))
//
//const payload = bcs.BitcoinAuthPayload.serialize({
//signature: sign,
//messagePrefix: concatBytes(bytes('utf8', input.messagePrefix), varintByteNum(messageLength)),
//messageInfo: bytes('utf8', input.messageInfo),
//publicKey: signer.getPublicKey().toBytes(),
//fromAddress: bytes('utf8', signer.getBitcoinAddress().toStr()),
//}).toBytes()
//
//return new Authenticator(BuiltinAuthValidator.BITCOIN, payload)
//}
//}

////region AccountAuthenticator
//
//// AccountAuthenticatorType single byte representing the spot in the enum from the Rust implementation
//type AccountAuthenticatorType uint8
//
//const (
//	AccountAuthenticatorEd25519      AccountAuthenticatorType = 0 // AccountAuthenticatorEd25519 is the authenticator type for ed25519 accounts
//	AccountAuthenticatorMultiEd25519 AccountAuthenticatorType = 1 // AccountAuthenticatorMultiEd25519 is the authenticator type for multi-ed25519 accounts
//	AccountAuthenticatorSingleSender AccountAuthenticatorType = 2 // AccountAuthenticatorSingleSender is the authenticator type for single-key accounts
//	AccountAuthenticatorMultiKey     AccountAuthenticatorType = 3 // AccountAuthenticatorMultiKey is the authenticator type for multi-key accounts
//)
//
//// AccountAuthenticator a generic authenticator type for a transaction
////
//// Implements:
////   - [AccountAuthenticatorImpl]
////   - [bcs.Marshaler]
////   - [bcs.Unmarshaler]
////   - [bcs.Struct]
//type AccountAuthenticator struct {
//	Variant AccountAuthenticatorType // Variant is the type of authenticator
//	Auth    AccountAuthenticatorImpl // Auth is the actual authenticator
//}
//
////region AccountAuthenticator AccountAuthenticatorImpl implementation
//
//// PubKey returns the public key of the authenticator
//func (ea *AccountAuthenticator) PubKey() PublicKey {
//	return ea.Auth.PublicKey()
//}
//
//// Signature returns the signature of the authenticator
//func (ea *AccountAuthenticator) Signature() Signature {
//	return ea.Auth.Signature()
//}
//
//// Verify returns true if the authenticator can be cryptographically verified
//func (ea *AccountAuthenticator) Verify(data []byte) bool {
//	return ea.Auth.Verify(data)
//}
//
////endregion

////region AccountAuthenticator bcs.Struct implementation
//
//// MarshalBCS serializes the [AccountAuthenticator] to the BCS format
////
//// Implements:
////   - [bcs.Marshaler]
//func (ea *AccountAuthenticator) MarshalBCS(ser *bcs.Serializer) {
//	ser.Uleb128(uint32(ea.Variant))
//	ea.Auth.MarshalBCS(ser)
//}
//
//// UnmarshalBCS deserializes the [AccountAuthenticator] from the BCS format
////
//// Implements:
////   - [bcs.Unmarshaler]
//func (ea *AccountAuthenticator) UnmarshalBCS(des *bcs.Deserializer) {
//	kindNum := des.Uleb128()
//	if des.Error() != nil {
//		return
//	}
//	ea.Variant = AccountAuthenticatorType(kindNum)
//	switch ea.Variant {
//	case AccountAuthenticatorEd25519:
//		ea.Auth = &Ed25519Authenticator{}
//	case AccountAuthenticatorMultiEd25519:
//		ea.Auth = &MultiEd25519Authenticator{}
//	case AccountAuthenticatorSingleSender:
//		ea.Auth = &SingleKeyAuthenticator{}
//	case AccountAuthenticatorMultiKey:
//		ea.Auth = &MultiKeyAuthenticator{}
//	default:
//		des.SetError(fmt.Errorf("unknown AccountAuthenticator kind: %d", kindNum))
//		return
//	}
//	ea.Auth.UnmarshalBCS(des)
//}
//
//func (ea *AccountAuthenticator) FromKeyAndSignature(key PublicKey, sig Signature) error {
//	switch key.(type) {
//	case *Ed25519PublicKey:
//		switch sig.(type) {
//		case *Ed25519Signature:
//			ea.Variant = AccountAuthenticatorEd25519
//			ea.Auth = &Ed25519Authenticator{
//				PubKey: key.(*Ed25519PublicKey),
//				Sig:    sig.(*Ed25519Signature),
//			}
//		default:
//			return errors.New("invalid signature type for Ed25519PublicKey")
//		}
//	case *MultiEd25519PublicKey:
//		switch sig.(type) {
//		case *MultiEd25519Signature:
//			ea.Variant = AccountAuthenticatorMultiEd25519
//			ea.Auth = &MultiEd25519Authenticator{
//				PubKey: key.(*MultiEd25519PublicKey),
//				Sig:    sig.(*MultiEd25519Signature),
//			}
//		default:
//			return errors.New("invalid signature type for MultiEd25519PublicKey")
//		}
//	case *AnyPublicKey:
//		switch sig.(type) {
//		case *AnySignature:
//			ea.Variant = AccountAuthenticatorSingleSender
//			ea.Auth = &SingleKeyAuthenticator{
//				PubKey: key.(*AnyPublicKey),
//				Sig:    sig.(*AnySignature),
//			}
//		default:
//			return errors.New("invalid signature type for AnyPublicKey")
//		}
//	case *MultiKey:
//		switch sig.(type) {
//		case *MultiKeySignature:
//			ea.Variant = AccountAuthenticatorMultiKey
//			ea.Auth = &MultiKeyAuthenticator{
//				PubKey: key.(*MultiKey),
//				Sig:    sig.(*MultiKeySignature),
//			}
//		default:
//			return errors.New("invalid signature type for MultiKey")
//		}
//	default:
//		return errors.New("Invalid key type")
//	}
//	return nil
//}

//endregion
//endregion
