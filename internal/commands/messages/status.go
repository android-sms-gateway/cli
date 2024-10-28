package messages

import (
	"fmt"

	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/urfave/cli/v2"
)

var status = &cli.Command{
	Name:      "status",
	Aliases:   []string{"state"},
	Usage:     "Get message status",
	Args:      true,
	ArgsUsage: "Message ID",
	Category:  "Messages",
	Action: func(c *cli.Context) error {
		id := c.Args().Get(0)
		if id == "" {
			return cli.Exit("Message ID is empty", codes.ParamsError)
		}

		client := metadata.GetClient(c.App.Metadata)
		renderer := metadata.GetRenderer(c.App.Metadata)

		res, err := client.GetState(c.Context, id)
		if err != nil {
			return cli.Exit(err.Error(), codes.ClientError)
		}

		s, err := renderer.MessageState(res)
		if err != nil {
			return cli.Exit(err.Error(), codes.OutputError)
		}
		fmt.Println(s)

		return nil
	},
}
