package main

import "sync"

var (
	lastBlock *Block = nil
	lbLock    sync.Mutex
)

var (
	isLeader     bool = false
	isLeaderLock sync.Mutex
)

var (
	currentSlot     int64
	currentSlotLock sync.Mutex
)

var (
	blockNo     int64
	blockNoLock sync.Mutex
)

func GetBlockNo() int64 {
	blockNoLock.Lock()
	defer blockNoLock.Unlock()
	return blockNo
}

func UpdateBlockNo(new int64) {
	blockNoLock.Lock()
	defer blockNoLock.Unlock()
	blockNo = new
}

func GetCurrentSlot() int64 {
	currentSlotLock.Lock()
	defer currentSlotLock.Unlock()
	return currentSlot
}

func UpdateCurrentSlot(new int64) {
	currentSlotLock.Lock()
	defer currentSlotLock.Unlock()
	currentSlot = new
}

func GetIsLeader() bool {
	isLeaderLock.Lock()
	defer isLeaderLock.Unlock()
	return isLeader
}

func UpdateIsLeader(b bool) {
	isLeaderLock.Lock()
	defer isLeaderLock.Unlock()
	isLeader = b
}

func GetLastBlock() *Block {
	lbLock.Lock()
	defer lbLock.Unlock()
	return lastBlock
}

func UpdateLastBlock(block *Block) {
	lbLock.Lock()
	defer lbLock.Unlock()
	lastBlock = block
}
