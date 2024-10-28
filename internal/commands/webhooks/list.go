package webhooks

import (
	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/urfave/cli/v2"
)

var list = &cli.Command{
	Category: "Webhooks",
	Name:     "list",
	Aliases:  []string{"l", "ls"},
	Usage:    "List webhooks",
	Action: func(c *cli.Context) error {
		// client := metadata.GetClient(c.App.Metadata)
		// renderer := metadata.GetRenderer(c.App.Metadata)

		return cli.Exit("Not implemented", codes.ParamsError)
	},
}
