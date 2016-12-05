package network

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/logger"
	"github.com/seantcanavan/utils"
)

type Network struct {
	endpoints map[string]string
	lgr       *logger.Logger
}

func NewNetwork() (*Network, error) {

	connectionAssetPath, assetPathErr := utils.AssetPath("connections.json")
	if assetPathErr != nil {
		return nil, assetPathErr

	}
	fileBytes, readErr := ioutil.ReadFile(connectionAssetPath)
	if readErr != nil {
		fmt.Println(fmt.Sprintf("Error reading in provided connections file from path: %v", connectionAssetPath))
		return nil, readErr
	}

	networkLogger, loggerError := logger.FromVolatilityValue("network_package")
	if loggerError != nil {
		return nil, loggerError
	}

	endpoints := make(map[string]string)

	jsonErr := json.Unmarshal(fileBytes, &endpoints)
	if jsonErr != nil {
		fmt.Println(fmt.Sprintf("Error unmarshalling the connections file into a map: %v", jsonErr))
		return nil, jsonErr
	}

	netw := &Network{}
	netw.lgr = networkLogger
	netw.endpoints = endpoints

	return netw, nil
}

// this was multi-threaded at first but it had reduced reliability when 20+
// network requests went out simultaneously. this will run in the background as
// it is so it won't really need to be multithreaded.
func (con *Network) IsInternetReachable() bool {

	numQueries := len(con.endpoints)
	errorCount := 0
	threshold := numQueries / 2

	for name, url := range con.endpoints {
		result, err := http.Get(url)
		if err != nil {
			con.lgr.LogMessage("Error querying internet endpoint: %v at: %v received: %v", name, url, err.Error())
			errorCount++
		} else {
			defer result.Body.Close()
			var bodyBuffer bytes.Buffer
			_, _ = io.Copy(&bodyBuffer, result.Body)
			con.lgr.LogMessage("Successfully queried internet endpoint: %v at: %v", name, url)
			// if we have a good ratio of errors to non-errors, we can afford to quit
			if !(errorCount > (numQueries / 2)) {
				break
			}
		}
	}

	con.lgr.LogMessage("Finished querying external APIs to test internet connectivity")
	con.lgr.LogMessage("Received %d errors back with a threshold of %d", errorCount, threshold)

	// if more than half the queries error out return false, else true
	return !(errorCount > (numQueries / 2))
}

func (con *Network) Run() {
	for 1 == 1 {
		interval := config.Cfg.NetQueryFrequencySeconds
		time.Sleep(time.Duration(interval) * time.Second)
		connected := con.IsInternetReachable()
		if !connected {
			//reboot machine
			con.lgr.LogMessage("Internet is unreachable. Rebooting the machine immediately.")
		} else {
			con.lgr.LogMessage("Internet is reachable. Sleeping for %d seconds before checking again", interval)
		}
	}
}
