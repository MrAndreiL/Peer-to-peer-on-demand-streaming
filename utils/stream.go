package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CreateStream(inputFile, outputDir, fileName string, segmentDuration int) string {
	// Firstly, create output directory to store playlist and stream segments.
	err := CreateDirectory(outputDir)
	if err != nil {
		fmt.Println("Error when creating output directory: " + err.Error())
	}

	// Execute command to create HLS playlist and streams.
	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-c:a", "libmp3lame", // baseline profile compatible with most devices.
		"-b:a", "128k",
		"-map", "0:0", // number streams from 0
		"-f", "segment", // length of each segment
		"-segment_time", "10", // all segments in the playlist
		"-segment_list", fmt.Sprintf("%s/%s.m3u8", outputDir, fileName),
		"-segment_format", "mpegts",
		outputDir+"/"+fileName+"%03d.ts",
	)

	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("failed to create HLS: %v\nOutput: %s", err, string(output))
	}
	playlist := outputDir + "/" + fileName
	return playlist
}

func Streamfy() {
	// Iterate through all existing raw media files.
	// Create streams for each media file.
	audioRawPath, err := filepath.Abs(AudioPath)
	if err != nil {
		fmt.Println("Error when retrieving absolute path: ", err.Error())
		os.Exit(1)
	}
	MediaToStream(audioRawPath)
	videoRawPath, err := filepath.Abs(VideoPath)
	if err != nil {
		fmt.Println("Error when retrieving absolute path: ", err.Error())
		os.Exit(1)
	}
	MediaToStream(videoRawPath)
}

func MediaToStream(pathToFile string) {
	filepath.Walk(pathToFile, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err.Error())
		}
		ext := filepath.Ext(path)
		if ext == ".mp4" || ext == ".mp3" {
			fileName := strings.TrimSuffix(info.Name(), ext)
			directory := filepath.Dir(pathToFile) + "/stream/" + fileName
			CreateDirectory(directory)
			playlist := CreateStream(path, directory, fileName, SegmentLength)
			WriteToFile(filepath.Dir(pathToFile), playlist+".m3u8")
			// DeleteFile(path)
		}
		return nil
	})
}

func DeleteFile(path string) {
	if err := os.Remove(path); err != nil {
		fmt.Println("Error deleting raw file")
		os.Exit(1)
	}
}

func WriteToFile(path, message string) {
	path += "/list.txt"
	fd, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("Failure at opening list.txt file: ", err.Error())
		os.Exit(1)
	}
	defer fd.Close()
	if _, err := fd.WriteString(message + "\n"); err != nil {
		fmt.Println("Error at writing to list.txt file: ", err.Error())
		os.Exit(1)
	}
}
