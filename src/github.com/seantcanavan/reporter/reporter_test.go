package reporter

import (
	"testing"

	"github.com/seantcanavan/config"
	"github.com/seantcanavan/profile"
)

var repr reporter.Reporter
var prof profiler.Profiler

func TestMain(m *testing.M) {
	flag.Parse()
	cfg, err := config.ConfigFromFile("../config/config.json")
	if err != nil {
		t.Error()
	}

	repr = NewReporter(cfg)
	prof = NewSysProfiler(repr)

	os.Exit(m.Run())
}

func TestSimpleEmail(t *testing.T) {
	err := repr.SendPlainEmail("test subject", []byte{"test body"})
	if err != nil {
		t.Error(err)
	}
}

func TestEmailAttachment(t *testing.T) {
	fileProfile, err := prof.ProfileAsFile()
	if err != nil {
		t.Error(err)
	}

	err = repr.SendEmailAttachment("test subject", []byte{"test body"}, fileProfile)
	if err != nil {
		t.Error(err)
	}
}

func TestEmailArchive(t *testing.T) {
	archiveProfile, err := prof.ProfileAsArchive()
	if err != nil {
		t.Error(err)
	}

	err = repr.SendEmailAttachment("test subject", []byte{"test body"}, archiveProfile)
	if err != nil {
		t.Error(err)
	}
}
