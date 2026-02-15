package messages

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func sendCmd() *cli.Command {
	const defaultDataPort = 53739

	return &cli.Command{
		Name:      "send",
		Usage:     "Send message",
		Args:      true,
		ArgsUsage: "Message content",
		Category:  "Messages",
		Flags: []cli.Flag{
			// Body fields
			&cli.StringFlag{
				Name:     "id",
				Usage:    "Message ID",
				Required: false,
			},
			&cli.StringFlag{
				Name:        "device-id",
				Aliases:     []string{"device", "deviceId"},
				Usage:       "Optional device ID for explicit selection",
				DefaultText: "auto",
			},
			&cli.StringSliceFlag{
				Name:     "phones",
				Aliases:  []string{"p", "phone"},
				Usage:    "Phone numbers (E.164 format, e.g. +19162255887)",
				Required: true,
			},
			&cli.IntFlag{
				Name:        "sim-number",
				Aliases:     []string{"simNumber", "sim"},
				Usage:       "SIM card index (one-based index, e.g. 1)",
				DefaultText: "see device settings",
			},
			&cli.BoolFlag{
				Name:    "delivery-report",
				Aliases: []string{"deliveryReport"},
				Usage:   "Enable delivery report",
				Value:   true,
			},
			&cli.IntFlag{
				Name:  "priority",
				Usage: "Priority, use >= 100 to bypass all limits and delays (-128 to 127)",
				Value: 0,
			},
			&cli.DurationFlag{
				Name:        "ttl",
				Usage:       "Time to live (duration, e.g. 1h30m)",
				DefaultText: "unlimited",
			},
			&cli.TimestampFlag{
				Name:     "valid-until",
				Aliases:  []string{"validUntil"},
				Usage:    "Valid until (RFC3339 format, e.g. 2006-01-02T15:04:05Z07:00)",
				Layout:   time.RFC3339,
				Timezone: time.Local,
			},

			// Data Message
			&cli.BoolFlag{
				Name:  "data",
				Usage: "Send data message instead of text, content should be base64",
				Value: false,
			},
			&cli.UintFlag{
				Name:    "data-port",
				Aliases: []string{"dataPort"},
				Usage:   "Destination port for data message (1 to 65535)",
				Value:   defaultDataPort,
			},

			// Query params
			&cli.BoolFlag{
				Name:    "skip-phone-validation",
				Aliases: []string{"skipPhoneValidation"},
				Usage:   "Skip phone number validation (default: false)",
				Value:   false,
			},
			&cli.UintFlag{
				Name:    "device-active-within",
				Aliases: []string{"deviceActiveWithin"},
				Usage:   "Filter devices active within the specified number of hours",
				Value:   0,
			},
		},
		Before: sendBefore,
		Action: sendAction,
	}
}

func sendBefore(c *cli.Context) error {
	sim := c.Int("sim-number")
	if sim < 0 || sim > 255 {
		return cli.Exit("SIM card index must be between 0 and 255 (0 for default)", codes.ParamsError)
	}

	ttl := c.Duration("ttl")
	validUntil := c.Timestamp("valid-until")
	if ttl != 0 && validUntil != nil {
		return cli.Exit("TTL and Valid Until flags are mutually exclusive", codes.ParamsError)
	}

	if ttl < 0 {
		return cli.Exit("TTL must be positive", codes.ParamsError)
	}

	if validUntil != nil && validUntil.Before(time.Now()) {
		return cli.Exit("Valid Until must be in the future", codes.ParamsError)
	}

	priority := c.Int("priority")
	if priority < int(smsgateway.PriorityMinimum) || priority > int(smsgateway.PriorityMaximum) {
		return cli.Exit(
			fmt.Sprintf(
				"Priority must be between %d and %d",
				smsgateway.PriorityMinimum,
				smsgateway.PriorityMaximum,
			),
			codes.ParamsError,
		)
	}

	isDataMessage := c.Bool("data")
	if !isDataMessage {
		return nil
	}

	dataPort := c.Uint("data-port")
	if dataPort < 1 || dataPort > 65535 {
		return cli.Exit("Data port must be between 1 and 65535", codes.ParamsError)
	}

	data := strings.TrimSpace(c.Args().Get(0))
	if data == "" {
		return cli.Exit("Message is empty", codes.ParamsError)
	}
	if _, err := base64.StdEncoding.DecodeString(data); err != nil {
		if _, err2 := base64.RawStdEncoding.DecodeString(data); err2 != nil {
			return cli.Exit("Invalid base64 data", codes.ParamsError)
		}
	}

	return nil
}

func sendAction(c *cli.Context) error {
	msg := c.Args().Get(0)
	if msg == "" {
		return cli.Exit("Message is empty", codes.ParamsError)
	}

	client := metadata.GetClient(c.App.Metadata)
	renderer := metadata.GetRenderer(c.App.Metadata)

	isDataMessage := c.Bool("data")
	var dataMessage *smsgateway.DataMessage
	var textMessage *smsgateway.TextMessage
	if isDataMessage {
		dataMessage = &smsgateway.DataMessage{
			Data: msg,
			Port: uint16(c.Uint("data-port")), //nolint:gosec // validated
		}
	} else {
		textMessage = &smsgateway.TextMessage{
			Text: msg,
		}
	}

	var simNumber *uint8
	if sim := c.Int("sim-number"); sim > 0 {
		simNumber = lo.ToPtr(uint8(sim)) //nolint:gosec // validated
	}

	var ttl *uint64
	if ttlRaw := uint64(c.Duration("ttl").Seconds()); ttlRaw > 0 {
		ttl = &ttlRaw
	}

	withDeliveryReport := c.Bool("delivery-report")

	req := smsgateway.Message{
		ID:       c.String("id"),
		DeviceID: c.String("device-id"),

		Message:     "",
		TextMessage: textMessage,
		DataMessage: dataMessage,

		PhoneNumbers: c.StringSlice("phones"),
		IsEncrypted:  false,

		SimNumber:          simNumber,
		WithDeliveryReport: &withDeliveryReport,
		Priority:           smsgateway.MessagePriority(c.Int("priority")), //nolint:gosec // validated

		TTL:        ttl,
		ValidUntil: c.Timestamp("valid-until"),
	}

	options := []smsgateway.SendOption{}

	if c.Bool("skip-phone-validation") {
		options = append(options, smsgateway.WithSkipPhoneValidation(true))
	}
	if v := c.Uint("device-active-within"); v > 0 {
		options = append(options, smsgateway.WithDeviceActiveWithin(v))
	}

	res, err := client.Send(c.Context, req, options...)
	if err != nil {
		return cli.Exit(err.Error(), codes.ClientError)
	}

	s, err := renderer.MessageState(res)
	if err != nil {
		return cli.Exit(err.Error(), codes.OutputError)
	}
	fmt.Fprintln(os.Stdout, s)

	return nil
}
