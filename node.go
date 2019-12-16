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

	mem := MemPool{
		ApiMemTxChan:  apiMemTxChan,
		MemPeerTxChan: memPeerTxChan,
		PeerMemTx:     peerMemTxChan,
		IncomingBlock: nil,
		TXPool:        nil,
	}
	go mem.Start()

	nodeLogger.Info("Node started")
	select {}
}
