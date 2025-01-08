package main

import (
	"time"

	"github.com/YassudAug/fullcycle-client-server-api/client"
	"github.com/YassudAug/fullcycle-client-server-api/server"
)

func main() {
	go server.Handler()

	time.Sleep(1 * time.Second)
	client.RequestDollarPriceBRL()
}
