package main

import (
	"github.com/ipfs/go-log"
)

type MemPool struct {
	ApiMemTxChan  chan *TransactionWrapper // transaction receive from api
	MemPeerTxChan chan *TransactionWrapper // mem -> peer: broadcast
	PeerMemTx     chan *TransactionWrapper // peer -> mem: tx received
	IncomingBlock chan *Block              //block receive from
	TXPool        []*TransactionWrapper
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
	for {
		select {
		case tx := <-m.ApiMemTxChan:
			memLogger.Info("receive from api")
			m.MemPeerTxChan <- tx
			m.TXPool = append(m.TXPool, tx)
		case tx := <-m.PeerMemTx:
			memLogger.Info("receive from peer")
			m.TXPool = append(m.TXPool, tx)
		}
	}
}
