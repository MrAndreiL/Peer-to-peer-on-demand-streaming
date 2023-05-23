package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/libp2p/go-reuseport"
)

var peerConnects = make(chan *Peer, 3)     // Peer connection management channel.
var peerDisconnects = make(chan *Peer, 3)  // Peer disconnection management channel.
var managementShutdown = make(chan string) // channel to signal management cleanup
var connections []*Peer

func Listen(protocol, address string) net.Listener {
	listener, err := reuseport.Listen(protocol, address)

	if err != nil {
		fmt.Println("Error when creating listener: ", err.Error())
		os.Exit(1)
	}
	return listener
}

func Dial(protocol, address string) net.Conn {
	conn, err := net.Dial(protocol, address)

	if err != nil {
		fmt.Println("Error when dialing: ", err.Error())
		return nil
	}
	return conn
}

func ReuseDial(protocol, fromAddr, toAddr string) net.Conn {
	conn, err := reuseport.Dial(protocol, fromAddr, toAddr)

	if err != nil {
		fmt.Println("Error when dialing: ", err.Error())
		return nil
	}
	return conn
}

func Accept(listener net.Listener) net.Conn {
	conn, err := listener.Accept()

	if err != nil {
		fmt.Println("Error when accepting request: ", err.Error())
		os.Exit(1)
	}
	return conn
}

func Close(connection net.Conn) {
	err := connection.Close()

	if err != nil {
		fmt.Println("Error when closing connection: ", err.Error())
	}
}

func BuffReadFromNetwork(connection net.Conn) string {
	data, err := bufio.NewReader(connection).ReadBytes('\n')

	if err != nil {
		fmt.Println("Error when reading from network: ", err.Error())
	}
	return string(data)
}

func Flush(writer *bufio.Writer) {
	err := writer.Flush()

	if err != nil {
		fmt.Println("Error when flushing data to network: ", err.Error())
	}
}

func BuffWriteToNetwork(connection net.Conn, message string) {
	writer := bufio.NewWriter(connection)

	_, err := writer.WriteString(message)
	if err != nil {
		fmt.Println("Error when writing to network: ", err.Error())
	}
	Flush(writer)
}

func ManageConnection(connection net.Conn) *Peer {
	// Send connect command.
	BuffWriteToNetwork(connection, "connect\n")
	// Receive private address.
	private := strings.Trim(BuffReadFromNetwork(connection), "\n")
	// Extract public address from connection struct.
	public := connection.RemoteAddr().String()
	// Generate network name for the waiting peer.
	networkName := StringWithCharset(IdLength, Charset)
	// Send network name to peer.
	BuffWriteToNetwork(connection, networkName+"\n")
	// Save peer in memory.
	return CreateNetworkPeer(public, private, networkName, connection)
}

var test = 1

func SendAllPeersToOne(connection net.Conn) {
	// Write the expected number of available peers.
	BuffWriteToNetwork(connection, "1\n")
	BuffReadFromNetwork(connection)
	// Send peer info.
	BuffWriteToNetwork(connection, connections[1].publicAddress+"\n")
	fmt.Println(BuffReadFromNetwork(connection))
}

func HandleConnection(connection net.Conn) {
	keepAlive := true
	var peer *Peer
	for keepAlive {
		message := strings.Trim(BuffReadFromNetwork(connection), "\n")
		if message == "connect" {
			fmt.Println("New peer connected!")
			peer = ManageConnection(connection)
			peerConnects <- peer
		}
		if message == "goodbye" {
			fmt.Println("Peer disconnected!")
			peerDisconnects <- peer
			BuffWriteToNetwork(connection, "goodbye\n")
			keepAlive = false
		}
		if message == "found" {
			SendAllPeersToOne(connection)
		}
		if message == "pair" {
			BuffWriteToNetwork(connection, connections[0].publicAddress+"\n")
			BuffReadFromNetwork(connection)
		}
	}
	Close(connection)
}

/*
func OpenServer() {
	// Open Audio streaming server.
	go ServeVideoAudio()
}
*/
/*
	func ServeVideoAudio() {
		AudioStream, err := filepath.Abs(AudioPathStream)
		if err != nil {
			fmt.Println("Error when getting abs path to stream directory: ", err.Error())
			os.Exit(1)
		}
		audioServer := http.NewServeMux()

		VideoStream, err := filepath.Abs(VideoPathStream)
		if err != nil {
			fmt.Println("Error when getting abs path to stream directory: ", err.Error())
			os.Exit(1)
		}
		videoServer := http.NewServeMux()
		videoServer.Handle("/", addHeaders(http.FileServer(http.Dir(VideoStream))))
		go func() {
			http.ListenAndServe(ServerHost+AudioPort, audioServer)
		}()
		go func() {
			http.ListenAndServe(ServerHost+VideoPort, videoServer)
		}()
	}

// add CORS support.

	func addHeaders(h http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			h.ServeHTTP(w, r)
		}
	}
*/
func FindPeerIndex(peer *Peer, connections []*Peer) int {
	for i, p := range connections {
		if p.networkName == peer.networkName {
			return i
		}
	}
	return -1
}

func RemovePeer(peer *Peer, connections []*Peer) []*Peer {
	index := FindPeerIndex(peer, connections)
	if index != -1 {
		ret := make([]*Peer, 0)
		ret = append(ret, connections[:index]...)
		return append(ret, connections[index+1:]...)
	}
	return connections
}

func MemoryManagementRoutine() {
	for {
		select {
		case peer := <-peerConnects:
			// Add to connections data structure.
			connections = append(connections, peer)
		case peer := <-peerDisconnects:
			// Delete a peer from data structure.
			connections = RemovePeer(peer, connections)
		case <-managementShutdown:
			// Close down management thread
			return
		}
	}
}

func PeerManagement() {
	for {
		if len(connections) >= 2 {
			// To A. -> TODO
			BuffWriteToNetwork(connections[0].connection, "found\n")

			BuffWriteToNetwork(connections[1].connection, "pair\n")
			break
		}
	}
}

func Serve(listener net.Listener) {
	// opens up routine that deals with concurrent memory management.
	go MemoryManagementRoutine()

	go PeerManagement()

	for {
		connection := Accept(listener)

		go HandleConnection(connection)
	}
}
