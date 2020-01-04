package main

import (
	"sync"
	"time"
)

var lock sync.RWMutex

const logFileName string = "logfile.log"

const dbFile = "blockchain.db"

const blocksBucket = "blockchain"

const sqliteFile = ":memory:"

var TopK int64

const blockTime = int64(1 * time.Second)

const nCandidate = 50

const peerIp = "127.0.0.1"

var parameter = []float64{1.0, 1.0, 1.0, 1.0, 1.0}

const dnsServer = "127.0.0.1"
