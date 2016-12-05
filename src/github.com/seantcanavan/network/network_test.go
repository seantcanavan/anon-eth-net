package network

import (
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/utils"
)

var netw *Network

func TestMain(m *testing.M) {

	configPath, configPathErr := utils.AssetPath("config.json")
	if configPathErr != nil {
		fmt.Println(configPathErr)
		return
	}

	cfgErr := config.FromFile(configPath)
	if cfgErr != nil {
		fmt.Println(cfgErr)
		return
	}

	network, netwErr := NewNetwork()
	if netwErr != nil {
		fmt.Println(netwErr)
		return
	}

	netw = network

	result := m.Run()
	netw.lgr.Flush()
	os.Exit(result)
}

func TestConnectionPass(t *testing.T) {

	reachable := netw.IsInternetReachable()

	if !reachable {
		t.Error("The internet is unreachable. The unit test is either broken, you're not connected to the internet, or a sufficient number of the APIs that this code base relies upon in order to test internet connectivity are no longer working. Please verify your internet connection is working first to fix the issue or update the collection of REST endpoints if necessary.")
	}
}
