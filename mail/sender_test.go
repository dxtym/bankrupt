package mail

import (
	"testing"

	"github.com/dxtym/bankrupt/utils"
	"github.com/stretchr/testify/require"
)

func TestEmailSender(t *testing.T) {
	config, err := utils.LoadConfig("..")
	require.NoError(t, err)

	sender := NewEmailSender(config.SenderName, config.SenderAddress, config.SenderPassword)
	require.NotNil(t, sender)	

	subject := "Test Email"
	content := `<h1>Test Email</h1>
	<p>Hello, world!</p>
	<p>Best regard,<br>Bankrupt</p>`
	to := []string{config.SenderAddress}
	attachments := []string{"test.txt"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachments)
	require.NoError(t, err)
}