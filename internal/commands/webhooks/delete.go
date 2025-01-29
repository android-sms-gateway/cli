package webhooks

import (
	"fmt"

	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/urfave/cli/v2"
)

var delete = &cli.Command{
	Category:  "Webhooks",
	Name:      "delete",
	Aliases:   []string{"d"},
	Usage:     "Delete webhook",
	Args:      true,
	ArgsUsage: "ID",
	Action: func(c *cli.Context) error {
		id := c.Args().Get(0)
		if id == "" {
			return cli.Exit("ID is empty", codes.ParamsError)
		}

		client := metadata.GetClient(c.App.Metadata)
		renderer := metadata.GetRenderer(c.App.Metadata)

		err := client.DeleteWebhook(c.Context, id)
		if err != nil {
			return cli.Exit(err.Error(), codes.ClientError)
		}

		b, err := renderer.Success()
		if err != nil {
			return cli.Exit(err.Error(), codes.OutputError)
		}
		fmt.Println(b)
		return nil
	},
}
