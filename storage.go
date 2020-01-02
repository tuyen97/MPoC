package main

import (
	"github.com/bvinc/go-sqlite-lite/sqlite3"
	log "github.com/sirupsen/logrus"
	"os"
)

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
	// if !SqliteExists() {
	conn, err := sqlite3.Open(sqliteFile)
	if err != nil {
		log.Error(err)
	}
	stmt, _ := conn.Prepare(`CREATE TABLE Balance(address TEXT PRIMARY KEY, amount INTEGER)`)
	err = stmt.Exec()
	if err != nil {
		log.Error("Cannot execute query")
	}
	_ = stmt.Reset()
	stmt, _ = conn.Prepare(`CREATE TABLE Stake(address TEXT PRIMARY KEY, amount INTEGER)`)
	_ = stmt.Exec()

	_ = stmt.Reset()
	stmt, _ = conn.Prepare(`CREATE TABLE Vote(address TEXT , amount INTEGER, voter TEXT, CONSTRAINT vote_pk PRIMARY KEY (address, voter))`)
	_ = stmt.Exec()

	_ = stmt.Reset()
	stmt, _ = conn.Prepare(`CREATE TABLE CountTX(address TEXT , amount INTEGER, epoch INTEGER, CONSTRAINT vote_pk PRIMARY KEY (address, epoch))`)
	_ = stmt.Exec()

	_ = stmt.Reset()
	stmt, _ = conn.Prepare(`CREATE TABLE Leader(address TEXT, epoch INTEGER, CONSTRAINT vote_pk PRIMARY KEY (address, epoch))`)
	_ = stmt.Exec()
	// conn.Close()

	// }
	// conn, _ := sqlite3.Open(sqliteFile)
	return conn
}
