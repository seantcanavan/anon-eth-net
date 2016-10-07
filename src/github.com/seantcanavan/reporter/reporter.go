package reporter

import (
	"bytes"
	"net/smtp"
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

func (r *Reporter) InitializeReporter() {
	r.emailServer = "smtp.gmail.com"
	r.emailPort = "587"
	r.emailAuth = smtp.PlainAuth("", r.GmailAddress, r.GmailPassword, r.emailServer)
}

// SendEmailReport will send the current slice of strings to yourself. The
// email username is pulled from config.json and the target is always yourself.
// This guarantees that the emails stay secluded and private.
func (r *Reporter) SendEmailUpdate(subject string, contents []string) error {
	return r.prepareAndSendEmail(subject, contents)
}

// SendCheckinEmail will send a simple email to the configured email address for
// the simple purposes of letting the owner know that the computer is alive and
// everything is still running. Contains a standard report of the computer's
// status like CPU usage, memory usage, disk usage, running processes, etc.
func (r *Reporter) SendCheckinEmail() {
	// computerProfile := profiler.GetProfile()
	// r.prepareAndSendEmail("Check In", computerProfile)
}

// SendUrgentStatusUpdate will immediately immediately send an email containing
// any error message(s) included and also included the provided log file for the
// owner to inspect. A standard report of the computer's health will also be
// provided.
func (r *Reporter) SendUrgentStatusUpdate(subject string, contents []string, logFilePath string) {
	// find a way to send attachments
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
func (r *Reporter) prepareAndSendEmail(subject string, contents []string) error {

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
