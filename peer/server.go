package main

import (
	"Peer-to-peer-on-demand-streaming/utils"
	"fmt"
)

func ManageServer() {
	listener := utils.Listen(utils.ProtocolTcp, utils.ServerHost+utils.Port)

	fmt.Println("Listening on:", utils.ServerHost+utils.Port)

	utils.Serve(listener)
}

func OpenServer() {
	go ManageServer()
}
