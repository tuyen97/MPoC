package main

import (
	"fmt"
	"testing"
)

func TestNewGenesis(t *testing.T) {
	ws, _ := NewWallets()
	balance := make(map[string]int)
	var bps []string
	for i := 0; i < 5; i++ {
		add := ws.CreateWallet()
		balance[add] = 1000000
		bps = append(bps, add)
	}
	ws.SaveToFile()
	genesis := NewGenesis(balance, bps)
	genesis.Save()
}

func TestGetGenesis(t *testing.T) {
	if DBExists() {
		_, genesis := GetGenesis()
		fmt.Println(genesis)
		if genesis == nil || len(genesis.BPs) == 0 {
			t.Error("cannot get genesis block")
		}
	}
}
