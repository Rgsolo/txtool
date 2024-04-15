package main

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"log"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func main() {
	// 构造交易输入
	prevTxHash, err := chainhash.NewHashFromStr("prev_tx_hash")
	if err != nil {
		log.Fatal(err)
	}
	prevOutPoint := wire.NewOutPoint(prevTxHash, 0)
	txIn := wire.NewTxIn(prevOutPoint, nil, nil)

	// 构造交易输出
	addr, err := btcutil.DecodeAddress("tb1q8k2d4q9z7f6c9a2n2t4mz8y6j5v4l8z9k2x8q", &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}
	script, err := txscript.PayToAddrScript(addr)
	if err != nil {
		log.Fatal(err)
	}
	txOut := wire.NewTxOut(100000, script)

	// 构造交易
	tx := wire.NewMsgTx(wire.TxVersion)
	tx.AddTxIn(txIn)
	tx.AddTxOut(txOut)

	// 计算交易费用
	feeRate := btcutil.Amount(10) // 每字节交易费用
	txSize := tx.SerializeSize()
	fee := feeRate * btcutil.Amount(txSize)

	connCfg := &rpcclient.ConnConfig{
		Host:         "",
		Endpoint:     "",
		User:         "",
		Pass:         "",
		DisableTLS:   false,
		HTTPPostMode: false,
	}
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	// 计算找零金额
	var totalInput btcutil.Amount
	for _, in := range tx.TxIn {
		hash := in.PreviousOutPoint.Hash
		prevTx, err := client.GetRawTransaction(&hash)
		if err != nil {
			log.Fatal(err)
		}
		prevTxOut := prevTx.MsgTx().TxOut[in.PreviousOutPoint.Index]
		totalInput += btcutil.Amount(prevTxOut.Value)
	}
	changeAmt := totalInput - btcutil.Amount(txOut.Value) - fee
	changeAddr, err := btcutil.DecodeAddress("tb1q8k2d4q9z7f6c9a2n2t4mz8y6j5v4l8z9k2x8q", &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}
	changeScript, err := txscript.PayToAddrScript(changeAddr)
	if err != nil {
		log.Fatal(err)
	}
	changeTxOut := wire.NewTxOut(int64(changeAmt), changeScript)
	tx.AddTxOut(changeTxOut)

	// 签名交易
	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)

	sigHashes := txscript.NewTxSigHashes(tx, prevOutFetcher)

	wif, err := btcutil.DecodeWIF("")
	if err != nil {
		log.Fatal(err)
	}
	privKey := wif.PrivKey
	pubKey := privKey.PubKey()
	sig, err := txscript.RawTxInWitnessSignature(tx, sigHashes, 0, int64(txOut.Value), nil, txscript.SigHashAll, privKey)
	if err != nil {
		log.Fatal(err)
	}
	witness := wire.TxWitness{}
	witness = append(witness, sig)
	witness = append(witness, pubKey.SerializeCompressed())
	txIn.Witness = witness

	// 广播交易
	hash, err := client.SendRawTransaction(tx, false)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction %v sent successfully\n", hash)
}
