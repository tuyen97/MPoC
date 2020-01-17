package main

import (
	"github.com/ipfs/go-log"
	"sync"
)

type MemPool struct {
	ApiMemTxChan    chan *Transaction   // api -> mem: transaction receive from api
	MemPeerTxChan   chan *Transaction   // mem -> peer: broadcast
	PeerMemTx       chan *Transaction   // peer -> mem: tx received
	BFMemChan       chan bool           //BlockFactory -> mem : notify chan
	MemBFChan       chan []*Transaction //mem -> BlockFactory: tx to forge block
	ReturnMemBFChan chan []*Transaction //bf->mem : delete tx that is added in block
	IncomingBlock   chan *Block         //block receive from
	TXPool          map[string]*Transaction
}
type Pool struct {
	sync.Mutex
	Data map[string]*Transaction
}

var mempool *MemPool = nil
var memLogger = log.Logger("mem")

//func GetOrInitMempool() *MemPool {
//	if mempool != nil {
//		return mempool
//	}
//	mempool = &MemPool{}
//	mempool.TxInp = make(chan *TransactionWrapper)
//	mempool.IncomingBlock = make(chan *Block)
//	mempool.Start()
//	return mempool
//}

func (m *MemPool) Start() {
	log.SetLogLevel("mem", "info")
	go func() {
		for {
			select {
			case tx := <-m.ApiMemTxChan:
				//memLogger.Infof("receive tx %s  from api", tx)
				m.MemPeerTxChan <- tx
				if GetIsLeader() {
					if m.TXPool[string(tx.ID)] == nil {
						m.TXPool[string(tx.ID)] = tx
					}
				}
			case tx := <-m.PeerMemTx:
				//memLogger.Info("receive from peer")
				if GetIsLeader() {
					if m.TXPool[string(tx.ID)] == nil {
						m.TXPool[string(tx.ID)] = tx
					}
				}
			case <-m.BFMemChan:
				memLogger.Info("Receive signal from bf")
				//m.MemBFChan <- m.TXPool.Data
				var data []*Transaction
				for _, value := range m.TXPool {
					data = append(data, value)
				}
				memLogger.Info("Push to BF")
				m.MemBFChan <- data
				memLogger.Info("Push to BF done")
				memLogger.Info("done signal from bf")
			case txs := <-m.ReturnMemBFChan:
				//memLogger.Info("Receive tx from bf")
				for _, tx := range txs {
					delete(m.TXPool, string(tx.ID))
				}
			}
		}
	}()
	memLogger.Info("mem started")
}
