package config

import (
	"testing"

	"github.com/seantcanavan/utils"
)

func TestSimpleConfigLoadFromFile(t *testing.T) {
	conf, err := ConfigFromFile("config.json")

	if err != nil {
		t.Errorf("generic unmarshal error: %v", conf)
	}

	if conf.CheckInFrequencySeconds != 3600 {
		t.Errorf("conf.CheckInFrequencySeconds did not unmarshal correctly: %v", conf.CheckInFrequencySeconds)
	}

	if conf.UpdateFrequencySeconds != 3600 {
		t.Errorf("conf.UpdateFrequencySeconds did not unmarshal correctly: %v", conf.UpdateFrequencySeconds)
	}

	if conf.RemoteUpdateURI != "https://github.com/seantcanavan/anon-eth-net.git" {
		t.Errorf("conf.RemoteUpdateURI did not unmarshal correctly: %v", conf.RemoteUpdateURI)
	}

	if conf.RemoteVersionURI != "https://raw.githubusercontent.com/seantcanavan/anon-eth-net/master/src/github.com/seantcanavan/main/version.no" {
		t.Errorf("conf.RemoteVersionURI did not unmarshal correctly: %v", conf.RemoteVersionURI)
	}

	if conf.LocalVersionURI != "main/version.no" {
		t.Errorf("conf.LocalVersionURI did not unmarshal correctly: %v", conf.LocalVersionURI)
	}

	if conf.LocalVersion != 0 {
		t.Errorf("conf.LocalVersion did not unmarshal correctly: %v", conf.LocalVersion)
	}

	if conf.MineEther != false {
		t.Errorf("conf.MineEther did not unmarshal correctly: %v", conf.MineEther)
	}

	if conf.GPUMine != false {
		t.Errorf("conf.GPUMine did not unmarshal correctly: %v", conf.MineEther)
	}

	if conf.CPUMine != true {
		t.Errorf("conf.CPUMine did not unmarshal correctly: %v", conf.CPUMine)
	}

	if conf.EtherWallet != "" {
		t.Errorf("conf.EtherWallet did not unmarshal correctly: %v", conf.EtherWallet)
	}

	fileLines, fileErr := utils.ReadLines(conf.CheckInGmailCredentialsFile)
	if fileErr != nil {
		t.Errorf("issue reading in email credentials file")
	}

	if conf.CheckInGmailAddress != fileLines[0] {
		t.Errorf("conf.CheckInGmailAddress did not load correctly: %v", conf.CheckInGmailAddress)
	}

	if conf.CheckInGmailPassword != fileLines[1] {
		t.Errorf("conf.CheckInGmailPassword did not load correctly: %v", conf.CheckInGmailPassword)
	}
}
