package config

import (
	"testing"

	"github.com/seantcanavan/utils"
)

func TestSimplecfgigLoadFromFile(t *testing.T) {
	cfg, err := ConfigFromFile("config.json")

	if err != nil {
		t.Errorf("generic unmarshal error: %v", cfg)
	}

	if cfg.CheckInFrequencySeconds != 3600 {
		t.Errorf("cfg.CheckInFrequencySeconds did not unmarshal correctly: %v", cfg.CheckInFrequencySeconds)
	}

	if cfg.UpdateFrequencySeconds != 3600 {
		t.Errorf("cfg.UpdateFrequencySeconds did not unmarshal correctly: %v", cfg.UpdateFrequencySeconds)
	}

	if cfg.RemoteUpdateURI != "https://github.com/seantcanavan/anon-eth-net.git" {
		t.Errorf("cfg.RemoteUpdateURI did not unmarshal correctly: %v", cfg.RemoteUpdateURI)
	}

	if cfg.RemoteVersionURI != "https://raw.githubusercontent.com/seantcanavan/anon-eth-net/master/src/github.com/seantcanavan/main/version.no" {
		t.Errorf("cfg.RemoteVersionURI did not unmarshal correctly: %v", cfg.RemoteVersionURI)
	}

	if cfg.LocalVersionURI != "main/version.no" {
		t.Errorf("cfg.LocalVersionURI did not unmarshal correctly: %v", cfg.LocalVersionURI)
	}

	if cfg.LocalVersion != 0 {
		t.Errorf("cfg.LocalVersion did not unmarshal correctly: %v", cfg.LocalVersion)
	}

	if cfg.MineEther != false {
		t.Errorf("cfg.MineEther did not unmarshal correctly: %v", cfg.MineEther)
	}

	if cfg.GPUMine != false {
		t.Errorf("cfg.GPUMine did not unmarshal correctly: %v", cfg.MineEther)
	}

	if cfg.CPUMine != true {
		t.Errorf("cfg.CPUMine did not unmarshal correctly: %v", cfg.CPUMine)
	}

	if cfg.EtherWallet != "" {
		t.Errorf("cfg.EtherWallet did not unmarshal correctly: %v", cfg.EtherWallet)
	}

	fileLines, fileErr := utils.ReadLines(cfg.CheckInGmailCredentialsFile)
	if fileErr != nil {
		t.Errorf("issue reading in email credentials file")
	}

	if cfg.CheckInGmailAddress != fileLines[0] {
		t.Errorf("cfg.CheckInGmailAddress did not load correctly: %v", cfg.CheckInGmailAddress)
	}

	if cfg.CheckInGmailPassword != fileLines[1] {
		t.Errorf("cfg.CheckInGmailPassword did not load correctly: %v", cfg.CheckInGmailPassword)
	}

	if cfg.FirstRunAfterUpdate != false {
		t.Errorf("cfg.FirstRunAfterUpDate did not load correctly: %v", cfg.FirstRunAfterUpdate)
	}

	if cfg.InitialStartup != false {
		t.Error("cfg.InitialStartup did not load correctly: %v", cfg.InitialStartup)
	}
}
