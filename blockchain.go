package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
	"os"
)

var logger *log.Logger = log.New()
var dbFile = "blockchain.db"
var bc *BlockChain = nil

const blocksBucket = "blocks"
const utxoBucket = "utxo"
const stakeBucket = "stake"

type BlockChain struct {
	DB  *bolt.DB //opening db
	tip []byte   //last block hash
}

func GetChain() *BlockChain {
	if bc != nil {
		return bc
	}
	bc := &BlockChain{}
	if DbExists() {
		db, _ := bolt.Open(dbFile, 0600, nil)
		_ = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))
			bc.tip = b.Get([]byte("l"))
			return nil
		})
		bc.DB = db
		return bc
	}
	return nil
}

func CreateChain(Creator string) {
	if DbExists() {
		logger.SetFormatter(&log.TextFormatter{ForceColors: true})
		logger.Error("Blockchain file already exist")
		return
	}
	var tip []byte
	genesis := GenesisBlock(Creator)
	db, _ := bolt.Open(dbFile, 0600, nil)
	_ = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			b, _ := tx.CreateBucket([]byte(blocksBucket))
			_ = b.Put(genesis.Hash, genesis.Serialize())
			_ = b.Put([]byte("l"), genesis.Hash)

			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	bc = &BlockChain{
		DB:  db,
		tip: tip,
	}
}

func (bc *BlockChain) AddBlock(block *Block) bool {
	//get best height
	bestHeight := 0
	lastBlock, _ := bc.GetBlock(bc.tip)
	bestHeight = lastBlock.Index

	//miss blocks, request
	if block.Index > bestHeight+1 {
		return false
	}

	_ = bc.DB.Update(func(tx *bolt.Tx) error {
		//get the bucket
		b := tx.Bucket([]byte(blocksBucket))
		hash := b.Get([]byte("l"))

		//get last block
		blockData := b.Get(hash)

		//if not store blockchain
		if blockData == nil {
			bc.tip = block.Hash
		} else {
			bl := *DeserializeBlock(blockData)
			if bl.Index < block.Index {
				bc.tip = block.Hash
			}
		}
		//put value by key
		_ = b.Put(block.Hash, block.Serialize())
		_ = b.Put([]byte("l"), block.Hash)
		return nil
	})

	return true
}

func (bc *BlockChain) GetBlock(hash []byte) (*Block, bool) {
	var block Block
	err := bc.DB.View(func(tx *bolt.Tx) error {
		//get the bucket
		b := tx.Bucket([]byte(blocksBucket))
		//get value by key
		blockData := b.Get(hash)
		if blockData == nil {
			return errors.New("Block is not found.")
		}

		block = *DeserializeBlock(blockData)

		return nil
	})
	if err != nil {
		return &block, false
	}
	return &block, true

}

func (bc *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{
		currentHash: bc.tip,
		db:          bc.DB,
	}
}

// FindTransaction finds a transaction by its ID
func (bc *BlockChain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()
	for {
		block := bci.Next()

		for _, tx := range block.Txs {
			if bytes.Compare(tx.ID, ID) == 0 {
				return tx, nil
			}
		}

		if len(block.PrevHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

//func GetBlockChain() (*BlockChain) {
//	if !dbExists() {
//
//	}
//}

func DbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

//DeserializeBlock deserialize block from byte array
func DeserializeBlock(d []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&block)
	return &block
}
