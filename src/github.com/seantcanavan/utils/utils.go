package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// FullDateString will return the current time formatted as a string.
func FullDateString() string {
	return time.Now().String()
}

// UnixDateString will return the current time in unix time format as a string.
func UnixDateString() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

// FullDateStringSafe returns the current time as a string with only file-name
// safe characters. Used to quickly and easily generate unique file names based
// off of the current system time.
func FullDateStringSafe() string {
	t := time.Now()
	return fmt.Sprintf("%v-%02d-%02d %02d_%02d_%02d.%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
}

// TimeStampFileName will generate a string to be used to uniquely name and
// identify a new file. The characters used in the file name are guaranteed to
// be safe for a file. The given fileBaseName will be used to append to the
// beginning of the file name to give some control over the file name.
func TimeStampFileName(fileBaseName string, fileExtension string) string {

	dts := FullDateStringSafe()

	var nameBuffer bytes.Buffer
	nameBuffer.WriteString(fileBaseName)
	nameBuffer.WriteString(" ")
	nameBuffer.WriteString(dts)
	nameBuffer.WriteString(fileExtension)

	return nameBuffer.String()
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

// DirectoryList returns a list of entries in the given directory. It filters
// optionally for the given filter.
func DirectoryList(directoryName string, filter string) []string {

	fileList := []string{}
	_ = filepath.Walk(directoryName, func(path string, f os.FileInfo, err error) error {
		if filter != "" {
			if strings.Contains(path, filter) {
				fileList = append(fileList, path)
			}
		}
		return nil
	})

	return fileList
}
