package types

import (
	"github.com/rooch-network/rooch-go-sdk/bcs"
)

type Identifier string

// ModuleId the identifier for a module e.g. 0x1::coin
type ModuleId struct {
	Address AccountAddress
	Name    string
}

func (mod *ModuleId) MarshalBCS(ser *bcs.Serializer) {
	mod.Address.MarshalBCS(ser)
	ser.WriteString(mod.Name)
}

func (mod *ModuleId) UnmarshalBCS(des *bcs.Deserializer) {
	mod.Address.UnmarshalBCS(des)
	mod.Name = des.ReadString()
}

type FunctionId struct {
	ModuleId     ModuleId   `json:"module_id"`
	FunctionName Identifier `json:"function_name"`
}

func (fi *FunctionId) MarshalBCS(ser *bcs.Serializer) {
	fi.ModuleId.MarshalBCS(ser)
	ser.WriteString(string(fi.FunctionName))
}

func (fi *FunctionId) UnmarshalBCS(des *bcs.Deserializer) {
	fi.ModuleId.UnmarshalBCS(des)
	fi.FunctionName = Identifier(des.ReadString())
}
