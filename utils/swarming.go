package utils

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func FilesPathsListing(fileName string) []string {
	dirPath := MappingFilePath[fileName]
	var paths []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}
		paths = append(paths, path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return paths
}

func NewMediaDirectory(name, mediaType string) string {
	name = strings.TrimSuffix(name, filepath.Ext(name))
	path := ""
	if mediaType == "audio" {
		path = AudioPathStream + name
	} else {
		path = VideoPathStream + name
	}
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	path, err := filepath.Abs(path + "/")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return (path + "/")
}

func fillString(returnString string, toLength int) string {
	for {
		lengthString := len(returnString)
		if lengthString < toLength {
			returnString = returnString + ":"
			continue
		}
		break
	}
	return returnString
}

func SendFiles(connection net.Conn, playlist []string) {
	for i := 0; i < len(playlist); i++ {
		file, err := os.Open(playlist[i])
		if err != nil {
			fmt.Println("Error when opening file: ", err.Error())
			os.Exit(1)
		}
		fileInfo, err := file.Stat()
		if err != nil {
			fmt.Println("Error when extracting file info: ", err.Error())
			os.Exit(1)
		}
		fileSize := fillString(fmt.Sprint(fileInfo.Size()), 10)
		fileName := fillString(fileInfo.Name(), 64)
		fmt.Println("Sending filename and filesize!")
		connection.Write([]byte(fileSize))
		connection.Write([]byte(fileName))
		file.Close()
	}
	for i := 0; i < len(playlist); i++ {
		file, err := os.Open(playlist[i])
		if err != nil {
			fmt.Println("Error when opening file: ", err.Error())
			os.Exit(1)
		}
		sendBuffer := make([]byte, BufferSize)
		fmt.Println("Start sending file!")
		for {
			_, err := file.Read(sendBuffer)
			if err == io.EOF {
				break
			}
			connection.Write(sendBuffer)
		}
		fmt.Println("File has been sent!")
		file.Close()
	}
	fmt.Println("Sending done!")
}

func ReceiveFiles(path string, connection net.Conn, length int) {
	var files []string
	var lengths []int64
	for i := 0; i < length; i++ {
		bufferFileName := make([]byte, 64)
		bufferFileSize := make([]byte, 10)

		connection.Read(bufferFileSize)
		fileSize, err := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		connection.Read(bufferFileName)
		fileName := strings.Trim(string(bufferFileName), ":")
		fmt.Println(fileName, fileSize)
		files = append(files, fileName)
		lengths = append(lengths, fileSize)
	}
	for i := 0; i < length; i++ {
		newFile, err := os.Create(path + files[i])
		if err != nil {
			fmt.Println("Error when creating new file: ", err.Error())
			os.Exit(1)
		}
		defer newFile.Close()
		var receivedBytes int64 = 0
		for {
			if (lengths[i] - receivedBytes) < BufferSize {
				io.CopyN(newFile, connection, (lengths[i] - receivedBytes))
				connection.Read(make([]byte, (receivedBytes+BufferSize)-lengths[i]))
				break
			}
			io.CopyN(newFile, connection, BufferSize)
			receivedBytes += BufferSize
		}
		fmt.Println("File received successfully!")
	}
	fmt.Println("Swarming accomplished!")
}

func ForwardSwarming(playlist []string, position int, connection net.Conn) {
	var lista []string
	for i := position; i < len(playlist); i++ {
		lista = append(lista, playlist[i])
	}
	SendFiles(connection, lista)
}

func BackwardsSwarming(playlist []string, position int, connection net.Conn) {
	var lista []string
	for i := position; i >= 1; i-- {
		lista = append(lista, playlist[i])
	}
	SendFiles(connection, lista)
}

func BidirectionalSwarming(playlist []string, position int, connection net.Conn) {
	var lista []string
	for i := position; i < len(playlist); i++ {
		lista = append(lista, playlist[i])
	}
	for i := position - 1; i >= 1; i-- {
		lista = append(lista, playlist[i])
	}
	SendFiles(connection, lista)
}
