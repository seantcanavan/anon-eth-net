package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Config struct {
	CheckInEmail            string        `json:CheckInEmail`            // The email address where this program will report regular status updates to
	CheckInFrequencySeconds time.Duration `json:CheckInFrequencySeconds` // The frequency with which this program will send status updates. In seconds.
	UpdateFrequencySeconds  time.Duration `json:UpdateFrequencySeconds`  // The frequency with which this program will attempt to update itself. In seconds.
	RemoteUpdateURI         string        `json:RemoteUpdateURI`         // The remote location where new source code can be obtained from for this program.
	RemoteVersionURI        string        `json:RemoteVersionURI`        // The remote URI where the latest version number of this program can be obtained from.
	LocalVersionURI         string        `json:LocalVersionURI`         // The local URI where the current running version of this program can be obtained from.
	LocalVersion            uint64        `json:LocalVersion`            // The local version of this program that is currently running.
	EtherWallet             string        `json:EtherWallet`             // The wallet address of whatever Ethereum wallet you want your mining contributions to go to.
	MineEther               bool          `json:MineEther`               // Whether or not this zombie should mine ether.
	GPUMine                 bool          `json:GPUMine`                 // If this zombie is mining ether, should it GPU mine?
	CPUMine                 bool          `json:CPUMine`                 // If this zombie is mining ether, should it CPU mine?
}

func GetConfigFromFile(fileName string) (*Config, error) {
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var conf Config

	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, err
	}

	fmt.Println("Loaded config from file:")
	fmt.Printf("%+v\n", conf)
	return &conf, nil
}
