package utils

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func Listen(protocol, address string) net.Listener {
	listener, err := net.Listen(protocol, address)

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
		os.Exit(1)
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
		os.Exit(1)
	}
}

func BuffReadFromNetwork(connection net.Conn) string {
	data, err := bufio.NewReader(connection).ReadString('\n')

	if err != nil {
		fmt.Println("Error when reading from network: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Received: ", data)
	return data
}

func Flush(writer *bufio.Writer) {
	err := writer.Flush()

	if err != nil {
		fmt.Println("Error when flushing data to network: ", err.Error())
		os.Exit(1)
	}
}

func BuffWriteToNetwork(connection net.Conn, message string) {
	writer := bufio.NewWriter(connection)

	_, err := writer.WriteString(message)
	if err != nil {
		fmt.Println("Error when writing to network: ", err.Error())
		os.Exit(1)
	}
	Flush(writer)
}

func HandleConnection(connection net.Conn) {
	message := BuffReadFromNetwork(connection)

	BuffWriteToNetwork(connection, message)

	Close(connection)
}

func Serve(listener net.Listener) {
	for {
		connection := Accept(listener)

		go HandleConnection(connection)
	}
}
