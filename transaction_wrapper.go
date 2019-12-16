package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type TransactionWrapper struct {
	ID        []byte
	TX        *Transaction
	Timestamp int64
}

func (t *TransactionWrapper) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(t)
	if err != nil {
		log.Println("serialize error")
		return []byte{}
	}
	return result.Bytes()
}

func (t *TransactionWrapper) SetId() {
	data := t.Serialize()
	hash := sha256.Sum256(data)
	t.ID = hash[:]
}

func DeserializeTxW(d []byte) *TransactionWrapper {
	var txw TransactionWrapper
	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&txw)
	return &txw
}
