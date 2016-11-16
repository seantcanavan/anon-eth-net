package config

import (
	"testing"

	"github.com/seantcanavan/utils"
)

func TestConfigFromFilePass(t *testing.T) {
	assetPath, assetErr := utils.AssetPath("config.json")
	if assetErr != nil {
		t.Error(assetErr)
	}

	err := FromFile(assetPath)

	if err != nil {
		t.Errorf("error loading config from file: %v", err)
	}

	gmailAsset, gmailAssetErr := utils.AssetPath(Cfg.CheckInGmailCredentialsFile)
	if gmailAssetErr != nil {
		t.Error(gmailAssetErr)
	}

	fileLines, fileErr := utils.ReadLines(gmailAsset)
	if fileErr != nil {
		t.Errorf("issue reading in email credentials file")
	}

	if Cfg.CheckInFrequencySeconds != 3600 {
		t.Errorf("Cfg.CheckInFrequencySeconds did not unmarshal correctly: %v", Cfg.CheckInFrequencySeconds)
	}

	if Cfg.UpdateFrequencySeconds != 3600 {
		t.Errorf("Cfg.UpdateFrequencySeconds did not unmarshal correctly: %v", Cfg.UpdateFrequencySeconds)
	}

	if Cfg.RemoteUpdateURI != "https://github.com/seantcanavan/anon-eth-net.git" {
		t.Errorf("Cfg.RemoteUpdateURI did not unmarshal correctly: %v", Cfg.RemoteUpdateURI)
	}

	if Cfg.RemoteVersionURI != "https://raw.githubusercontent.com/seantcanavan/anon-eth-net/master/src/github.com/seantcanavan/assets/version.no" {
		t.Errorf("Cfg.RemoteVersionURI did not unmarshal correctly: %v", Cfg.RemoteVersionURI)
	}

	if Cfg.LocalVersionURI != "version.no" {
		t.Errorf("Cfg.LocalVersionURI did not unmarshal correctly: %v", Cfg.LocalVersionURI)
	}

	if Cfg.LogVolatility != 3 {
		t.Errorf("Cfg.LogVolatility did not unmarshal correctly: %v", Cfg.LogVolatility)
	}

	if Cfg.DeviceName != "LG Smart Fridge" {
		t.Errorf("Cfg.DeviceName did not unmarshal correctly: %v", Cfg.DeviceName)
	}

	if Cfg.DeviceId != "519a2a15-afad-4c1a-94a3-71660c83504b" {
		t.Errorf("Cfg.DeviceId did not unmarshal correctly: %v", Cfg.DeviceId)
	}

	if Cfg.InitialStartup != false {
		t.Errorf("Cfg.InitialStartup did not unmarshal correctly: %v", Cfg.InitialStartup)
	}

	if Cfg.FirstRunAfterUpdate != false {
		t.Errorf("Cfg.FirstRunAfterUpDate did not unmarshal correctly: %v", Cfg.FirstRunAfterUpdate)
	}

	// these values are loaded dynamically from the CheckInGmailCredentialsFile
	if Cfg.CheckInGmailAddress != fileLines[0] {
		t.Errorf("Cfg.CheckInGmailAddress did not load correctly: %v", Cfg.CheckInGmailAddress)
	}

	if Cfg.CheckInGmailPassword != fileLines[1] {
		t.Errorf("Cfg.CheckInGmailPassword did not load correctly: %v", Cfg.CheckInGmailPassword)
	}

}
