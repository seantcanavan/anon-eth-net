package profiler

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/seantcanavan/config"
)

func TestReportAsBytes(t *testing.T) {
	bytes := ReportAsBytes()
	fmt.Println(string(bytes))
}

func TestReportAsFile(t *testing.T) {

	fileName, err := ReportAsFile(false, 0)
	if err != nil {
		t.Error(err)
	}

	bytes, readErr := ioutil.ReadFile(fileName)
	if readErr != nil {
		t.Error(readErr)
	}

	fmt.Println(string(bytes))
}

func TestReportAsEmailBody(t *testing.T) {

	conf, err := config.ConfigFromFile("../config/config.json")
	if err != nil {
		t.Error(err)
	}

	err = SendReportAsEmail(conf)

	if err != nil {
		t.Error(err)
	}
}
