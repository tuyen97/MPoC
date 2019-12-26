package main

type Candidate struct {
	Address          string
	TotalVote        int
	NumberOfVote     int
	TotalTxLastEpoch int
	LeaderCount      int
	TotalReward      int
}

type ByScore []Candidate

func (c Candidate) CalculateScore() float64 {
	return float64(c.TotalVote)*parameter[0] +
		float64(c.NumberOfVote)*parameter[1] +
		float64(c.TotalTxLastEpoch)*parameter[2] +
		float64(c.LeaderCount)*parameter[3] +
		float64(c.TotalReward)*parameter[4]
}

func (bps ByScore) Len() int {
	return len(bps)
}

func (bps ByScore) Less(i, j int) bool {
	return bps[i].CalculateScore() < bps[j].CalculateScore()
}
func (bps ByScore) Swap(i, j int) {
	bps[i], bps[j] = bps[j], bps[i]
}
