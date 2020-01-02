package main

import (
	"fmt"
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/ipfs/go-log"
	"sort"
)

var idx *Index = nil

type Index struct {
	Conn *sqlite3.Conn
}

func ReIndex() {

}

var indexLogger = log.Logger("bf")

func (i *Index) GetTxCountForAddr(addr string, epoch int) int {
	//update balance
	stmt, _ := i.Conn.Prepare("SELECT * FROM CountTX WHERE address=? AND epoch=?")
	_ = stmt.Exec(addr, epoch)
	_ = stmt.Reset()
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			break
		}
		if !hasRow {
			// The query is finished
			break
		}

		// Use Scan to access column data from a row
		var count int
		var address string
		var epoch int
		err = stmt.Scan(&address, &count, &epoch)
		if err != nil {
			break
		}
		return count
	}
	return 0
}

func (i *Index) UpdateLeader(leaders []string, epoch int) {
	stmt, _ := i.Conn.Prepare(`INSERT INTO Leader VALUES (?,?)`)
	for _, leader := range leaders {
		err := stmt.Exec(leader, epoch)
		if err != nil {
			indexLogger.Errorf("Cannot insert leader of epoch: %d", epoch)
		}
		stmt.Reset()
	}
	stmt.Close()
}

func (i *Index) Update(block *Block) {
	epoch := block.Index / int(TopK)
	if block.Index%int(TopK) == 0 {
		i.UpdateLeader(block.BPs, epoch)
	}
	countTx := make(map[string]int)
	for _, tx := range block.Txs {
		if tx.Type == 0 {
			addr := tx.Sender
			if countTx[addr] != 0 {
				countTx[addr] += 1
			} else {
				countTx[addr] = 1
			}
		}
		//update vote and stake
		i.ExecuteTransaction(&tx, epoch)
	}
	//update count tx
	for addr, count := range countTx {
		currCount := i.GetTxCountForAddr(addr, epoch)

		//insert count
		if currCount == 0 {
			stmt, _ := i.Conn.Prepare("INSERT INTO CountTX VALUES (?,?,?)")
			err := stmt.Exec(addr, count, epoch)
			if err != nil {
				indexLogger.Error("Cannot insert tx count")
			}
			stmt.Close()
			//update tx count for address
		} else {
			stmt, _ := i.Conn.Prepare("UPDATE CountTX SET amount=? WHERE address=? AND epoch=?")
			err := stmt.Exec(currCount+count, addr, epoch)
			if err != nil {
				indexLogger.Error("Cannot insert tx count")
			}
			stmt.Close()
		}
	}

}

//ExecuteTransaction: update database with information in block
//Information include:
//- TX count from given address on given epoch
//- Current stakes of all participants
//- Current votes
//- Leaders of this epoch
func (i *Index) ExecuteTransaction(tx *Transaction, epoch int) bool {
	switch tx.Type {
	//staking
	case 1:
		//get address
		addr := tx.Sender
		amount := tx.StakeAmount
		balance := i.GetBalance(addr)
		if balance < amount {
			return false
		}
		//update balance
		stmt, _ := i.Conn.Prepare("UPDATE Balance SET amount=? WHERE address=?")
		_ = stmt.Exec(balance-amount, addr)
		_ = stmt.Reset()

		//update stake
		stake := i.GetStake(addr)
		//not yet staked
		if stake == -1 {
			stmt, _ = i.Conn.Prepare("INSERT INTO Stake VALUES (?,?)")
			_ = stmt.Exec(addr, amount)
		} else {
			stmt, _ = i.Conn.Prepare("UPDATE Stake SET amount=? WHERE address=?")
			_ = stmt.Exec(stake+amount, addr)
		}
	case 2:
		voter := tx.Sender
		//delete old vote
		stmt, _ := i.Conn.Prepare("DELETE FROM Vote WHERE voter=?")
		_ = stmt.Exec(voter)
		_ = stmt.Reset()

		//add new vote
		stake := i.GetStake(voter)
		//log.SetFormatter(&log.TextFormatter{ForceColors: true})
		// indexLogger.Infof("%s stake %d", voter, stake)
		if stake < 0 {
			return false
		}
		for _, candidate := range tx.Candidate {
			//stmt, _ = i.Conn.Prepare("DELETE FROM Vote WHERE address=? AND  voter=?")
			//_ = stmt.Exec(candidate, voter)
			//_ = stmt.Reset()
			stmt, _ = i.Conn.Prepare("INSERT INTO Vote VALUES (?,?,?)")
			// indexLogger.Infof("%s receive %d from %s", candidate, stake, voter)
			_ = stmt.Exec(candidate, stake, voter)
			_ = stmt.Reset()

		}
	}
	return true
}

func (i *Index) Init() {
	i.Conn = InitSqlite()
	b, genesis := GetGenesis()
	if !b {
		indexLogger.Error("genesis not exist")
	}

	for k, v := range genesis.Balance {
		stmt, _ := i.Conn.Prepare(`INSERT INTO Balance VALUES (?, ?)`)
		err := stmt.Exec(k, v)
		stmt.Reset()
		if err != nil {
			indexLogger.Error("cannot execute statement")
		}
	}
}

func (i *Index) GetBalance(addr string) int {
	stmt, _ := i.Conn.Prepare("SELECT amount FROM Balance where address=?")
	_ = stmt.Exec(addr)
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			break
		}
		if !hasRow {
			// The query is finished
			break
		}

		// Use Scan to access column data from a row
		var amount int
		err = stmt.Scan(&amount)
		if err != nil {
			break
		}
		return amount
	}
	return -1
}

