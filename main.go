package main

import (
	"os"
)

//<api_port> <peer_port> <address>
func main() {
	node := Node{}
	node.Init(os.Args[3])
	node.Start(os.Args[1], os.Args[2])
	select {}
}
