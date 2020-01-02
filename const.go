package main

const logFileName string = "logfile.log"

const dbFile = "blockchain.db"

const blocksBucket = "blockchain"

const sqliteFile = ":memory:"

var TopK int64

var blockTime int64

const nCandidate = 50

const peerIp = "0.0.0.0"

var parameter = []float64{1.0, 1.0, 1.0, 1.0, 1.0}
