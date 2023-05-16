package utils

import (
	"fmt"
	"net"
	"strings"
)

var networkName string

func SendConnectionInfo(connection net.Conn) {
	BuffWriteToNetwork(connection, connection.LocalAddr().String()+"\n")
	networkName = strings.Trim(BuffReadFromNetwork(connection), "\n")
	fmt.Println("Connected to cluster with name:", networkName)
}

func OpenPeer() {
	connection := Dial(ProtocolTcp, ServerHost+SuperPeerPort)

	fmt.Println("Dialing the superpeer...")
	BuffWriteToNetwork(connection, "connect\n")

	keepAlive := true
	for keepAlive {
		message := BuffReadFromNetwork(connection)
		switch strings.Trim(message, "\n") {
		case "connect":
			fmt.Println("Connecting to cluster...")
			SendConnectionInfo(connection)
		default:
			keepAlive = false
		}
	}
	Close(connection)
}
