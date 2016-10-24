package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nu7hatch/gouuid"
	"github.com/seantcanavan/config"
)

var cfg = &config.Config{}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("One and only argument accepted: the file path to config.json to initialize the program.")
		os.Exit(1)
	}

	if os.Args[1] == "h" || os.Args[1] == "help" || os.Args[1] == "?" {
		fmt.Println("One and only argument accepted: the file path to config.json to initialize the program. Please refer to the sample provided config.json file in ./config/config.json to get started.")
	}

	if _, err := os.Stat(os.Args[1]); err == nil {
		if loadedConfig, configError := config.ConfigFromFile(os.Args[1]); configError == nil {
			cfg = loadedConfig
		} else {
			fmt.Println(fmt.Sprintf("Could not successfully load config from file: %v", os.Args[1]))
			os.Exit(1)
		}
	} else {
		fmt.Println(fmt.Sprintf("Config file you passed in could not be found: %v", os.Args[1]))
		os.Exit(1)
	}

	if cfg.InitialStartup {
		initialStartup()
	}

	if cfg.FirstRunAfterUpdate {
		firstRunAfterUpdate()
	}

	go func() {
		if err := waitForUpdates(); err != nil {
			fmt.Println(fmt.Sprintf("Error occurred while processing updates: %v", err.Error()))
		} else {
			fmt.Println("waitForUpdates() gracefully excited. Well played.")
		}
	}()

	if cfg.MineEther {
		go func() {
			if err := mine(); err != nil {
				fmt.Println(fmt.Sprintf("Error occurred during the mining process: %v", err.Error()))
			} else {
				fmt.Println("mine() gracefully excited. Well played.")
			}
		}()
	}
}

func initialStartup() {
	uuid, err := uuid.NewV4()
	if err != nil {
		cfg.DeviceId = uuid.String()
	}
}

func firstRunAfterUpdate() {

}

// waitForUpdates will continuously check for updated versions of the software
// and update to a newer version if found. Successive version checks will take
// place after a given number of seconds and compare the remote build number
// to the local build number to see if an update is required.
func waitForUpdates() error {
	for 1 == 1 {
		fmt.Println(fmt.Sprintf("waiting for updates. sleeping %v seconds", cfg.CheckInFrequencySeconds))
		time.Sleep(cfg.CheckInFrequencySeconds * time.Second)
		local, localError := localVersion(cfg.LocalVersionURI)
		remote, remoteError := remoteVersion(cfg.RemoteVersionURI)

		if localError != nil {
			return localError
		} else if remoteError != nil {
			return remoteError
		}

		if remote > local {
			fmt.Println(fmt.Sprintf("localVersion: %v", local))
			fmt.Println(fmt.Sprintf("remoteVersion: %v", remote))
			fmt.Println("Newer remote version available. Performing update.")
			doUpdate()
		}
	}
	return nil
}

// mine will kick off all the programs in order to start mining ethereum.
func mine() error {
	for 1 == 1 {
		fmt.Println("mining...")
		time.Sleep(10 * time.Second)
	}
	return nil
}

// getCurrentVersion will grab the version of this program from the local given
// file path where the version number should reside as a whole integer number.
// The default project structure is to have this file be named 'version.no' and
// be placed within the main package.
func localVersion(versionFilePath string) (uint64, error) {
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
func remoteVersion(versionFilePath string) (uint64, error) {
	var s string // hold the value from the http GET
	resp, getError := http.Get(cfg.RemoteVersionURI)
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

func doUpdate() error {
	fmt.Println("updating...")
	return nil
}
