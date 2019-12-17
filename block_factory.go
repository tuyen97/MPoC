package main

import (
	"github.com/ipfs/go-log"
	"time"
)

type BlockFactory struct {
	Address   string
	BFMemChan chan bool
	MemBFChan chan []*TransactionWrapper
}

var bfLogger = log.Logger("bf")

const blockTime = int64(1 * time.Second)

func (b *BlockFactory) ticker() {
	_, g := GetGenesis()
	sinceGenesis := (time.Now().UnixNano() - g.Timestamp) / blockTime
	time.Sleep(time.Duration(sinceGenesis))
	ticker := time.NewTicker(time.Duration(1 * time.Second))
	for {
		for now := range ticker.C {
			//blockNo := (now.UnixNano() - g.Timestamp) / blockTime
			currentSlot := (now.UnixNano() - g.Timestamp) % blockTime
			if currentSlot == 0 {
				topK := GetOrInitIndex().GetTopKVote(3)
				if len(topK) == 0 {
					//leaders := GetInitialBPs()
				}
			}
			b.BFMemChan <- true
		}
	}
}

func (b *BlockFactory) ServeInternal() {
	for {
		select {
		case txs := <-b.MemBFChan:
			bfLogger.Infof("Got %d transaction", len(txs))
		}
	}
}
func (b *BlockFactory) Start() {
	log.SetLogLevel("bf", "info")
	go b.ticker()
	go b.ServeInternal()
	bfLogger.Infof("BF started")
}

//func main() {
//	b := BlockFactory{}
//	b.Start()
//}
