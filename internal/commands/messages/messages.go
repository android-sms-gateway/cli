package messages

import (
	"github.com/android-sms-gateway/cli/internal/commands/messages/batch"
	"github.com/urfave/cli/v2"
)

func Commands() []*cli.Command {
	return []*cli.Command{
		sendCmd(),
		statusCmd(),
		batch.Commands(),
	}
}
