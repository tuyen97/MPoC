package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"log"
	"strconv"
)

type Transaction struct {
	ID          []byte   // sha256 hash
	Sender      string   // sender public key
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

func (tx *Transaction) prepareData() []byte {
	t := []byte(strconv.Itoa(tx.Type))
	var option []byte
	switch tx.Type {
	//normal tx
	case 0:
		option = []byte(tx.Data)
		//staking tx
	case 1:
		option = bytes.Join([][]byte{[]byte(tx.Sender), []byte(strconv.Itoa(tx.StakeAmount))}, []byte{})
		//voting tx
	case 2:
		for _, b := range tx.Candidate {
			option = bytes.Join([][]byte{[]byte(b)}, []byte{})
		}
		option = bytes.Join([][]byte{[]byte(tx.Sender), option}, []byte{})
	}
	dataToSign := bytes.Join([][]byte{t, option}, []byte{})
	return dataToSign
}

func DeserializeTx(d []byte) *Transaction {
	var tx Transaction
	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&tx)
	return &tx
}

func (t *Transaction) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(t)
	if err != nil {
		log.Println("serialize error")
		return []byte{}
	}
	return result.Bytes()
}

func (tx *Transaction) String() string {
	b, _ := json.Marshal(tx)
	return string(b)
}
