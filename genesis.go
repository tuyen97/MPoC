package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"github.com/boltdb/bolt"
	"log"
	"time"
)

type GenesisBlock struct {
	Hash      []byte
	Timestamp int64
	Balance   map[string]int
	BPs       []string
}

func (b *GenesisBlock) SetHash() {
	hash := sha256.Sum256(b.Serialize())
	b.Hash = hash[:]
}

func (b *GenesisBlock) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		log.Println("serialize error")
		return []byte{}
	}
	return result.Bytes()
}

func NewGenesis(balance map[string]int, bps []string) *GenesisBlock {
	g := GenesisBlock{
		Timestamp: time.Now().UnixNano(),
		Balance:   balance,
		BPs:       bps,
	}
	g.SetHash()
	return &g
}

func (g *GenesisBlock) Save() bool {
	if DBExists() {
		return false
	}
	db, _ := bolt.Open(dbFile, 0600, nil)
	defer db.Close()
	_ = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		if b == nil {
			b, _ = tx.CreateBucket([]byte(blocksBucket))
		}
		_ = b.Put(g.Hash, g.Serialize())
		_ = b.Put([]byte("g"), g.Hash)
		return nil
	})
	return true
}

//DeserializeBlock deserialize block from byte array
func DeserializeGenesis(d []byte) *GenesisBlock {
	var block GenesisBlock
	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&block)
	return &block
}

func GetGenesis() (bool, *GenesisBlock) {
	if !DBExists() {
		return false, &GenesisBlock{}
	}
	db, _ := bolt.Open(dbFile, 0600, nil)
	defer db.Close()
	var genesis *GenesisBlock

	_ = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		hash := b.Get([]byte("g"))
		data := b.Get(hash)
		genesis = DeserializeGenesis(data)
		return nil
	})
	return true, genesis
}

func (g *GenesisBlock) String() string {
	b, _ := json.Marshal(g)
	return string(b)
}

func GetInitialBPs() []string {
	b, g := GetGenesis()
	if !b {
		log.Println("Genesis does not exist")
	}
	return g.BPs
}
