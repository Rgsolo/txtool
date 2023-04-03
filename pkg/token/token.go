package token

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

//go:embed erc20.abi.json
var erc20Abi []byte

//go:embed erc721.abi.json
var erc721Abi []byte

//go:embed erc1155.abi.json
var erc1155Abi []byte

type DecodedCallData struct {
	Signature string
	Name      string
	Inputs    []DecodedArgument
}

type DecodedArgument struct {
	SolType abi.Argument
	Value   interface{}
}

var (
	Erc20   = initABI(erc20Abi)
	Erc721  = initABI(erc721Abi)
	Erc1155 = initABI(erc1155Abi)
)

func initABI(bytesJSON []byte) *abi.ABI {
	json, err := abi.JSON(bytes.NewReader(bytesJSON))
	if err != nil {
		panic(err)
	}
	return &json
}

func ParseCallData(input []byte, abiSpec *abi.ABI) (*DecodedCallData, error) {
	if err := validateCallData(input); err != nil {
		return nil, err
	}
	argumentsData := input[4:]

	method, err := abiSpec.MethodById(input)
	if err != nil {
		return nil, err
	}
	values, err := method.Inputs.UnpackValues(argumentsData)
	if err != nil {
		return nil, fmt.Errorf("signature %q matches, but arguments mismatch: %v", method.String(), err)
	}

	return createDecodedCallData(method, values), nil
}

func validateCallData(input []byte) error {
	if len(input) < 4 {
		return fmt.Errorf("invalid call data, incomplete method signature (%d bytes < 4)", len(input))
	}
	argumentsData := input[4:]
	if len(argumentsData)%32 != 0 {
		return fmt.Errorf("invalid call data; length should be a multiple of 32 bytes (was %d)", len(argumentsData))
	}
	return nil
}

func createDecodedCallData(method *abi.Method, values []interface{}) *DecodedCallData {
	decoded := DecodedCallData{Signature: method.Sig, Name: method.RawName}
	for i := 0; i < len(method.Inputs); i++ {
		decoded.Inputs = append(decoded.Inputs, DecodedArgument{
			SolType: method.Inputs[i],
			Value:   values[i],
		})
	}
	return &decoded
}
