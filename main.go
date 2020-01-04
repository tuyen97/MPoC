package main

import "os"
import "time"

//<api_port> <peer_port> <address>
func main() {
	flagDns := os.Args[1]
	if flagDns == "1" {
		s := Server{}
		s.Start()
	} else {
		time.Sleep(3 * time.Second)
		node := Node{}
		node.Init(os.Args[4])
		node.Start(os.Args[2], os.Args[3])
	}
	select {}
}
