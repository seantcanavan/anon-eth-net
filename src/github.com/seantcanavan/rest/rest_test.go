package rest

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	// "time"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/utils"
)

func TestMain(m *testing.M) {

	assetPath, assetErr := utils.AssetPath("config.json")
	if assetErr != nil {
		fmt.Println(assetErr)
		return
	}

	cfgErr := config.FromFile(assetPath)
	if cfgErr != nil {
		fmt.Println(cfgErr)
		return
	}
	os.Exit(m.Run())
}

func TestSimpleRestBringupPass(t *testing.T) {
	fmt.Println(config.Cfg.LogVolatility)
	restHandler, restErr := NewRestHandler()
	if restErr != nil {
		t.Error(restErr)
	}

	port := restHandler.Port
	response, err := http.Get("http://localhost:" + strconv.Itoa(port) + "/checkin/")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(response)
}
