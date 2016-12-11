package reporter

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/logger"
)

func TestMain(m *testing.M) {

	logErr := logger.StandardLogger("reporter_test")
	if logErr != nil {
		fmt.Println(fmt.Sprintf("Could not initialize logger: %v", logErr))
		return
	}

	configErr := config.FromFile()
	if configErr != nil {
		fmt.Println(configErr)
		return
	}

	result := m.Run()
	// flush logs
	os.Exit(result)
}

func TestSimpleEmail(t *testing.T) {
	err := SendPlainEmail("TestSimpleEmail", []byte("TestSimpleEmail"))
	if err != nil {
		t.Error(err)
	}
}

func TestEmailAttachment(t *testing.T) {
	testName := "bklajkjlkja.txt"
	ioutil.WriteFile(testName, []byte("TestEmailAttachment"), 0744)
	defer os.Remove(testName)

	filePtr, openErr := os.Open(testName)
	if openErr != nil {
		t.Error(openErr)
	}

	err := SendAttachment("TestEmailAttachment", []byte("TestEmailAttachment"), filePtr)
	if err != nil {
		t.Error(err)
	}
}
