package config

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	configErr := FromFile()
	if configErr != nil {
		fmt.Println(fmt.Sprintf("Could not initialize config: %v", configErr))
		return
	}

	result := m.Run()
	os.Exit(result)
}

func TestConfigFromFilePass(t *testing.T) {

	// ---------- verify required values unmarshalled correctly ----------
	if Cfg.CheckInGmailAddress == "" {
		t.Errorf("Cfg.CheckInGmailAddress did not unmarshal correctly: %v", Cfg.CheckInGmailAddress)
	}

	if Cfg.CheckInGmailPassword == "" {
		t.Errorf("Cfg.CheckInGmailPassword did not unmarshal correctly: %v", Cfg.CheckInGmailPassword)
	}

	if Cfg.CheckInFrequencySeconds != 3600 {
		t.Errorf("Cfg.CheckInFrequencySeconds did not unmarshal correctly: %v", Cfg.CheckInFrequencySeconds)
	}

	if Cfg.NetQueryFrequencySeconds != 3600 {
		t.Errorf("Cfg.NetQueryFrequencySeconds did not unmarshal correctly: %v", Cfg.NetQueryFrequencySeconds)
	}

	if Cfg.LogVolatility != 3 {
		t.Errorf("Cfg.LogVolatility did not unmarshal correctly: %v", Cfg.LogVolatility)
	}

	// ---------- verify optional values unmarshalled correctly ----------
	if Cfg.DeviceName != "My Little Raspberry Pi" {
		t.Errorf("Cfg.DeviceName did not unmarshal correctly: %v", Cfg.DeviceName)
	}

	if Cfg.DeviceId != "aa3b2cd6-7a95-4b2a-af33-7eb953d730a9" {
		t.Errorf("Cfg.DeviceId did not unmarshal correctly: %v", Cfg.DeviceId)
	}

	// ---------- verify default values unmarshalled correctly ----------
	if Cfg.InitialStartup != "yes" {
		t.Errorf("Cfg.InitialStartup did not unmarshal correctly: %v", Cfg.InitialStartup)
	}

	if Cfg.FirstRunAfterUpdate != "no" {
		t.Errorf("Cfg.FirstRunAfterUpDate did not unmarshal correctly: %v", Cfg.FirstRunAfterUpdate)
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
}
