package network

import (
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
)

var netw *Network

func TestMain(m *testing.M) {

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
