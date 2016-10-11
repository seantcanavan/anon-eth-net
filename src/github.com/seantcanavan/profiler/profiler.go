package profiler

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/logger"
	"github.com/seantcanavan/reporter"
)

const CPU_AND_DISK_UTIL = "CPU and Disk Utilization"
const MEMORY_UTILIZATION = "Memory Utilization"
const DISK_FREE_SPACE = "Disk Free Space"
const PROCESS_LIST = "Running Processes"
const UPTIME = "Uptime"
const NETWORK_DETAILS = "Network Details"
const END_COMMAND_DIVIDER = "------------------------------------------------------------"
const SEPARATING_SEQUENCE = "\n\n"
const FILE_CLEANUP_DELAY = 180

// ReportAsStrings will return a system report generated
// entirely in memory. The full report will be returned as a concatenated slice
// of strings.
func ReportAsStrings() []string {

	var bytesRepr bytes.Buffer
	var stringSliceRepr []string

	bytesRepr.Write(ReportAsBytes())
	scanner := bufio.NewScanner(&bytesRepr)

	for scanner.Scan() {
    	stringSliceRepr = append(stringSliceRepr, scanner.Text())
	}

	return stringSliceRepr
}

// ReportAsBytes will return a system report generated
// entirely in memory. The full report will be returned as an array of bytes.
func ReportAsBytes() []byte {

	var inMemoryBuffer bytes.Buffer

	inMemoryBuffer.Write(cpuAndDiskUtilization())
	inMemoryBuffer.WriteString(SEPARATING_SEQUENCE)
	inMemoryBuffer.Write(memoryUtilization())
	inMemoryBuffer.WriteString(SEPARATING_SEQUENCE)
	inMemoryBuffer.Write(diskFreeSpace())
	inMemoryBuffer.WriteString(SEPARATING_SEQUENCE)
	inMemoryBuffer.Write(runningProcessList())
	inMemoryBuffer.WriteString(SEPARATING_SEQUENCE)
	inMemoryBuffer.Write(uptime())
	inMemoryBuffer.WriteString(SEPARATING_SEQUENCE)
	inMemoryBuffer.Write(networkDetails())
	inMemoryBuffer.WriteString(SEPARATING_SEQUENCE)

	return inMemoryBuffer.Bytes()
}

// ReportAsFile will return a system report that has been saved to a
// file and return the name of that file. If the requested report is
// transient - it will be deleted after automatically after two minutes.
func ReportAsFile(transient bool) (string, error) {

	bytes := ReportAsBytes()
	fileName := logger.LogFileHandle("sysreport")

	ioutil.WriteFile(fileName, bytes, 0744)

	if transient {
		go cleanupFile(handle)
	}

	return fileName, nil
}

func cleanupFile(filename string) error {
	time.Sleep(FILE_CLEANUP_DELAY * time.Second)
	return os.Remove(fileName)

}

// ReportAsArchive will generate an individual file for each
// system resources and then compress them all together. Returns a pointer to
// the compressed file containing all of the reports inside of it. Deletes the
// individual files after they've been compressed to clean up disk space. If the
// requested report is transient it will be deleted after a minute automatically.
func ReportAsArchive(transient bool) (*os.File, error) {

	reports := make(map[string][]byte)
	tarBuffer := new(bytes.Buffer)
	tarWrite := tar.NewWriter(tarBuffer)

	reports.add(logger.LogFileHandle("memory_util"), memoryUtilization())
	reports.add(logger.LogFileHandle("network_details"), networkDetails())
	reports.add(logger.LogFileHandle("process_list"), processList())
	reports.add(logger.LogFileHandle("cpu_and_disk"), cpuAndDiskUtilization())
	reports.add(logger.LogFileHandle("disk_space"), diskFreeSpace())
	reports.add(logger.LogFileHandle("uptime"), uptime())

	for file, contents := range reports {
		tarHeader := &tar.Header{
			Name: file,
			Mode: 0600,
			Size: int64(len(contents)),
		}
		if err := tarWriter.WriteHeader(tarHeader); err != nil {
			logger.
		}
		if _, err := tarWriter.Write(contents); err != nil {
			return err
		}
	}



	return nil, nil
}

func beautifyTitle(title string) []byte {
	var titleBuffer bytes.Buffer

	titleBuffer.WriteString("-------------------- ")
	titleBuffer.WriteString(title)
	titleBuffer.WriteString(" ")
	for titleBuffer.Len() < 60 {
		titleBuffer.WriteString("-")
	}

	return titleBuffer.Bytes()
}

func cpuAndDiskUtilization() []byte {

	return execCommand(CPU_AND_DISK_UTIL, "iostat",  "1", "10")
}

func memoryUtilization() []byte {

	return execCommand(MEMORY_UTILIZATION, "free")
}

func diskFreeSpace() []byte {

	return execCommand(DISK_FREE_SPACE, "df", "-h")
}

func processList() []byte {

	return execCommand(PROCESS_LIST, "ps", "-AlF")
}

func uptime() []byte {

	return execCommand(UPTIME, "uptime")
}

func networkDetails() []byte {

	var networkBuffer bytes.Buffer
	networkCommands := []string{"netstat -r", "netstat -i", "netstat -se", "netstat -vea"}

	for i := range networkCommands {
		splitCmd := strings.Split(networkCommands[i], " ")
		networkBuffer.Write(execCommand(NETWORK_DETAILS + ": " + networkCommands[i], splitCmd[0], splitCmd[1]))
	}

	return networkBuffer.Bytes()
}

func execCommand(header string, command string, args ...string) []byte {

	var cmdBuffer bytes.Buffer
	cmdBuffer.Write(beautifyTitle(header))
	cmdBuffer.WriteString("\n")

	if out, err := exec.Command(command, args...).Output(); err != nil {
		cmdBuffer.WriteString(fmt.Sprintf("Command failed: \n%v \n\nError: \n%v \n\nOutput: \n%v", command, err, os.Stderr))
	} else {
		cmdBuffer.Write(out)
	}

	cmdBuffer.WriteString(END_COMMAND_DIVIDER)
	return cmdBuffer.Bytes()
}

func cleanupFile(systemReport *os.File) error {
	time.Sleep(60 * time.Second)
	return os.Remove(systemReport.Name())
}
