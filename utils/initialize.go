package utils

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"strings"
)

func CreateMapping() map[string]map[string]string {
	// Return a mapping with the following structure:
	// {"path to playlist file" -> {name.mp3/mp4 -> hashedCode}}
	mapping := make(map[string]map[string]string)
	mapping = AudioListing(mapping)
	mapping = VideoListing(mapping)
	return mapping
}

func AudioListing(m map[string]map[string]string) map[string]map[string]string {
	// Appends to the mapping according to the
	// required structure for audio.
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
		rel := make(map[string]string)
		rel[name] = hashcode
		m[line] = rel
		fmt.Println(rel)
	}
	return m
}

func VideoListing(m map[string]map[string]string) map[string]map[string]string {
	// Appends to the mapping according to the
	// required structure for video.
	videoListPath, err := filepath.Abs(VideoList)
	if err != nil {
		fmt.Println("Error when retrieving absolute path: ", err.Error())
		os.Exit(1)
	}
	// Open file and read content.
	lines := ReadFileLine(videoListPath)
	// Get .mp4 and hashed versions.
	for _, line := range lines {
		name, hashcode := NameHashCode(line, ".mp4")
		rel := make(map[string]string)
		rel[name] = hashcode
		m[line] = rel
		fmt.Println(rel)
	}
	return m
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

func HashCode(name, size string) string {
	h := fnv.New32a()
	s := name + size
	h.Write([]byte(s))
	return fmt.Sprint(h.Sum32())
}
