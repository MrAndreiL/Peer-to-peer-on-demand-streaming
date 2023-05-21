package main

import (
	"Peer-to-peer-on-demand-streaming/utils"
)

func main() {
	// 1. Create initial files to save a list of all
	// available media files.
	utils.CreateFilesIfNotExists()
	// 2. Transform all media files into stream files
	// and save them into list.txt files.
	utils.Streamfy()
	// 3. Open list.txt files and create a mapping
	// to be distributed across the network.
	mappingFileHash := utils.CreateFileHashMapping()
	mappingFilePath := utils.CreateFilePathMapping()
	mappingFileLength := utils.CreateFileLengthMapping()
	utils.SetFileMappings(mappingFileHash, mappingFilePath, mappingFileLength)
	// 4. Open Peer.
	// If set, give superpeer privileges.
	if utils.SuperPeer {
		go utils.OpenSuperPeer()
	} else {
		go utils.OpenPeer()
	}
	// 5. Wait on goroutines.
	for {

	}
}
