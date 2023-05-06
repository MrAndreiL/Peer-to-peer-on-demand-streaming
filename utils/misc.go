package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

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
