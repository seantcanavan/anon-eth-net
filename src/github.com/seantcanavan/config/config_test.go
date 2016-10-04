package config

import (
	"testing"
)

func TestSimpleConfigLoadFromFile(t *testing.T) {
	conf, err := GetConfigFromFile("config.json")

	if err != nil {
		t.Errorf("generic unmarshal error: %v", conf)
	}

	if conf.CheckInEmail != "seantcanavan@gmail.com" {
		t.Error("conf.CheckInEmail did not unmarshal correctly")
	}

	if conf.CheckInFrequencySeconds != 3600 {
		t.Error("conf.CheckInFrequencySeconds did not unmarshal correctly")
	}

	if conf.UpdateFrequencySeconds != 3600 {
		t.Error("conf.UpdateFrequencySeconds did not unmarshal correctly")
	}

	if conf.RemoteUpdateURI != "https://github.com/seantcanavan/anon-eth-net.git" {
		t.Error("conf.RemoteUpdateURI did not unmarshal correctly")
	}

	if conf.RemoteVersionURI != "https://raw.githubusercontent.com/seantcanavan/anon-eth-net/master/src/github.com/seantcanavan/main/version.no" {
		t.Error("conf.RemoteVersionURI did not unmarshal correctly")
	}

	if conf.LocalVersionURI != "main/version.no" {
		t.Error("conf.LocalVersionURI did not unmarshal correctly")
	}

	if conf.LocalVersion != 0 {
		t.Error("conf.LocalVersion did not unmarshal correctly")
	}

	if conf.MineEther != false {
		t.Error("conf.MineEther did not unmarshal correctly")
	}

	if conf.GPUMine != false {
		t.Error("conf.GPUMine did not unmarshal correctly")
	}

	if conf.CPUMine != true {
		t.Error("conf.CPUMine did not unmarshal correctly")
	}

	if conf.EtherWallet != "" {
		t.Error("conf.EtherWallet did not unmarshal correctly")
	}
}
