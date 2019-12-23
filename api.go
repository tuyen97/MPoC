package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ipfs/go-log"
	"net/http"
	"time"
)

var api_logger = log.Logger("api")

type Api struct {
	memTxChan chan *Transaction
}

type TXRequest struct {
	Data string
}

type StakeRequest struct {
	Sender string
	Amount int
}

type VoteRequest struct {
	Voter      string
	Candidates []string
}

func (a *Api) GetLastBlock(w http.ResponseWriter, r *http.Request) {
	b, err := GetLastBlock()
	if err != nil {
		fmt.Fprintf(w, "cannot get last block")
	} else {
		fmt.Fprint(w, b)
	}

}

func (a *Api) VoteFunc(w http.ResponseWriter, r *http.Request) {
	var v VoteRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&v)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Cannot decode body")
	}
	api_logger.Info(v.Candidates)
	tx := Transaction{
		ID:          nil,
		Sender:      v.Voter,
		Type:        2,
		StakeAmount: 0,
		Candidate:   v.Candidates,
		Data:        "txr.Data",
		Timestamp:   time.Now().UnixNano(),
	}
	tx.SetId()
	a.memTxChan <- &tx
	_, _ = fmt.Fprintf(w, "success")

}

func (a *Api) StakeFunc(w http.ResponseWriter, r *http.Request) {
	var s StakeRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&s)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Cannot decode body")
	}
	tx := Transaction{
		ID:          nil,
		Sender:      s.Sender,
		Type:        1,
		StakeAmount: s.Amount,
		Candidate:   nil,
		Data:        "txr.Data",
		Timestamp:   time.Now().UnixNano(),
	}
	tx.SetId()
	a.memTxChan <- &tx
	_, _ = fmt.Fprintf(w, "success")
}

func (a *Api) IndexFunc(w http.ResponseWriter, r *http.Request) {
	//var txr TXRequest
	//decoder := json.NewDecoder(r.Body)
	//err := decoder.Decode(&txr)
	//if err != nil {
	//	logger.Errorf("Cannot decode body")
	//}
	tx := Transaction{
		ID:          nil,
		Sender:      "",
		Type:        0,
		StakeAmount: 0,
		Candidate:   nil,
		Data:        "txr.Data",
		Timestamp:   time.Now().UnixNano(),
	}
	tx.SetId()
	a.memTxChan <- &tx
	_, _ = fmt.Fprintf(w, "success")
}

func (api *Api) Start(port string) {
	//log.SetAllLoggers(logging.INFO)
	log.SetLogLevel("api", "info")
	router := mux.NewRouter()
	router.HandleFunc("/", api.IndexFunc).Methods("GET")
	router.HandleFunc("/stake", api.StakeFunc).Methods("POST")
	router.HandleFunc("/vote", api.VoteFunc).Methods("POST")
	router.HandleFunc("/lastblock", api.GetLastBlock).Methods("GET")
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), router)
	api_logger.Info("Server started")
}
