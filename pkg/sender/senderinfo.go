package sender

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
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

	nextNonce, err := s.client.NonceAt(context.Background(), from, nil)
	if err != nil {
		return fmt.Errorf("failed to get next nonce for sender: %w", err)
	}

	balance, err := s.client.BalanceAt(context.Background(), from, nil)
	if err != nil {
		return fmt.Errorf("failed to get sender balance: %w", err)
	}

	receipt, err := s.client.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil && !errors.Is(err, ethereum.NotFound) {
		return fmt.Errorf("failed to get transaction receipt: %w", err)
	}

	printTransactionInfo(tx, from, nextNonce, balance, receipt)

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

func printTransactionInfo(tx *types.Transaction, from common.Address, nextNonce uint64, balance *big.Int, receipt *types.Receipt) {
	fmt.Println("\n📋 Transaction Information")
	fmt.Println("· Hash: ", tx.Hash())
	fmt.Println("· Transaction Type: ", getTransactionTypeString(tx.Type()))
	fmt.Printf("· Nonce: %d \n", tx.Nonce())
	fmt.Printf("· Value: %s \n", decimal.NewFromBigInt(tx.Value(), -18))

	var fee *big.Int
	if tx.Type() == types.DynamicFeeTxType {
		fmt.Println("· Max Priority Fee Per Gas: ", tx.GasTipCap().String())
		fmt.Println("· Max Fee Per Gas: ", tx.GasFeeCap().String())
		fee = new(big.Int).Mul(tx.GasFeeCap(), big.NewInt(int64(tx.Gas())))
	} else {
		fmt.Printf("· Gas Price: %s Gwei\n", decimal.NewFromBigInt(tx.GasPrice(), -9))
		fee = new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(tx.Gas())))
	}
	fmt.Println("· Gas Limit: ", tx.Gas())
	fmt.Printf("· Fee: %s\n", decimal.NewFromBigInt(fee, -18))

	if len(tx.Data()) != 0 {
		printContractInfo(tx.Data())
	}

	// 输出发送者信息，无论交易是否上链都要输出
	fmt.Println("\n📬 Sender Information")
	fmt.Println("· Sender: ", from.Hex())
	fmt.Println("· Next Nonce: ", nextNonce)
	fmt.Printf("· Balance: %s\n", decimal.NewFromBigInt(balance, -18))

	// 检查是否有交易回执，判断交易是否已经上链
	if receipt != nil {
		fmt.Println("\n⛓️ Transaction has been mined!")
		fmt.Println("· Block Number: ", receipt.BlockNumber.Uint64())
		status := "Success"
		if receipt.Status == 0 {
			status = "Failed"
		}
		fmt.Println("· Transaction Status: ", status)
		return // 已经上链，无需再检查 nonce 和余额信息，提前返回
	}

	// 如果交易没有上链，检查 nonce 和余额
	if nextNonce > tx.Nonce() {
		fmt.Printf("\n⚠️  Nonce %d is already used\n", tx.Nonce())
	} else {
		fmt.Printf("\n✔️  Nonce is available\n")

		// 计算总费用（交易费 + 交易金额）
		totalCost := new(big.Int).Add(fee, tx.Value())

		// 检查余额是否足够支付交易和费用
		if balance.Cmp(totalCost) < 0 {
			// 余额不足，计算差额
			shortfall := new(big.Int).Sub(totalCost, balance)
			fmt.Printf("⚠️  Balance is not enough. Shortfall: %s\n", decimal.NewFromBigInt(shortfall, -18))
		} else {
			fmt.Println("✔️  Balance is sufficient")
		}
	}
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
