package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func RegisterAddr(addr string) error {
	req, err := json.Marshal(map[string]string{
		"Addr": addr,
	})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s:13000/register", dnsServer), "application/json", bytes.NewBuffer(req))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return nil
}
func QueryDns() Peers {
	resp, err := http.Get(fmt.Sprintf("http://%s:13000/query", dnsServer))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var p Peers
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&p)
	return p
}
