package updater

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/logger"
)

func TestMain(m *testing.M) {

	logErr := logger.StandardLogger("updater_test")
	if logErr != nil {
		fmt.Println(fmt.Sprintf("Could not initialize logger: %v", logErr))
		return
	}

	configErr := config.FromFile()
	if configErr != nil {
		fmt.Println(configErr)
		os.Exit(1)
	}

	result := m.Run()
	os.Exit(result)
}

func TestVersionCompare(t *testing.T) {

	update, updateErr := UpdateNecessary()
	if updateErr != nil {
		t.Error(updateErr)
	}

	fmt.Println(fmt.Sprintf("update necessary: %v", update))
}

func TestRun(t *testing.T) {

	config.Cfg.UpdateFrequencySeconds = 2
	Run()
	time.Sleep(time.Second * 6)

}
