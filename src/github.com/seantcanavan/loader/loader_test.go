package loader

import (
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/utils"
)

func TestMain(m *testing.M) {

	configErr := config.FromFile()
	if configErr != nil {
		fmt.Println(fmt.Sprintf("could not initialize config: %v", configErr))
	}

	result := m.Run()
	os.Exit(result)
}

func TestProcessMap(t *testing.T) {

	loaderAssetPath, assetErr := utils.SysAssetPath("loader_test.json")
	if assetErr != nil {
		t.Error(assetErr)
	}

	loader, loaderErr := NewLoader(loaderAssetPath)

	if loaderErr != nil {
		t.Error(loaderErr)
	}

	loader.StartAsynchronous()
}
