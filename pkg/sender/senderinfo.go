package sender

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rgsolo/txtool/pkg/token"
	"github.com/shopspring/decimal"
)

type Sender struct {
	client *ethclient.Client
}

func NewSender(rpcURL string) *Sender {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum client: %v", err)
	}
	return &Sender{client: client}
}

func (s *Sender) GetSenderInfo(tx *types.Transaction) error {
	from, err := getTransactionSender(tx)
	if err != nil {
		return err
	}

	nextNonce, err := s.client.PendingNonceAt(context.Background(), from)
	if err != nil {
		return fmt.Errorf("failed to get next nonce for sender: %w", err)
	}

	balance, err := s.client.BalanceAt(context.Background(), from, nil)
	if err != nil {
		return fmt.Errorf("failed to get sender balance: %w", err)
	}

	printTransactionInfo(tx, from, nextNonce, balance)

	return nil
}

func (s *Sender) Send(tx *types.Transaction) error {
	return s.client.SendTransaction(context.Background(), tx)
}

func getTransactionSender(tx *types.Transaction) (common.Address, error) {
	from, err := types.Sender(types.NewLondonSigner(tx.ChainId()), tx)
	if err != nil {
		return common.Address{}, fmt.Errorf("\nfailed to get transaction sender: %w", err)
	}
	return from, nil
}

func printTransactionInfo(tx *types.Transaction, from common.Address, nextNonce uint64, balance *big.Int) {
	fmt.Println("\n📋 transaction information")
	fmt.Println("· hash: ", tx.Hash())
	fmt.Println("· Transaction Type: ", getTransactionTypeString(tx.Type()))
	fmt.Printf("· nonce: %d \n", tx.Nonce())

	var fee *big.Int
	if tx.Type() == types.DynamicFeeTxType {
		fmt.Println("· Max Priority Fee Per Gas: ", tx.GasTipCap().String())
		fmt.Println("· Max Fee Per Gas: ", tx.GasFeeCap().String())
		fee = new(big.Int).Mul(tx.GasFeeCap(), big.NewInt(int64(tx.Gas())))
	} else {
		fmt.Printf("· Gas Price: %s Gwei\n", decimal.NewFromBigInt(tx.GasPrice(), -9))
		fee = new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))
	}
	fmt.Println("· gasLimit: ", tx.Gas())
	fmt.Printf("· fee: %s\n", decimal.NewFromBigInt(fee, -18))

	if len(tx.Data()) != 0 {
		printContractInfo(tx.Data())
	}

	fmt.Println("\n📬 sender information")
	fmt.Println("· sender: ", from.Hex())
	fmt.Println("· next nonce :", nextNonce)
	fmt.Printf("· balance : %s\n", decimal.NewFromBigInt(balance, -18))

	if nextNonce > tx.Nonce() {
		fmt.Printf("\n⚠️  nonce %d is already used ", tx.Nonce())
	} else {
		fmt.Printf("\n✔️  nonce is available")
		checkNonceAndBalance(tx, fee, balance)
	}
}

func calculateTotalCost(tx *types.Transaction) *big.Int {
	var fee *big.Int
	if tx.Type() == types.DynamicFeeTxType {
		fmt.Println("🌱Max Priority Fee Per Gas: ", tx.GasTipCap().String())
		fmt.Println("🌱Max Fee Per Gas: ", tx.GasFeeCap().String())
		fee = new(big.Int).Mul(tx.GasFeeCap(), big.NewInt(int64(tx.Gas())))
	} else {
		fmt.Println("🌱Gas Price: ", tx.GasPrice().String())
		fee = new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))
	}
	fmt.Println("🌱gasLimit: ", tx.Gas())
	return fee
}

func checkNonceAndBalance(tx *types.Transaction, fee *big.Int, balance *big.Int) {
	totalCost := new(big.Int).Add(fee, tx.Value())

	if balance.Cmp(totalCost) < 0 {
		fmt.Println("\n⚠️ balance is not enough")
	} else {
		fmt.Println("\n✔️ balance is enough~")
	}
}

func printContractInfo(data []byte) {
	fmt.Println("\n📥 input data")

	contractType, methodSig, args, err := token.ParseTransactionData(data)
	if err != nil {
		fmt.Println("\n🙁not supported contract")
		return
	}

	fmt.Println("· contractType: ", contractType)
	fmt.Println("· methodSig: ", methodSig)
	fmt.Println("· args: ", args)

	// contractAbis := []*abi.ABI{token.Erc20, token.Erc721, token.Erc1155}
	// contractNames := []string{"erc20", "erc721", "erc1155"}

	// var decodedData *token.DecodedCallData
	// var err error

	// for i, contractAbi := range contractAbis {
	// 	decodedData, err = token.ParseCallData(data, contractAbi)
	// 	if err == nil {
	// 		fmt.Printf("· %s: %s \n", contractNames[i], decodedData.Signature)
	// 		break
	// 	}
	// }

	// for _, input := range decodedData.Inputs {
	// 	fmt.Printf("· %s[%s]: %s \n", input.SolType.Name, input.SolType.Type, input.Value)
	// }
}

func getTransactionTypeString(txType uint8) string {
	switch txType {
	case uint8(types.LegacyTxType):
		return "Legacy"
	case uint8(types.DynamicFeeTxType):
		return "EIP-1559 (Dynamic Fee)"
	default:
		return "Unknown"
	}
}
