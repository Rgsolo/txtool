package main

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"log"
	"testing"
)

func TestSign(t *testing.T) {
	tx := wire.NewMsgTx(wire.TxVersion)

	var signBuilder SignsBuilder
	// signBuilder
	changeAddr, err := btcutil.DecodeAddress("tb1q8k2d4q9z7f6c9a2n2t4mz8y6j5v4l8z9k2x8q", &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}
	changeScript, err := txscript.PayToAddrScript(changeAddr)
	if err != nil {
		log.Fatal(err)
	}
	txOut := wire.NewTxOut(int64(546), changeScript)

	signBuilder.outs = append(signBuilder.outs, *txOut)

	// 构造交易输入
	prevTxHash, err := chainhash.NewHashFromStr("prev_tx_hash")
	if err != nil {
		log.Fatal(err)
	}
	prevOutPoint := wire.NewOutPoint(prevTxHash, 0)
	txIn := wire.NewTxIn(prevOutPoint, nil, nil)
	tx.AddTxIn(txIn)
	prevOutFetcher, err := signBuilder.txPrevOutputFetcher(tx)
	if err != nil {
		t.Fatal()
	}

	sigHashes := txscript.NewTxSigHashes(tx, prevOutFetcher)
	err = spendTaprootHash("", prevOutFetcher, tx.TxIn[0], tx, sigHashes, 0, "")
	if err != nil {
		t.Fatal()
	}
}
