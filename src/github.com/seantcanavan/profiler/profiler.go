package profiler

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"
)

const CPU_AND_DISK_UTIL = "CPU and Disk Utilization"
const MEMORY_UTILIZATION = "Memory Utilization"
const DISK_FREE_SPACE = "Disk Free Space"
const PROCESS_LIST = "Running Processes"
const UPTIME = "Uptime"
const NETWORK_DETAILS = "Network Details"
const END_COMMAND_DIVIDER = "(------------------------------------------------)"

type Profiler struct {

}


// GetInMemorySystemReportAsStringSlice will return a system report generated
// entirely in memory. The full report will be returned as a concatenated slice
// of strings.
func GetInMemorySystemReportAsStringSlice() {

}

// GetInMemorySystemReportAsByteArray will return a system report generated
// entirely in memory. The full report will be returned as an array of bytes.
func GetInMemorySystemReportAsByteArray() {

}

// GetInFileSystemReport will return a system report that has been saved to a
// file and return the pointer to that file. If the requested report is
// transient - it will be deleted after a minute automatically.
func GetInFileSystemReport(transient bool) *os.File {
	return nil
}

// GetCompressedSystemReport will generate an individual file for each
// system resources and then compress them all together. Returns a pointer to
// the compressed file containing all of the reports inside of it. Deletes the
// individual files after they've been compressed to clean up disk space. If the
// requested report is transient it will be deleted after a minute automatically.
func GetCompressedSystemReport(transient bool) *os.File {
	return nil
}

func beautifyTitle(title string) []byte {
	var titleBuffer bytes.Buffer
	titleBuffer.WriteString("--------------------  ")
	titleBuffer.WriteString(title)
	titleBuffer.WriteString("  --------------------")
	return titleBuffer.Bytes()
}

func (p *Profiler) getCPUAndDiskUtilization() []byte {

	return p.execCommand(CPU_AND_DISK_UTIL, "iostat 1 10")
}

func (p *Profiler) getMemoryUtilization() []byte {

	return p.execCommand(MEMORY_UTILIZATION, "free")
}

func (p *Profiler) getDiskFreeSpace() []byte {

	return p.execCommand(DISK_FREE_SPACE, "df -h")
}

func (p *Profiler) getRunningProcessList() []byte {

	return p.execCommand(PROCESS_LIST, "ps -AlF")
}

func (p *Profiler) getUptime() []byte {

	return p.execCommand(UPTIME, "uptime")
}

func (p *Profiler) getNetworkDetails() []byte {

	return p.execCommand(NETWORK_DETAILS, "netstat -vea")
	// TODO(Canavan) add all these commands to the network details
	// netstat -r
	// netstat -i
	// netstat -se
	// full socket list
	// netstat -vea
}

func (p *Profiler) execCommand(header string, command string) []byte {

	cmd := exec.Command(command)

	var cmdBuffer bytes.Buffer
	cmdBuffer.Write(beautifyTitle(header))

	if out, err := cmd.Output(); err == nil {
		cmdBuffer.Write(out)
	} else {
		cmdBuffer.WriteString(fmt.Sprintf("Command failed: %v", command))
	}

	cmdBuffer.WriteString(END_COMMAND_DIVIDER)
	return cmdBuffer.Bytes()
}

func (p *Profiler) cleanupFile(systemReport *os.File) error {
	time.Sleep(60 * time.Second)
	return os.Remove(systemReport.Name())
}
