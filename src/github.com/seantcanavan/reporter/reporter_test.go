package reporter

import (
	"testing"

	"github.com/seantcanavan/config"
)

func TestSimpleEmail(t *testing.T) {
	cfg, _ := config.GetConfigFromFile("../config/config.json")

    r := &Reporter{
		GmailAddress:  cfg.CheckInGmailAddress,
		GmailPassword: cfg.CheckInGmailPassword,
		DeviceName:    cfg.DeviceName,
		DeviceId:      cfg.DeviceId,
	}

	r.InitializeReporter()

	r.SendEmailUpdate("test subject", []string{"test body"})
}
