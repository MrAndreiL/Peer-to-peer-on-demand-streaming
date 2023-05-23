package utils

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
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
	var swarm []net.Conn
	// Announce the superpeer to start sending.
	BuffWriteToNetwork(connection, fmt.Sprint(p)+"\n")
	for i := 0; i < p; i++ {
		// Read public address.
		address := strings.Trim(BuffReadFromNetwork(connection), "\n")
		fmt.Println(address)
		BuffWriteToNetwork(connection, "ok\n")
		fmt.Println("Starting hole punching from: ", connection.LocalAddr().String())
		HolePunching(address, listener, connection)
		// Upon successful hole punching, save the connection in a list.
		punchConn := punchConnection
		swarm = append(swarm, punchConn)
	}
	// After hole punching, start swarm protocol.
	Swarm(swarm, "lane.mp4", "video")
}

func Swarm(swarm []net.Conn, fileName, mediaType string) {
	// 1. Ask first peer in connection for the total
	// number of files to be transferred.
	BuffWriteToNetwork(swarm[0], "number\n")
	BuffReadFromNetwork(swarm[0])
	BuffWriteToNetwork(swarm[0], fileName+"\n")
	fileLength := strings.Trim(BuffReadFromNetwork(swarm[0]), "\n")
	fmt.Println(fileLength)
	// Create directory in appropiate place.
	pathToNewDir := NewMediaDirectory(fileName, mediaType)
	fmt.Println(pathToNewDir)
	// 2. Initiate full transfer.
	for i := 0; i < len(swarm); i++ {
		go Transfer(swarm[i], pathToNewDir, mediaType, fileName, fileLength, i)
	}
	Streamfy()
	// 3. Upon successful swarming, open HLS server.
	playlistFile := strings.Trim(fileName, filepath.Ext(fileName)) + ".m3u8"
	StartStream(pathToNewDir, playlistFile)
}

func Transfer(connection net.Conn, path, mediaType, fileName, length string, index int) {
	// 1. Send "transfer" initiation.
	BuffWriteToNetwork(connection, "transfer\n")
	BuffReadFromNetwork(connection)
	BuffWriteToNetwork(connection, fileName+"\n")
	BuffReadFromNetwork(connection)
	// 2. Send mode type for transfer (forwards/backwards/both).
	if index == 0 {
		// first peer forwards = full transfer.
		BuffWriteToNetwork(connection, "forward\n")
		BuffReadFromNetwork(connection)
		// send start position.
		BuffWriteToNetwork(connection, "1\n")
	} else if index == 1 {
		// seconds peer backwards = full transfer.
		BuffWriteToNetwork(connection, "backwards\n")
		BuffReadFromNetwork(connection)
		// send start position.
		BuffWriteToNetwork(connection, length+"\n")
	} else {
		// All other peers receive both mode.
		BuffWriteToNetwork(connection, "both\n")
		BuffReadFromNetwork(connection)
		// send start position.
		l, err := strconv.Atoi(length)
		if err != nil {
			fmt.Println("Error when converting: ", err.Error())
			os.Exit(1)
		}
		l = l + 1
		randomPosition := rand.Intn(l-2) + 2
		BuffWriteToNetwork(connection, fmt.Sprint(randomPosition)+"\n")
	}
	// 3. Start swarming.
	fmt.Println("Start swarming")
	l, err := strconv.Atoi(length)
	if err != nil {
		fmt.Println("Error when converting: ", err.Error())
		os.Exit(1)
	}
	ReceiveFiles(path, connection, l)
	fmt.Println("Swarming successful!")
	connection.Close()
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
}

func Swarming(connection net.Conn) {
	// Receive swarming request from peer.
	keepAlive := true
	for keepAlive {
		message := strings.Trim(BuffReadFromNetwork(connection), "\n")
		if message == "number" {
			BuffWriteToNetwork(connection, "ok\n")
			// Receive media file name to be transferred.
			fileName := strings.Trim(BuffReadFromNetwork(connection), "\n")
			// Send back the number of files required for full transfer.
			BuffWriteToNetwork(connection, MappingFileLength[fileName]+"\n")
		}
		if message == "transfer" {
			fmt.Println("Preparing to transfer HLS files")
			BuffWriteToNetwork(connection, "ok\n")
			// Recevie media file requested.
			fileName := strings.Trim(BuffReadFromNetwork(connection), "\n")
			fmt.Println(fileName)
			BuffWriteToNetwork(connection, "ok\n")
			playlist := FilesPathsListing(fileName)
			// Receive mode type for transfer (forward/backwards/both)
			mode := strings.Trim(BuffReadFromNetwork(connection), "\n")
			fmt.Println(mode)
			BuffWriteToNetwork(connection, "ok\n")
			position := strings.Trim(BuffReadFromNetwork(connection), "\n")
			pos, _ := strconv.Atoi(position)
			if mode == "forward" {
				ForwardSwarming(playlist, pos, connection)
			} else if mode == "backwards" {
				BackwardsSwarming(playlist, pos, connection)
			} else {
				BidirectionalSwarming(playlist, pos, connection)
			}
			connection.Close()
			keepAlive = false
		}
	}
}

func OpenPeer() {
	fmt.Println("Application starting...")
	fmt.Println("Listing all available superpeers...")
	fmt.Println("1)" + Server1Host + SuperPeerPort)
	fmt.Println("2)" + Server2Host + SuperPeerPort)
	address := Clustering()
	connection := ReuseDial(ProtocolTcp, ":", address)
	// Open for listening.
	listener := Listen(ProtocolTcp, connection.LocalAddr().String())
	keepAlive := true
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
			go Swarming(punchConnection)
		}
	}
	Close(connection)
}

func Clustering() string {
	// Dialing first cluster superpeer.

	fmt.Println("Dialing first superpeer...")
	connection := ReuseDial(ProtocolTcp, ":", Server1Host+SuperPeerPort)
	BuffWriteToNetwork(connection, "clustering\n")
	l := strings.Trim(BuffReadFromNetwork(connection), "\n")
	fmt.Println("First cluster has " + l + " connected peers.")
	connection.Close()
	load1, err := strconv.Atoi(l)
	if err != nil {
		fmt.Println("Failure when converting: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Dialing second superpeer...")
	connection = ReuseDial(ProtocolTcp, ":", Server2Host+SuperPeerPort)
	BuffWriteToNetwork(connection, "clustering\n")
	l = strings.Trim(BuffReadFromNetwork(connection), "\n")
	fmt.Println("Second cluster has " + l + " connected peers.")
	connection.Close()
	load2, err := strconv.Atoi(l)
	if err != nil {
		fmt.Println("Failure when converting: ", err.Error())
		os.Exit(1)
	}

	if load1 >= load2 {
		fmt.Println("Connecting to first superpeer cluster...")
		return (Server1Host + SuperPeerPort)
	} else {
		fmt.Println("Connecting to second superpeer cluster...")
		return (Server2Host + SuperPeerPort)
	}
}
