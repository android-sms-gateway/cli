package providers

import (
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/urfave/cli/v2"
)

// MessageContentProvider fills the content portion of a smsgateway.Message.
// The caller creates the Message struct and passes it to the provider.
// The provider sets only the relevant content field (TextMessage, DataMessage,
// or MmsMessage). The caller is responsible for envelope fields (ID, PhoneNumbers, etc.).
type MessageContentProvider interface {
	// PrepareMessage fills the content field of msg from CLI context parameters.
	// It parses required flags and arguments from the context and sets
	// the appropriate content field on the provided Message.
	PrepareMessage(c *cli.Context, msg *smsgateway.Message) error
}
