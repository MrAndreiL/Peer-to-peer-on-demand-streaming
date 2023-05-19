package utils

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var networkName string
var punchConnection net.Conn = nil

func ConnectionLogin(connection net.Conn) {
	// Send private address.
	BuffWriteToNetwork(connection, connection.LocalAddr().String()+"\n")
	// Receive network name.
	networkName := strings.Trim(BuffReadFromNetwork(connection), "\n")
	fmt.Println("Connected with network name: ", networkName)
}

func AcceptIncomingConnection(listener net.Listener) {
	punchConnection = Accept(listener)
	fmt.Println("accepted: ", punchConnection.RemoteAddr().String())
	fmt.Println(punchConnection.LocalAddr().String())
}

func SearchForConnection(connection net.Conn, listener net.Listener) {
	// Receive the number of peers to connect.
	BuffWriteToNetwork(connection, "found\n")
	peers := strings.Trim(BuffReadFromNetwork(connection), "\n")
	fmt.Println("Found ", peers, " peers available to pair...")
	// Convert from string to int the number of available peers.
	p, err := strconv.Atoi(peers)
	if err != nil {
		fmt.Println("Cannot convert the received number of peers.")
		os.Exit(1)
	}
	// Announce the superpeer to start sending.
	BuffWriteToNetwork(connection, fmt.Sprint(p)+"\n")
	for i := 0; i < p; i++ {
		// Read public address.
		address := strings.Trim(BuffReadFromNetwork(connection), "\n")
		fmt.Println(address)
		BuffWriteToNetwork(connection, "ok\n")
		fmt.Println("Starting hole punching from: ", connection.LocalAddr().String())
		HolePunching(address, listener, connection)
	}
	fmt.Println(BuffReadFromNetwork(punchConnection))
}

func HolePunching(address string, listener net.Listener, connection net.Conn) {
	punchConnection = nil
	go AcceptIncomingConnection(listener)

	for punchConnection == nil {
		// Dial the other peer.
		fmt.Println("Dialing")
		connect := ReuseDial(ProtocolTcp, connection.LocalAddr().String(), address)
		if connect != nil {
			punchConnection = connect
		}
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Hole punched, connection established!")
	BuffWriteToNetwork(punchConnection, "mesaj\n")
}

func OpenPeer() {
	connection := Dial(ProtocolTcp, PeerHost+SuperPeerPort)
	// Open for listening.
	listener := Listen(ProtocolTcp, connection.LocalAddr().String())
	keepAlive := true
	fmt.Println("Dialing the superpeer...")
	BuffWriteToNetwork(connection, "connect\n")
	for keepAlive {
		message := strings.Trim(BuffReadFromNetwork(connection), "\n")
		if message == "connect" {
			fmt.Println("Connecting to cluster")
			ConnectionLogin(connection)
		}
		if message == "goodbye" {
			fmt.Println("Disconnecting from network...")
			keepAlive = false
		}
		if message == "found" {
			fmt.Println("pairing...")
			SearchForConnection(connection, listener)
		}
		if message == "pair" {
			fmt.Println("pairing...")
			BuffWriteToNetwork(connection, "pair\n")
			address := strings.Trim(BuffReadFromNetwork(connection), "\n")
			fmt.Println("A: ", address)
			BuffWriteToNetwork(connection, "ok\n")
			fmt.Println("Starting hole punching...")
			HolePunching(address, listener, connection)
			fmt.Println(BuffReadFromNetwork(punchConnection))
		}
	}
	Close(connection)
}
