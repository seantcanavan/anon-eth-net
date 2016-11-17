package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/utils"
)

var path string
var now int64
var nowString string
var protocol string
var host string
var port int
var portString string

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

	now = time.Now().Unix()
	nowString = strconv.FormatInt(now, 10)

	restHandler, restErr := NewRestHandler()
	if restErr != nil {
		fmt.Println(restErr)
		return
	}

	port = restHandler.Port
	portString = strconv.Itoa(restHandler.Port)
	protocol = "http"
	host = "localhost"

	os.Exit(m.Run())
}

func TestCheckinHandlerPass(t *testing.T) {

	path = buildRestPath(protocol, host, portString, CHECKIN_REST_PATH, nowString, "samplegmail")

	fmt.Println(fmt.Sprintf("TestAllRestEndPoints: http.Get -> : %v", path))

	response, err := http.Get(path)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
	}

	response, err = http.Post(path, "application/json", bytes.NewBuffer([]byte("method not supported")))
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusMethodNotAllowed, response.StatusCode))
	}
}

func TestExecuteHandlerPass(t *testing.T) {
	path = buildRestPath(protocol, host, portString, EXECUTE_REST_PATH, nowString, "python")

	response, err = http.Get(path)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusMethodNotAllowed, response.StatusCode))
	}

	response, err = http.Post(path, "text/plain", bytes.NewBuffer([]byte("print(\"python script woah!\"")))
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
	}

	path = buildRestPath(protocol, host, portString, EXECUTE_REST_PATH, nowString, "binary")

	response, err = http.Post(path, "application/octet-stream", bytes.NewBuffer([]byte("this will surely fail")))
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
	}

	path = buildRestPath(protocol, host, portString, EXECUTE_REST_PATH, nowString, "script")

	response, err = http.Post(path, "text/plain", bytes.NewBuffer([]byte("echo hello world")))
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
	}
}

// uncomment to test manually - don't want to reboot the computer every time
// this test is executed
// func TestRebootHandler(t *testing.T) {
// 	path = buildRestPath(protocol, host, portString, REBOOT_REST_PATH, nowString, "10")
// }

func TestLogHandlerPass(t *testing.T) {
	path = buildRestPath(protocol, host, portString, LOG_REST_PATH, nowString, "samplegmail")
}

func TestUpdateHandlerPass(t *testing.T) {
	path = buildRestPath(protocol, host, portString, UPDATE_REST_PATH, nowString, "remoteUpdateURI")
}

func TestConfigHandlerPass(t *testing.T) {
	path = buildRestPath(protocol, host, portString, CONFIG_REST_PATH, nowString)
}
