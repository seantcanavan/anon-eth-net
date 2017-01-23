package network

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/logger"
)

var netw *Network

func TestMain(m *testing.M) {

	logErr := logger.StandardLogger("network_test")
	if logErr != nil {
		fmt.Println(fmt.Sprintf("Could not initialize logger: %v", logErr))
		return
	}

	configErr := config.FromFile()
	if configErr != nil {
		fmt.Println(configErr)
		return
	}

	network, netwErr := NewNetwork()
	if netwErr != nil {
		fmt.Println(netwErr)
		return
	}

	netw = network

	result := m.Run()
	os.Exit(result)
}

func TestConnectionPass(t *testing.T) {
	reachable := netw.IsInternetReachable()

	if !reachable {
		t.Error("The internet is unreachable. The unit test is either broken, you're not connected to the internet, or a sufficient number of the APIs that this code base relies upon in order to test internet connectivity are no longer working. Please verify your internet connection is working first to fix the issue or update the collection of REST endpoints if necessary.")
	}
}

func TestRunPass(t *testing.T) {
	config.Cfg.NetQueryFrequencySeconds = 1
	netw.Run()
	time.Sleep(time.Second * 5)
}
