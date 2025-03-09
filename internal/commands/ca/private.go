package ca

import (
	"github.com/android-sms-gateway/cli/internal/commands/ca/common"
	"github.com/android-sms-gateway/client-go/ca"
)

var private = common.NewIPCertificateCommand(
	"private",
	"Issue a new certificate for Private server",
	[]string{"p"},
	ca.CSRTypePrivateServer,
)
