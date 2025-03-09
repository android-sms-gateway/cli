package ca

import (
	"github.com/android-sms-gateway/cli/internal/commands/ca/common"
	"github.com/android-sms-gateway/client-go/ca"
)

var webhooks = common.NewIPCertificateCommand(
	"webhooks",
	"Issue a new certificate for receiving webhooks to local IP address",
	[]string{"wh"},
	ca.CSRTypeWebhook,
)
