package utils

import "fmt"

func OpenSuperPeer() {
	listener := Listen(ProtocolTcp, Server1Host+SuperPeerPort)

	fmt.Println("SuperPeer Opened at...", Server1Host+SuperPeerPort)

	Serve(listener)
}
