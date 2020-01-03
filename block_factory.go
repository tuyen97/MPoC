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
	MemBFChan       chan map[string]*Transaction
	ReturnBFMemChan chan map[string]*Transaction
}

var bfLogger = log.Logger("bf")
var lastBlock *Block = nil
var currentSlot int64
var blockNo int64
var index *Index

//only store tx when isLeader = true
var isLeader bool = false

func (b *BlockFactory) ticker() {
	_, g := GetGenesis()
	//sleep 5 block before start
	sinceGenesis := blockTime - ((time.Now().UnixNano() - g.Timestamp) % blockTime)
	bfLogger.Infof("Sleep %d", sinceGenesis)
	time.Sleep(time.Duration(sinceGenesis))
	ticker := time.NewTicker(time.Duration(blockTime))
	for {
		for now := range ticker.C {
			//blockNo := (now.UnixNano() - g.Timestamp) / blockTime
			currentSlot = (now.UnixNano() - g.Timestamp) % (TopK * blockTime) / 1000000000

			blockNo = (time.Now().UnixNano() - g.Timestamp) / blockTime
			// bfLogger.Infof("Current slot: %d", currentSlot)
			if i, _ := strconv.Atoi(b.Address); i < 21 {
				bfLogger.Infof("Slot now: %d", int(currentSlot))
			}

			// lastblock, err := GetLastBlock()
			//am i the current bp?
			// bfLogger.Infof("I am %s", b.Address)
			if lastBlock != nil {
				// bfLogger.Info("last block ", lastBlock)
				if SliceExists(lastBlock.BPs, b.Address) {
					lock.Lock()
					isLeader = true
					lock.Unlock()
				} else {
					lock.Lock()
					isLeader = false
					lock.Unlock()
				}
				if b.Address == lastBlock.BPs[currentSlot] {
					bfLogger.Infof("I am %s", b.Address)
					b.BFMemChan <- true
				}
				//use initial bps
			} else {
				// bfLogger.Info("genesis ", g)
				if SliceExists(g.BPs, b.Address) {
					lock.Lock()
					isLeader = true
					lock.Unlock()
				} else {
					lock.Lock()
					isLeader = false
					lock.Unlock()
				}
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
			//bfLogger.Infof("Got %d transaction", len(txs))

			// blockNo := (time.Now().UnixNano() - g.Timestamp) / blockTime
			// currentSlot := (time.Now().UnixNano() - g.Timestamp) % (TopK * blockTime) / 1000000000
			var tnx []Transaction
			for _, tx := range txs {
				tnx = append(tnx, *tx)

			}

			var bps []string
			// bl, err := GetLastBlock()
			//end of epoch -> recalculate bps
			fmt.Printf("for: %d\n", currentSlot)
			if currentSlot == 0 {
				topk := index.GetTopKBP(int(TopK))
				bfLogger.Infof("top k: %s", topk)
				bps = topk
			} else {
				if lastBlock == nil {
					bps = GetInitialBPs()
				} else {
					bps = lastBlock.BPs
				}
			}

			var prevHash []byte
			if lastBlock == nil {
				prevHash = []byte("genesis")
			} else {
				prevHash = lastBlock.Hash
			}

			block := Block{
				Hash:      nil,
				PrevHash:  prevHash,
				Index:     int(blockNo),
				Timestamp: time.Now().UnixNano(),
				Creator:   b.Address,
				Txs:       tnx,
				BPs:       bps,
			}
			//bfLogger.Info("Bps ", bps)
			block.SetHash()
			//bfLogger.Infof("New block produced %d", int(block.Index))

			//update database
			index.Update(&block)

			//bfLogger.Info("replace last block")
			lastBlock = &block
			go func() { b.ReturnBFMemChan <- txs }()
			go func() { b.BFPeerChan <- &block }()
		case block := <-b.PeerBFChan:
			txs := make(map[string]*Transaction)
			//logrus.Infof("Got %d transaction, bps: %s", len(block.Txs), block.BPs)
			// bfLogger.Infof("Got %d transaction", len(block.Txs))
			for _, tx := range block.Txs {
				txs[string(tx.ID)] = &tx
			}
			if lastBlock != nil {
				if block.Index > lastBlock.Index {
					index.Update(block)
				}
			} else {
				index.Update(block)
			}
			//bfLogger.Info("replace last block")
			lastBlock = block
			// block.Save()
			b.ReturnBFMemChan <- txs
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
