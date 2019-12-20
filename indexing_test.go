package main

//func TestGetOrInitIndex(t *testing.T) {
//	_ = GetOrInitIndex()
//}
//
//func TestIndex_GetBalance(t *testing.T) {
//	index := GetOrInitIndex()
//	balance := index.GetBalance("1NdjTWZJrGsk9V7zDXj31kFTeAKyFvB34")
//	fmt.Println(balance)
//	if balance < -1 {
//		t.Error("get balance fail")
//	}
//}
//
//func TestIndex_Update(t *testing.T) {
//	ws, _ := NewWallets()
//	tx := Transaction{
//		ID:          nil,
//		Signature:   nil,
//		Sender:      ws.Wallets["16N2sT4C5JTfoqtpMUu79qkDod18K4E51J"].PublicKey,
//		Type:        1,
//		StakeAmount: 100,
//		Candidate:   nil,
//		Data:        "",
//	}
//
//	index := GetOrInitIndex()
//	index.Update(&tx)
//	if index.GetBalance("16N2sT4C5JTfoqtpMUu79qkDod18K4E51J") >= 1000000 {
//		t.Error("cannot execute tx")
//	}
//	candidate := []string{"1NdjTWZJrGsk9V7zDXj31kFTeAKyFvB34i", "1MqU2j2JcWR2S2gv48UuPZNh1BtVvfN8Be"}
//	vote := Transaction{
//		ID:          nil,
//		Signature:   nil,
//		Sender:      ws.Wallets["16N2sT4C5JTfoqtpMUu79qkDod18K4E51J"].PublicKey,
//		Type:        2,
//		StakeAmount: 0,
//		Candidate:   candidate,
//		Data:        "",
//	}
//	index.Update(&vote)
//	//fmt.Println(index.GetTotalVote("1DFpCwvyCpc6EargD6mtYfjGrpuDyfTNUp"))
//	if index.GetTotalVote("1NdjTWZJrGsk9V7zDXj31kFTeAKyFvB34i") <= 0 {
//		t.Error("Vote failed")
//	}
//}
//
//func TestIndex_GetTopKVote(t *testing.T) {
//	idx := GetOrInitIndex()
//	topk := idx.GetTopKVote(2)
//	fmt.Println(topk)
//}
