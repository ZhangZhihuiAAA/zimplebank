package mail

import (
	"testing"

	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWith126(t *testing.T) {
    if testing.Short() {
        t.Skip()
    }

    config, err := util.LoadConfig("..")
    require.NoError(t, err)

    sender := NewOneTwoSixSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

    subject := "A test email"
    content := `
    <h1>Hello world</h1>
    <p>This is a test message.</p>
    `
    to := []string{"ZhangZhihuiAAA@126.com"}
    attachedFiles := []string{"../README.md"}

    err = sender.Send(subject, content, to, nil, nil, attachedFiles)
    require.NoError(t, err)
}