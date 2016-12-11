package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// the root asset directory where all the external files used are stored
const ASSET_ROOT_DIR = "../assets/"

// GetAssetPath will return the relative path to the file represented by
// assetName otherwise it will return an error if the file doesn't exist.
func AssetPath(assetName string) (string, error) {

	relativePath := ASSET_ROOT_DIR + assetName

	if _, err := os.Stat(relativePath); os.IsNotExist(err) {
		return "", err
	}

	return relativePath, nil
}

// GetSysAssetPath will return the relative path to the file represented by
// assetName but also add in the GOOS after the filename and before the
// extension. This allows loading system-specific files with one command instead
// of a complicated switch statement every time.
func SysAssetPath(assetName string) (string, error) {

	var relativeName bytes.Buffer
	var extIndex int

	fileExt := filepath.Ext(assetName)

	if fileExt != "" {
		// if there is an extension, insert right before it
		extIndex = strings.Index(assetName, fileExt)
	} else {
		// if there is no extension, insert at the end of the name
		extIndex = len(assetName)
	}

	relativeName.WriteString(ASSET_ROOT_DIR)
	relativeName.WriteString(assetName[0:extIndex])

	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		relativeName.WriteString("_")
		relativeName.WriteString(runtime.GOOS)
	default:
		return "", fmt.Errorf("Invalid GOOS for asset string: %v", runtime.GOOS)
	}

	relativeName.WriteString(assetName[extIndex:])
	path := relativeName.String()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("Relative file does not exist: %v", path)
	}

	return relativeName.String(), nil
}

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
	return fmt.Sprintf("[%v-%02d-%02d][%02d_%02d_%02d.%02d]",
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
	nameBuffer.WriteString("_")
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

// credit to: https://gist.github.com/jniltinho/9788121
func ExternalIPAddress() (string, error) {

	var ipBuffer bytes.Buffer

	resp, err := http.Get("http://myexternalip.com/raw")
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	_, copyErr := io.Copy(&ipBuffer, resp.Body)
	if copyErr != nil {
		return "", copyErr
	}

	return strings.Trim(ipBuffer.String(), "\n"), nil
}
