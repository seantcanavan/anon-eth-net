package reporter

import (
	"bytes"
	"net/smtp"
	"strings"

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

func (r *Reporter) InitializeReporter(cfg *config.Config) {
	r.GmailAddress = cfg.CheckInGmailAddress
	r.GmailPassword = cfg.CheckInGmailPassword
	r.DeviceName = cfg.DeviceName
	r.DeviceId = cfg.DeviceId
	r.emailServer = "smtp.gmail.com"
	r.emailPort = "587"
	r.emailAuth = smtp.PlainAuth("", r.GmailAddress, r.GmailPassword, r.emailServer)
}

// SendPlainEmail will send the content of the byte array as the body of an
// email along with the provided subject. The default sender and receiver are
// defined by InitializeReporter() which in turn can be defined via a
// config.json file. A sample is provided in the config package folder.
func (r *Reporter) SendPlainEmail(subject string, contents []byte) error {
	return r.prepareAndSendEmail(subject, strings.Split(string(contents), "\n"), nil)
}

// SendEmailAttachment will send the content of the byte array as the body of an
// email along with the provided subject. A pointer to a file is also accepted
// to be uploaded and sent as an attachment. The default sender and receiver are
// defined by InitializeReporter() which in turn can be defined via a
// config.json file. A sample is provided in the config package folder.
func (r *Reporter) SendEmailAttachment(subject string, contents []byte, attachment *os.File) error {
	return r.prepareAndSendEmail(subject, contents, attachment)
}

// SendCheckinEmail will send a simple email to the configured email address for
// the simple purposes of letting the owner know that the computer is alive and
// everything is still running. Contains a standard report of the computer's
// status like CPU usage, memory usage, disk usage, running processes, etc.
func (r *Reporter) SendCheckinEmail() {
	// r.prepareAndSendEmail("Check In", computerProfile)
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

// sendEmail will accept the subject and contents of the email as strings
// and concatenate them all into a nice clean buffer before sending off the
// email. Emails are sent asynchronously.
func (r *Reporter) prepareAndSendEmail(subject string, contents []string attachment *os.File) error {

	var messageBuffer bytes.Buffer
	messageBuffer.WriteString("From: ")
	messageBuffer.WriteString(r.GmailAddress)
	messageBuffer.WriteString("\n")
	messageBuffer.WriteString("To: ")
	messageBuffer.WriteString(r.GmailAddress)
	messageBuffer.WriteString("\n")
	messageBuffer.WriteString("Subject: ")
	messageBuffer.WriteString(r.generateSubject(subject))
	messageBuffer.WriteString("\n\n")

	for _, currentLine := range contents {
		messageBuffer.WriteString(currentLine)
	}

	return r.sendEmail(messageBuffer.Bytes())
}

func (r *Reporter) sendEmail(messageContents []byte) error {

	return smtp.SendMail(r.emailServer+":"+r.emailPort,
		r.emailAuth, r.GmailAddress, []string{r.GmailAddress}, messageContents)
}
