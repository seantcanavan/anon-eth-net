package reporter

import (
	"bytes"
	"net/smtp"
	"os"
	"time"

	"github.com/jordan-wright/email"
	"github.com/seantcanavan/config"
)

const EMAIL_SERVER = "smtp.gmail.com"
const EMAIL_PORT = "587"
const MAX_EMAIL_TIMEOUT_ATTEMPTS = 5

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

	if attachmentPtr != nil {
		jwEmail.AttachFile(attachmentPtr.Name())
	}

	emailAuth := smtp.PlainAuth("", config.Cfg.CheckInGmailAddress, config.Cfg.CheckInGmailPassword, EMAIL_SERVER)
	count := 0
	var emailErr error

	for count < MAX_EMAIL_TIMEOUT_ATTEMPTS {
		emailErr = jwEmail.Send(EMAIL_SERVER+":"+EMAIL_PORT, emailAuth)
		if emailErr == nil {
			break
		}
		count++
		time.Sleep(time.Second * 5)
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
