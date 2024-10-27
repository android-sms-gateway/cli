package webhooks

import (
	"log"

	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	{
		Name:        "webhook",
		Description: "Manage webhooks",
		Aliases:     []string{"wh"},
		Subcommands: []*cli.Command{
			{
				Name:      "register",
				Aliases:   []string{"r"},
				Args:      true,
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
					},
				},
				Action: func(c *cli.Context) error {

					u := c.Args().Get(0)
					if u == "" {
						return cli.Exit("URL is empty", 1)
					}

					id := c.String("id")
					log.Printf("Registering webhook: %s with ID: %s", u, id)

					return nil
				},
			},
		},
	},
}
