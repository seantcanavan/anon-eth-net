package profiler

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/seantcanavan/loader"
	"github.com/seantcanavan/reporter"
	"github.com/seantcanavan/utils"
)

const END_COMMAND_DIVIDER = "------------------------------------------------------------"
const SEPARATING_SEQUENCE = "\n\n"
const PROFILE_FILE_EXTENSION = ".rep"
const PROFILE_EMAIL_SUBJECT = "System Profile"
const SYS_PROFILE_FILE_NAME = "profile_file"
const SYS_PROFILE_ARCHIVE_NAME = "profile_archive"

type SysProfiler struct {
	repr *reporter.Reporter
}

func NewSysProfiler(repr *reporter.Reporter) *SysProfiler {
	SysProfiler := SysProfiler{}
	SysProfiler.repr = repr
	return &SysProfiler
}

// ProfileAsArchive will generate an individual file for each
// system resources and then compress them all together. Returns a pointer to
// the compressed file containing all of the profile pieces inside of it.
// Automatically deletes the individual files after they've been compressed to
// clean up disk space.
func (sp *SysProfiler) ProfileAsArchive() (*os.File, error) {

	var profileLoader *loader.Loader

	tarBall, err := os.Create(utils.TimeStampFileName(SYS_PROFILE_ARCHIVE_NAME, ".tar"))
	if err != nil {
		_ = tarBall.Close()
		_ = os.Remove(tarBall.Name())
		return nil, err
	}

	switch runtime.GOOS {
	case "windows":
		profileLoader, err = loader.NewLoader("profiler_loader_windows.json")
	case "darwin":
		profileLoader, err = loader.NewLoader("profiler_loader_darwin.json")
	case "linux":
		profileLoader, err = loader.NewLoader("profiler_loader_linux.json")
	default:
		return nil, fmt.Errorf("Could not create profile for unsupported operating system: %v", runtime.GOOS)
	}

	if err != nil {
		return nil, fmt.Errorf("Loader returned error while trying to generate Profile: %v", err)
	}

	profilerProcesses := profileLoader.StartSynchronous()

	gzipWriter := gzip.NewWriter(tarBall)
	tarWriter := tar.NewWriter(gzipWriter)

	defer gzipWriter.Close()
	defer tarWriter.Close()

	for _, currentProcess := range profilerProcesses {
		var logContents []byte
		logContents, err = currentProcess.Lgr.CurrentLogContents()
		if err != nil {
			break
		}

		var logName string
		logName, err = currentProcess.Lgr.CurrentLogName()
		if err != nil {
			break
		}

		tarHeader := &tar.Header{
			Name: logName,
			Mode: 0600,
			Size: int64(len(logContents)),
		}
		if err := tarWriter.WriteHeader(tarHeader); err != nil {
			break
		}
		if _, err := tarWriter.Write(logContents); err != nil {
			break
		}

		err = os.Remove(logName)
		if err != nil {
			break
		}
	}

	if err != nil {
		_ = tarBall.Close()
		_ = os.Remove(tarBall.Name())
		return nil, err
	}

	return tarBall, nil
}

// SendArchiveReportAsAttachment will generate each individual piece of the
// system profile inside its own file. It will then gzip and tarball the
// resulting pieces into a single archive for compressing and convenience
// purposes. The original pieces will be automatically cleaned up the archive
// is generated.
func (sp *SysProfiler) SendArchiveProfileAsAttachment() (*os.File, error) {
	filePtr, err := sp.ProfileAsArchive()
	if err != nil {
		return nil, err
	}
	return filePtr, sp.repr.SendAttachment(generateEmailSubject(), generateEmailBody(), filePtr)
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
	var buf bytes.Buffer
	buf.WriteString(PROFILE_EMAIL_SUBJECT)
	buf.WriteString(" ")
	buf.WriteString(utils.FullDateString())
	return string(buf.Bytes())
}

func generateEmailBody() []byte {
	var buf bytes.Buffer
	buf.WriteString("A full system profile is attached.")
	return buf.Bytes()
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
