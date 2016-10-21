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
			fmt.Println("Could not successfully load config from file: %v", os.Args[1])
			os.Exit(1)
		}
	} else {
		fmt.Println("Config file you passed in could not be found: %v", os.Args[1])
		os.Exit(1)
	}

	cfg.LocalVersion = getCurrentVersion(cfg.LocalVersionURI)

	if cfg.InitialStartup {
		initialStartup()
	}

	if cfg.FirstRunAfterUpdate {
		firstRunAfterUpdate()
	}

	go func() {
		if err := waitForUpdates(); err != nil {
			fmt.Println("Error occurred while processing updates: %v", err.Error())
		} else {
			fmt.Println("waitForUpdates() gracefully excited. Well played.")
		}
	}()

	if cfg.MineEther {
		go func() {
			if err := mine(); err != nil {
				fmt.Println("Error occurred during the mining process: %v", err.Error())
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
		fmt.Println("waiting for updates. sleeping %v seconds", cfg.CheckInFrequencySeconds)
		time.Sleep(cfg.CheckInFrequencySeconds * time.Second)
		var s string // hold the value from the http GET
		if resp, err := http.Get(cfg.RemoteVersionURI); err != nil {
			return err
			defer resp.Body.Close()
			if body, err := ioutil.ReadAll(resp.Body); err == nil {
				s = string(body[:])
			} else {
				return err
			}
		}

		if remoteVersion, err := strconv.ParseUint(s, 10, 64); err == nil {
			if remoteVersion > cfg.LocalVersion {
				fmt.Println("localVersion: %v", cfg.LocalVersion)
				fmt.Println("remoteVersion: %v", remoteVersion)
				fmt.Println("Newer remote version available. Performing update.")
				doUpdate()
			}
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

// getCurrentVersion will query the local version file for the build
// number and return it.
func getCurrentVersion(versionFilePath string) uint64 {
	if bytes, err := ioutil.ReadFile(versionFilePath); err == nil {
		s := string(bytes)
		s = strings.Trim(s, "\n")
		if compiledVersion, castError := strconv.ParseUint(s, 10, 64); castError == nil {
			return compiledVersion
		}
	}
	return 0
}

func doUpdate() error {
	fmt.Println("updating...")
	return nil
}
