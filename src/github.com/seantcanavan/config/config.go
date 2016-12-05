// The config package represents a static instance of the Config struct which is
// globally accessible once properly initialized.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/seantcanavan/utils"
	"github.com/nu7hatch/gouuid"
)

var Cfg *Config

// Config represents a set of public configuration values used throughout the
// program to help anon-eth-net execute in a manner that the user expects. All
// values can be configured via the config.json file and changes to the config
// can be exported to the same file as well. This helps the program change state
// over time as well.
// (R) means the value is required to be provided by the user
// (O) means the value is optional and not required to be provided by the user
// (D) means the value is default value already set and should only be
// changed after careful consideration.
type Config struct {
	CheckInGmailAddress  		string 		  `json:"CheckInGmailAddress"`	       // (R) the gmail address to send updates to and receive updates from. parsed from line 1 of CheckInEmailCredentialsFile
	CheckInGmailPassword 		string 		  `json:"CheckInGmailPassword"`	       // (R) the password for the gmail account. parsed from line 2 of CheckInEmailCredentialsFile
	CheckInFrequencySeconds     time.Duration `json:"CheckInFrequencySeconds"`     // (R) The frequency with which this program will send status updates. In seconds.
	NetQueryFrequencySeconds    time.Duration `json:"NetQueryFrequencySeconds"`    // (R) The frequency with which this program will attempt to connect to the outside world to verify internet connectivity
	LogVolatility               int           `json:"LogVolatility"`               // (R) How quickly or slowly logs are pruned from the local disk. More volatility means logs last less time. Use 0 for most conservative logging strategy, 3 for least conservative.

	DeviceName                  string        `json:"DeviceName"`                  // (O) The canonical DeviceName for the machine currently executing this program.
	DeviceId                    string        `json:"DeviceId"`                    // (O) The unique ID for the machine currently executing this program.

	InitialStartup              string        `json:"InitialStartup"`              // (D) Whether or not this is the first time that the program is starting.
	FirstRunAfterUpdate         string        `json:"FirstRunAfterUpdate"`         // (D) Whether or not this is the first time that the program is running after an update has been executed.
	UpdateFrequencySeconds      int           `json:"UpdateFrequencySeconds"`      // (D) The frequency with which this program will attempt to update itself. In seconds.
	RemoteUpdateURI             string        `json:"RemoteUpdateURI"`             // (D) The remote location where new source code can be obtained from for this program.
	RemoteVersionURI            string        `json:"RemoteVersionURI"`            // (D) The remote URI where the latest version number of this program can be obtained from.
	LocalVersion                uint64        `json:"LocalVersion"`                // (D) The local version of this program that is currently running.

}

// COnfigJSONParametersExplained() returns a nicely formatted string which
// describes all the public variables available to the user for configuration.
func ConfigJSONParametersExplained() string {
	return `
	CheckInGmailAddress  		string 		  json:"CheckInGmailAddress"	     // (R) the gmail address to send updates to and receive updates from. parsed from line 1 of CheckInEmailCredentialsFile
	CheckInGmailPassword 		string 		  json:"CheckInGmailPassword"	     // (R) the password for the gmail account. parsed from line 2 of CheckInEmailCredentialsFile
	CheckInFrequencySeconds     time.Duration json:"CheckInFrequencySeconds"     // (R) The frequency with which this program will send status updates. In seconds.
	NetQueryFrequencySeconds    time.Duration json:"NetQueryFrequencySeconds"    // (R) The frequency with which this program will attempt to connect to the outside world to verify internet connectivity
	LogVolatility               int           json:"LogVolatility"               // (R) How quickly or slowly logs are pruned from the local disk. More volatility means logs last less time. Use 0 for most conservative logging strategy, 3 for least conservative.

	DeviceName                  string        json:"DeviceName"                  // (O) The canonical DeviceName for the machine currently executing this program.
	DeviceId                    string        json:"DeviceId"                    // (O) The unique ID for the machine currently executing this program.

	InitialStartup              string        json:"InitialStartup"              // (D) Whether or not this is the first time that the program is starting.
	FirstRunAfterUpdate         string        json:"FirstRunAfterUpdate"         // (D) Whether or not this is the first time that the program is running after an update has been executed.
	UpdateFrequencySeconds      int           json:"UpdateFrequencySeconds"      // (D) The frequency with which this program will attempt to update itself. In seconds.
	RemoteUpdateURI             string        json:"RemoteUpdateURI"             // (D) The remote location where new source code can be obtained from for this program.
	RemoteVersionURI            string        json:"RemoteVersionURI"            // (D) The remote URI where the latest version number of this program can be obtained from.
	LocalVersion                uint64        json:"LocalVersion"                // (D) The local version of this program that is currently running.
`
}

