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

	"github.com/seantcanavan/logger"
)

const CPU_AND_DISK_UTIL = "CPU and Disk Utilization"
const MEMORY_UTILIZATION = "Memory Utilization"
const DISK_FREE_SPACE = "Disk Free Space"
const PROCESS_LIST = "Running Processes"
const UPTIME = "Uptime"
const NETWORK_DETAILS = "Network Details"
const END_COMMAND_DIVIDER = "------------------------------------------------------------"
const SEPARATING_SEQUENCE = "\n\n"

// GetInMemorySystemReportAsStringSlice will return a system report generated
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

// GetInMemorySystemReportAsByteArray will return a system report generated
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

// GetInFileSystemReport will return a system report that has been saved to a
// file and return the pointer to that file. If the requested report is
// transient - it will be deleted after  automatically.
func ReportAsFile(transient bool, duration time.Duration) (string, error) {

	bytes := ReportAsBytes()
	fileName := logger.LogFileHandle("sysreport")

	if transient {
		go func() {
			time.Sleep(duration * time.Second)
			//TODO(Canavan) handle this error somehow
			_ = os.Remove(fileName)
		}()
	}

	ioutil.WriteFile(fileName, bytes, 0744)
	return fileName, nil
}

// GetCompressedSystemReport will generate an individual file for each
// system resources and then compress them all together. Returns a pointer to
// the compressed file containing all of the reports inside of it. Deletes the
// individual files after they've been compressed to clean up disk space. If the
// requested report is transient it will be deleted after a minute automatically.
func ReportAsArchive(transient bool) (*os.File, error) {

	memoryUtilizationFileName := logger.LogFileHandle("memory_util")
	networkDetailsFileName := logger.LogFileHandle("network_details")
	processListFileName := logger.LogFileHandle("process_list")
	cpuAndDiskFileName := logger.LogFileHandle("cpu_and_disk")
	diskSpaceFileName := logger.LogFileHandle("disk_space")
	uptimeFileName := logger.LogFileHandle("uptime")

	mem := memoryUtilization()
	net := networkDetails()
	proc := runningProcessList()
	cpudisk := cpuAndDiskUtilization()
	freedisk := diskFreeSpace()
	uptime := uptime()

	ioutil.WriteFile(memoryUtilizationFileName, mem, 0744)
	ioutil.WriteFile(networkDetailsFileName, net, 0744)
	ioutil.WriteFile(processListFileName, proc, 0744)
	ioutil.WriteFile(cpuAndDiskFileName, cpudisk, 0744)
	ioutil.WriteFile(diskSpaceFileName, freedisk, 0744)
	ioutil.WriteFile(uptimeFileName, uptime, 0744)



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

func runningProcessList() []byte {

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
