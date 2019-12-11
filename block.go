package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"strconv"
	"time"
)

type Block struct {
	Hash      []byte
	PrevHash  []byte
	Index     int
	Timestamp int64
	Creator   string
	Txs       []Transaction
}

func GenesisBlock(Creator string) *Block {
	timestamp := time.Now().UnixNano()
	coinbase := NewCoinBaseTX(Creator, "base")
	block := &Block{
		Hash:      nil,
		PrevHash:  nil,
		Index:     0,
		Timestamp: timestamp,
		Creator:   Creator,
		Txs:       []Transaction{*coinbase},
	}
	block.SetHash()
	return block
}

//SetHash : set block hash
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	index := []byte(strconv.Itoa(b.Index))
	headers := bytes.Join([][]byte{b.PrevHash, b.HashTransactions(), index, timestamp, []byte(b.Creator)}, []byte{})
	hash := sha256.Sum256(headers)
	b.Hash = hash[:]
}

func (b *Block) HashTransactions() []byte {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Txs {
		txHashes = append(txHashes, tx.ID)
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))
	return txHash[:]
}

// Serialize serialize block to byte
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Println("serialize error")
		return []byte{}
	}
	return result.Bytes()
}
