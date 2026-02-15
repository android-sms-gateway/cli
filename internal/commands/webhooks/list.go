package webhooks

import (
	"fmt"
	"os"

	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/urfave/cli/v2"
)

func listCmd() *cli.Command {
	return &cli.Command{
		Category: "Webhooks",
		Name:     "list",
		Aliases:  []string{"l", "ls"},
		Usage:    "List webhooks",
		Action: func(c *cli.Context) error {
			client := metadata.GetClient(c.App.Metadata)
			renderer := metadata.GetRenderer(c.App.Metadata)

			res, err := client.ListWebhooks(c.Context)
			if err != nil {
				return cli.Exit(err.Error(), codes.ClientError)
			}

			b, err := renderer.Webhooks(res)
			if err != nil {
				return cli.Exit(err.Error(), codes.OutputError)
			}
			fmt.Fprintln(os.Stdout, b)

			return nil
		},
	}
}
