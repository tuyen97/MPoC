package main

import (
	"github.com/ipfs/go-log"
	"os"
)

type Node struct {
	mempool      MemPool
	api          Api
	blockFactory BlockFactory
	peer         Peer
}

var nodeLogger = log.Logger("node")

func main() {

	memPeerTxChan := make(chan *TransactionWrapper)
	peerMemTxChan := make(chan *TransactionWrapper)
	peer := Peer{
		MemPeerTxChan: memPeerTxChan,
		PeerMemTxChan: peerMemTxChan,
	}
	peer.Start(os.Args[1])

	apiMemTxChan := make(chan *TransactionWrapper)
	api := Api{memTxChan: apiMemTxChan}
	api.Start(os.Args[2])

	txPool := make(map[string]*TransactionWrapper)
	bfMemChan := make(chan bool)
	memBfChan := make(chan map[string]*TransactionWrapper)
	returnMemBFChan := make(chan map[string]*TransactionWrapper)
	mem := MemPool{
		TXPool:          txPool,
		ApiMemTxChan:    apiMemTxChan,
		MemPeerTxChan:   memPeerTxChan,
		PeerMemTx:       peerMemTxChan,
		BFMemChan:       bfMemChan,
		MemBFChan:       memBfChan,
		ReturnMemBFChan: returnMemBFChan,
	}
	go mem.Start()

	bf := BlockFactory{ReturnBFMemChan: returnMemBFChan, BFMemChan: bfMemChan, MemBFChan: memBfChan, Address: os.Args[3]}
	go bf.Start()
	nodeLogger.Info("Node started")
	select {}
}
