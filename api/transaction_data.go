// Copyright (c) RoochNetwork
// SPDX-License-Identifier: Apache-2.0

package api

import (
	"fmt"
	"github.com/rooch-network/rooch-go-sdk/bcs"
	"github.com/rooch-network/rooch-go-sdk/types"
	"github.com/rooch-network/rooch-go-sdk/utils"

	"github.com/rooch-network/rooch-go-sdk/address"
	"strings"
)

const DEFAULT_GAS = uint64(50000000)

// CallScript represents a script to be executed
type CallScript struct {
	Code     string   `json:"code"`
	Args     []Args   `json:"args"`
	TypeArgs []string `json:"type_args"`
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
	//*FunctionArgs          // Embed FunctionArgs
	Address  string   `json:"address,omitempty"`
	Module   string   `json:"module,omitempty"`
	Function string   `json:"function,omitempty"`
	Target   string   `json:"target,omitempty"`
	Args     []Args   `json:"args,omitempty"`
	TypeArgs []string `json:"type_args,omitempty"`
}

//// NewCallFunctionArgs creates a new CallFunctionArgs instance
//func NewCallFunctionArgs(funcArgs *FunctionArgs, args []Args, typeArgs []string) *CallFunction {
//	return &CallFunctionArgs{
//		Address:  funcArgs.Address,
//		Module:   funcArgs.Module,
//		Function: funcArgs.Address,
//		Args:     args,
//		TypeArgs: typeArgs,
//	}
//}

//type CallFunctionArgs struct {
//	Target   string
//	Address  string
//	Module   string
//	Function string
//	Args     []bcs.Args
//	TypeArgs []string
//}

type CallFunction struct {
	Address  string
	Module   string
	Function string
	Args     []Args
	TypeArgs []string
}

func NewCallFunction(input CallFunctionArgs) *CallFunction {
	var pkg, mod, fn string
	if input.Target != "" {
		parts := strings.Split(input.Target, "::")
		pkg, mod, fn = parts[0], parts[1], parts[2]
	} else {
		pkg, mod, fn = input.Address, input.Module, input.Function
	}

	if input.Args == nil {
		input.Args = make([]Args, 0)
	}
	if input.TypeArgs == nil {
		input.TypeArgs = make([]string, 0)
	}

	return &CallFunction{
		Address:  pkg,
		Module:   mod,
		Function: fn,
		Args:     input.Args,
		TypeArgs: input.TypeArgs,
	}
}

func (c *CallFunction) FunctionId() string {
	return fmt.Sprintf("%s::%s::%s",
		address.NormalizeRoochAddress(c.Address, true),
		c.Module,
		c.Function)
}

func (c *CallFunction) EncodeArgs() []string {
	result := make([]string, len(c.Args))
	for i, arg := range c.Args {
		result[i] = arg.EncodeWithHex()
	}
	return result
}

func (c *CallFunction) EncodeArgsToByteArrays() [][]uint8 {
	result := make([][]uint8, len(c.Args))
	for i, arg := range c.Args {
		result[i] = arg.Encode()
	}
	return result
}

type MoveActionType interface {
	isActionType()
}

//type MoveActionType int

const (
	MoveActionVariantFunction int = 1
	MoveActionVariantScript   int = 2
)

//type MoveActionImpl interface {
//	MoveActionType() MoveActionType // This is specifically to ensure that wrong types don't end up here
//}

//func (*CallFunction) MoveActionType() MoveActionType { return MoveActionVariantScript }
//func (*CallScript) MoveActionType() MoveActionType   { return MoveActionVariantFunction }

func (*CallFunction) isActionType() {}
func (*CallScript) isActionType()   {}

type MoveAction struct {
	Scheme int
	Val    MoveActionType
}

func NewCallFunctionAction(input CallFunctionArgs) *MoveAction {
	return &MoveAction{
		Scheme: MoveActionVariantFunction,
		Val:    NewCallFunction(input),
	}
}

func NewCallScriptAction(input *CallScript) *MoveAction {
	return &MoveAction{
		Scheme: MoveActionVariantScript,
		Val:    input,
	}
}

type TransactionData struct {
	Sender         *types.RoochAddress
	SequenceNumber *uint64
	ChainId        *uint64
	MaxGas         uint64
	Action         *MoveAction
}

func NewTransactionData(
	action *MoveAction,
	sender string,
	sequenceNumber uint64,
	chainId uint64,
	maxGas uint64,
) (*TransactionData, error) {
	if maxGas == 0 {
		maxGas = DEFAULT_GAS
	}

	var senderAddr *types.RoochAddress
	if sender != "" {
		addr, err := address.NewRoochAddress(sender)
		if err != nil {
			return nil, err
		}
		senderAddr = addr
	}

	return &TransactionData{
		Sender:         senderAddr,
		SequenceNumber: &sequenceNumber,
		ChainId:        &chainId,
		MaxGas:         maxGas,
		Action:         action,
	}, nil
}

func (t *TransactionData) Encode() []byte {
	call := t.Action.Val.(*CallFunction)

	data := bcs.RoochTransactionData{
		Sender:         *t.Sender,
		SequenceNumber: *t.SequenceNumber,
		ChainId:        *t.ChainId,
		MaxGas:         t.MaxGas,
		Action: bcs.TransactionAction{
			Kind: "CallFunction",
			FunctionId: bcs.FunctionId{
				ModuleId: bcs.ModuleId{
					Address: call.Address,
					Name:    call.Module,
				},
				Name: call.Function,
			},
			Args:     call.EncodeArgsToByteArrays(),
			TypeArgs: call.TypeArgs,
		},
	}

	return data.Serialize()
}

func (t *TransactionData) Hash() []byte {
	return utils.Sha3256(t.Encode())
}
