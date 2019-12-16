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
	memTxChan chan *TransactionWrapper
}

type TXRequest struct {
	Data string
}

func (a *Api) IndexFunc(w http.ResponseWriter, r *http.Request) {
	var txr TXRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&txr)
	if err != nil {
		logger.Errorf("Cannot decode body")
	}
	tx := Transaction{
		ID:          nil,
		Signature:   nil,
		Sender:      nil,
		Type:        0,
		StakeAmount: 0,
		Candidate:   nil,
		Data:        txr.Data,
	}
	txw := TransactionWrapper{
		ID:        nil,
		TX:        &tx,
		Timestamp: time.Now().UnixNano(),
	}
	a.memTxChan <- &txw
	_, _ = fmt.Fprintf(w, "success")
}

func (api *Api) Start(port string) {
	//log.SetAllLoggers(logging.INFO)
	log.SetLogLevel("api", "info")
	router := mux.NewRouter()
	router.HandleFunc("/", api.IndexFunc).Methods("POST")
	go http.ListenAndServe(fmt.Sprintf("127.0.0.1:%s", port), router)
	api_logger.Info("Server started")
}
