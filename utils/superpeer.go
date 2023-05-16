package utils

import "fmt"

func OpenSuperPeer() {
	listener := Listen(ProtocolTcp, ServerHost+SuperPeerPort)

	fmt.Println("SuperPeer Opened at...", ServerHost+SuperPeerPort)

	Serve(listener)
}