// FromFile will generate a config struct from the local standard config file
// which is located inside of the assets folder as 'config.json'. It will be
// fully configured based off of the values in the json.
func FromFile() error {

	configAssetPath, assetErr := utils.AssetPath("config.json")
	if assetErr != nil {
		return assetErr
	}

	// read in the pre-existing config file
	bytes, loadErr := ioutil.ReadFile(configAssetPath)
	if loadErr != nil {
		fmt.Println(fmt.Sprintf("Error reading in provided config file from path: %v", configAssetPath))
		return loadErr
	}

	newConfig := &Config{}

	// unmarshal the JSON directly into a config struct instance
	jsonErr := json.Unmarshal(bytes, &newConfig)
	if jsonErr != nil {
		fmt.Println(fmt.Sprintf("Error unmarshalling the config file into a struct: %v", jsonErr))
		return jsonErr
	}

	// check if a manual email login file was provided to secretly override the defaults
	emailAssetPath, emailAssetErr := utils.AssetPath("emaillogin.txt")
	if emailAssetErr == nil {
		fileLines, readErr := utils.ReadLines(emailAssetPath)
		if readErr != nil {
			return readErr
		}
		newConfig.CheckInGmailAddress = fileLines[0]
		newConfig.CheckInGmailPassword = fileLines[1]
	}

	// verify all the required values are correctly setup by the user
	if newConfig.CheckInGmailAddress == "" {
		return errors.New("Cannot use empty gmail address when starting up. Please update the config.json asset with an appropriate value and restart.")
	}

	if newConfig.CheckInGmailPassword == "" {
		return errors.New("Cannot use empty email password when starting up. Please update the config.json asset with an appropriate value and restart.")
	}

	if newConfig.CheckInFrequencySeconds == 0 {
		return errors.New("Cannot use an empty or zero value for check in frequency . Please update the config.json asset with an appropriate value and restart.")
	}

	if newConfig.NetQueryFrequencySeconds == 0 {
		return errors.New("Cannot use an empty or zero value for internet query frequency. Please update the config.json asset with an appropriate value and restart.")
	}

	// verify all the optional values are correctly set to a default, if necessary
	if newConfig.DeviceName == "" {
		randInt := rand.Int()
		newConfig.DeviceName = "device_" + strconv.Itoa(randInt)
	}

	if newConfig.DeviceId == "" {
		// if the DeviceId hasn't been set by the user - let's give them a nice UUID
		uuid, err := uuid.NewV4()
		if err != nil {
			return err
		}
		// update the UUID if it doesn't exist
		newConfig.DeviceId = uuid.String()
	}

	if newConfig.InitialStartup == "" {
		newConfig.InitialStartup = "yes"
	}

	if newConfig.FirstRunAfterUpdate == "" {
		newConfig.FirstRunAfterUpdate = "no"
	}

	if newConfig.UpdateFrequencySeconds == 0 {
		newConfig.UpdateFrequencySeconds = 3600
	}

	if newConfig.RemoteUpdateURI == "" {
		newConfig.RemoteUpdateURI = "https://github.com/seantcanavan/anon-eth-net.git"
	}

	if newConfig.RemoteVersionURI == "" {
		newConfig.RemoteVersionURI = "https://raw.githubusercontent.com/seantcanavan/anon-eth-net/master/src/github.com/seantcanavan/assets/version.no"
	}

	// load the local version number from the local asset
	localVersionAsset, assetErr := utils.AssetPath("version.no")

	bytes, err := ioutil.ReadFile(localVersionAsset)
	if err != nil {
		return err
	}

	s := string(bytes)
	s = strings.Trim(s, "\n")
	localVersion, castError := strconv.ParseUint(s, 10, 64)
	if castError != nil {
		return castError
	}

	newConfig.LocalVersion = localVersion

	fmt.Println("Loaded config from file:")
	fmt.Println(fmt.Sprintf("%+v\n", newConfig))

	Cfg = newConfig
	return nil
}

// ToFile will save the current instance of config to the local standard config
// file which is located inside of the assets folder as 'config.json'. This will
// help preserver changes to the configuration between settings.
func ToFile() error {

	configAssetPath, assetErr := utils.AssetPath("config.json")
	if assetErr != nil {
		return assetErr
	}

	bytes, marshalError := json.Marshal(Cfg)
	if marshalError != nil {
		return marshalError
	}

	writeError := ioutil.WriteFile(configAssetPath, bytes, 0644)
	if writeError != nil {
		return writeError
	}
	return nil
}
