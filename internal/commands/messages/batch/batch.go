package batch

import "github.com/urfave/cli/v2"

func Commands() *cli.Command {
	return &cli.Command{
		Name:     "batch",
		Usage:    "Bulk message operations",
		Category: "Messages",
		Subcommands: []*cli.Command{
			batchSendCmd(),
		},
	}
}
