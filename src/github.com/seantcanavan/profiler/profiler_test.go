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
	cfg, err := config.ConfigFromFile("../config/config.json")
	if err != nil {
		return
	}

	repr = reporter.NewReporter(cfg)
	prof = NewSysProfiler(repr)

	os.Exit(m.Run())
}

func TestSendFileArchiveAsAttachment(t *testing.T) {
	err := prof.SendArchiveProfileAsAttachment()
	if err != nil {
		t.Error(err)
	}
}
