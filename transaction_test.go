package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestNewCoinBaseTX(t *testing.T) {
	coinbase := NewCoinBaseTX("15LA5YCCy7BXqbsUM8sHFWVW66dcDXBCFr", "reward")
	if coinbase.IsCoinBase() != true {
		t.Errorf("cannot check coinbase")
	}

	if coinbase.Output[0].Amount != 10 {
		t.Errorf("reward amount invalid")
	}
	//fmt.Println(fmt.Sprintf("%s",coinbase.Output[0].PubKeyHash))
	pubkeydecode := Base58Decode([]byte("1EX5Ps4rygpjSVLSpK9Bmxnkrydn5GAwa6"))
	pubkeydecode = pubkeydecode[1 : len(pubkeydecode)-4]
	if bytes.Compare(coinbase.Output[0].PubKeyHash, pubkeydecode) != 0 {
		t.Errorf("reward to address failed")
	}
}

func TestTransaction_Verify(t *testing.T) {
	wallets, _ := NewWallets()
	coinbase := NewCoinBaseTX("1CWt1eoHjH42K79YzeDcvwx25cPDsFvo15", "reward")
	txin := TXInput{
		TXID:      coinbase.ID,
		Vout:      0,
		PubKey:    wallets.Wallets["1CWt1eoHjH42K79YzeDcvwx25cPDsFvo15"].PublicKey,
		Signature: nil,
	}
	txout := NewTXOutput(3, "1M8xqMqNGCf4ZLLRDqjgud6fsQykzqSUsf")
	tx := Transaction{
		ID:     nil,
		Type:   0,
		Input:  []TXInput{txin},
		Output: []TXOutput{*txout},
	}
	prevTx := make(map[string]Transaction)
	prevTx[hex.EncodeToString(coinbase.ID)] = *coinbase
	tx.Sign(wallets.Wallets["1CWt1eoHjH42K79YzeDcvwx25cPDsFvo15"].PrivateKey, prevTx)
	tx.SetId()
	if !tx.Verify(prevTx) {
		t.Errorf("verify error")
	}
}
