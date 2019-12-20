package main

//func TestTransaction_Verify(t *testing.T) {
//	wallets, _ := NewWallets()
//	pubkey := wallets.Wallets["16N2sT4C5JTfoqtpMUu79qkDod18K4E51J"].PublicKey
//
//	normaltx := Transaction{
//		ID:          nil,
//		Signature:   nil,
//		Sender:      pubkey,
//		Type:        0,
//		StakeAmount: 0,
//		Candidate:   nil,
//		Data:        "normal tx",
//	}
//	normaltx.Sign(wallets.Wallets["16N2sT4C5JTfoqtpMUu79qkDod18K4E51J"].PrivateKey)
//	normaltx.SetId()
//	if !normaltx.Verify() {
//		t.Error("Failed to sign and verify normal tx")
//	}
//
//	staketx := Transaction{
//		ID:          nil,
//		Signature:   nil,
//		Sender:      pubkey,
//		Type:        1,
//		StakeAmount: 10,
//		Candidate:   nil,
//		Data:        "",
//	}
//	staketx.Sign(wallets.Wallets["16N2sT4C5JTfoqtpMUu79qkDod18K4E51J"].PrivateKey)
//	staketx.SetId()
//	if !staketx.Verify() {
//		t.Error("Failed to sign and verify stake tx")
//	}
//}
