package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/nu7hatch/gouuid"
	"github.com/seantcanavan/config"
	"github.com/seantcanavan/loader"
	"github.com/seantcanavan/logger"
	"github.com/seantcanavan/updater"
)

// SeanLogger for the main package and any errors it encounters while executing
var lgr *logger.Logger

// Updater for the main package which will track local and remote versions of the code
var udr *updater.Updater

// Loader for the main package which will execute all of the third party processes
var ldr *loader.Loader

func main() {

	// check if the user is confused at first and is so, print the config.json
	// with code documentation included. and and then exit.
	if os.Args[1] == "h" || os.Args[1] == "help" || os.Args[1] == "?" {
		fmt.Println("Does not require any command line arguments. Refer to the default config/config.json file for all the parameters required for AEN to execute successfully.")
		fmt.Println(config.ConfigJSONParametersExplained())
		os.Exit(1)
	}

	// load the main config file from JSON which we'll use throughout execution.
	// we're exiting under an error condition since it's so early in execution.
	if _, err := os.Stat(os.Args[1]); err == nil {
		if configError := config.ConfigFromFile(config.LOCAL_EXTERNAL_PATH); configError != nil {
			fmt.Println(fmt.Sprintf("Could not successfully load config from file: %v", os.Args[1]))
			os.Exit(1)
		}
	} else {
		fmt.Println(fmt.Sprintf("Config file you passed in could not be found: %v", os.Args[1]))
		os.Exit(1)
	}

	// generate a SeanLogger instance with the predefined volatility value and
	// name it after the main_package to differentiate it from other packages
	mainLogger, loggerError := logger.FromVolatilityValue("main_package")
	if loggerError != nil {
		fmt.Println("Couldn't create logger... Exiting...")
		os.Exit(1)
	}

	// generate a loader instance with the appropriate JSON file based on the
	// system architecture
	var mainLoader *loader.Loader
	var loaderError error

	switch runtime.GOOS {
	case "windows":
		mainLoader, loaderError = loader.NewLoader("main_loader_windows.json")
	case "darwin":
		mainLoader, loaderError = loader.NewLoader("main_loader_darwin.json")
	case "linux":
		mainLoader, loaderError = loader.NewLoader("main_loader_linux.json")
	default:
		fmt.Println(fmt.Sprintf("Could not create loader for unsupported operating system: %v", runtime.GOOS))
		os.Exit(1)
	}

	if loaderError != nil {
		fmt.Println("Couldn't create loader... Exiting...")
		os.Exit(1)
	}

	mainUpdater, updaterError := updater.NewUpdater()
	if updaterError != nil {
		fmt.Println("Couldn't create updater... Exiting...")
		os.Exit(1)
	}

	// maintain a local logging reference for anything kicked off by the main package
	lgr = mainLogger
	// maintain a local loader reference for executing processes
	udr = mainUpdater
	// maintain a local updater reference for updating the main program
	ldr = mainLoader

	// if this is our first time ever starting up - run the initial config
	if config.Cfg.InitialStartup {
		err := initialStartup()
		if err != nil {
			lgr.LogMessage(err.Error())
		}
	}

	// if this is our first time starting up after an update - run the update config
	if config.Cfg.FirstRunAfterUpdate {
		err := firstRunAfterUpdate()
		if err != nil {
			lgr.LogMessage(err.Error())
		}
	}

	// kick off the watchdog loop that will regularly watch for and execute updates when available
	go func() {
		for 1 == 1 {
			if err := udr.Run(); err != nil {
				lgr.LogMessage("Error occurred in the Updater module: %v", err.Error())
			} else {
				lgr.LogMessage("Updater has exited gracefully. Well played.")
			}
		}
	}()

	go func() {
		for 1 == 1 {
			processes := ldr.StartAsynchronous()

			for _, resultProcess := range processes {
				//process results here
				logContents, logError := resultProcess.Lgr.CurrentLogContents()
				if logError == nil {
					lgr.LogMessage(string(logContents))
				}
			}
		}
	}()
}

// initialStartup will be executed only when this program is running for the
// first time on a new host.
func initialStartup() error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}

	// update the UUID if it doesn't exist
	config.Cfg.DeviceId = uuid.String()

	// we're finishing the first run!
	config.Cfg.InitialStartup = false

	// more stuff later!

	// push the UUID back to the file for next time
	return config.ConfigToFile(config.LOCAL_EXTERNAL_PATH)
}

// firstRunAfterUpdate will be executed only when this program is running for
// the first time after an update has recently been applied. This gives us the
// opportunity to do some post-update cleanup to make sure everything is in
// working order.
func firstRunAfterUpdate() error {
	config.Cfg.FirstRunAfterUpdate = false
	return config.ConfigToFile(config.LOCAL_EXTERNAL_PATH)
}
