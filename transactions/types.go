package transactions

import (
	"fmt"
)

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