func (i *Index) GetStake(addr string) int {
	stmt, _ := i.Conn.Prepare("SELECT amount FROM Stake where address=?")
	_ = stmt.Exec(addr)
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			break
		}
		if !hasRow {
			// The query is finished
			break
		}

		// Use Scan to access column data from a row
		var amount int
		err = stmt.Scan(&amount)
		if err != nil {
			break
		}
		return amount
	}
	return -1
}

func (i *Index) GetTopKVote(k int) []string {
	sql := fmt.Sprintf("select address,sum(amount) as total from Vote group by address order by total desc limit %d", k)
	stmt, _ := i.Conn.Prepare(sql)
	_ = stmt.Exec()
	topK := []string{}
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			break
		}
		if !hasRow {
			// The query is finished
			break
		}

		// Use Scan to access column data from a row
		var address string
		var total int
		err = stmt.Scan(&address, &total)
		topK = append(topK, address)
		if err != nil {
			break
		}
	}
	//not enough vote, use initial setup
	if len(topK) < k {
		return GetInitialBPs()
	}
	return topK
}

func (i *Index) GetTotalVote(addr string) int {
	stmt, _ := i.Conn.Prepare("SELECT amount FROM Vote where address=?")
	_ = stmt.Exec(addr)
	totalVote := 0
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			break
		}
		if !hasRow {
			// The query is finished
			break
		}

		// Use Scan to access column data from a row
		var i int
		err = stmt.Scan(&i)
		if err != nil {
			break
		}
		totalVote += i
	}
	return totalVote
}

func (i *Index) GetNumberOfVote(addr string) int {
	stmt, _ := i.Conn.Prepare(`SELECT count(*) FROM Vote WHERE address=?`)
	err := stmt.Exec(addr)
	if err != nil {
		indexLogger.Error(err)
	}
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			break
		}
		if !hasRow {
			// The query is finished
			break
		}

		// Use Scan to access column data from a row
		var total int
		err = stmt.Scan(&total)
		if err != nil {
			break
		}
		return total
	}
	return 0
}

func (i *Index) GetTXLastEpoch(addr string) int {
	stmt, _ := i.Conn.Prepare(`select amount from CountTX where epoch=(select max(epoch) from CountTX) AND address=?`)
	err := stmt.Exec(addr)
	if err != nil {
		indexLogger.Error(err)
	}
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			break
		}
		if !hasRow {
			// The query is finished
			break
		}

		// Use Scan to access column data from a row
		var total int
		err = stmt.Scan(&total)
		if err != nil {
			break
		}
		return total
	}
	return 0

}

func (i *Index) GetLeaderCount(addr string) int {
	stmt, _ := i.Conn.Prepare(`SELECT count(*) FROM Leader WHERE address=?`)
	err := stmt.Exec(addr)
	if err != nil {
		indexLogger.Error(err)
	}
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			break
		}
		if !hasRow {
			// The query is finished
			break
		}

		// Use Scan to access column data from a row
		var total int
		err = stmt.Scan(&total)
		if err != nil {
			break
		}
		return total
	}
	return 0

}

func (i *Index) GetTotalReward(addr string) int {
	stmt, _ := i.Conn.Prepare(`SELECT count(amount) FROM CountTX WHERE address=?`)
	err := stmt.Exec(addr)
	if err != nil {
		indexLogger.Error(err)
	}
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			break
		}
		if !hasRow {
			// The query is finished
			break
		}

		// Use Scan to access column data from a row
		var total int
		err = stmt.Scan(&total)
		if err != nil {
			break
		}
		return total
	}
	return 0

}

func (i *Index) GetTopNCandidate() []string {
	sql := fmt.Sprintf("select address from CountTX where epoch= (select max(epoch) from CountTX) order by amount asc limit %d", nCandidate)
	stmt, _ := i.Conn.Prepare(sql)
	_ = stmt.Exec()
	var topK []string
	for {
		hasRow, err := stmt.Step()
		if err != nil {
			break
		}
		if !hasRow {
			// The query is finished
			break
		}

		// Use Scan to access column data from a row
		var address string
		err = stmt.Scan(&address)
		topK = append(topK, address)
		if err != nil {
			break
		}
	}
	//not enough vote, use initial setup
	if len(topK) < nCandidate {
		return []string{}
	}
	return topK

}

func (i *Index) GetTopKBP(k int) []string {
	//get top n candidate
	candidates := i.GetTopNCandidate()
	if len(candidates) < nCandidate {
		return GetInitialBPs()
	}
	var bps []Candidate
	count := 0
	for _, addr := range candidates {
		indexLogger.Infof("Calc %s", addr)
		totalVote := i.GetTotalVote(addr)
		numberVote := i.GetNumberOfVote(addr)
		txInLastEpoch := i.GetTXLastEpoch(addr)
		leaderCount := i.GetLeaderCount(addr)
		totalReward := i.GetTotalReward(addr)
		candidate := Candidate{
			Address:          addr,
			TotalVote:        totalVote,
			NumberOfVote:     numberVote,
			TotalTxLastEpoch: txInLastEpoch,
			LeaderCount:      leaderCount,
			TotalReward:      totalReward,
		}
		bps = append(bps, candidate)
		count++
		if count == 3 {
			break
		}
	}
	sort.Sort(ByScore(bps))
	var adds []string
	for _, c := range bps {
		adds = append(adds, c.Address)
	}
	return adds
}

func GetOrInitIndex() *Index {
	if idx == nil {
		//conn, err := sqlite3.Open(sqliteFile)
		idx = &Index{}
		idx.Init()
		return idx
	}
	return idx
}
