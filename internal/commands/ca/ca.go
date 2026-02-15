package ca

import (
	"github.com/android-sms-gateway/cli/internal/commands/ca/common"
	"github.com/android-sms-gateway/client-go/ca"
	"github.com/urfave/cli/v2"
)

func Commands() []*cli.Command {
	return []*cli.Command{
		common.NewIPCertificateCommand(
			"webhooks",
			"Issue a new certificate for receiving webhooks to local IP address",
			[]string{"wh"},
			ca.CSRTypeWebhook,
		),
		common.NewIPCertificateCommand(
			"private",
			"Issue a new certificate for Private server",
			[]string{"p"},
			ca.CSRTypePrivateServer,
		),
	}
}
