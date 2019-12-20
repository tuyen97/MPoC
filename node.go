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
	node.peer = peer

	apiMemTxChan := make(chan *Transaction)
	api := Api{memTxChan: apiMemTxChan}
	node.api = api

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
