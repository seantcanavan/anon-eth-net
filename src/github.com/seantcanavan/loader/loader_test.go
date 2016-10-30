package loader

import (
	// "fmt"
	"testing"
	// "time"

	"github.com/seantcanavan/config"
)

func TestProcessMap(t *testing.T) {
	config.ConfigFromFile(config.LOCAL_EXTERNAL_PATH)
	loader, err := NewLoader("loader.json")
	if err != nil {
		t.Error(err)
	}

	loader.StartSynchronous()
}
