package main

import (
	"crypto/ecdsa"
)

// Block/ Tx are signable
type Signable interface {
	Sign(privKey ecdsa.PrivateKey)
	Verify(pub ecdsa.PublicKey) bool
}
