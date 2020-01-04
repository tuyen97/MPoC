package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Register struct {
	Addr string
}
type Peers struct {
	Infos []string
}

func (p *Peers) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}

var PeerInfos []string

type Server struct {
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var re Register
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&re)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Cannot decode body")
	}
	if !SliceExists(PeerInfos, re.Addr) {
		PeerInfos = append(PeerInfos, re.Addr)
	}
	_, _ = fmt.Fprintf(w, "Success")
}

func handleQuery(w http.ResponseWriter, r *http.Request) {
	res := Peers{Infos: PeerInfos}
	fmt.Fprintf(w, res.String())
}

func (s *Server) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/register", handleRegister).Methods("POST")
	router.HandleFunc("/query", handleQuery)
	go http.ListenAndServe("0.0.0.0:13000", router)
	fmt.Println("server started")
}
