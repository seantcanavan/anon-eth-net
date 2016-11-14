package loader

import (
	// "fmt"
	// "runtime"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/utils"
)

func TestProcessMap(t *testing.T) {

	assetPath, assetErr := utils.AssetPath("config.json")
	if assetErr != nil {
		t.Error(assetErr)
	}

	cfgErr := config.FromFile(assetPath)
	if cfgErr != nil {
		t.Error(cfgErr)
	}

	loaderAssetPath, assetErr := utils.SysAssetPath("loader_test.json")
	if assetErr != nil {
		t.Error(assetErr)
	}

	loader, loaderErr := NewLoader(loaderAssetPath)

	if loaderErr != nil {
		t.Error(loaderErr)
	}

	loader.StartSynchronous()
}
