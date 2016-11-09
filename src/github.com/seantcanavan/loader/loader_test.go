package loader

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/seantcanavan/config"
)

func TestProcessMap(t *testing.T) {
	config.ConfigFromFile("loader_config.json")

	var loader *Loader
	var loaderError error

	switch runtime.GOOS {
	case "windows":
		loader, loaderError = NewLoader("loader_test_windows.json")
	case "darwin":
		loader, loaderError = NewLoader("loader_test_darwin.json")
	case "linux":
		loader, loaderError = NewLoader("loader_test_linux.json")
	default:
		t.Error(fmt.Errorf("Could not create loader for unsupported operating system: %v", runtime.GOOS))
	}

	if loaderError != nil {
		t.Error(loaderError)
	}

	loader.StartSynchronous()
}
