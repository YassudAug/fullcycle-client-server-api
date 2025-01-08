package main

import (
	"github.com/YassudAug/fullcycle-client-server-api/client"
	"github.com/YassudAug/fullcycle-client-server-api/server"
)

func main() {
	server.Handler()
	client.RequestDollarPriceBRL()
}
