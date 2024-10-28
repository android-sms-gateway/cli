package webhooks

import (
	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/urfave/cli/v2"
)

var register = &cli.Command{
	Category:  "Webhooks",
	Name:      "register",
	Aliases:   []string{"r"},
	ArgsUsage: "URL",
	Usage:     "Register webhook",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "id",
			Usage:    "ID",
			Required: false,
		},
		&cli.StringFlag{
			Name:     "event",
			Aliases:  []string{"e"},
			Usage:    "Event",
			Required: true,
			Action: func(c *cli.Context, event string) error {
				if !smsgateway.IsValidWebhookEvent(event) {
					return cli.Exit("Invalid event", codes.ParamsError)
				}

				return nil
			},
		},
	},
	Action: func(c *cli.Context) error {
		url := c.Args().Get(0)
		if url == "" {
			return cli.Exit("URL is empty", codes.ParamsError)
		}

		// client := metadata.GetClient(c.App.Metadata)
		// renderer := metadata.GetRenderer(c.App.Metadata)

		// req := smsgateway.Webhook{
		// 	ID:    c.String("id"),
		// 	URL:   url,
		// 	Event: c.String("event"),
		// }

		return cli.Exit("Not implemented", codes.ParamsError)
	},
}
