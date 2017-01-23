package network

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/seantcanavan/anon-eth-net/config"
	"github.com/seantcanavan/anon-eth-net/loader"
	"github.com/seantcanavan/anon-eth-net/logger"
	"github.com/seantcanavan/anon-eth-net/utils"
)

type Network struct {
	endpoints map[string]string
}

func NewNetwork() (*Network, error) {

	connectionAssetPath, assetPathErr := utils.AssetPath("connections.json")
	if assetPathErr != nil {
		return nil, assetPathErr
	}

	logger.Lgr.LogMessage("Successfully loaded connections from file: %v", connectionAssetPath)

	fileBytes, readErr := ioutil.ReadFile(connectionAssetPath)
	if readErr != nil {
		return nil, readErr
	}

	logger.Lgr.LogMessage("Successfully read connections file: %v", connectionAssetPath)

	loadedEndpoints := make(map[string]string)

	jsonErr := json.Unmarshal(fileBytes, &loadedEndpoints)
	if jsonErr != nil {
		return nil, jsonErr
	}

	logger.Lgr.LogMessage("Successfully unmarshalled JSON endpoint data into map: %+v", loadedEndpoints)

	netw := &Network{endpoints: loadedEndpoints}
	return netw, nil
}

// this was multi-threaded at first but it had reduced reliability when 20+
// network requests went out simultaneously. this will run in the background as
// it is so it won't really need to be multithreaded.
func (con *Network) IsInternetReachable() bool {

	numQueries := len(con.endpoints)
	errorCount := 0
	threshold := numQueries / 2

	logger.Lgr.LogMessage("Checking internet connectivity with threshold: %d", threshold)

	for name, url := range con.endpoints {

		result, err := http.Get(url)
		if err != nil {
			logger.Lgr.LogMessage("Error querying internet endpoint: %v at: %v received: %v", name, url, err.Error())
			errorCount++
		} else {
			defer result.Body.Close()
			var bodyBuffer bytes.Buffer
			_, _ = io.Copy(&bodyBuffer, result.Body)
			logger.Lgr.LogMessage("Successfully queried internet endpoint: %v at: %v", name, url)
			// if we have a good ratio of errors to non-errors, we can afford to quit
			if !(errorCount > (numQueries / 2)) {
				break
			}
		}
	}

	logger.Lgr.LogMessage("Finished querying external APIs to test internet connectivity")
	logger.Lgr.LogMessage("Received %d errors back with a threshold of %d", errorCount, threshold)

	// if more than half the queries error out return false, else true
	return !(errorCount > (numQueries / 2))
}

// Run will ensure that the network manager is always running and verifying at
// a set interval that this machine can speak to others via the internet
func (con *Network) Run() {

	go func() {

		for 1 == 1 {

			interval := config.Cfg.NetQueryFrequencySeconds

			logger.Lgr.LogMessage("Network manager will sleep for %d seconds before querying the internet", interval)

			time.Sleep(time.Duration(interval) * time.Second)

			connected := con.IsInternetReachable()

			if !connected {
				//reboot machine
				logger.Lgr.LogMessage("Internet is unreachable. Rebooting the machine immediately.")
				rebootAssetPath, assetErr := utils.SysAssetPath("reboot_loader.json")
				if assetErr != nil {
					logger.Lgr.LogMessage("Unable to load the reboot loader asset: %v", assetErr.Error())
				} else {
					rebootLoader, loaderErr := loader.NewLoader(rebootAssetPath)
					if loaderErr != nil {
						logger.Lgr.LogMessage("Unable to instantiate new loader from asset: %v with error: %v", rebootAssetPath, loaderErr.Error())
					} else {
						_ = rebootLoader.StartSynchronous()
					}
				}
			} else {
				logger.Lgr.LogMessage("Internet is reachable. Sleeping for %d seconds before checking again", interval)
			}
		}

	}()

}
