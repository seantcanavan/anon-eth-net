package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

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

func TestSimpleRestBringUpPass(t *testing.T) {
	now := time.Now().Unix()
	fmt.Println(config.Cfg.LogVolatility)
	restHandler, restErr := NewRestHandler()
	if restErr != nil {
		t.Error(restErr)
	}

	port := restHandler.Port

	var addressBuf bytes.Buffer

	addressBuf.WriteString("http://localhost:")
	addressBuf.WriteString(strconv.Itoa(port))
	addressBuf.WriteString("/checkin/")
	addressBuf.WriteString(strconv.FormatInt(now, 10))
	addressBuf.WriteString("/")
	addressBuf.WriteString(strings.Split(config.Cfg.CheckInGmailAddress, "@")[0])

	fmt.Println(fmt.Sprintf("TestSimpleRestBringUpPass: %v", addressBuf.String()))

	response, err := http.Get(addressBuf.String())
	if err != nil {
		t.Error(err)
	}
	fmt.Println(response)
}
