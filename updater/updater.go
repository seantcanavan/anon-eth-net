package updater

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/seantcanavan/anon-eth-net/config"
	"github.com/seantcanavan/anon-eth-net/logger"
)

// Run will continuously check for updated versions of the software
// and update to a newer version if found. Successive version checks will take
// place after a given number of seconds and compare the remote build number
// to the local build number to see if an update is required.
func Run() {

	go func() {

		for 1 == 1 {

			logger.Lgr.LogMessage("waiting for updates. sleeping %v", config.Cfg.UpdateFrequencySeconds)
			time.Sleep(time.Duration(config.Cfg.UpdateFrequencySeconds) * time.Second)

			local := config.Cfg.LocalVersion
			remote, remoteErr := remoteVersion()

			if remoteErr != nil {
				logger.Lgr.LogMessage("Error retrieving the remote version: %v", remoteErr.Error())
				continue
			}

			if remote > local {
				logger.Lgr.LogMessage("localVersion: %v", local)
				logger.Lgr.LogMessage("remoteVersion: %v", remote)
				logger.Lgr.LogMessage("Newer remote version available. Performing update.")
				doUpdate()
			}
		}
	}()
}

// UpdateNecessary will look at the remotely defined version number as well as
// the locally defined version number and compare the two. Based on the result
// it will recommend a course of action. It will return True is the remote
// version is higher (newer) than the local version.
func UpdateNecessary() (bool, error) {

	localVersion := config.Cfg.LocalVersion

	remoteVersion, remoteErr := remoteVersion()
	if remoteErr != nil {
		return false, remoteErr
	}

	if localVersion > remoteVersion {
		logger.Lgr.LogMessage("Your version, %v, is higher than the remote: %v. Push your changes!", localVersion, remoteVersion)
	}

	if localVersion == remoteVersion {
		logger.Lgr.LogMessage("Your version, %v, equals the remote: %v. Do some work!", localVersion, remoteVersion)
	}

	if localVersion < remoteVersion {
		logger.Lgr.LogMessage("Your version, %v, is lower than the remote: %v. Pull the latest code and build it!", localVersion, remoteVersion)
	}

	return remoteVersion > localVersion, nil

}

// remoteVersion will grab the version of this program from the remote given
// file path where the version number should reside as a whole integer number.
// The default project structure is to have this file be named 'version.no' and
// queried directly via the github.com API.
func remoteVersion() (uint64, error) {

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

	logger.Lgr.LogMessage("Successfully retrieved remote version: %v", remoteVersion)
	return remoteVersion, nil
}

// doUpdate will hopefully someday actually perform the update
func doUpdate() error {
	logger.Lgr.LogMessage("performing an update")
	return nil
}
