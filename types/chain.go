package types

//export enum MultiChainID {
// Bitcoin = 0,
// Ether = 60,
// Sui = 784,
// Nostr = 1237,
// Rooch = 20230101,
// }
//

type MultiChainIDVariant uint64

const (
	MultiChainIDVariantBitcoin MultiChainIDVariant = 0
	MultiChainIDVariantEther   MultiChainIDVariant = 60 // Deprecated
	MultiChainIDVariantSui     MultiChainIDVariant = 784
	MultiChainIDVariantNostr   MultiChainIDVariant = 1237
	MultiChainIDVariantRooch   MultiChainIDVariant = 20230101
)

//func (mcid *MultiChainIDVariant) MarshalBCS(ser *bcs.Serializer) {
//	//ser.U64(uint64(mcid))
//	ser.U64(uint64(mcid))
//}

//func (mcid *MultiChainIDVariant) UnmarshalBCS(des *bcs.Deserializer) {
//	//fc.FunctionId.UnmarshalBCS(des)
//	//fc.Function = des.ReadString()
//	//fc.ArgTypes = bcs.DeserializeSequence[TypeTag](des)
//	//alen := des.Uleb128()
//	//fc.Args = make([][]byte, alen)
//	//for i := range alen {
//	//	fc.Args[i] = des.ReadBytes()
//	//}
//
//	fc.FunctionId.UnmarshalBCS(des)
//	//fc.Code = des.ReadBytes()
//	fc.TyArgs = bcs.DeserializeSequence[TypeTag](des)
//	alen := des.Uleb128()
//	fc.Args = make([][]byte, alen)
//	for i := range alen {
//		fc.Args[i] = des.ReadBytes()
//	}
//}
