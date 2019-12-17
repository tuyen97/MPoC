package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"math/big"
	"strconv"
)

type Transaction struct {
	ID          []byte   // sha256 hash
	Signature   []byte   //ECDSA signature with secp256k1
	Sender      []byte   // sender public key
	Type        int      // 0:normal, 1:staking, 2: voting
	StakeAmount int      // use for staking
	Candidate   []string // list of candidate address, use for voting
	Data        string   // use for normal tx
	Timestamp   int64
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

func (tx *Transaction) Sign(key ecdsa.PrivateKey) {
	dataToSign := tx.prepareData()
	r, s, _ := ecdsa.Sign(rand.Reader, &key, dataToSign)
	signature := append(r.Bytes(), s.Bytes()...)
	tx.Signature = signature[:]
}

func (tx *Transaction) prepareData() []byte {
	t := []byte(strconv.Itoa(tx.Type))
	var option []byte
	switch tx.Type {
	//normal tx
	case 0:
		option = []byte(tx.Data)
		//staking tx
	case 1:
		option = bytes.Join([][]byte{tx.Sender, []byte(strconv.Itoa(tx.StakeAmount))}, []byte{})
		//voting tx
	case 2:
		for _, b := range tx.Candidate {
			option = bytes.Join([][]byte{[]byte(b)}, []byte{})
		}
		option = bytes.Join([][]byte{tx.Sender, option}, []byte{})
	}
	dataToSign := bytes.Join([][]byte{t, option}, []byte{})
	return dataToSign
}

func (tx *Transaction) Verify() bool {
	curve := elliptic.P256()
	r := big.Int{}
	s := big.Int{}
	sigLen := len(tx.Signature)
	r.SetBytes(tx.Signature[:(sigLen / 2)])
	s.SetBytes(tx.Signature[(sigLen / 2):])

	x := big.Int{}
	y := big.Int{}
	keyLen := len(tx.Sender)
	x.SetBytes(tx.Sender[:(keyLen / 2)])
	y.SetBytes(tx.Sender[(keyLen / 2):])

	dataToVerify := tx.prepareData()
	rawPubKey := ecdsa.PublicKey{curve, &x, &y}
	if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
		return false
	}
	return true
}

//func (tx *Transaction) Verify()  {
//
//}
//
////Verify transaction inputs with previous transaction
//func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
//	if tx.IsCoinBase() {
//		return true
//	}
//	txCopy := tx.TrimmedCopy()
//	curve := elliptic.P256()
//
//	for inID, vin := range tx.Input {
//		prevTx := prevTXs[hex.EncodeToString(vin.TXID)]
//		txCopy.Input[inID].Signature = nil
//		txCopy.Input[inID].PubKey = prevTx.Output[vin.Vout].PubKeyHash
//		txCopy.SetId()
//		txCopy.Input[inID].PubKey = nil
//
//		r := big.Int{}
//		s := big.Int{}
//		sigLen := len(vin.Signature)
//		r.SetBytes(vin.Signature[:(sigLen / 2)])
//		s.SetBytes(vin.Signature[(sigLen / 2):])
//
//		x := big.Int{}
//		y := big.Int{}
//		keyLen := len(vin.PubKey)
//		x.SetBytes(vin.PubKey[:(keyLen / 2)])
//		y.SetBytes(vin.PubKey[(keyLen / 2):])
//
//		dataToVerify := fmt.Sprintf("%x\n", txCopy)
//		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
//		if ecdsa.Verify(&rawPubKey, []byte(dataToVerify), &r, &s) == false {
//			return false
//		}
//		txCopy.Input[inID].PubKey = nil
//	}
//
//	return true
//}
//
//// String returns a human-readable representation of a transaction
//func (tx Transaction) String() string {
//	var lines []string
//
//	lines = append(lines, fmt.Sprintf("--- Transaction %x:", tx.ID))
//	lines = append(lines, fmt.Sprintf("     Type: %d", tx.Type))
//	for i, input := range tx.Input {
//
//		lines = append(lines, fmt.Sprintf("     Input %d:", i))
//		lines = append(lines, fmt.Sprintf("       TXID:      %x", input.TXID))
//		lines = append(lines, fmt.Sprintf("       Out:       %d", input.Vout))
//		lines = append(lines, fmt.Sprintf("       Signature: %x", input.Signature))
//		lines = append(lines, fmt.Sprintf("       PubKey:    %x", input.PubKey))
//	}
//
//	for i, output := range tx.Output {
//		lines = append(lines, fmt.Sprintf("     Output %d:", i))
//		lines = append(lines, fmt.Sprintf("       Value:  %d", output.Amount))
//		lines = append(lines, fmt.Sprintf("       Script: %x", output.PubKeyHash))
//	}
//
//	return strings.Join(lines, "\n")
//}
