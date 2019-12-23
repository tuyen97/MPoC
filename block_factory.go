package main

import (
	"fmt"
	"time"

	"github.com/ipfs/go-log"
)

type BlockFactory struct {
	Address         string
	PeerBFChan      chan *Block
	BFPeerChan      chan *Block
	BFMemChan       chan bool
	MemBFChan       chan map[string]*Transaction
	ReturnBFMemChan chan map[string]*Transaction
}

var bfLogger = log.Logger("bf")

const TopK = 3

const blockTime = int64(1 * time.Second)

func (b *BlockFactory) ticker() {
	_, g := GetGenesis()
	//sleep 5 block before start
	sinceGenesis := (time.Now().UnixNano()-g.Timestamp)%blockTime + 5*blockTime
	time.Sleep(time.Duration(sinceGenesis))
	ticker := time.NewTicker(time.Duration(1 * time.Second))
	for {
		for now := range ticker.C {
			//blockNo := (now.UnixNano() - g.Timestamp) / blockTime
			currentSlot := (now.UnixNano() - g.Timestamp) % (3 * blockTime) / 1000000000
			//bfLogger.Infof("Slot now: %d", int(currentSlot))
			lastblock, err := GetLastBlock()
			//am i the current bp?
			if err == nil {
				if b.Address == lastblock.BPs[currentSlot] {
					b.BFMemChan <- true
				}
				//use initial bps
			} else {
				if b.Address == g.BPs[currentSlot] {
					b.BFMemChan <- true
				}
			}
		}
	}
}

func (b *BlockFactory) ServeInternal() {
	for {
		select {
		//gather txs to produce block
		case txs := <-b.MemBFChan:
			bfLogger.Infof("Got %d transaction", len(txs))
			_, g := GetGenesis()
			index := GetOrInitIndex()

			blockNo := (time.Now().UnixNano() - g.Timestamp) / blockTime
			currentSlot := (time.Now().UnixNano() - g.Timestamp) % (3 * blockTime) / 1000000000
			var tnx []Transaction
			for _, tx := range txs {
				//execute transaction
				if index.Update(tx) {
					tnx = append(tnx, *tx)
				}

			}

			var bps []string
			//end of epoch -> recalculate bps
			fmt.Println("c:", currentSlot)
			if currentSlot == TopK-1 {
				topk := index.GetTopKVote(TopK)
				fmt.Println("top k:", topk)
				bps = topk
			}

			block := Block{
				Hash:      nil,
				PrevHash:  nil,
				Index:     int(blockNo),
				Timestamp: time.Now().UnixNano(),
				Creator:   b.Address,
				Txs:       tnx,
				BPs:       bps,
			}
			block.SetHash()
			bfLogger.Infof("New block produced %d", int(block.Index))
			block.Save()
			b.ReturnBFMemChan <- txs
			b.BFPeerChan <- &block
		case block := <-b.PeerBFChan:
			txs := make(map[string]*Transaction)
			for _, tx := range block.Txs {
				txs[string(tx.ID)] = &tx
			}
			b.ReturnBFMemChan <- txs
		}
	}
}
func (b *BlockFactory) Start() {
	log.SetLogLevel("bf", "info")
	logger.Infof("i am %s", b.Address)
	go b.ticker()
	go b.ServeInternal()
	bfLogger.Infof("BF started")
}

//func main() {
//	b := BlockFactory{}
//	b.Start()
//}
