package messages

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/android-sms-gateway/cli/internal/commands/flags"
	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/urfave/cli/v2"
)

func sendCmd() *cli.Command {
	const defaultDataPort = 53739

	fl := []cli.Flag{
		// Body fields
		&cli.StringFlag{
			Name:     "id",
			Category: "Body",
			Usage:    "Message ID",
			Required: false,
		},
		&cli.StringSliceFlag{
			Name:     "phones",
			Category: "Body",
			Aliases:  []string{"p", "phone"},
			Usage:    "Phone numbers (E.164 format, e.g. +19162255887)",
			Required: true,
		},

		// Data Message
		&cli.BoolFlag{
			Name:     "data",
			Category: "Data Message",
			Usage:    "Send data message instead of text, content should be base64 encoded",
			Value:    false,
		},
		&cli.UintFlag{
			Name:     "data-port",
			Category: "Data Message",
			Aliases:  []string{"dataPort"},
			Usage:    "Destination port for data message (1 to 65535)",
			Value:    defaultDataPort,
		},
	}
	fl = append(fl, flags.Send()...)

	return &cli.Command{
		Name:      "send",
		Usage:     "Send message",
		Args:      true,
		ArgsUsage: "Message content",
		Category:  "Messages",
		Flags:     fl,
		Before:    sendBefore,
		Action:    sendAction,
	}
}

func sendBefore(c *cli.Context) error {
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
	sendFlags, err := flags.NewSendFlags(c)
	if err != nil {
		return cli.Exit(err.Error(), codes.ParamsError)
	}

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

	req := smsgateway.Message{
		ID:           c.String("id"),
		Message:      "",
		TextMessage:  textMessage,
		DataMessage:  dataMessage,
		PhoneNumbers: c.StringSlice("phones"),
		IsEncrypted:  false,

		DeviceID:           "",
		SimNumber:          nil,
		WithDeliveryReport: nil,
		Priority:           0,
		TTL:                nil,
		ValidUntil:         nil,
		ScheduleAt:         nil,
	}
	req = sendFlags.Merge(req)

	options := sendFlags.Option()

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
