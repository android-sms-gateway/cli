package providers

import (
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/urfave/cli/v2"
)

// DataContentProvider fills the DataMessage field of a Message.
// It parses the data and port from CLI flags.
type DataContentProvider struct{}

// PrepareMessage fills the DataMessage field of msg.
// It reads the data from the first CLI argument and the port from the --data-port flag.
func (p DataContentProvider) PrepareMessage(c *cli.Context, msg *smsgateway.Message) error {
	data := c.Args().Get(0)
	port := uint16(c.Uint("data-port")) //nolint:gosec // validated upstream by sendBefore
	msg.DataMessage = &smsgateway.DataMessage{
		Data: data,
		Port: port,
	}
	return nil
}
