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
	"github.com/seantcanavan/logger"
)

// the local instance of SeanLogger which will only log messages from the main package
var log *logger.SeanLogger

func main() {

	if os.Args[1] == "h" || os.Args[1] == "help" || os.Args[1] == "?" {
		// fmt.Println("One and only argument accepted: the file path to config.json to initialize the program. Please refer to the sample provided config.json file in ./config/config.json to get started.")
		fmt.Println("Does not require any command line arguments. Refer to the default config/config.json file for all the parameters required for AEN to execute successfully.")
		fmt.Println(config.ConfigJSONParametersExplained())
	}

	// if len(os.Args) != 2 {
	// 	fmt.Println("One and only argument accepted: the file path to config.json to initialize the program.")
	// 	os.Exit(1)
	// }

	if _, err := os.Stat(os.Args[1]); err == nil {
		if configError := config.ConfigFromFile(config.LOCAL_EXTERNAL_PATH); configError != nil {
			fmt.Println(fmt.Sprintf("Could not successfully load config from file: %v", os.Args[1]))
			os.Exit(1)
		}
	} else {
		fmt.Println(fmt.Sprintf("Config file you passed in could not be found: %v", os.Args[1]))
		os.Exit(1)
	}

	mainLogger, loggerError := logger.LoggerFromConservativeValue(config.Cfg.LoggingVolatility, "main_package")
	if loggerError != nil {
		fmt.Println("Couldn't create log... Exiting...")
		os.Exit(1)
	}
	// maintain a local logging reference for anything kicked off by the main package
	log = mainLogger
	// if this is our first time ever starting up - run the initial config
	if config.Cfg.InitialStartup {
		initialStartup()
	}
	// if this is our first time starting up after an update - run the update / cleanup config
	if config.Cfg.FirstRunAfterUpdate {
		firstRunAfterUpdate()
	}
	// kick off the watchdog loop that will regularly watch for and execute updates when available
	go func() {
		if err := waitForUpdates(); err != nil {
			log.LogMessage(fmt.Sprintf("Error occurred while processing updates: %v", err.Error()))
		} else {
			log.LogMessage("waitForUpdates() gracefully excited. Well played.")
		}
	}()
	// if the user elects to, mine ethereum as well in our free time
	if config.Cfg.MineEther {
		go func() {
			if err := mine(); err != nil {
				log.LogMessage(fmt.Sprintf("Error occurred during the mining process: %v", err.Error()))
			} else {
				log.LogMessage("mine() gracefully excited. Well played.")
			}
		}()
	}
}

// initialStartup will be executed only when this program is running for the
// first time on a new host.
func initialStartup() {
	uuid, err := uuid.NewV4()
	if err != nil {
		// update the UUID if it doesn't exist
		config.Cfg.DeviceId = uuid.String()
	}
	// we're finishing the first run!
	config.Cfg.InitialStartup = false
	// push the UUID back to the file for next time
	config.ConfigToFile(config.LOCAL_EXTERNAL_PATH)
}

// firstRunAfterUpdate will be executed only when this program is running for
// the first time after an update has recently been applied. This gives us the
// opportunity to do some post-update cleanup to make sure everything is in
// working order.
func firstRunAfterUpdate() {
	config.Cfg.FirstRunAfterUpdate = false
	config.ConfigToFile(config.LOCAL_EXTERNAL_PATH)
}

// waitForUpdates will continuously check for updated versions of the software
// and update to a newer version if found. Successive version checks will take
// place after a given number of seconds and compare the remote build number
// to the local build number to see if an update is required.
func waitForUpdates() error {
	for 1 == 1 {
		log.LogMessage(fmt.Sprintf("waiting for updates. sleeping %v seconds", config.Cfg.CheckInFrequencySeconds))
		time.Sleep(config.Cfg.CheckInFrequencySeconds * time.Second)
		local, localError := localVersion(config.Cfg.LocalVersionURI)
		remote, remoteError := remoteVersion(config.Cfg.RemoteVersionURI)

		if localError != nil {
			return localError
		} else if remoteError != nil {
			return remoteError
		}

		if remote > local {
			log.LogMessage(fmt.Sprintf("localVersion: %v", local))
			log.LogMessage(fmt.Sprintf("remoteVersion: %v", remote))
			log.LogMessage("Newer remote version available. Performing update.")
			doUpdate()
		}
	}
	return nil
}

// mine will kick off all the programs in order to start mining ethereum.
func mine() error {
	for 1 == 1 {
		log.LogMessage("mining...")
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

func doUpdate() error {
	log.LogMessage("updating...")
	return nil
}
