package main

import (
	"github.com/ipfs/go-log"
	"time"
)

type BlockFactory struct {
	Address         string
	PeerBFChan      chan *Block
	BFMemChan       chan bool
	MemBFChan       chan map[string]*Transaction
	ReturnBFMemChan chan map[string]*Transaction
}

var bfLogger = log.Logger("bf")

const TopK = 3

const blockTime = int64(1 * time.Second)

func (b *BlockFactory) ticker() {
	_, g := GetGenesis()
	sinceGenesis := (time.Now().UnixNano() - g.Timestamp) / blockTime
	time.Sleep(time.Duration(sinceGenesis))
	ticker := time.NewTicker(time.Duration(1 * time.Second))
	for {
		for now := range ticker.C {
			//blockNo := (now.UnixNano() - g.Timestamp) / blockTime
			currentSlot := (now.UnixNano() - g.Timestamp) % (3 * blockTime) / 1000000000
			//bfLogger.Infof("Slot now: %d", int(currentSlot))
			lastblock, err := GetLastBlock()
			//am i the current bp?
			if err != nil {
				if b.Address == lastblock.BPs[currentSlot] {
					b.BFMemChan <- true
				}
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
			blockNo := (time.Now().UnixNano() - g.Timestamp) / blockTime
			currentSlot := (time.Now().UnixNano() - g.Timestamp) % blockTime
			var tnx []Transaction
			for _, tx := range txs {
				tnx = append(tnx, *tx)
			}

			var bps []string
			lastblock, err := GetLastBlock()
			//end of epoch -> recalculate bps
			if currentSlot == TopK-1 {
				topk := GetOrInitIndex().GetTopKVote(TopK)
				//not enough bps
				if len(topk) < TopK {
					//if have block
					if err != nil {
						bps = lastblock.BPs
					} else {
						bps = GetInitialBPs()
					}
				} else {
					//bps = new calculated value
					bps = topk
				}
			}

			block := Block{
				Hash:      nil,
				PrevHash:  nil,
				Index:     int(blockNo),
				Timestamp: time.Now().UnixNano(),
				Creator:   []byte(b.Address),
				Txs:       tnx,
				BPs:       bps,
			}
			block.SetHash()
			bfLogger.Infof("New block produced %d", int(block.Index))
			b.ReturnBFMemChan <- txs
		case block := <-b.PeerBFChan:

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
