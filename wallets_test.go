package main

import (
	"fmt"
	"testing"
)

func TestWallets_LoadFromFile(t *testing.T) {
	wallets, _ := NewWallets()
	wallet := wallets.CreateWallet()
	wallets.SaveToFile()

	ws, _ := NewWallets()
	if len(ws.Wallets[wallet].PublicKey) > 0 {
		fmt.Println("ok")
	} else {
		fmt.Errorf("error")
	}

}
