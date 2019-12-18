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

	memPeerTxChan := make(chan *Transaction)
	peerMemTxChan := make(chan *Transaction)
	bfPeerChan := make(chan *Block)
	peerBfBlockChan := make(chan *Block)
	peer := Peer{
		MemPeerTxChan:   memPeerTxChan,
		PeerMemTxChan:   peerMemTxChan,
		BFPeerBlockChan: bfPeerChan,
		PeerBFBlockChan: peerBfBlockChan,
	}
	peer.Start(os.Args[1])

	apiMemTxChan := make(chan *Transaction)
	api := Api{memTxChan: apiMemTxChan}
	api.Start(os.Args[2])

	txPool := make(map[string]*Transaction)
	bfMemChan := make(chan bool)
	memBfChan := make(chan map[string]*Transaction)
	returnMemBFChan := make(chan map[string]*Transaction)
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

	bf := BlockFactory{
		BFPeerChan:      bfPeerChan,
		PeerBFChan:      peerBfBlockChan,
		ReturnBFMemChan: returnMemBFChan,
		BFMemChan:       bfMemChan,
		MemBFChan:       memBfChan,
		Address:         os.Args[3]}
	go bf.Start()
	nodeLogger.Info("Node started")
	select {}
}
