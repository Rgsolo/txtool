package main

import (
	"encoding/hex"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"log"
	"testing"
)

func TestHexKeys(t *testing.T) {
	decodeString, _ := hex.DecodeString("")
	key, _ := btcec.PrivKeyFromBytes(decodeString)
	pubKey := key.PubKey()

	// 1BEkoAtvdGvpZaFCmuVqxXUb8YJ2wVewBf
	address, err := btcutil.NewAddressPubKey(pubKey.SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("P2PKH地址: %s", address.EncodeAddress())

	script, err := btcutil.NewAddressScriptHash(pubKey.SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}
	t.Logf("P2SH地址: %s", script.EncodeAddress())

	// 生成Bech32地址
	witnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(key.PubKey().SerializeCompressed()), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}
	t.Logf("Bech32地址: %s", witnessPubKeyHash.EncodeAddress())

	wif, err := btcutil.NewWIF(key, &chaincfg.MainNetParams, true)
	if err != nil {
		log.Fatal(err)
	}
	t.Log(wif.String())
}

func TestWIF(t *testing.T) {
	wif, err := btcutil.DecodeWIF("")
	if err != nil {
		log.Fatal(err)
	}
	key := wif.PrivKey
	// 生成Bech32地址
	witnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(key.PubKey().SerializeCompressed()), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}
	t.Logf("Bech32地址: %s", witnessPubKeyHash.EncodeAddress())

}

func TestSeed(t *testing.T) {
	masterKey := generateMasterKey(t)

	//derivationPath := "m/86'/0'/0'/0"
	// Derive the extended key using the derivation path
	extendedKey, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 84) // m/86'
	if err != nil {
		log.Fatal(err)
	}
	extendedKey, err = extendedKey.Derive(hdkeychain.HardenedKeyStart + 0) // m/86'/0'
	if err != nil {
		log.Fatal(err)
	}
	extendedKey, err = extendedKey.Derive(hdkeychain.HardenedKeyStart + 0) // m/86'/0'/0'
	if err != nil {
		log.Fatal(err)
	}
	extendedKey, err = extendedKey.Derive(0) // m/86'/0'/0'/0
	if err != nil {
		log.Fatal(err)
	}

	extendedKey, err = extendedKey.Derive(0) // m/86'/0'/0'/0/0
	if err != nil {
		log.Fatal(err)
	}

	// Convert the extended key to xprv format
	xprv := extendedKey.String()
	if err != nil {
		log.Fatal(err)
	}
	t.Logf("xprv: %s", xprv)

	// 派生公钥
	publicKey, err := extendedKey.ECPubKey()
	if err != nil {
		log.Fatal(err)
	}

	// 将地址格式化为BC1地址
	bc1Address, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(publicKey.SerializeCompressed()), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}
	// bc1qvlq8je98ra2pvr4yu9p3c8ct6d7qjep6ghg8v3
	t.Logf("P2WPKH: %s", bc1Address.EncodeAddress())

	// bc1p0ehfkl3rhssrj6pwgp5jn480vwr3sz9cw6nlut4v22rcc035fm9qmh2msv
}

// derivationPath := "m/86'/0'/0'/0"
var DerivationPath = accounts.DerivationPath{hdkeychain.HardenedKeyStart + 86, hdkeychain.HardenedKeyStart + 0, hdkeychain.HardenedKeyStart + 0, 0, 0}

func TestP2TR(t *testing.T) {
	extendedKey := generateMasterKey(t)
	var err error
	//derivationPath := "m/86'/0'/0'/0"
	for _, u := range DerivationPath {
		extendedKey, err = extendedKey.Derive(u)
	}
	// Convert the extended key to xprv format
	xprv := extendedKey.String()
	if err != nil {
		log.Fatal(err)
	}
	// xprvA37JQJs6QA23xeHBeWuMBy6VWNzfJiipXwC46xU5ejmgskwjt8oaWrfdFE9GMWeEsZRGBEgqiPCvX8wijGq8LRtrn5R5NhWWyggzx7H1JRG
	t.Logf("xprv: %s", xprv)

	// 派生公钥
	publicKey, err := extendedKey.ECPubKey()
	if err != nil {
		log.Fatal(err)
	}

	ecdsa := publicKey.ToECDSA()
	ethAddress := crypto.PubkeyToAddress(*ecdsa)
	// eth address: 0xD15214f915D57F5b3789e3e4b3607F0Fc3e031Bc
	t.Logf("eth address: %s", ethAddress)

	tapKey := txscript.ComputeTaprootKeyNoScript(publicKey)

	address, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(tapKey), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}

	// bc1pphfe5s4tldvyysffek5nh8gem0523w2r4qj6f8ttzd6s85wzq49qmsfc8s
	t.Logf("P2TR: %s", address)
}

func generateMasterKey(t *testing.T) *hdkeychain.ExtendedKey {
	seed := bip39.NewSeed("animal forest core rookie rice can design term detail company trouble network", "")

	// 使用BIP39种子生成主私钥
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		t.Fatal(err)
	}
	return masterKey
}
