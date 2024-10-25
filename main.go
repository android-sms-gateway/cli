package main

import (
	"log"
	"os"

	"github.com/android-sms-gateway/cli/internal/config"
	"github.com/urfave/cli/v2"
)

var (
	version = "0.0.0"
)

func main() {
	app := &cli.App{
		Name:     "sms-gate-cli",
		Usage:    "CLI interface for working with SMS Gateway for Android™",
		Version:  version,
		Commands: []*cli.Command{},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "endpoint",
				DefaultText: config.DefaultEndpoint,
				Usage:       "Endpoint",
				Aliases:     []string{"e"},
				EnvVars: []string{
					"ASG_ENDPOINT",
				},
			},
			&cli.StringFlag{
				Name:    "username",
				Aliases: []string{"u"},
				Usage:   "SMS Gateway username",
				EnvVars: []string{
					"ASG_USERNAME",
				},
				Required: true,
			},
			&cli.StringFlag{
				Name:    "password",
				Aliases: []string{"p"},
				Usage:   "SMS Gateway password",
				EnvVars: []string{
					"ASG_PASSWORD",
				},
				Required: true,
			},
		},
		Authors: []*cli.Author{
			{
				Name:  "Aleksandr Soloshenko",
				Email: "support@sms-gate.app",
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
