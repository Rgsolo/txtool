package token

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

func TestToken(t *testing.T) {
	d := common.Hex2Bytes("f242432a0000000000000000000000001ce1a3f7ed42c4d688822e800cae1b5620fe117700000000000000000000000010be62d6cf64c21a0d03fa340533ce8d800bab9f0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000")
	callData, err := ParseCallData(d, Erc1155)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(callData.Signature)
	for _, input := range callData.Inputs {
		fmt.Println(input.Value)
	}
}

func TestAbi(t *testing.T) {
	fmt.Println()

	json, err := abi.JSON(bytes.NewBuffer(erc20Abi))
	if err != nil {
		return
	}
	fmt.Println(json.Events["Transfer"].ID.String())

	fmt.Println(json.Methods["transferTokenOwnership"].Sig)
	fmt.Println(json.Methods["transferTokenOwnership"].Name)
	fmt.Println(json.Methods["transferTokenOwnership"].RawName)
	fmt.Println(json.Methods["transferTokenOwnership"].)

	fmt.Println(Erc721.Events["Transfer"].ID.String())
}

func TestName(t *testing.T) {
	client, err := rpc.Dial("tcp", "10.95.2.48:8001")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	err = client.Call("chain.getChainInfo", nil, &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}

func Test1(t *testing.T) {
	_, err := net.Dial("tcp", "10.95.2.48:8001")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	//var reply string
	//err = conn.Write("chain.getChainInfo", nil, &reply)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fmt.Println(reply)
}
