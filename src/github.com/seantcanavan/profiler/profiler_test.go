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
	cfg, err := config.ConfigFromFile("../config/config.json")
	if err != nil {
		return
	}

	repr = reporter.NewReporter(cfg)
	prof = NewSysProfiler(repr)

	os.Exit(m.Run())
}

func TestProfileAsBytes(t *testing.T) {
	bytes := prof.ProfileAsBytes()
	fmt.Println(string(bytes))
}

func TestProfileAsFile(t *testing.T) {
	_, err := prof.ProfileAsFile()
	if err != nil {
		t.Error(err)
	}
}

func TestProfileAsArchive(t *testing.T) {
	_, err := prof.ProfileAsArchive()
	if err != nil {
		t.Error(err)
	}
}

func TestSendByteProfileAsEmail(t *testing.T) {
	err := prof.SendByteProfileAsEmail()
	if err != nil {
		t.Error(err)
	}
}

func TestSendFileProfileAsAttachment(t *testing.T) {
	err := prof.SendFileProfileAsAttachment()
	if err != nil {
		t.Error(err)
	}
}

func TestSendFileArchiveAsAttachment(t *testing.T) {
	err := prof.SendArchiveProfileAsAttachment()
	if err != nil {
		t.Error(err)
	}
}
