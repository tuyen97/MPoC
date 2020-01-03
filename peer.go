package main

import (
	"context"
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	connmgr "github.com/libp2p/go-libp2p-connmgr"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/libp2p/go-libp2p-core/peer"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	yamux "github.com/libp2p/go-libp2p-yamux"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"time"
)

var logger = log.Logger("peer")

type Peer struct {
	MemPeerTxChan   chan *Transaction // mem -> peer :broadcast
	PeerMemTxChan   chan *Transaction //peer -> mem: receive from broadcast
	PeerBFBlockChan chan *Block       // peer -> bf: receive from broadcast
	BFPeerBlockChan chan *Block       // bf -> peer: block
}

type DiscoveryNotifiee struct {
	peerchan chan peer.AddrInfo
}

func (d *DiscoveryNotifiee) HandlePeerFound(peerInfo peer.AddrInfo) {
	d.peerchan <- peerInfo
}

func (p *Peer) broadcastTx(pub *pubsub.Topic, ctx context.Context, tx Transaction) {
	err := pub.Publish(ctx, tx.Serialize())
	if err != nil {
		logger.Error("Cannot broadcast tx")
	}
}
func (p *Peer) broadcastBlock(pub *pubsub.Topic, ctx context.Context, block *Block) {
	err := pub.Publish(ctx, block.Serialize())
	if err != nil {
		logger.Error("Cannot broadcast Block")
	}
}

func RegisterTopic(pubsub *pubsub.PubSub, ctx context.Context, topic string, handler func(sub *pubsub.Subscription, ctx context.Context)) *pubsub.Topic {
	t, err := pubsub.Join(topic)
	if err != nil {
		logger.Error("Cannot join topic genesis")
	}
	sub, _ := t.Subscribe()
	go handler(sub, ctx)
	return t
}

func handleGenesis(sub *pubsub.Subscription, ctx context.Context) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			//fmt.Println("Error (sub.Next): %v", err)
			panic(err)
		}

		fmt.Printf("%s: %s\n", msg.GetFrom(), string(msg.GetData()))
	}

}

func (p *Peer) handleIncomingTx(sub *pubsub.Subscription, ctx context.Context) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			//fmt.Println("Error (sub.Next): %v", err)
			panic(err)
		}
		tx := DeserializeTx(msg.GetData())
		//logger.Info("Receive tx")
		p.PeerMemTxChan <- tx
	}

}

func (p *Peer) handleIncomingBlock(sub *pubsub.Subscription, ctx context.Context) {
	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			//fmt.Println("Error (sub.Next): %v", err)
			panic(err)
		}
		block := DeserializeBlock(msg.GetData())
		// logger.Info("Receive block")
		p.PeerBFBlockChan <- block
	}
}

func (p *Peer) Start(port string) {
	//log.SetAllLoggers(logging.WARNING)
	log.SetLogLevel("peer", "info")
	ctx := context.Background()
	//new host
	connMgr := connmgr.NewConnManager(4, 12, 1*time.Second)
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%s", peerIp, port)),
		libp2p.ConnectionManager(connMgr),
		libp2p.Muxer("/yamux/1.0.0", yamux.DefaultTransport),
	)
	if err != nil {
		logger.Error(err)
	}
	logger.Info("i am ", host.Addrs())

	//new mdns
	service, err := discovery.NewMdnsService(ctx, host, 2*time.Second, "test")
	peerchan := make(chan peer.AddrInfo)

	//register notifiee
	service.RegisterNotifee(&DiscoveryNotifiee{peerchan: peerchan})
	go func() {
		for {
			select {
			case addrinfo := <-peerchan:
				if host.Network().Connectedness(addrinfo.ID) != network.Connected {
					if err := host.Connect(ctx, addrinfo); err != nil {
						logger.Error("Connection failed:", err)
					}
				}
			}
		}
	}()

	//new pubsub
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		logger.Error("Cannot init pubsub")
	}
	pub := make(map[string]*pubsub.Topic)
	//register topics
	RegisterTopic(ps, ctx, "genesis", handleGenesis)
	txTopic := RegisterTopic(ps, ctx, "tx", p.handleIncomingTx)
	pub["tx"] = txTopic
	blockTopic := RegisterTopic(ps, ctx, "block", p.handleIncomingBlock)
	pub["block"] = blockTopic

	//serve other component in node
	p.ServeInternal(pub, ctx)

	logger.Info("Peer started")
}

func (p *Peer) ServeInternal(pub map[string]*pubsub.Topic, ctx context.Context) {
	go func() {
		for {
			select {
			case tx := <-p.MemPeerTxChan:
				//logger.Info("Broadcast tx")
				p.broadcastTx(pub["tx"], ctx, *tx)
			case block := <-p.BFPeerBlockChan:
				logger.Info("Broadcast block")
				p.broadcastBlock(pub["block"], ctx, block)
			}
		}
	}()

}
