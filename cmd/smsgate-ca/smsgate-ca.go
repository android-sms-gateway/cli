package main

import (
	"fmt"
	"os"
	"time"

	"github.com/android-sms-gateway/cli/internal/commands/ca"
	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	client "github.com/android-sms-gateway/client-go/ca"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var (
	version = "0.0.0"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error loading .env file: %v\n", err)
		os.Exit(codes.ParamsError)
	}

	const defaultTimeout = 30 * time.Second

	cmds := ca.Commands()

	app := &cli.App{
		Name:     "smsgate-ca",
		Usage:    "CLI interface for interacting with Certificate Authority of SMSGate",
		Version:  version,
		Commands: cmds,
		Flags: []cli.Flag{
			&cli.DurationFlag{
				Name:     "timeout",
				Aliases:  []string{"t"},
				Usage:    "Request timeout",
				Required: false,
				Value:    defaultTimeout,
				EnvVars: []string{
					"ASG_CA_TIMEOUT",
				},
			},
		},
		Authors: []*cli.Author{
			{
				Name:  "Aleksandr Soloshenko",
				Email: "support@sms-gate.app",
			},
		},
		Before: func(c *cli.Context) error {
			c.App.Metadata[metadata.CAClientKey] = client.NewClient()

			timeout := c.Duration("timeout")
			if timeout <= 0 {
				return cli.Exit(
					"Timeout must be greater than 0",
					codes.ParamsError,
				)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(codes.ParamsError)
	}
}
