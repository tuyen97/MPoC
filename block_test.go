package main

import (
	"fmt"
	"testing"
)

func TestGenesisBlock(t *testing.T) {
	genesis := GenesisBlock("15LA5YCCy7BXqbsUM8sHFWVW66dcDXBCFr")
	fmt.Println(genesis.Hash)
}
