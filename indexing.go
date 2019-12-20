package main

import (
	"fmt"
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	"github.com/ipfs/go-log"
)

var idx *Index = nil

type Balance struct {
	Address string
	Amount  int
}

type Stake struct {
	Address string
	Amount  int
}

type TotalVote struct {
	Address string
	Amount  int
}

type Index struct {
	Conn *sqlite3.Conn
}

func ReIndex() {

}

var indexLogger = log.Logger("bf")

func (i *Index) Update(tx *Transaction) bool {
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
		indexLogger.Infof("%s stake %d", voter, stake)
		if stake < 0 {
			return false
		}
		for _, candidate := range tx.Candidate {
			stmt, _ = i.Conn.Prepare("DELETE FROM Vote WHERE address=? AND  voter=?")
			_ = stmt.Exec(candidate, voter)
			_ = stmt.Reset()
			stmt, _ = i.Conn.Prepare("INSERT INTO Vote VALUES (?,?,?)")
			indexLogger.Infof("%s receive %d from %s", candidate, stake, voter)
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

func GetOrInitIndex() *Index {
	if idx == nil {
		//conn, err := sqlite3.Open(sqliteFile)
		idx = &Index{}
		idx.Init()
		return idx
	}
	return idx
}
