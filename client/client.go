package main

import (
	"Peer-to-peer-on-demand-streaming/utils"
)

func main() {
	connection := utils.Dial(utils.ProtocolTcp, utils.ServerHost+utils.Port)

	utils.BuffWriteToNetwork(connection, "Hello World\n")

	utils.BuffReadFromNetwork(connection)

	utils.Close(connection)
}
