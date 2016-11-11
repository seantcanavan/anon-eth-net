package profiler

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/reporter"
)

func TestMain(m *testing.M) {
	flag.Parse()
	err := config.ConfigFromFile("profiler_config.json")
	if err != nil {
		return
	}

	reporter.NewReporter()
	os.Exit(m.Run())
}

func TestSendFileArchiveAsAttachment(t *testing.T) {
	filePtr, err := SendArchiveProfileAsAttachment()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("Archive successfully created: %v", filePtr.Name()))
}
