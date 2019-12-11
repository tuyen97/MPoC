package main

import (
	"fmt"
	"testing"
)

func TestCreateChain(t *testing.T) {
	CreateChain("15LA5YCCy7BXqbsUM8sHFWVW66dcDXBCFr")
}

func TestGetChain(t *testing.T) {
	if DbExists() {
		bc := GetChain()
		if bc == nil {
			t.Error("Db File exist but cannot get chain")
		}
	} else {
		bc := GetChain()
		if bc != nil {
			t.Error("DB FIle not exist")
		}
	}
}

func TestBlockChain_FindTransaction(t *testing.T) {
	if DbExists() {
		bc := GetChain()
		block, b := bc.GetBlock(bc.tip)
		if !b {
			t.Error("cannot get last block")
		}
		_, err := bc.FindTransaction(block.Txs[0].ID)
		if err != nil {
			t.Error("cannot get tx")
		}
	}
}

func TestBlockChain_GetBlock(t *testing.T) {
	if DbExists() {
		bc := GetChain()
		block, b := bc.GetBlock(bc.tip)
		if !b {
			t.Error("cannot get last block")
		}
		fmt.Println(block.Timestamp)
	}

}
