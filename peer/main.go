package main

import (
	"Peer-to-peer-on-demand-streaming/utils"
)

func main() {
	// 1. If the peer is a superpeer, then have a server open
	// at any time (used as randezvous, connection, or discovery).
	if utils.SuperPeer {
		utils.OpenSuperPeer()
	} else {
		utils.OpenPeer()
	}
	// 1. Create initial files to save a list of all
	// available media files.
	// utils.CreateFilesIfNotExists()
	// 2. Transform all media files into stream files
	// and save them into list.txt files.
	// utils.Streamfy()
	// 3. Open list.txt files and create a mapping
	// to be distributed across the network.
	// mapping := utils.CreateMapping()
	// mapping := utils.CreateMapping()
	// fmt.Println(mapping)
	// utils.OpenServer()

	// for {

	// }
}
