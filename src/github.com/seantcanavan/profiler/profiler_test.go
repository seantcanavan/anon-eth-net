package profiler

import (
	"flag"
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
	err := prof.SendArchiveProfileAsAttachment()
	if err != nil {
		t.Error(err)
	}
}
