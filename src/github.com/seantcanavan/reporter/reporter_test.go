package reporter

import (
	"testing"

	"github.com/seantcanavan/config"
)

func TestSimpleEmail(t *testing.T) {

    r := &Reporter{}
	cfg, err := config.GetConfigFromFile("../config/config.json")
	if err != nil {
		return err
	}

	r.InitializeReporter(cfg)

	r.SendEmailUpdate("test subject", []string{"test body"})
}
