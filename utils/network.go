package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/libp2p/go-reuseport"
)

var peerConnects = make(chan *Peer, 3)     // Peer connection management channel.
var peerDisconnects = make(chan *Peer, 3)  // Peer disconnection management channel.
var managementShutdown = make(chan string) // channel to signal management cleanup
var connections []*Peer
var mediaList = make(chan []string, 3)

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
	BuffWriteToNetwork(connection, connections[0].publicAddress+"\n")
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
			BuffWriteToNetwork(connection, connections[len(connections)-1].publicAddress+"\n")
			BuffReadFromNetwork(connection)
		}
		if message == "clustering" {
			BuffWriteToNetwork(connection, fmt.Sprint(len(connections))+"\n")
			keepAlive = false
		}
		if message == "search" {
			ClusterSearch(connection)
		}
		if message == "listing" {
			BuffWriteToNetwork(connection, "ok\n")
			v := strings.Trim(BuffReadFromNetwork(connection), "\n")
			nr, err := strconv.Atoi(v)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			// Read all nr lines and put into array.
			var media []string
			for i := 0; i < nr; i++ {
				val := strings.Trim(BuffReadFromNetwork(connection), "\n")
				media = append(media, val)
				fmt.Println(media)
				BuffWriteToNetwork(connection, "ok\n")
			}
			fmt.Println("Sending into channel")
			mediaList <- media
		}
		if message == "find" {
			PeerManagement()
		}
	}
	Close(connection)
}

func InList(element string, list []string) bool {
	if len(list) == 0 {
		return false
	}
	for i := 0; i < len(list); i++ {
		if list[i] == element {
			return true
		}
	}
	return false
}

func ClusterSearch(connection net.Conn) {
	fmt.Println("Cluster searching...")
	// Send "list" to all peers which are not current connection.
	var response []string
	for i := 0; i < len(connections); i++ {
		if connections[i].connection != connection {
			BuffWriteToNetwork(connections[i].connection, "list\n")
			fmt.Println("Test")
			// Wait on channel.
			select {
			case media := <-mediaList:
				for i := 0; i < len(media); i++ {
					if InList(media[i], response) == false {
						response = append(response, media[i])
					}
				}
				break
			}
		}
	}
	BuffWriteToNetwork(connection, "listing\n")
	BuffReadFromNetwork(connection)
	BuffWriteToNetwork(connection, fmt.Sprint(len(response))+"\n")
	BuffReadFromNetwork(connection)
	fmt.Println("sending")
	for i := 0; i < len(response); i++ {
		BuffWriteToNetwork(connection, response[i]+"\n")
		BuffReadFromNetwork(connection)
	}
	fmt.Println("out")
}

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
			fmt.Println("test")
			BuffWriteToNetwork(connections[len(connections)-1].connection, "found\n")

			BuffWriteToNetwork(connections[0].connection, "pair\n")
			break
		}
	}
}
func Serve(listener net.Listener) {
	// opens up routine that deals with concurrent memory management.
	go MemoryManagementRoutine()

	// go PeerManagement()

	for {
		connection := Accept(listener)

		go HandleConnection(connection)
	}
}
