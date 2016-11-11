package reporter

import (
	"bytes"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
	"github.com/seantcanavan/config"
)

var Rpr *Reporter

type Reporter struct {
	gmailAddress  string
	gmailPassword string
	deviceName    string
	deviceId      string
	emailServer   string
	emailPort     string
	emailAuth     smtp.Auth
}

func NewReporter() error {
	r := Reporter{}
	r.gmailAddress = config.Cfg.CheckInGmailAddress
	r.gmailPassword = config.Cfg.CheckInGmailPassword
	r.deviceName = config.Cfg.DeviceName
	r.deviceId = config.Cfg.DeviceId
	r.emailServer = "smtp.gmail.com"
	r.emailPort = "587"
	r.emailAuth = smtp.PlainAuth("", r.gmailAddress, r.gmailPassword, r.emailServer)
	Rpr = &r
	return nil
}

// SendPlainEmail will send the content of the byte array as the body of an
// email along with the provided subject. The default sender and receiver are
// defined by NewReporter() which in turn can be defined via a
// config.json file. A sample is provided in the config package folder.
func (r *Reporter) SendPlainEmail(subject string, contents []byte) error {
	return r.SendAttachment(subject, contents, nil)
}

//SendAttachment will send the content of the byte array as the body of an email
// along with the provided subject. The device ID is automatically added to the
// email subject line in order to help differentiate emails from multiple
// devices to the same address. The sender and receiver are defined by
// NewReporter() which in turn can be defined via config.json file.
func (r *Reporter) SendAttachment(subject string, contents []byte, attachmentPtr *os.File) error {
	jwEmail := &email.Email{
		To:      []string{r.gmailAddress},
		From:    r.gmailAddress,
		Subject: r.generateSubject(subject),
		Text:    contents,
	}

	if attachmentPtr != nil {
		jwEmail.AttachFile(attachmentPtr.Name())
	}

	return jwEmail.Send(r.emailServer+":"+r.emailPort, r.emailAuth)
}

// generateSubject will append the device ID to the beginning of the email
// subject for easier sorting / searching through the list of emails to help
// keep track of emails by device.
func (r *Reporter) generateSubject(subject string) string {
	var subjectBuffer bytes.Buffer
	subjectBuffer.WriteString("[")
	subjectBuffer.WriteString(r.deviceId)
	subjectBuffer.WriteString("] ")
	subjectBuffer.WriteString(subject)
	return subjectBuffer.String()
}
