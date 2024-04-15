package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func spendTaprootHash(
	inscriptionScript string,
	txPrevOutputFetcher *txscript.MultiPrevOutFetcher,
	txIn *wire.TxIn,
	originTx *wire.MsgTx,
	sigHashes *txscript.TxSigHashes,
	idx int,
	priKey string,
) error {
	decodeString, _ := hex.DecodeString(priKey)
	privateKey, _ := btcec.PrivKeyFromBytes(decodeString)
	txout := txPrevOutputFetcher.FetchPrevOutput(originTx.TxIn[idx].PreviousOutPoint)

	if inscriptionScript != "" {
		inscriptionScriptByte, err := hex.DecodeString(inscriptionScript)
		if err != nil {
			return err
		}
		proof := &txscript.TapscriptProof{
			TapLeaf:  txscript.NewBaseTapLeaf(schnorr.SerializePubKey(privateKey.PubKey())),
			RootNode: txscript.NewBaseTapLeaf(inscriptionScriptByte),
		}

		controlBlock := proof.ToControlBlock(privateKey.PubKey())
		controlBlockWitness, err := controlBlock.ToBytes()
		if err != nil {
			return err
		}

		witnessArray, err := txscript.CalcTapscriptSignaturehash(txscript.NewTxSigHashes(originTx, txPrevOutputFetcher),
			txscript.SigHashDefault, originTx, idx, txPrevOutputFetcher, txscript.NewBaseTapLeaf(inscriptionScriptByte))
		if err != nil {
			return err
		}
		signature, err := schnorr.Sign(privateKey, witnessArray)
		if err != nil {
			return err
		}
		witnessList := wire.TxWitness{signature.Serialize(), inscriptionScriptByte, controlBlockWitness}
		fmt.Printf("Taproot witness:%x\n", witnessList)

		txIn.Witness = witnessList
		return nil
	}

	// TODO: check txout.Value > 546
	witness, err := txscript.TaprootWitnessSignature(originTx,
		sigHashes,
		idx,
		txout.Value,
		txout.PkScript,
		txscript.SigHashDefault,
		privateKey)
	if err != nil {
		return err
	}

	txIn.Witness = witness
	return nil
}

type SignsBuilder struct {
	outs []wire.TxOut
}

// txPrevOutputFetcher add all PreviousOutPoint
func (sb SignsBuilder) txPrevOutputFetcher(tx *wire.MsgTx) (*txscript.MultiPrevOutFetcher, error) {
	if len(tx.TxIn) != len(sb.outs) {
		return nil, errors.New("signs builder length is not equal tx.TxIns")
	}
	inputFetcher := txscript.NewMultiPrevOutFetcher(nil)
	for i, s := range sb.outs {
		inputFetcher.AddPrevOut(tx.TxIn[i].PreviousOutPoint, &wire.TxOut{
			Value:    s.Value,
			PkScript: s.PkScript,
		})
	}

	return inputFetcher, nil
}
