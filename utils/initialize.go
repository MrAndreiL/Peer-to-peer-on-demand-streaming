package utils

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func CreateFileHashMapping() map[string]string {
	mapping := make(map[string]string)
	// Audio mapping first.
	audioListPath, err := filepath.Abs(AudioList)
	if err != nil {
		fmt.Println("Error when retrieving absolute path: ", err.Error())
		os.Exit(1)
	}
	// Open file and read content.
	lines := ReadFileLine(audioListPath)
	// Get .mp3 and hashed versions.
	for _, line := range lines {
		name, hashcode := NameHashCode(line, ".mp3")
		mapping[name] = hashcode
	}
	// Video mapping second.
	videoListPath, err := filepath.Abs(VideoList)
	if err != nil {
		fmt.Println("Error when retrieving absolute path: ", err.Error())
		os.Exit(1)
	}
	// Open file and read content.
	lines = ReadFileLine(videoListPath)
	// Get .mp4 and hashed versions.
	for _, line := range lines {
		name, hashcode := NameHashCode(line, ".mp4")
		mapping[name] = hashcode
	}
	return mapping
}

func CreateFilePathMapping() map[string]string {
	mapping := make(map[string]string)
	// Audio mapping first.
	audioListPath, err := filepath.Abs(AudioList)
	if err != nil {
		fmt.Println("Error when retrieving absolute path: ", err.Error())
		os.Exit(1)
	}
	// Open file and read content.
	lines := ReadFileLine(audioListPath)
	// Get .mp3 and hashed versions.
	for _, line := range lines {
		name, path := NameFilePath(line, ".mp3")
		mapping[name] = path
	}
	// Video mapping second.
	videoListPath, err := filepath.Abs(VideoList)
	if err != nil {
		fmt.Println("Error when retrieving absolute path: ", err.Error())
		os.Exit(1)
	}
	// Open file and read content.
	lines = ReadFileLine(videoListPath)
	// Get .mp4 and hashed versions.
	for _, line := range lines {
		name, path := NameFilePath(line, ".mp4")
		mapping[name] = path
	}
	return mapping
}

func CreateFileLengthMapping() map[string]string {
	mapping := make(map[string]string)
	// Audio mapping first.
	audioListPath, err := filepath.Abs(AudioList)
	if err != nil {
		fmt.Println("Error when retrieving absolute path: ", err.Error())
		os.Exit(1)
	}
	// Open file and read content.
	lines := ReadFileLine(audioListPath)
	// Get .mp3 and hashed versions.
	for _, line := range lines {
		name, length := NameFileLength(line, ".mp3")
		mapping[name] = length
	}
	// Video mapping second.
	videoListPath, err := filepath.Abs(VideoList)
	if err != nil {
		fmt.Println("Error when retrieving absolute path: ", err.Error())
		os.Exit(1)
	}
	// Open file and read content.
	lines = ReadFileLine(videoListPath)
	// Get .mp4 and hashed versions.
	for _, line := range lines {
		name, length := NameFileLength(line, ".mp4")
		mapping[name] = length
	}
	return mapping
}

func ReadFileLine(path string) []string {
	readFile, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file: ", err.Error())
		os.Exit(1)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	lines := make([]string, 0)
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lines = append(lines, line)
	}
	return lines
}

func NameFilePath(line, extension string) (string, string) {
	fd, err := os.Stat(line)
	if err != nil {
		fmt.Println("Error when getting file status: ", err.Error())
		os.Exit(1)
	}
	name := strings.TrimSuffix(fd.Name(), filepath.Ext(fd.Name())) + extension
	path := strings.TrimSuffix(line, fd.Name())
	return name, path
}

func NameHashCode(line, extension string) (string, string) {
	fd, err := os.Stat(line)
	if err != nil {
		fmt.Println("Error when getting file status: ", err.Error())
		os.Exit(1)
	}
	name := strings.TrimSuffix(fd.Name(), filepath.Ext(fd.Name())) + extension
	hashCode := HashCode(fd.Name(), fmt.Sprint(fd.Size()))
	return name, hashCode
}

func NameFileLength(line, extension string) (string, string) {
	fd, err := os.Stat(line)
	if err != nil {
		fmt.Println("Error when getting file status: ", err.Error())
		os.Exit(1)
	}
	name := strings.TrimSuffix(fd.Name(), filepath.Ext(fd.Name())) + extension
	path := strings.TrimSuffix(line, fd.Name())
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Error when getting the number of files: ", err.Error())
		os.Exit(1)
	}
	return name, fmt.Sprint(len(files))
}

func HashCode(name, size string) string {
	h := fnv.New32a()
	s := name + size
	h.Write([]byte(s))
	return fmt.Sprint(h.Sum32())
}
