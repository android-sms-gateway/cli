package providers

import (
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/urfave/cli/v2"
)

// TextContentProvider fills the TextMessage field of a Message.
// It parses the message text from the first CLI argument.
type TextContentProvider struct{}

// PrepareMessage fills the TextMessage field of msg.
// It reads the message text from the first CLI argument.
func (p TextContentProvider) PrepareMessage(c *cli.Context, msg *smsgateway.Message) error {
	text := c.Args().Get(0)
	msg.TextMessage = &smsgateway.TextMessage{
		Text: text,
	}
	return nil
}
