package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/seantcanavan/utils"
)

const LOCAL_EXTERNAL_PATH = "../config/config.json"
const LOCAL_INTERNAL_PATH = "config.json"

var Cfg *Config

type Config struct {
	CheckInGmailCredentialsFile string        `json:"CheckInGmailCredentialsFile"` // (R) The email address where this program will report regular status updates to
	CheckInFrequencySeconds     time.Duration `json:"CheckInFrequencySeconds"`     // (R) The frequency with which this program will send status updates. In seconds.
	UpdateFrequencySeconds      int           `json:"UpdateFrequencySeconds"`      // (R) The frequency with which this program will attempt to update itself. In seconds.
	RemoteUpdateURI             string        `json:"RemoteUpdateURI"`             // (R) The remote location where new source code can be obtained from for this program.
	RemoteVersionURI            string        `json:"RemoteVersionURI"`            // (R) The remote URI where the latest version number of this program can be obtained from.
	LocalVersionURI             string        `json:"LocalVersionURI"`             // (N/A) The local URI where the current running version of this program can be obtained from.
	LocalVersion                uint64        `json:"LocalVersion"`                // (N/A) The local version of this program that is currently running.
	LogVolatility               int           `json:"LogVolatility"`               // (R) How quickly or slowly logs are pruned from the local disk. More volatility means logs last less time. Use 0 for most conservative logging strategy, 3 for least conservative.
	DeviceName                  string        `json:"DeviceName"`                  // (O) The canonical DeviceName for the machine currently executing this program.
	DeviceId                    string        `json:"DeviceId"`                    // (O) The unique ID for the machine currently executing this program.
	InitialStartup              bool          `json:"InitialStartup"`              // (N/A) Whether or not this is the first time that the program is starting.
	FirstRunAfterUpdate         bool          `json:"FirstRunAfterUpdate"`         // (N/A) Whether or not this is the first time that the program is running after an update has been executed.

	// You may manually set these values in your code if it remains private to you only. Otherwise you may configure the email credentials file instead for the values to be pulled from.
	CheckInGmailAddress  string // (O) the gmail address to send updates to and receive updates from. parsed from line 1 of CheckInEmailCredentialsFile
	CheckInGmailPassword string // (O) the password for the gmail account. parsed from line 2 of CheckInEmailCredentialsFile
}

func ConfigJSONParametersExplained() string {
	return `
	CheckInGmailCredentialsFile string        json:"CheckInGmailCredentialsFile" // (R) The email address where this program will report regular status updates to
	CheckInFrequencySeconds     time.Duration json:"CheckInFrequencySeconds"     // (R) The frequency with which this program will send status updates. In seconds.
	UpdateFrequencySeconds      int           json:"UpdateFrequencySeconds"      // (R) The frequency with which this program will attempt to update itself. In seconds.
	RemoteUpdateURI             string        json:"RemoteUpdateURI"             // (R) The remote location where new source code can be obtained from for this program.
	RemoteVersionURI            string        json:"RemoteVersionURI"            // (R) The remote URI where the latest version number of this program can be obtained from.
	LocalVersionURI             string        json:"LocalVersionURI"             // (N/A) The local URI where the current running version of this program can be obtained from.
	LocalVersion                uint64        json:"LocalVersion"                // (N/A) The local version of this program that is currently running.
	LogVolatility               int           json:"LogVolatility"               // (R) How quickly or slowly logs are pruned from the local disk. More volatility means logs last less time. Use 0 for most conservative logging strategy, 3 for least conservative.
	DeviceName                  string        json:"DeviceName"                  // (O) The canonical DeviceName for the machine currently executing this program.
	DeviceId                    string        json:"DeviceId"                    // (O) The unique ID for the machine currently executing this program.
	InitialStartup              bool          json:"InitialStartup"              // (N/A) Whether or not this is the first time that the program is starting.
	FirstRunAfterUpdate         bool          json:"FirstRunAfterUpdate"         // (N/A) Whether or not this is the first time that the program is running after an update has been executed.

	// You may manually set these values in your code if it remains private to you only. Otherwise you may configure the email credentials file instead for the values to be pulled from.
	CheckInGmailAddress  string // (O) the gmail address to send updates to and receive updates from. parsed from line 1 of CheckInEmailCredentialsFile
	CheckInGmailPassword string // (O) the password for the gmail account. parsed from line 2 of CheckInEmailCredentialsFile
`
}

// ConfigFromFile will generate a config struct from the local standard config
//path which should point to a valid JSON file containing key and value pairs
//for all the fields in the config struct that you wish to use throughout the
// code. A sample config file is provided, config.json, which tells you all of
// the possible fields to include as well as hopefully enough of a clue as to
// what their use is.
func ConfigFromFile(filePath string) error {

	bytes, loadErr := ioutil.ReadFile(filePath)
	if loadErr != nil {
		return loadErr
	}

	var newConfig Config

	jsonErr := json.Unmarshal(bytes, &newConfig)
	if jsonErr != nil {
		fmt.Println(fmt.Sprintf("Error unmarshalling the config file into a struct: %v", jsonErr))
		return jsonErr
	}

	fileLines, emailErr := utils.ReadLines(newConfig.CheckInGmailCredentialsFile)
	if emailErr != nil {
		fmt.Println(fmt.Sprintf("Email credentials file issue: %v", emailErr))
		return emailErr
	}

	newConfig.CheckInGmailAddress = fileLines[0]
	newConfig.CheckInGmailPassword = fileLines[1]

	fmt.Println("Loaded config from file:")
	fmt.Println(fmt.Sprintf("%+v\n", newConfig))

	Cfg = &newConfig
	return nil
}

// ConfigToFile will overwrite the local config file with the given instance of
// config. This is useful if the program's variables change dynamically because
// they can be saved permanently to disk for next startup / shutdown / restart.
func ConfigToFile(filePath string) error {
	bytes, marshalError := json.Marshal(Cfg)
	if marshalError != nil {
		return marshalError
	}

	writeError := ioutil.WriteFile(filePath, bytes, 0644)
	if writeError != nil {
		return writeError
	}
	return nil
}
