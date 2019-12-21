package main

import (
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	log "github.com/sirupsen/logrus"
	"os"
)

const dbFile = "blockchain.db"
const blocksBucket = "blockchain"
const sqliteFile = "data.sqlite3"

func DBExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

func SqliteExists() bool {
	if _, err := os.Stat(sqliteFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func InitSqlite() *sqlite3.Conn {
	if !SqliteExists() {
		conn, _ := sqlite3.Open(sqliteFile)
		stmt, _ := conn.Prepare(`CREATE TABLE Balance(address TEXT PRIMARY KEY, amount INTEGER)`)
		err := stmt.Exec()
		if err != nil {
			log.Error("Cannot execute query")
		}
		//_ = stmt.Reset()
		stmt, _ = conn.Prepare(`CREATE TABLE Stake(address TEXT PRIMARY KEY, amount INTEGER)`)
		_ = stmt.Exec()

		//_ = stmt.Reset()
		stmt, _ = conn.Prepare(`CREATE TABLE Vote(address TEXT , amount INTEGER, voter TEXT, CONSTRAINT vote_pk PRIMARY KEY (address, voter))`)
		_ = stmt.Exec()
		conn.Close()
	}
	conn, _ := sqlite3.Open(sqliteFile)
	return conn
}
