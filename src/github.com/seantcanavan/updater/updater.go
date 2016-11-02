package updater

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/logger"
)

type Updater struct {
	lgr *logger.Logger
}

func NewUpdater() (*Updater, error) {
	localLogger, loggerError := logger.FromVolatilityValue("updater_package")
	if loggerError != nil {
		return nil, loggerError
	}

	udr := Updater{lgr: localLogger}
	return &udr, nil
}

// Run will continuously check for updated versions of the software
// and update to a newer version if found. Successive version checks will take
// place after a given number of seconds and compare the remote build number
// to the local build number to see if an update is required.
func (udr *Updater) Run() error {
	udr.lgr.LogMessage("waiting for updates. sleeping %v seconds", config.Cfg.CheckInFrequencySeconds)
	time.Sleep(config.Cfg.CheckInFrequencySeconds * time.Second)
	local, localError := udr.localVersion(config.Cfg.LocalVersionURI)
	remote, remoteError := udr.remoteVersion(config.Cfg.RemoteVersionURI)

	if localError != nil {
		return localError
	} else if remoteError != nil {
		return remoteError
	}

	if remote > local {
		udr.lgr.LogMessage("localVersion: %v", local)
		udr.lgr.LogMessage("remoteVersion: %v", remote)
		udr.lgr.LogMessage("Newer remote version available. Performing update.")
		udr.doUpdate()
	}
	return nil
}

// getCurrentVersion will grab the version of this program from the local given
// file path where the version number should reside as a whole integer number.
// The default project structure is to have this file be named 'version.no' and
// be placed within the main package.
func (udr *Updater) localVersion(versionFilePath string) (uint64, error) {
	bytes, err := ioutil.ReadFile(versionFilePath)
	if err != nil {
		return 0, err
	}

	s := string(bytes)
	s = strings.Trim(s, "\n")
	compiledVersion, castError := strconv.ParseUint(s, 10, 64)
	if castError != nil {
		return 0, castError
	}
	return compiledVersion, nil
}

// getRemoteVersion will grab the version of this program from the remote given
// file path where the version number should reside as a whole integer number.
// The default project structure is to have this file be named 'version.no' and
// queried directly via the github.com API.
func (udr *Updater) remoteVersion(versionFilePath string) (uint64, error) {
	var s string // hold the value from the http GET
	resp, getError := http.Get(config.Cfg.RemoteVersionURI)
	if getError != nil {
		return 0, getError
	}

	defer resp.Body.Close()
	body, readError := ioutil.ReadAll(resp.Body)
	if readError != nil {
		return 0, readError
	}

	s = string(body[:])
	s = strings.Trim(s, "\n")

	remoteVersion, castError := strconv.ParseUint(s, 10, 64)
	if castError != nil {
		return 0, castError
	}

	return remoteVersion, nil
}

func (udr *Updater) doUpdate() error {
	udr.lgr.LogMessage("updating...")
	return nil
}
