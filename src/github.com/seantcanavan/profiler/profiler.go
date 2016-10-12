package profiler

import (
	"archive/tar"
	"compress/gzip"
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/seantcanavan/reporter"
	"github.com/seantcanavan/utils"
)

const CPU_AND_DISK_UTIL = "CPU and Disk Utilization"
const MEMORY_UTILIZATION = "Memory Utilization"
const DISK_FREE_SPACE = "Disk Free Space"
const PROCESS_LIST = "Running Processes"
const UPTIME = "Uptime"
const NETWORK_DETAILS = "Network Details"
const KERNEL_VERSION = "Kernel Version"
const END_COMMAND_DIVIDER = "------------------------------------------------------------"
const SEPARATING_SEQUENCE = "\n\n"
const FILE_CLEANUP_DELAY = 360
const PROFILE_FILE_EXTENSION = ".rep"
const PROFILE_EMAIL_SUBJECT = "System Profile"

type SysProfiler struct {
	repr *reporter.Reporter
}

func NewSysProfiler(repr *reporter.Reporter) *SysProfiler {
	SysProfiler := SysProfiler{}
	SysProfiler.repr = repr
	return &SysProfiler
}

// ProfileAsStrings will return a profile of the current executing system
// generated entirely in memory. The full profile will be returned as a
// concatenated slice of strings.
func (sp *SysProfiler) ProfileAsStrings() []string {

	var bytesRepr bytes.Buffer
	var stringSliceRepr []string

	bytesRepr.Write(sp.ProfileAsBytes())
	scanner := bufio.NewScanner(&bytesRepr)

	for scanner.Scan() {
    	stringSliceRepr = append(stringSliceRepr, scanner.Text())
	}

	return stringSliceRepr
}

// ProfileAsBytes will return a profile of the current executing system
// generated entirely in memory. The full profile  will be returned as an array
// of bytes.
func (sp *SysProfiler) ProfileAsBytes() []byte {

	var inMemoryBuffer bytes.Buffer

	inMemoryBuffer.Write(kernelVersion())
	inMemoryBuffer.Write(cpuAndDiskUtilization())
	inMemoryBuffer.Write(memoryUtilization())
	inMemoryBuffer.Write(diskFreeSpace())
	inMemoryBuffer.Write(processList())
	inMemoryBuffer.Write(uptime())
	inMemoryBuffer.Write(networkDetails())

	return inMemoryBuffer.Bytes()
}

// ProfileAsFile will return a profile of the current executing system and save
// the full profile to a local file and return the name of that file.
func (sp *SysProfiler) ProfileAsFile() (*os.File, error) {

	bytes := sp.ProfileAsBytes()
	fileName := utils.TimeStampFileName("sys_profile", PROFILE_FILE_EXTENSION)

	ioutil.WriteFile(fileName, bytes, 0744)
	filePtr, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return filePtr, nil
}

// ProfileAsArchive will generate an individual file for each
// system resources and then compress them all together. Returns a pointer to
// the compressed file containing all of the profile pieces inside of it.
// Automatically deletes the individual files after they've been compressed to
// clean up disk space.
func (sp *SysProfiler) ProfileAsArchive() (*os.File, error) {

	tarBall, err := os.Create(utils.TimeStampFileName("sys_archive", ".tar"))
	if err != nil {
		return nil, err
	}

	gzipWriter := gzip.NewWriter(tarBall)
	tarWriter := tar.NewWriter(gzipWriter)
	pieces := make(map[string][]byte)

	defer gzipWriter.Close()
	defer tarWriter.Close()

	pieces[utils.TimeStampFileName("kernel_version", PROFILE_FILE_EXTENSION)] = kernelVersion()
	pieces[utils.TimeStampFileName("memory_util", PROFILE_FILE_EXTENSION)] = memoryUtilization()
	pieces[utils.TimeStampFileName("network_details", PROFILE_FILE_EXTENSION)] = networkDetails()
	pieces[utils.TimeStampFileName("process_list", PROFILE_FILE_EXTENSION)] = processList()
	pieces[utils.TimeStampFileName("cpu_and_disk", PROFILE_FILE_EXTENSION)] = cpuAndDiskUtilization()
	pieces[utils.TimeStampFileName("disk_space", PROFILE_FILE_EXTENSION)] = diskFreeSpace()
	pieces[utils.TimeStampFileName("uptime", PROFILE_FILE_EXTENSION)] = uptime()

	for file, contents := range pieces {
		tarHeader := &tar.Header{
			Name: file,
			Mode: 0600,
			Size: int64(len(contents)),
		}
		if err := tarWriter.WriteHeader(tarHeader); err != nil {
			return nil, err
		}
		if _, err := tarWriter.Write(contents); err != nil {
			return nil, err
		}
	}

	return tarBall, nil
}

// SendByteProfileAsEmail will generate a full system profile in memory as a
// byte array and then stuff those bytes directly into the message of an email.
func (sp *SysProfiler) SendByteProfileAsEmail() error {
	bytes := sp.ProfileAsBytes()
	return sp.repr.SendPlainEmail(generateEmailSubject(), bytes)
}

// SendFileProfileAsAttachment will generate a full system profile on disk as a
// file and then attach the file directly to an email before sending.
func (sp *SysProfiler) SendFileProfileAsAttachment() error {
	filePtr, err := sp.ProfileAsFile()
	if err != nil {
		return err
	}
	return sp.repr.SendEmailAttachment(generateEmailSubject(), generateEmailBody(), filePtr)
}

// SendArchiveReportAsAttachment will generate each individual piece of the
// system profile inside its own file. It will then gzip and tarball the
// resulting pieces into a single archive for compressing and convenience
// purposes. The original pieces will be automatically cleaned up the archive
// is generated.
func (sp *SysProfiler) SendArchiveProfileAsAttachment() error {
	filePtr, err := sp.ProfileAsArchive()
	if err != nil {
		return err
	}
	return sp.repr.SendEmailAttachment(generateEmailSubject(), generateEmailBody(), filePtr)
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

func generateEmailSubject() string {
	return PROFILE_EMAIL_SUBJECT + utils.FullDateString()
}

func generateEmailBody() []byte {
	var buf bytes.Buffer
	buf.WriteString("A full system profile is attached.")
	return buf.Bytes()
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

func kernelVersion() []byte {
	return execCommand(KERNEL_VERSION, "uname", "-r")
}

func networkDetails() []byte {

	var networkBuffer bytes.Buffer
	networkCommands := []string{"netstat -r", "netstat -i", "netstat -s", "netstat -va"}

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
	cmdBuffer.WriteString(SEPARATING_SEQUENCE)
	return cmdBuffer.Bytes()
}
