package types

import (
	"fmt"
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

//type FunctionCall struct {
//	FunctionId FunctionId `json:"function_id"`
//	TyArgs     []TypeTag  `json:"ty_args"`
//	Args       [][]byte   `json:"args"`
//}

//func (fc *FunctionCall) MoveActionType() MoveActionVariant {
//	return MoveActionVariantFunction
//}

//endregion

//region EntryFunction bcs.Struct

//// Args and TypeTag interfaces as placeholders - implement according to your BCS package
//type Args interface{}
//type TypeTag interface{}

// CallScript represents a script to be executed
type CallScript struct {
	Code     string    `json:"code"`
	Args     []Args    `json:"args"`
	TypeArgs []TypeTag `json:"typeArgs"`
}

// FunctionArgs represents either a full address/module/function path or a direct target
type FunctionArgs struct {
	Address  string `json:"address,omitempty"`
	Module   string `json:"module,omitempty"`
	Function string `json:"function,omitempty"`
	Target   string `json:"target,omitempty"`
}

// NewFunctionArgs creates a new FunctionArgs instance
func NewFunctionArgs(input map[string]interface{}) (*FunctionArgs, error) {
	if target, ok := input["target"].(string); ok {
		return &FunctionArgs{
			Target: target,
		}, nil
	}

	// Extract address, module, function
	address, ok1 := input["address"].(string)
	module, ok2 := input["module"].(string)
	function, ok3 := input["function"].(string)

	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("missing required fields")
	}

	return &FunctionArgs{
		Address:  address,
		Module:   module,
		Function: function,
	}, nil
}

// CallFunctionArgs represents arguments for calling a function
type CallFunctionArgs struct {
	*FunctionArgs           // Embed FunctionArgs
	Args          []Args    `json:"args,omitempty"`
	TypeArgs      []TypeTag `json:"typeArgs,omitempty"`
}

// NewCallFunctionArgs creates a new CallFunctionArgs instance
func NewCallFunctionArgs(funcArgs *FunctionArgs, args []Args, typeArgs []TypeTag) *CallFunctionArgs {
	return &CallFunctionArgs{
		FunctionArgs: funcArgs,
		Args:         args,
		TypeArgs:     typeArgs,
	}
}
