package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ipfs/go-log"
	logrus "github.com/sirupsen/logrus"
)

type BlockFactory struct {
	Address         string
	PeerBFChan      chan *Block
	BFPeerChan      chan *Block
	BFMemChan       chan bool
	MemBFChan       chan []*Transaction
	ReturnBFMemChan chan []*Transaction
}

var bfLogger = log.Logger("bf")

var index *Index

//only store tx when isLeader = true

func (b *BlockFactory) ticker() {
	_, g := GetGenesis()
	//sleep 5 block before start
	sinceGenesis := 5*blockTime - ((time.Now().UnixNano() - g.Timestamp) % blockTime)
	bfLogger.Infof("Sleep %d", sinceGenesis)
	time.Sleep(time.Duration(sinceGenesis))
	ticker := time.NewTicker(time.Duration(blockTime))
	for {
		for now := range ticker.C {
			//blockNo := (now.UnixNano() - g.Timestamp) / blockTime
			UpdateCurrentSlot((now.UnixNano() - g.Timestamp) % (TopK * blockTime) / 1000000000)

			UpdateBlockNo((time.Now().UnixNano() - g.Timestamp) / blockTime)
			// bfLogger.Infof("Current slot: %d", currentSlot)
			if i, _ := strconv.Atoi(b.Address); i < 21 {
				bfLogger.Infof("Slot now: %d", int(GetCurrentSlot()))

			}

			// lastblock, err := GetLastBlock()
			//am i the current bp?
			// bfLogger.Infof("I am %s", b.Address)
			if GetLastBlock() != nil {
				// bfLogger.Info("last block ", lastBlock)
				if SliceExists(GetLastBlock().BPs, b.Address) {
					UpdateIsLeader(true)
				} else {
					UpdateIsLeader(false)
				}
				if b.Address == GetLastBlock().BPs[GetCurrentSlot()] {
					b.BFMemChan <- true
					bfLogger.Infof("I am %s", b.Address)
				}
				//use initial bps
			} else {
				// bfLogger.Info("genesis ", g)
				if SliceExists(g.BPs, b.Address) {
					UpdateIsLeader(true)
				} else {
					UpdateIsLeader(false)
				}
				if b.Address == g.BPs[GetCurrentSlot()] {
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
			// blockNo := (time.Now().UnixNano() - g.Timestamp) / blockTime
			// currentSlot := (time.Now().UnixNano() - g.Timestamp) % (TopK * blockTime) / 1000000000
			var tnx []*Transaction
			for _, tx := range txs {
				tnx = append(tnx, tx)

			}

			var bps []string
			// bl, err := GetLastBlock()
			//end of epoch -> recalculate bps
			bfLogger.Infof("for: %d\n", GetCurrentSlot())
			if GetCurrentSlot() == 0 {
				topk := index.GetTopKVote(int(TopK))
				bfLogger.Infof("top k: %s", topk)
				bps = topk
			} else {
				if GetLastBlock() == nil {
					bps = GetInitialBPs()
				} else {
					bps = GetLastBlock().BPs
				}
			}

			var prevHash []byte
			if GetLastBlock() == nil {
				prevHash = []byte("genesis")
			} else {
				prevHash = GetLastBlock().Hash
			}

			block := Block{
				Hash:      nil,
				PrevHash:  prevHash,
				Index:     int(GetBlockNo()),
				Timestamp: time.Now().UnixNano(),
				Creator:   b.Address,
				Txs:       tnx,
				BPs:       bps,
			}
			//bfLogger.Info("Bps ", bps)
			block.SetHash()
			// bfLogger.Infof("New block produced %d", int(block.Index))

			//update database
			index.Update(&block)

			// bfLogger.Info("replace last block")
			UpdateLastBlock(&block)
			b.ReturnBFMemChan <- txs
			b.BFPeerChan <- &block
			bfLogger.Info("done pr")
		case block := <-b.PeerBFChan:
			logrus.Infof("Got %d transaction, timestamp: %d", len(block.Txs), block.Timestamp)
			// bfLogger.Infof("Got %d transaction", len(block.Txs))
			if GetLastBlock() != nil {
				if block.Index > GetLastBlock().Index {
					index.Update(block)
				}
			} else {
				index.Update(block)
			}
			//bfLogger.Info("replace last block")
			UpdateLastBlock(block)
			// block.Save()
			b.ReturnBFMemChan <- block.Txs
		}
	}
}
func (b *BlockFactory) init() {
	_, g := GetGenesis()
	TopK = int64(len(g.BPs))
	index = GetOrInitIndex()
}
func (b *BlockFactory) Start() {
	log.SetLogLevel("bf", "info")
	f, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		// Cannot open log file. Logging to stderr
		fmt.Println(err)
	} else {
		logrus.SetOutput(f)
	}

	// logger.Infof("i am %s", b.Address)
	b.init()
	go b.ticker()
	go b.ServeInternal()
	bfLogger.Infof("BF started")
}

//func main() {
//	b := BlockFactory{}
//	b.Start()
//}
