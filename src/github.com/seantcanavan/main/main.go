package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/loader"
	"github.com/seantcanavan/logger"
	"github.com/seantcanavan/network"
	"github.com/seantcanavan/updater"
)

// Logger for the main package and any errors it encounters while executing
var lgr *logger.Logger

// Updater for the main package which will track local and remote versions of the code
var udr *updater.Updater

// Loader for the main package which will execute all of the third party processes
var ldr *loader.Loader

// Connection for the main package which will constantly monitor the outgoing internet connection
var net *network.Network

func main() {

	// check if the user is confused at first and is so, print the config.json
	// with code documentation included. and and then exit.
	if os.Args[1] == "h" || os.Args[1] == "help" || os.Args[1] == "?" {
		fmt.Println("Does not require any command line arguments. Refer to the default assets/config.json file for all the parameters required for AEN to execute successfully.")
		fmt.Println(config.ConfigJSONParametersExplained())
		os.Exit(1)
	}

	// load the main config file from JSON which we'll use throughout execution.
	// we're exiting under an error condition since it's so early in execution.
	configErr := config.FromFile()
	if configErr != nil {
		fmt.Println(fmt.Sprintf("Could not successfully load config. Received error %v. Make sure the JSON is well formed and the values are correct for each variable. Revert to the standard config on github if this problem persists.", configErr))
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
		fmt.Println(fmt.Sprintf("Could not create loader for unsupported operating system: %v. Please choose from one of the selected supported operating systems to continue. Refer to the README.md for the list.", runtime.GOOS))
		os.Exit(1)
	}

	if loaderError != nil {
		fmt.Println("Couldn't create the loader for executing external processes... Exiting...")
		os.Exit(1)
	}

	// generate a Logger instance with the predefined volatility value and
	// name it after the main_package to differentiate it from other packages
	mainLogger, loggerError := logger.FromVolatilityValue("main_package")
	if loggerError != nil {
		fmt.Println("Couldn't create the logger for logging local activity to disk... Exiting...")
		os.Exit(1)
	}

	mainUpdater, updaterError := updater.NewUpdater()
	if updaterError != nil {
		fmt.Println("Couldn't create the automatic updater for AEN... Exiting...")
		os.Exit(1)
	}

	mainNetwork, networkErr := network.NewNetwork()
	if networkErr != nil {
		fmt.Println("Couldn't create the network monitor... Exiting...")
		os.Exit(1)
	}

	// maintain a local logging reference for anything kicked off by the main package
	lgr = mainLogger
	// maintain a local loader reference for executing processes
	udr = mainUpdater
	// maintain a local updater reference for updating the main program
	ldr = mainLoader
	// maintain a local network reference for checking internet connection
	net = mainNetwork

	// if this is our first time ever starting up - run the initial config
	if config.Cfg.InitialStartup == "yes" {
		err := initialStartup()
		if err != nil {
			lgr.LogMessage(err.Error())
		}
	}

	// if this is our first time starting up after an update - run the update config
	if config.Cfg.FirstRunAfterUpdate == "yes" {
		err := firstRunAfterUpdate()
		if err != nil {
			lgr.LogMessage(err.Error())
		}
	}

	// kick off the updater loop
	go func() {
		lgr.LogMessage("Initializing the updater")
		for 1 == 1 {
			if err := udr.Run(); err != nil {
				lgr.LogMessage("Error occurred in the Updater module: %v", err.Error())
			} else {
				lgr.LogMessage("Updater has exited gracefully. Well played.")
			}
		}
	}()

	// kick off the process loader loop that will execute things like miners
	go func() {
		lgr.LogMessage("Initializing the loader")
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

	// kick off the network monitor loop to monitor internet connectivity
	go func() {
		lgr.LogMessage("Initializing the network monitor")
		net.Run()
	}()

	lgr.LogMessage("Backing up the latest config changes before exiting")
	config.ToFile()
}

// initialStartup will be executed only when this program is running for the
// first time on a new host.
func initialStartup() error {
	// more stuff later!

	// we're finishing the first run!
	config.Cfg.InitialStartup = "no"

	return config.ToFile()
}

// firstRunAfterUpdate will be executed only when this program is running for
// the first time after an update has recently been applied. This gives us the
// opportunity to do some post-update cleanup to make sure everything is in
// working order.
func firstRunAfterUpdate() error {
	// more stuff later!

	//we're finishing the run after an update
	config.Cfg.FirstRunAfterUpdate = "no"

	return config.ToFile()
}
