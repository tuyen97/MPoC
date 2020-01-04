package main

import (
	"github.com/ipfs/go-log"
)

type Node struct {
	mempool      MemPool
	api          Api
	blockFactory BlockFactory
	peer         Peer
}

var nodeLogger = log.Logger("node")

func (node *Node) Start(api, peer string) {
	log.SetLogLevel("node", "info")
	node.api.Start(api)
	node.blockFactory.Start()
	node.mempool.Start()
	node.peer.Start(peer)
	nodeLogger.Info("Node started")
}

func (node *Node) Init(addr string) {

	memPeerTxChan := make(chan *Transaction, 1000)
	peerMemTxChan := make(chan *Transaction, 1000)
	bfPeerChan := make(chan *Block, 1000)
	peerBfBlockChan := make(chan *Block, 1000)
	peer := Peer{
		MemPeerTxChan:   memPeerTxChan,
		PeerMemTxChan:   peerMemTxChan,
		BFPeerBlockChan: bfPeerChan,
		PeerBFBlockChan: peerBfBlockChan,
	}
	node.peer = peer

	apiMemTxChan := make(chan *Transaction, 1000)
	api := Api{memTxChan: apiMemTxChan}
	node.api = api

	txPool := Pool{Data: make(map[string]*Transaction)}
	bfMemChan := make(chan bool)
	memBfChan := make(chan map[string]*Transaction, 100)
	returnMemBFChan := make(chan map[string]*Transaction, 100)
	mem := MemPool{
		TXPool:          txPool,
		ApiMemTxChan:    apiMemTxChan,
		MemPeerTxChan:   memPeerTxChan,
		PeerMemTx:       peerMemTxChan,
		BFMemChan:       bfMemChan,
		MemBFChan:       memBfChan,
		ReturnMemBFChan: returnMemBFChan,
	}
	node.mempool = mem

	bf := BlockFactory{
		BFPeerChan:      bfPeerChan,
		PeerBFChan:      peerBfBlockChan,
		ReturnBFMemChan: returnMemBFChan,
		BFMemChan:       bfMemChan,
		MemBFChan:       memBfChan,
		Address:         addr}
	node.blockFactory = bf
}
