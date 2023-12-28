package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

const (
    // gmailSmtpAuthAddress   = "smtp.gmail.com"
    // gmailSmtpServerAddress = "smtp.gmail.com:587"
    smtpAuthAddress = "smtp.126.com"
    smtpServerAddress = "smtp.126.com:25"
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

    smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
    return e.Send(smtpServerAddress, smtpAuth)
}
