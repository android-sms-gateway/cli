package webhooks

import (
	"github.com/urfave/cli/v2"
)

var Commands = []*cli.Command{
	{
		Category: "Webhooks",
		Name:     "webhooks",
		Aliases:  []string{"w", "wh"},
		Usage:    "Manage webhooks",
		Subcommands: []*cli.Command{
			register,
			delete,
			list,
		},
	},
}
