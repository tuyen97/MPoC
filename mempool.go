package main

import (
	"github.com/ipfs/go-log"
)

type MemPool struct {
	ApiMemTxChan    chan *Transaction            // api -> mem: transaction receive from api
	MemPeerTxChan   chan *Transaction            // mem -> peer: broadcast
	PeerMemTx       chan *Transaction            // peer -> mem: tx received
	BFMemChan       chan bool                    //BlockFactory -> mem : notify chan
	MemBFChan       chan map[string]*Transaction //mem -> BlockFactory: tx to forge block
	ReturnMemBFChan chan map[string]*Transaction //bf->mem : delete tx that is added in block
	IncomingBlock   chan *Block                  //block receive from
	TXPool          map[string]*Transaction
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
				memLogger.Infof("receive tx %s  from api", tx)
				m.MemPeerTxChan <- tx
				if m.TXPool[string(tx.ID)] == nil {
					m.TXPool[string(tx.ID)] = tx
				}
			case tx := <-m.PeerMemTx:
				//memLogger.Info("receive from peer")
				if m.TXPool[string(tx.ID)] == nil {
					memLogger.Infof("add tx : %s", tx)
					m.TXPool[string(tx.ID)] = tx
				}
			case <-m.BFMemChan:
				//memLogger.Info("Receive signal from bf")
				m.MemBFChan <- m.TXPool
			case txs := <-m.ReturnMemBFChan:
				//memLogger.Info("Receive tx from bf")
				for key, _ := range txs {
					delete(m.TXPool, key)
				}

			}
		}
	}()

	memLogger.Info("mem started")

}
