package profiler

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"time"
	// "runtime"

	"github.com/seantcanavan/anon-eth-net/config"
	"github.com/seantcanavan/anon-eth-net/loader"
	"github.com/seantcanavan/anon-eth-net/logger"
	"github.com/seantcanavan/anon-eth-net/reporter"
	"github.com/seantcanavan/anon-eth-net/utils"
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

	logger.Lgr.LogMessage("Successfully created new archive tar: %v", tarBall.Name())

	loaderAssetPath, assetErr := utils.SysAssetPath("profiler_loader.json")
	if assetErr != nil {
		return nil, assetErr
	}

	logger.Lgr.LogMessage("Successfully loaded profile_loader asset: %v", loaderAssetPath)

	profileLoader, err = loader.NewLoader(loaderAssetPath)
	if err != nil {
		return nil, fmt.Errorf("Loader returned error while trying to generate Profile: %v", err)
	}

	logger.Lgr.LogMessage("Successfully created new profile loader instance with: %v", loaderAssetPath)

	profilerProcesses := profileLoader.StartAsynchronous()

	logger.Lgr.LogMessage("Finished executing profiler loader processes to get full system report details")

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

		logger.Lgr.LogMessage("Successfully wrote %v to the tarball and deleted it from disk", logName)
	}

	if err != nil {
		_ = tarBall.Close()
		_ = os.Remove(tarBall.Name())
		logger.Lgr.LogMessage("Error during tarball creation. Cleaned up tarball but process logs will remain")
		return nil, err
	}

	logger.Lgr.LogMessage("Successfully created tarball of log files from profile process loader")

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

	logger.Lgr.LogMessage("Successfully created system profile archive. Will attempt to email now")

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

func Run() {

	// kick off the system profiler loop to send out system profiles at the specified interval
	go func() {

		for 1 == 1 {

			logger.Lgr.LogMessage("Sleeping for %d seconds before sending a system profile", config.Cfg.CheckInFrequencySeconds)
			time.Sleep(time.Duration(config.Cfg.CheckInFrequencySeconds) * time.Second)
			logger.Lgr.LogMessage("Sending archive to provided email after sleeping %d seconds", config.Cfg.CheckInFrequencySeconds)
			SendArchiveProfileAsAttachment()

		}

	}()

}
