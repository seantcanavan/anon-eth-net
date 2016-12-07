package updater

import (
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/utils"
)

var udr *Updater

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

	updater, udrError := NewUpdater()
	if udrError != nil {
		fmt.Println(udrError)
		return
	}

	udr = updater

	result := m.Run()
	os.Exit(result)
}

func TestVersionCompare(t *testing.T) {

	update, updateErr := udr.UpdateNecessary()
	if updateErr != nil {
		t.Error(updateErr)
	}

	fmt.Println(fmt.Sprintf("update necessary: %v", update))
}
