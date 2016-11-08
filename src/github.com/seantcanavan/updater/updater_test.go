package updater

import (
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
)

func TestMain(m *testing.M) {
	cfgError := config.ConfigFromFile("updater_config.json")

	if cfgError != nil {
		fmt.Println("test init failure")
		fmt.Println(cfgError)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestVersionCompare(t *testing.T) {
	udr, udrError := NewUpdater()

	if udrError != nil {
		t.Error(udrError)
	}

	localVersion, localError := udr.localVersion(config.Cfg.LocalVersionURI)
	remoteVersion, remoteError := udr.remoteVersion(config.Cfg.RemoteVersionURI)

	if localError != nil {
		t.Error(localError)
	} else if remoteError != nil {
		t.Error(remoteError)
	}

	fmt.Println(fmt.Sprintf("localVersion: %v", localVersion))
	fmt.Println(fmt.Sprintf("remoteVersion: %v", remoteVersion))

	if localVersion > remoteVersion {
		fmt.Println("Your version is higher than the remote. Push your changes!")
	}

	if localVersion == remoteVersion {
		fmt.Println("Your version equals the remote version. Do some work!")
	}

	if localVersion < remoteVersion {
		fmt.Println("Your version is lower than the remote. Pull the latest code and build it!")
	}
}
