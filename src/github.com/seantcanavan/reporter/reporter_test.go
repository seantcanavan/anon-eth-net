package reporter

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/utils"
)

func TestMain(m *testing.M) {

	flag.Parse()

	assetPath, assetErr := utils.AssetPath("config.json")
	if assetErr != nil {
		fmt.Println(assetErr)
		return
	}

	err := config.FromFile(assetPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Exit(m.Run())
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
