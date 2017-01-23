package reporter

import (
	"bytes"
	"net/smtp"
	"os"
	"time"

	"github.com/jordan-wright/email"
	"github.com/seantcanavan/anon-eth-net/config"
	"github.com/seantcanavan/anon-eth-net/logger"
)

const EMAIL_SERVER = "smtp.gmail.com"
const EMAIL_PORT = "587"
const MAX_EMAIL_TIMEOUT_ATTEMPTS = 5
const SUCCESSIVE_EMAIL_ATTEMPTS_DELAY = 5

// SendPlainEmail will send the content of the byte array as the body of an
// email along with the provided subject. The default sender and receiver are
// defined by NewReporter() which in turn can be defined via a
// config.json file. A sample is provided in the config package folder.
func SendPlainEmail(subject string, contents []byte) error {
	return SendAttachment(subject, contents, nil)
}

//SendAttachment will send the content of the byte array as the body of an email
// along with the provided subject. The device ID is automatically added to the
// email subject line in order to help differentiate emails from multiple
// devices to the same address. The sender and receiver are defined by
// NewReporter() which in turn can be defined via config.json file.
func SendAttachment(subject string, contents []byte, attachmentPtr *os.File) error {
	jwEmail := &email.Email{
		To:      []string{config.Cfg.CheckInGmailAddress},
		From:    config.Cfg.CheckInGmailAddress,
		Subject: generateSubject(subject),
		Text:    contents,
	}

	logger.Lgr.LogMessage("Successfully created new jwemail instance to: %v", config.Cfg.CheckInGmailAddress)

	if attachmentPtr != nil {
		jwEmail.AttachFile(attachmentPtr.Name())
		logger.Lgr.LogMessage("Successfully attached file: %v", attachmentPtr.Name())
	}

	emailAuth := smtp.PlainAuth("", config.Cfg.CheckInGmailAddress, config.Cfg.CheckInGmailPassword, EMAIL_SERVER)

	logger.Lgr.LogMessage("Successfully generated SMTP email auth: %+v", emailAuth)

	count := 0
	var emailErr error

	for count < MAX_EMAIL_TIMEOUT_ATTEMPTS {
		emailErr = jwEmail.Send(EMAIL_SERVER+":"+EMAIL_PORT, emailAuth)
		if emailErr == nil {
			logger.Lgr.LogMessage("Successfully sent out email to: %v", config.Cfg.CheckInGmailAddress)
			break
		}
		count++
		logger.Lgr.LogMessage("Unsuccessfully sent out email to: %v. Sleeping for %d", config.Cfg.CheckInGmailAddress, SUCCESSIVE_EMAIL_ATTEMPTS_DELAY)
		time.Sleep(time.Second * SUCCESSIVE_EMAIL_ATTEMPTS_DELAY)
	}

	return emailErr
}

// generateSubject will append the device ID to the beginning of the email
// subject for easier sorting / searching through the list of emails to help
// keep track of emails by device.
func generateSubject(subject string) string {
	var subjectBuffer bytes.Buffer
	subjectBuffer.WriteString("[")
	subjectBuffer.WriteString(config.Cfg.DeviceId)
	subjectBuffer.WriteString("] ")
	subjectBuffer.WriteString(subject)
	return subjectBuffer.String()
}
