package main

import (
	"context"
	"fmt"
	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/peer"
	_ "github.com/libp2p/go-libp2p-pubsub"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/p2p/discovery"
	"time"
)

var logger = log.Logger("peer")

type Peer struct {
	MemPeerTxChan chan *TransactionWrapper // mem -> peer :broadcast
	PeerMemTxChan chan *TransactionWrapper //peer -> mem: receive from broadcast
}

type DiscoveryNotifiee struct {
	peerchan chan peer.AddrInfo
}

func (d *DiscoveryNotifiee) HandlePeerFound(peerInfo peer.AddrInfo) {
	d.peerchan <- peerInfo
}

func (p *Peer) broadcastTx(pub *pubsub.Topic, ctx context.Context, tx TransactionWrapper) {
	err := pub.Publish(ctx, tx.Serialize())
	if err != nil {
		logger.Error("Cannot broadcast tx")
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
		tx := DeserializeTxW(msg.GetData())
		p.PeerMemTxChan <- tx
	}

}

func (p *Peer) Start(port string) {
	//log.SetAllLoggers(logging.WARNING)
	log.SetLogLevel("peer", "info")
	ctx := context.Background()
	//new host
	host, err := libp2p.New(ctx, libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", port)))
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
				if err := host.Connect(ctx, addrinfo); err != nil {
					logger.Error("Connection failed:", err)
				}
			}
		}
	}()

	//new pubsub
	pubsub, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		logger.Error("Cannot init pubsub")
	}
	//register topics
	RegisterTopic(pubsub, ctx, "genesis", handleGenesis)
	txTopic := RegisterTopic(pubsub, ctx, "tx", p.handleIncomingTx)

	//serve other component in node
	p.ServeInternal(txTopic, ctx)

	logger.Info("Peer started")
}

func (p *Peer) ServeInternal(pub *pubsub.Topic, ctx context.Context) {
	go func() {
		for {
			select {
			case tx := <-p.MemPeerTxChan:
				logger.Info("Broadcast tx")
				p.broadcastTx(pub, ctx, *tx)
			}
		}
	}()

}

//func main() {
//	//log.SetAllLoggers(logging.WARNING)
//	//log.SetLogLevel("rendezvous", "info")
//
//	ctx := context.Background()
//	//new host
//	host, err := libp2p.New(ctx, libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/%s/tcp/%s", os.Args[1], os.Args[2])))
//	if err != nil {
//		logger.Error(err)
//	}
//
//	fmt.Println(host.Addrs())
//
//	//new pubsub
//	pubsub, err := pubsub.NewGossipSub(ctx, host)
//	if err != nil {
//		logger.Error("Cannot init pubsub")
//	}
//	topic, err := pubsub.Join("genesis")
//	if err != nil {
//		logger.Error("Cannot join topic")
//	}
//
//	sub, _ := topic.Subscribe()
//	go func() {
//		for {
//			msg, err := sub.Next(ctx)
//			if err != nil {
//				//fmt.Println("Error (sub.Next): %v", err)
//				panic(err)
//			}
//
//			fmt.Printf("%s: %s\n", msg.GetFrom(), string(msg.GetData()))
//		}
//	}()
//
//	//new mdns
//	service, err := discovery.NewMdnsService(ctx, host, 2*time.Second, "test")
//
//	peerchan := make(chan peer.AddrInfo)
//	//register notifiee
//	service.RegisterNotifee(&DiscoveryNotifiee{peerchan: peerchan})
//	go func() {
//		for {
//			select {
//			case addrinfo := <-peerchan:
//				logger.Infof("Connecting to %s", addrinfo.String())
//				logger.Infof("topic peers: %s", topic.ListPeers())
//				if err := host.Connect(ctx, addrinfo); err != nil {
//					logger.Error("Connection failed:", err)
//				}
//			}
//		}
//	}()
//
//	select {}
//}
