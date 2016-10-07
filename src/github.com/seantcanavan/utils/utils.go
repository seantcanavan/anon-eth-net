package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func FullDateString() string {
	return time.Now().String()
}

func UnixDateString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// FullDateStringSafe returns the current time as a string with only file-name safe characters
func FullDateStringSafe() string {
	t := time.Now()
	return fmt.Sprintf("%v-%02d-%02d %02d_%02d_%02d.%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

// ReadLines reads in a file by path and returns a slice of strings
// credit to: https://stackoverflow.com/a/18479916/584947
func ReadLines(path string) ([]string, error) {

	file, err := os.Open(path)

	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

func DirectoryList(directoryName string, filter string) []string {

	fileList := []string{}
	_ = filepath.Walk(directoryName, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, filter) {
			fileList = append(fileList, path)
		}
		return nil
	})

	return fileList
}
