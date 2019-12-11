package main

import "bytes"

type TXInput struct {
	TXID      []byte
	Vout      int
	PubKey    []byte
	Signature []byte
}

//UsesKey checks that an input uses a specific key to unlock an output
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
