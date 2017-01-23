package main

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/loader"
	"github.com/seantcanavan/logger"
	"github.com/seantcanavan/network"
	"github.com/seantcanavan/profiler"
	"github.com/seantcanavan/rest"
	"github.com/seantcanavan/updater"
	"github.com/seantcanavan/utils"
)

func main() {

	//------------------ CHECK FOR COMMAND LINE HELP ARGUMENTS ------------------
	if len(os.Args) > 1 {
		fmt.Println("Does not require any command line arguments. Refer to the default ./assets/config.json file for all the parameters required for anon-eth-net to execute successfully.")
		fmt.Println(config.ConfigJSONParametersExplained())
		os.Exit(1)
	}

	//------------------ GENERATE THE LOGGING FILE FOR THE MAIN PACKAGE ------------------
	loggerErr := logger.StandardLogger("main_package")
	if loggerErr != nil {
		fmt.Println(loggerErr)
		fmt.Println("Couldn't create the logger for logging local activity to disk... Exiting...")
		os.Exit(1)
	}

	//------------------ LOAD THE CONFIG.JSON ASSET AND UNMARSHAL THE VALUES ------------------
	configErr := config.FromFile()
	if configErr != nil {
		fmt.Println(fmt.Sprintf("Could not successfully load config. Received error %v. Make sure the JSON is well formed and the values are correct for each variable. Revert to the standard config on github if this problem persists.", configErr))
		os.Exit(1)
	}

	//------------------ CREATE LOADER INSTANCE TO RUN PROCESSES LOCALLY BASED ON GOOS ------------------
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

	//------------------ CREATE REST INSTANCE TO ENABLE COMMUNIATION VIA REST ------------------
	certPath, certPathErr := utils.AssetPath("server.cert")
	if certPathErr != nil {
		fmt.Println(certPathErr)
		return
	}

	certValue, certReadErr := ioutil.ReadFile(certPath)
	if certReadErr != nil {
		fmt.Println(certReadErr)
		return
	}

	mainRest, restErr := rest.NewRestHandler()
	if restErr != nil {
		fmt.Println(restErr)
		return
	}

	rootAuthorities := x509.NewCertPool()
	if ok := rootAuthorities.AppendCertsFromPEM([]byte(certValue)); !ok {
		fmt.Println("Unable to append certificate to set of root certificate authorities")
		return
	}

	//------------------ IF THIS IS OUR FIRST TIME STARTING UP EVER, TAKE APPROPRIATE ACTIONS ------------------
	if config.Cfg.InitialStartup == "yes" {
		err := initialStartup()
		if err != nil {
			logger.Lgr.LogMessage(err.Error())
		}
	}

	//------------------ IF THIS IS OUR FIRST TIME STARTING UP AFTER AN UPDATE, TAKE APPROPRIATE ACTIONS ------------------
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
	mainLoader.Run()

	// kick off the network monitor loop to monitor internet connectivity
	logger.Lgr.LogMessage("Initializing the network monitor")
	mainNetwork.Run()

	// kick off the REST endpoints
	logger.Lgr.LogMessage("Initializing the REST interface")
	mainRest.StartupRestServer()

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
