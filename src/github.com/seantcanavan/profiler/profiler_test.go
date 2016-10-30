package profiler

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/reporter"
)

var repr *reporter.Reporter
var prof *SysProfiler

func TestMain(m *testing.M) {
	flag.Parse()
	err := config.ConfigFromFile(config.LOCAL_EXTERNAL_PATH)
	if err != nil {
		return
	}

	repr = reporter.NewReporter()
	prof = NewSysProfiler(repr)

	os.Exit(m.Run())
}

func TestSendFileArchiveAsAttachment(t *testing.T) {
	filePtr, err := prof.SendArchiveProfileAsAttachment()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fmt.Sprintf("Archive successfully created: %v", filePtr.Name()))
}
