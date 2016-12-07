package updater

import (
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
)

var udr *Updater

func TestMain(m *testing.M) {

	configErr := config.FromFile()

	if configErr != nil {
		fmt.Println(configErr)
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
