package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"errors"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
)

type Block struct {
	Hash      []byte
	PrevHash  []byte
	Index     int
	Timestamp int64
	Creator   string
	Txs       []Transaction
	BPs       []string
}

//SetHash : set block hash
func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	index := []byte(strconv.Itoa(b.Index))
	bps := []byte{}
	for _, bp := range b.BPs {
		bps = bytes.Join([][]byte{[]byte(bp)}, []byte{})
	}
	headers := bytes.Join([][]byte{b.PrevHash, b.HashTransactions(), index, timestamp, []byte(b.Creator), bps}, []byte{})
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

func (bl *Block) Save() bool {
	//if !DBExists() {
	//	return false
	//}
	lastBlock, err := GetLastBlock()
	bestHeight := 0
	if err == nil {
		bestHeight = lastBlock.Index
	}

	if bl.Index > bestHeight {
		db, _ := bolt.Open(dbFile, 0600, nil)
		defer db.Close()
		_ = db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			if b == nil {
				b, _ = tx.CreateBucket([]byte(blocksBucket))
			}
			_ = b.Put(bl.Hash, bl.Serialize())
			_ = b.Put([]byte("l"), bl.Hash)
			return nil
		})
		return true
	}

	return false
}

//
func GetLastBlock() (*Block, error) {
	if !DBExists() {
		return &Block{}, errors.New("Blockchain empty")
	}
	db, _ := bolt.Open(dbFile, 0600, nil)
	defer db.Close()
	var block *Block

	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		hash := b.Get([]byte("l"))
		data := b.Get(hash)
		block = DeserializeBlock(data)
		return nil
	})
	if len(block.BPs) == 0 {
		return &Block{}, errors.New("Blockchain empty")
	}
	return block, nil
}

//DeserializeBlock deserialize block from byte array
func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&block)
	return &block
}

func (block *Block) String() string {
	b, _ := json.Marshal(block)
	return string(b)
}
