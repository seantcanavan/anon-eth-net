package profiler

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/reporter"
)

func TestReportAsBytes(t *testing.T) {
	bytes := ReportAsBytes()
	fmt.Println(string(bytes))
}

func TestReportAsFile(t *testing.T) {
	fileName, err := ReportAsFile(true)
	if err != nil {
		t.Error(err)
	}
}

func TestReportAsEmailBody(t *testing.T) {
	conf, err := config.ConfigFromFile("../config/config.json")
	if err != nil {
		t.Error(err)
	}

	rep := reporter.Reporter{}
	rep.InitializeReporter(conf)

	err = ReportAsBytes()
	if err != nil {
		t.Error(err)
	}


}
