package webhooks

import (
	"github.com/urfave/cli/v2"
)

func Commands() []*cli.Command {
	return []*cli.Command{
		{
			Category: "Webhooks",
			Name:     "webhooks",
			Aliases:  []string{"w", "wh"},
			Usage:    "Manage webhooks",
			Subcommands: []*cli.Command{
				registerCmd(),
				deleteCmd(),
				listCmd(),
			},
		},
	}
}
