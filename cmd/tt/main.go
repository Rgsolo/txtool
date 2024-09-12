package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rgsolo/txtool/pkg/chaininfo"
	"github.com/rgsolo/txtool/pkg/sender"
	"github.com/tidwall/pretty"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("Usage: tt <signed_tx> | tt send <signed_tx> | tt <chain_name>")
		return
	}

	if len(args) == 1 && !strings.HasPrefix(args[0], "0x") {
		// 处理链名称输入情况
		chainName := args[0]
		chain, err := chaininfo.GetChainInfoByName(chainName)
		if err != nil {
			fmt.Printf("Error getting chain info for '%s': %v\n", chainName, err)
			return
		}
		printChain(chain)
		return
	}

	isSend, signedTx := parseArguments(args)

	transaction, err := decodeTransaction(signedTx)
	if err != nil {
		fmt.Printf("failed to decode raw transaction: %v\n", err)
		return
	}

	printTransactionJson(transaction)

	chainInfo, err := chaininfo.GetChainInfo(int(transaction.ChainId().Int64()))
	if err != nil {
		fmt.Println("Error getting chain info:", err)
		return
	}

	printChainInfo(chainInfo)

	newSender := sender.NewSender(chainInfo.RPC[0])

	err = newSender.GetSenderInfo(transaction)
	if err != nil {
		fmt.Println("Error getting sender info:", err)
		return
	}

	if isSend {
		err = newSender.Send(transaction)
		if err != nil {
			fmt.Println("Error sending transaction:", err)
			return
		}
		fmt.Printf("send success: %s", transaction.Hash().String())
	}
}

func parseArguments(args []string) (bool, string) {
	isSend := false
	signedTx := args[0]
	if len(args) == 2 {
		signedTx = args[1]
		isSend = true
	}
	return isSend, signedTx
}

func decodeTransaction(signedTx string) (*types.Transaction, error) {
	decode, err := hexutil.Decode(signedTx)
	if err != nil {
		return nil, err
	}

	transaction := new(types.Transaction)
	err = transaction.UnmarshalBinary(decode)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func printTransactionJson(transaction *types.Transaction) {
	fmt.Println("\n🔖 Metadata")
	j, _ := transaction.MarshalJSON()
	coloredJson := pretty.Color(pretty.Pretty(j), pretty.TerminalStyle)
	fmt.Println(string(coloredJson))
}

func printChainInfo(chain *chaininfo.ChainInfo) {
	fmt.Println("\n🔗 chain information")
	fmt.Println("· chain name: ", chain.Name)
	fmt.Println("· chain ID: ", chain.ChainID)
	fmt.Println("· chain url: ", chain.RPC[0])
}

func printChain(chain *chaininfo.ChainInfo) {
	// 将链的信息格式化为漂亮的 JSON
	prettyJSON, err := json.MarshalIndent(chain, "", "  ")
	if err != nil {
		fmt.Printf("Failed to generate pretty JSON: %v\n", err)
		return
	}

	// 输出格式化的 JSON
	fmt.Println(string(prettyJSON))
}
