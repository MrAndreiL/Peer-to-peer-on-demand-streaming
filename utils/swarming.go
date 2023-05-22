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

func SendFile(filePath string, connection net.Conn) {
	// 1. Open file to be sent.
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error when opening file: ", err.Error())
		os.Exit(1)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error when extracting file info: ", err.Error())
		os.Exit(1)
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Println("Sending filename and filesize!")
	connection.Write([]byte(fileSize))
	connection.Write([]byte(fileName))
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
}

func ReceiveFile(path string, connection net.Conn) string {
	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)

	connection.Read(bufferFileSize)
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)

	connection.Read(bufferFileName)
	fileName := strings.Trim(string(bufferFileName), ":")

	newFile, err := os.Create(path + fileName)
	if err != nil {
		fmt.Println("Error when creating new file: ", err.Error())
		os.Exit(1)
	}
	defer newFile.Close()

	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < BufferSize {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BufferSize)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BufferSize)
		receivedBytes += BufferSize
	}
	fmt.Println("File received successfully!")
	return path + fileName
}
