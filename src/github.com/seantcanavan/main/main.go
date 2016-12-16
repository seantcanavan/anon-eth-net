package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/loader"
	"github.com/seantcanavan/logger"
	"github.com/seantcanavan/network"
	"github.com/seantcanavan/profiler"
	"github.com/seantcanavan/updater"
	"github.com/seantcanavan/utils"
)

// Loader for the main package which will execute all of the third party processes
var ldr *loader.Loader

// Connection for the main package which will constantly monitor the outgoing internet connection
var net *network.Network

func main() {

	// check if the user is confused at first and is so, print the config.json
	// with code documentation included. and and then exit.
	if len(os.Args) > 1 {
		if os.Args[1] == "h" || os.Args[1] == "help" || os.Args[1] == "?" {
			fmt.Println("Does not require any command line arguments. Refer to the default ./assets/config.json file for all the parameters required for anon-eth-net to execute successfully.")
			fmt.Println(config.ConfigJSONParametersExplained())
			os.Exit(1)
		}
	}

	// generate a Logger instance with the predefined volatility value and
	// name it after the main_package to differentiate it from other packages
	loggerErr := logger.StandardLogger("main_package")
	if loggerErr != nil {
		fmt.Println(loggerErr)
		fmt.Println("Couldn't create the logger for logging local activity to disk... Exiting...")
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
	var loaderErr error

	switch runtime.GOOS {
	case "windows", "darwin", "linux":
		loaderAssetPath, assetErr := utils.SysAssetPath("main_loader.json")
		if assetErr != nil {
			fmt.Println(assetErr)
			fmt.Println("Could not successfully load main_loader.json... Exiting...")
			os.Exit(1)
		}
		mainLoader, loaderErr = loader.NewLoader(loaderAssetPath)
	default:
		fmt.Println(fmt.Sprintf("Could not create loader for unsupported operating system: %v. Please choose from one of the selected supported operating systems to continue. Refer to the README.md for the list.", runtime.GOOS))
		os.Exit(1)
	}

	if loaderErr != nil {
		fmt.Println(loaderErr)
		fmt.Println("Couldn't create the loader for executing external processes... Exiting...")
		os.Exit(1)
	}

	mainNetwork, networkErr := network.NewNetwork()
	if networkErr != nil {
		fmt.Println(networkErr)
		fmt.Println("Couldn't create the network monitor... Exiting...")
		os.Exit(1)
	}

	// maintain a local updater reference for updating the main program
	ldr = mainLoader
	// maintain a local network reference for checking internet connection
	net = mainNetwork

	// if this is our first time ever starting up - run the initial config
	if config.Cfg.InitialStartup == "yes" {
		err := initialStartup()
		if err != nil {
			logger.Lgr.LogMessage(err.Error())
		}
	}

	// if this is our first time starting up after an update - run the update config
	if config.Cfg.FirstRunAfterUpdate == "yes" {
		err := firstRunAfterUpdate()
		if err != nil {
			logger.Lgr.LogMessage(err.Error())
		}
	}

	// kick off the profiler loop
	logger.Lgr.LogMessage("Initializing the profiler")
	profiler.Run()

	// kick off the updater loop
	logger.Lgr.LogMessage("Initializing the updater")
	updater.Run()

	// kick off the process loader loop that will execute things like miners
	logger.Lgr.LogMessage("Initializing the loader")
	ldr.Run()

	// kick off the network monitor loop to monitor internet connectivity
	logger.Lgr.LogMessage("Initializing the network monitor")
	net.Run()


	// create a channel to listen to type os.Signal on with depth = 1
	sigs := make(chan os.Signal, 1)
	// create a channel to listen to type bool on with depth = 1
	done := make(chan bool, 1)

	// redirect the signals SIGINT and SIGTERM to the channel 'sigs'
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// wait in a separate go routine for the SIGINT and SIGTERM signals
	go func() {
		// when the OS sends SIGINT or SIGTERM to us, save the signal to value to 'signal'
		signal := <-sigs
		logger.Lgr.LogMessage("Received interrupting signal: %v", signal)
		// push 'true' to the 'done' channel when we've successfully received SIGINT or SIGTERM
		done <- true
	}()

	logger.Lgr.LogMessage("Executing... Press CTRL+C to exit. Browse local log files to keep an eye on each individual component.")
	// block until we receive SIGINT or SIGTERM and 'true' is pushed down the 'done' pipe
	<-done
	logger.Lgr.LogMessage("Clean exit after a CTRL+C interrupt.")
	logger.Lgr.LogMessage("Backing up the latest config changes before exiting")
	config.ToFile()
	logger.Lgr.LogMessage("Fin")
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
