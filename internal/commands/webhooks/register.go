package webhooks

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func registerCmd() *cli.Command {
	return &cli.Command{
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
				Name: "device-id",
				Usage: "Optional device ID for explicit selection. If not set, account-wide " +
					"webhook will be registered.",
				Required: false,
			},
			&cli.StringFlag{
				Name:    "event",
				Aliases: []string{"e"},
				Usage: "Event, one of: " + strings.Join(
					smsgateway.WebhookEventTypes(),
					", ",
				),
				Required: true,
				Action: func(_ *cli.Context, event string) error {
					if !smsgateway.IsValidWebhookEvent(event) {
						return cli.Exit("Invalid event", codes.ParamsError)
					}

					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			targetURL := strings.TrimSpace(c.Args().Get(0))
			if targetURL == "" {
				return cli.Exit("URL is empty", codes.ParamsError)
			}

			// accept only absolute http/https URLs
			parsed, err := url.Parse(targetURL)
			if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
				return cli.Exit("invalid URL", codes.ParamsError)
			}

			var deviceID *string
			if did := c.String("device-id"); did != "" {
				deviceID = lo.ToPtr(did)
			}

			client := metadata.GetClient(c.App.Metadata)
			renderer := metadata.GetRenderer(c.App.Metadata)

			req := smsgateway.Webhook{
				ID:       c.String("id"),
				URL:      targetURL,
				Event:    c.String("event"),
				DeviceID: deviceID,
			}

			res, err := client.RegisterWebhook(c.Context, req)
			if err != nil {
				return cli.Exit(err.Error(), codes.ClientError)
			}

			b, err := renderer.Webhook(res)
			if err != nil {
				return cli.Exit(err.Error(), codes.OutputError)
			}
			fmt.Fprintln(os.Stdout, b)

			return nil
		},
	}
}
