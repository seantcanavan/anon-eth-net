package rest

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
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
var port string
var publicKey string
var transport *http.Transport
var client *http.Client

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

	certPath, certPathErr := utils.AssetPath("server.cert")
	if certPathErr != nil {
		fmt.Println(certPathErr)
		return
	}

	certValue, certReadErr := ioutil.ReadFile(certPath)
	if certReadErr != nil {
		fmt.Println(certReadErr)
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
	protocol = "https"
	host = "localhost"

	rootAuthorities := x509.NewCertPool()
	if ok := rootAuthorities.AppendCertsFromPEM([]byte(certValue)); !ok {
		fmt.Println("Unable to append certificate to set of root certificate authorities")
		return
	}

	transport = &http.Transport{TLSClientConfig: &tls.Config{RootCAs: rootAuthorities}}
	client = &http.Client{Transport: transport}

	os.Exit(m.Run())
}

func TestCheckinHandlerPass(t *testing.T) {

	path = buildRestPath(protocol, host, port, CHECKIN_REST_PATH, nowString, "samplegmail")

	fmt.Println(fmt.Sprintf("TestCheckinHandlerPass: client.Get -> : %v", path))

	response, err := client.Get(path)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
	}

	fmt.Println(fmt.Sprintf("TestCheckinHandlerPass: client.Post -> %v", path))

	response, err = client.Post(path, "application/json", bytes.NewBuffer([]byte("method not supported")))
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusMethodNotAllowed, response.StatusCode))
	}
}

func TestExecuteHandlerPass(t *testing.T) {
	path = buildRestPath(protocol, host, port, EXECUTE_REST_PATH, nowString, "python")

	fmt.Println(fmt.Sprintf("TestExecuteHandlerPass: client.Get -> %v", path))

	response, err := client.Get(path)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusMethodNotAllowed, response.StatusCode))
	}

	fmt.Println(fmt.Sprintf("TestExecuteHandlerPass: client.Post -> %v", path))

	response, err = client.Post(path, "text/plain", bytes.NewBuffer([]byte("print(\"python script woah!\")")))
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
	}

	// path = buildRestPath(protocol, host, port, EXECUTE_REST_PATH, nowString, "binary")

	// fmt.Println(fmt.Sprintf("TestExecuteHandlerPass: client.Post -> %v", path))

	// response, err = client.Post(path, "application/octet-stream", bytes.NewBuffer([]byte("this will surely fail")))
	// if err != nil {
	// 	t.Error(err)
	// }

	// if response.StatusCode != http.StatusOK {
	// 	t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
	// }

	path = buildRestPath(protocol, host, port, EXECUTE_REST_PATH, nowString, "script")

	fmt.Println(fmt.Sprintf("TestExecuteHandlerPass: client.Post -> %v", path))

	response, err = client.Post(path, "text/plain", bytes.NewBuffer([]byte("echo \"hello world\"")))
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
// 	path = buildRestPath(protocol, host, port, REBOOT_REST_PATH, nowString, "10")

// 	fmt.Println(fmt.Sprintf("TestRebootHandler: client.Post -> %v", path))

// 	response, err := client.Post(path, "text/plain", bytes.NewBuffer([]byte("welcome to my house")))
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if response.StatusCode != http.StatusMethodNotAllowed {
// 		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusMethodNotAllowed, response.StatusCode))
// 	}

// 	fmt.Println(fmt.Sprintf("TestRebootHandler: client.Get -> %v", path))

// 	response, err = client.Get(path)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if response.StatusCode != http.StatusOK {
// 		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
// 	}

// 	fmt.Println("The computer should restart now after a short delay...")
// }

func TestLogHandlerPass(t *testing.T) {
	path = buildRestPath(protocol, host, port, LOG_REST_PATH, nowString, "samplegmail")

	fmt.Println(fmt.Sprintf("TestLogHandlerPass: client.Get -> %v", path))

	response, err := client.Get(path)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
	}

	fmt.Println(fmt.Sprintf("TestLoginHandlerPass: client.Post -> %v", path))

	response, err = client.Post(path, "text/plain", bytes.NewBuffer([]byte("welcome to my house")))
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusMethodNotAllowed, response.StatusCode))
	}
}

func TestUpdateHandlerPass(t *testing.T) {
	path = buildRestPath(protocol, host, port, UPDATE_REST_PATH, nowString)

	fmt.Println(fmt.Sprintf("TestUpdateHandlerPass: client.Get -> %v", path))

	response, err := client.Get(path)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
	}

	fmt.Println(fmt.Sprintf("TestUpdateHandlerPass: client.Post -> %v", path))

	response, err = client.Post(path, "text/plain", bytes.NewBuffer([]byte("welcome to my house")))
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
	}

	fmt.Println(fmt.Sprintf("TestUpdateHandlerPass: client.Do -> %v", path))

	response, err = client.Head(path)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusMethodNotAllowed, response.StatusCode))
	}
}

func TestConfigHandlerPass(t *testing.T) {
	path = buildRestPath(protocol, host, port, CONFIG_REST_PATH, nowString)

	fmt.Println(fmt.Sprintf("TestConfigHandlerPass: client.Get -> %v", path))

	response, err := client.Get(path)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(fmt.Sprintf("TestConfigHandlerPass: client.Post -> %v", path))

	response, err = client.Post(path, "text/plain", bytes.NewBuffer([]byte("welcome to my house")))
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusOK, response.StatusCode))
	}

	fmt.Println(fmt.Sprintf("TestConfigHandlerPass: client.Do -> %v", path))

	response, err = client.Head(path)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusMethodNotAllowed {
		t.Error(fmt.Errorf("expected: %v, got: %v", http.StatusMethodNotAllowed, response.StatusCode))
	}
}
