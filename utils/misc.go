package utils

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"time"
)

var MappingFileHash map[string]string
var MappingFilePath map[string]string
var MappingFileLength map[string]string

type Peer struct {
	publicAddress  string
	privateAddress string
	networkName    string
	connection     net.Conn
}

func CreateNetworkPeer(public, private, name string, conn net.Conn) *Peer {
	var peer Peer
	peer.publicAddress = public
	peer.privateAddress = private
	peer.networkName = name
	peer.connection = conn
	return &peer
}

func SetFileMappings(fileHash, filePath, fileLength map[string]string) {
	MappingFileHash = fileHash
	MappingFilePath = filePath
	MappingFileLength = fileLength
}

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length+1)
	first := false
	for i := range b {
		if first == false {
			b[i] = '#'
			first = true
		} else {
			b[i] = charset[seededRand.Intn(len(charset))]
		}
	}
	return string(b)
}

func CreateDirectory(dir string) error {
	// Create output directory. Test if it already exists.
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory")
	}
	return nil
}

func CreateFilesIfNotExists() {
	// Creates audio and video list files to store
	// locally available media files.
	audioListPath, err := filepath.Abs(AudioList)
	if err != nil {
		fmt.Println("Error when retrieving absolute path: ", err.Error())
		os.Exit(1)
	}
	_, err = os.Stat(audioListPath)
	if os.IsExist(err) {
		return
	}
	fda, err := os.Create(audioListPath)
	if err != nil {
		fmt.Println("Could not create audio list file.")
		os.Exit(1)
	}
	defer fda.Close()
	videoListPath, err := filepath.Abs(VideoList)
	if err != nil {
		fmt.Println("Error when retrieving absolute path: ", err.Error())
		os.Exit(1)
	}
	_, err = os.Stat(videoListPath)
	if os.IsExist(err) {
		return
	}
	fdv, err := os.Create(videoListPath)
	if err != nil {
		fmt.Println("Could not create video list file")
		os.Exit(1)
	}
	defer fdv.Close()
}
