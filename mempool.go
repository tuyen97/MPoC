package main

import (
	"github.com/ipfs/go-log"
	"sync"
)

type MemPool struct {
	ApiMemTxChan    chan *Transaction            // api -> mem: transaction receive from api
	MemPeerTxChan   chan *Transaction            // mem -> peer: broadcast
	PeerMemTx       chan *Transaction            // peer -> mem: tx received
	BFMemChan       chan bool                    //BlockFactory -> mem : notify chan
	MemBFChan       chan map[string]*Transaction //mem -> BlockFactory: tx to forge block
	ReturnMemBFChan chan map[string]*Transaction //bf->mem : delete tx that is added in block
	IncomingBlock   chan *Block                  //block receive from
	TXPool          Pool
}
type Pool struct {
	sync.RWMutex
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
				lock.RLock()
				if isLeader {
					m.TXPool.Lock()
					if m.TXPool.Data[string(tx.ID)] == nil {
						m.TXPool.Data[string(tx.ID)] = tx
					}
					m.TXPool.Unlock()
				}
				lock.RUnlock()
			case tx := <-m.PeerMemTx:
				//memLogger.Info("receive from peer")
				lock.RLock()
				if isLeader {
					m.TXPool.Lock()
					if m.TXPool.Data[string(tx.ID)] == nil {
						//memLogger.Infof("add tx : %s", tx)
						m.TXPool.Data[string(tx.ID)] = tx
					}
					m.TXPool.Unlock()
				}
				lock.RUnlock()
			case <-m.BFMemChan:
				memLogger.Info("Receive signal from bf")
				m.TXPool.RLock()
				//m.MemBFChan <- m.TXPool.Data
				data := make(map[string]*Transaction)
				for key, value := range m.TXPool.Data {
					data[key] = value
				}
				memLogger.Info("Push to BF")
				select {
				case m.MemBFChan <- data:
					memLogger.Info("Push to BF done")
				default:
					memLogger.Info("Can't push to BF")
				}
				m.TXPool.RUnlock()
				memLogger.Info("done signal from bf")
			case txs := <-m.ReturnMemBFChan:
				//memLogger.Info("Receive tx from bf")
				m.TXPool.Lock()
				for key, _ := range txs {
					delete(m.TXPool.Data, key)
				}
				m.TXPool.Unlock()
			}
		}
	}()
	memLogger.Info("mem started")
}
