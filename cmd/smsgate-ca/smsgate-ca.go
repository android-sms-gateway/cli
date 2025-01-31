package main

import (
	"fmt"
	"os"

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

	app := &cli.App{
		Name:     "smsgate-ca",
		Usage:    "CLI interface for interacting with Certificate Authority of SMS Gateway for Android™",
		Version:  version,
		Commands: make([]*cli.Command, 0, len(ca.Commands)),
		Authors: []*cli.Author{
			{
				Name:  "Aleksandr Soloshenko",
				Email: "support@sms-gate.app",
			},
		},
		Before: func(c *cli.Context) error {
			c.App.Metadata[metadata.CAClientKey] = client.NewClient()
			return nil
		},
	}

	app.Commands = append(app.Commands, ca.Commands...)

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(codes.ParamsError)
	}
}
