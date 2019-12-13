package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type Wallet struct {
	PublicKey  []byte
	PrivateKey ecdsa.PrivateKey
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, _ := ecdsa.GenerateKey(curve, rand.Reader)
	pubkey := append(private.PublicKey.X.Bytes(), private.Y.Bytes()...)
	return *private, pubkey
}

//NewWallet create new key pair
func NewWallet() *Wallet {
	priv, pub := NewKeyPair()
	wallet := Wallet{pub, priv}
	return &wallet
}

func (wallet Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(wallet.PublicKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)

	address := Base58Encode(fullPayload)
	return address
}

func GetAddress(pubkey []byte) []byte {
	pubKeyHash := HashPubKey(pubkey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)

	address := Base58Encode(fullPayload)
	return address
}
