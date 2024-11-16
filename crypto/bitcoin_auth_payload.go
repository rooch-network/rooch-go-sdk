package crypto

import "github.com/rooch-network/rooch-go-sdk/bcs"

type BitcoinAuthPayload struct {
	Signature     []byte
	MessagePrefix []byte
	MessageInfo   string
	PublicKey     []byte
	FromAddress   []byte
}

func (bap *BitcoinAuthPayload) MarshalBCS(ser *bcs.Serializer) {
	ser.WriteBytes(bap.Signature)
	ser.WriteBytes(bap.MessagePrefix)
	ser.WriteString(bap.MessageInfo)
	ser.WriteBytes(bap.PublicKey)
	ser.WriteBytes(bap.FromAddress)
}

func (bap *BitcoinAuthPayload) UnmarshalBCS(des *bcs.Deserializer) {
	bap.Signature = des.ReadBytes()
	bap.MessagePrefix = des.ReadBytes()
	bap.MessageInfo = des.ReadString()
	bap.PublicKey = des.ReadBytes()
	bap.FromAddress = des.ReadBytes()
}
