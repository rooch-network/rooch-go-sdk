package api

import (
	"fmt"
	"github.com/rooch-network/rooch-go-sdk/address"
	"github.com/rooch-network/rooch-go-sdk/types"
)

func ParseStructTypeArgs(str string, normalizeAddress bool) ([]types.TypeTag, error) {
	tokens := SplitGenericParameters(str)
	result := make([]types.TypeTag, len(tokens))
	for i, token := range tokens {
		typetag, err := ParseTypeTagFromStr(token, normalizeAddress)
		if err != nil {
			return nil, err
		}
		result[i] = typetag
	}
	return result, nil
}

func ParseTypeTagFromStr(str string, normalizeAddress bool) (types.TypeTag, error) {
	switch str {
	case "address":
		return types.TypeTag{&types.AddressTag{}}, nil
	case "bool":
		return types.TypeTag{&types.BoolTag{}}, nil
	case "u8":
		return types.TypeTag{&types.U8Tag{}}, nil
	case "u16":
		return types.TypeTag{&types.U16Tag{}}, nil
	case "u32":
		return types.TypeTag{&types.U32Tag{}}, nil
	case "u64":
		return types.TypeTag{&types.U64Tag{}}, nil
	case "u128":
		return types.TypeTag{&types.U128Tag{}}, nil
	case "u256":
		return types.TypeTag{&types.U256Tag{}}, nil
	case "signer":
		return types.TypeTag{&types.SignerTag{}}, nil
	}

	if matches := VectorRegex.FindStringSubmatch(str); matches != nil {
		typeParam, err := ParseTypeTagFromStr(matches[1], normalizeAddress)
		if err != nil {
			return types.TypeTag{}, err
		}
		return types.TypeTag{
			&types.VectorTag{typeParam}}, nil
	}

	if matches := StructRegex.FindStringSubmatch(str); matches != nil {
		addr := matches[1]
		if normalizeAddress {
			addr = address.NormalizeRoochAddress(addr, true)
		}
		rooch_address, err := address.NewRoochAddress(addr)
		if err != nil {
			return types.TypeTag{}, err
		}

		var typeParams []types.TypeTag
		if matches[5] != "" {
			parsedTypeParams, err := ParseStructTypeArgs(matches[5], normalizeAddress)
			if err != nil {
				return types.TypeTag{}, err
			}
			typeParams = parsedTypeParams
		}

		return types.TypeTag{
			&types.StructTag{
				Address:    *rooch_address,
				Module:     matches[2],
				Name:       matches[3],
				TypeParams: typeParams,
			},
		}, nil
	}

	return types.TypeTag{}, fmt.Errorf("unknown type tag str: %s", str)
}
