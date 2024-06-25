package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
	smtpAuthAddress = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachments []string,
	) error
}

type Sender struct {
	name        string
	fromAddress   string
	fromPassword string
}

func NewEmailSender(name, fromAddress, fromPassword string) EmailSender {
	return &Sender{
		name: name,
		fromAddress: fromAddress,
		fromPassword: fromPassword,
	}
}

func (s *Sender) SendEmail(
	subject string,
	content string,
	to []string,
	cc []string,
	bcc []string,
	attachments []string,
) error {
	e := email.NewEmail()
	e.From = fmt.Sprintf("%s <%s>", s.name, s.fromAddress)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc
	e.Subject = subject
	e.HTML = []byte(content)

	for _, file := range attachments {
		e.AttachFile(file)
	}

	return e.Send(smtpServerAddress, smtp.PlainAuth("", s.fromAddress, s.fromPassword, smtpAuthAddress))
}