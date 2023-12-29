package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
    // GMAIL_SMTP_AUTH_ADDRESS   = "smtp.gmail.com"
    // GMAIL_SMTP_SERVER_ADDRESS = "smtp.gmail.com:587"
    SMTP_AUTH_ADDRESS   = "smtp.126.com"
    SMTP_SERVER_ADDRESS = "smtp.126.com:25"
)

type Sender interface {
    Send(
        subject string,
        content string,
        to []string,
        cc []string,
        bcc []string,
        attachedFiles []string,
    ) error
}

type OneTwoSixSender struct {
    name              string
    fromEmailAddress  string
    fromEmailPassword string
}

func NewOneTwoSixSender(name string, fromEmailAddress string, fromEmailPassword string) Sender {
    return OneTwoSixSender{
        name:              name,
        fromEmailAddress:  fromEmailAddress,
        fromEmailPassword: fromEmailPassword,
    }
}

func (sender OneTwoSixSender) Send(
    subject string,
    content string,
    to []string,
    cc []string,
    bcc []string,
    attachedFiles []string,
) error {
    e := email.NewEmail()
    e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
    e.Subject = subject
    e.HTML = []byte(content)
    e.To = to
    e.Cc = cc
    e.Bcc = bcc

    for _, f := range attachedFiles {
        _, err := e.AttachFile(f)
        if err != nil {
            return fmt.Errorf("failed to attach file %s: %w", f, err)
        }
    }

    smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, SMTP_AUTH_ADDRESS)
    return e.Send(SMTP_SERVER_ADDRESS, smtpAuth)
}
