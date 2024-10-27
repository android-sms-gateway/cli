package main

import (
	"log"
	"os"

	"github.com/android-sms-gateway/cli/internal/commands/messages"
	"github.com/android-sms-gateway/cli/internal/commands/webhooks"
	"github.com/android-sms-gateway/cli/internal/config"
	"github.com/android-sms-gateway/cli/internal/core/client"
	"github.com/android-sms-gateway/cli/internal/core/output"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var (
	version = "0.0.0"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:     "sms-gate-cli",
		Usage:    "CLI interface for working with SMS Gateway for Android™",
		Version:  version,
		Commands: make([]*cli.Command, 0, len(messages.Commands)+len(webhooks.Commands)),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "endpoint",
				DefaultText: config.DefaultEndpoint,
				Usage:       "Endpoint",
				Category:    "Configuration",
				Aliases:     []string{"e"},
				Value:       config.DefaultEndpoint,
				EnvVars: []string{
					"ASG_ENDPOINT",
				},
			},
			&cli.StringFlag{
				Name:     "username",
				Aliases:  []string{"u"},
				Usage:    "Username",
				Category: "Configuration",
				EnvVars: []string{
					"ASG_USERNAME",
				},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "password",
				Aliases:  []string{"p"},
				Usage:    "Password",
				Category: "Configuration",
				EnvVars: []string{
					"ASG_PASSWORD",
				},
				Required: true,
			},

			&cli.StringFlag{
				Name:     "format",
				Category: "Output",
				Usage:    "Output format. Supported: text, json, raw",
				Required: false,
				Value:    string(output.Text),
				Aliases:  []string{"f"},
			},
		},
		Authors: []*cli.Author{
			{
				Name:  "Aleksandr Soloshenko",
				Email: "support@sms-gate.app",
			},
		},
		Before: func(c *cli.Context) error {
			renderer, err := output.New(output.Format(c.String("format")))
			if err != nil {
				return err
			}

			c.App.Metadata[metadata.RendererKey] = renderer
			c.App.Metadata[metadata.ClientKey] = client.New(c.String("username"), c.String("password"), c.String("endpoint"))
			return nil
		},
	}

	app.Commands = append(app.Commands, messages.Commands...)
	// app.Commands = append(app.Commands, webhooks.Commands...)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
