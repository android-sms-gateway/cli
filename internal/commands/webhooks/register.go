package webhooks

import (
	"fmt"
	"strings"

	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/capcom6/go-helpers/slices"
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
			Usage:    "ID, optional",
			Required: false,
		},
		&cli.StringFlag{
			Name:    "event",
			Aliases: []string{"e"},
			Usage: "Event, one of: " + strings.Join(
				slices.Map(smsgateway.WebhookEventTypes(), func(e smsgateway.WebhookEvent) string { return string(e) }),
				", ",
			),
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

		client := metadata.GetClient(c.App.Metadata)
		renderer := metadata.GetRenderer(c.App.Metadata)

		req := smsgateway.Webhook{
			ID:    c.String("id"),
			URL:   url,
			Event: c.String("event"),
		}

		res, err := client.RegisterWebhook(c.Context, req)
		if err != nil {
			return cli.Exit(err.Error(), codes.ClientError)
		}

		b, err := renderer.Webhook(res)
		if err != nil {
			return cli.Exit(err.Error(), codes.OutputError)
		}
		fmt.Println(b)

		return nil
	},
}
