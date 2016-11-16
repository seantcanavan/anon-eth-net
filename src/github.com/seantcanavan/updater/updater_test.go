package updater

import (
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/utils"
)

func TestMain(m *testing.M) {
	assetPath, assetErr := utils.AssetPath("config.json")
	if assetErr != nil {
		fmt.Println(assetErr)
		return
	}

	cfgError := config.FromFile(assetPath)

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

	localAsset, assetErr := utils.AssetPath(config.Cfg.LocalVersionURI)
	if assetErr != nil {
		t.Error(assetErr)
	}

	localVersion, localError := udr.localVersion(localAsset)
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
