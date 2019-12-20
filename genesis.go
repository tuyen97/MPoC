package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"
	"os"
	"time"
)

type GenesisBlock struct {
	Timestamp int64
	Balance   map[string]int
	BPs       []string
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
	return &g
}

//DeserializeBlock deserialize block from byte array
func DeserializeGenesis(d []byte) *GenesisBlock {
	var block GenesisBlock
	decoder := gob.NewDecoder(bytes.NewReader(d))
	decoder.Decode(&block)
	return &block
}

func GetGenesis() (bool, *GenesisBlock) {
	file, err := os.Open("/home/tuyen/Desktop/genesis.json")
	if err != nil {
		return false, &GenesisBlock{}
	}
	decoder := json.NewDecoder(file)
	var genesis GenesisBlock
	err = decoder.Decode(&genesis)
	if err != nil {
		log.Panic(err)
	}
	return true, &genesis
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
