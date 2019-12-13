package main

import (
	"context"
	"github.com/ipfs/go-log"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/sirupsen/logrus"
	logging "github.com/whyrusleeping/go-logging"
)

var logger = log.Logger("rendezvous")

func main() {
	log.SetAllLoggers(logging.WARNING)
	log.SetLogLevel("rendezvous", "info")
	host, err := libp2p.New(context.Background(), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/9000"))
	if err != nil {
		logrus.Error(err)
	}
	logger.Info("i am ", host.Addrs())
	logger.Info("bootstrap ", dht.DefaultBootstrapPeers)
}
