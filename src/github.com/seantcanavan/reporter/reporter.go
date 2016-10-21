package reporter

import (
	"bytes"
	"net/smtp"
	"os"

	"github.com/jordan-wright/email"
	"github.com/seantcanavan/config"
)

type Reporter struct {
	GmailAddress  string
	GmailPassword string
	DeviceName    string
	DeviceId      string
	emailServer   string
	emailPort     string
	emailAuth     smtp.Auth
}

func NewReporter(cfg *config.Config) *Reporter {
	r := Reporter{}
	r.GmailAddress = cfg.CheckInGmailAddress
	r.GmailPassword = cfg.CheckInGmailPassword
	r.DeviceName = cfg.DeviceName
	r.DeviceId = cfg.DeviceId
	r.emailServer = "smtp.gmail.com"
	r.emailPort = "587"
	r.emailAuth = smtp.PlainAuth("", r.GmailAddress, r.GmailPassword, r.emailServer)
	return &r
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
		To:      []string{r.GmailAddress},
		From:    r.GmailAddress,
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
	subjectBuffer.WriteString(r.DeviceId)
	subjectBuffer.WriteString("] ")
	subjectBuffer.WriteString(subject)
	return subjectBuffer.String()
}
