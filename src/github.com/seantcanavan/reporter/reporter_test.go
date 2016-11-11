package reporter

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"

	"github.com/seantcanavan/config"
)

var repr *Reporter

func TestMain(m *testing.M) {
	flag.Parse()
	err := config.ConfigFromFile(config.LOCAL_EXTERNAL_PATH)
	if err != nil {
		return
	}

	err = NewReporter()
	if err != nil {
		return
	}
	os.Exit(m.Run())
}

func TestSimpleEmail(t *testing.T) {
	err := repr.SendPlainEmail("TestSimpleEmail", []byte("TestSimpleEmail"))
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

	err := repr.SendAttachment("TestEmailAttachment", []byte("TestEmailAttachment"), filePtr)
	if err != nil {
		t.Error(err)
	}
}
