package messages

import (
	"fmt"
	"time"

	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/urfave/cli/v2"
)

var send = &cli.Command{
	Name:      "send",
	Usage:     "Send message",
	Args:      true,
	ArgsUsage: "Message content",
	Category:  "Messages",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "id",
			Usage:    "Message ID",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "phones",
			Aliases:  []string{"p"},
			Usage:    "Phone numbers",
			Required: true,
		},
		&cli.IntFlag{
			Name:  "sim",
			Usage: "SIM card index (1-3)",
		},
		&cli.DurationFlag{
			Name:        "ttl",
			Usage:       "Time to live",
			DefaultText: "unlimited",
		},
		&cli.TimestampFlag{
			Name:   "validUntil",
			Usage:  "Valid until",
			Layout: time.DateTime,
		},
	},
	Before: func(c *cli.Context) error {
		ttl := c.Duration("ttl")
		validUntil := c.Timestamp("validUntil")
		if ttl > 0 && validUntil != nil {
			return cli.Exit("TTL and Valid Until flags are mutually exclusive", 1)
		}

		return nil
	},
	Action: func(c *cli.Context) error {
		msg := c.Args().Get(0)
		if msg == "" {
			return cli.Exit("Message is empty", codes.ParamsError)
		}

		client := metadata.GetClient(c.App.Metadata)
		renderer := metadata.GetRenderer(c.App.Metadata)

		req := smsgateway.Message{
			ID:           c.String("id"),
			Message:      msg,
			PhoneNumbers: c.StringSlice("phones"),
		}

		if sim := uint8(c.Int("sim")); sim > 0 {
			req.SimNumber = &sim
		}
		if ttl := uint64(c.Duration("ttl").Seconds()); ttl > 0 {
			req.TTL = &ttl
		}
		if validUntil := c.Timestamp("validUntil"); validUntil != nil {
			req.ValidUntil = validUntil
		}

		res, err := client.Send(c.Context, req)
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
