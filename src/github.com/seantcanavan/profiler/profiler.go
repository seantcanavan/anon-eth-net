package profiler

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"os/exec"
	// "runtime"

	"github.com/seantcanavan/loader"
	"github.com/seantcanavan/reporter"
	"github.com/seantcanavan/utils"
)

// the text to print after we finish executing a command so successive commands have some kind of visual break
const END_COMMAND_DIVIDER = "------------------------------------------------------------\n\n"

// prof is short for profile and is the file extension we use after generating full system profiles
const PROFILE_FILE_EXTENSION = ".prof"

// the subject of the email that we use when sending out the profile
const PROFILE_EMAIL_SUBJECT = "System Profile"

// the base name of the archive file that we save all our individual reports into
const SYS_PROFILE_ARCHIVE_NAME = "profile_archive"

// ProfileAsArchive will generate an individual file for each
// system resources and then compress them all together. Returns a pointer to
// the compressed file containing all of the profile pieces inside of it.
// Automatically deletes the individual files after they've been compressed to
// clean up disk space.
func ProfileAsArchive() (*os.File, error) {

	var profileLoader *loader.Loader

	tarBall, err := os.Create(utils.TimeStampFileName(SYS_PROFILE_ARCHIVE_NAME, ".tar"))
	if err != nil {
		_ = tarBall.Close()
		_ = os.Remove(tarBall.Name())
		return nil, err
	}

	loaderAssetPath, assetErr := utils.SysAssetPath("profiler_loader.json")
	if assetErr != nil {
		return nil, assetErr
	}

	profileLoader, err = loader.NewLoader(loaderAssetPath)
	if err != nil {
		return nil, fmt.Errorf("Loader returned error while trying to generate Profile: %v", err)
	}

	profilerProcesses := profileLoader.StartAsynchronous()

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
func SendArchiveProfileAsAttachment() (*os.File, error) {
	filePtr, err := ProfileAsArchive()
	if err != nil {
		return nil, err
	}
	return filePtr, reporter.SendAttachment(generateEmailSubject(), generateEmailBody(), filePtr)
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
	return cmdBuffer.Bytes()
}
