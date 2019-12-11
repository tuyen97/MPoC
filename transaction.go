package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

type Transaction struct {
	//Receiver string
	ID     []byte
	Type   int // 0:exchanging, 1:staking, 2: voting
	Input  []TXInput
	Output []TXOutput
}

//SetId set id
func (tx *Transaction) SetId() {
	var encoded bytes.Buffer
	var hash [32]byte
	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}
func (tx Transaction) IsCoinBase() bool {
	return len(tx.Input) == 1 && len(tx.Input[0].TXID) == 0 && tx.Input[0].Vout == -1
}

//NewCoinBaseTX creates a new coinbase transaction
func NewCoinBaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}
	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(10, to)
	tx := Transaction{nil, 0, []TXInput{txin}, []TXOutput{*txout}}
	tx.SetId()
	return &tx
}

//TrimmedCopy create a copy of this transaction take only th transaction id and vout
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput

	for _, input := range tx.Input {
		inputs = append(inputs, TXInput{input.TXID, input.Vout, nil, nil})
	}
	txCopy := Transaction{tx.ID, tx.Type, inputs, tx.Output}
	return txCopy

}

// Sign signs each input of a Transaction, proof that input is accept to use (by get the pubkey of source output), the output is belong to new user(by add the pubkey to new output)
func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinBase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Input {
		prevTx := prevTXs[hex.EncodeToString(vin.TXID)]
		//prev output => current input
		txCopy.Input[inID].Signature = nil
		//copy pubkey, proof owner of this input
		txCopy.Input[inID].PubKey = prevTx.Output[vin.Vout].PubKeyHash
		dataToSign := fmt.Sprintf("%x\n", txCopy)

		// //set id
		// txCopy.SetId()
		//set pubkey nil again
		txCopy.Input[inID].PubKey = nil

		r, s, _ := ecdsa.Sign(rand.Reader, &privKey, []byte(dataToSign))
		signature := append(r.Bytes(), s.Bytes()...)
		//signature of each input is created separately, differ other input in one transaction
		tx.Input[inID].Signature = signature
		txCopy.Input[inID].PubKey = nil
	}
}

//Verify transaction inputs with previous transaction
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinBase() {
		return true
	}
	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, vin := range tx.Input {
		prevTx := prevTXs[hex.EncodeToString(vin.TXID)]
		txCopy.Input[inID].Signature = nil
		txCopy.Input[inID].PubKey = prevTx.Output[vin.Vout].PubKeyHash
		txCopy.SetId()
		txCopy.Input[inID].PubKey = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		dataToVerify := fmt.Sprintf("%x\n", txCopy)
		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
			return false
		}
		txCopy.Input[inID].PubKey = nil
	}

	return true
}

// String returns a human-readable representation of a transaction
func (tx Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
	lines = append(lines, fmt.Sprintf("     Type: %d", tx.Type))
	for i, input := range tx.Input {

		lines = append(lines, fmt.Sprintf("     Input %d:", i))
		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.TXID))
		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
	}

	for i, output := range tx.Output {
		lines = append(lines, fmt.Sprintf("     Output %d:", i))
		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Amount))
		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
	}

	return strings.Join(lines, "\n")
}
