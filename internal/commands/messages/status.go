package messages

import (
	"log"

	"github.com/urfave/cli/v2"
)

var status = &cli.Command{
	Name:      "status",
	Usage:     "Get message status",
	Args:      true,
	ArgsUsage: "Message ID",
	Category:  "Messages",
	Action: func(c *cli.Context) error {
		id := c.Args().Get(0)
		if id == "" {
			return cli.Exit("Message ID is empty", 1)
		}

		log.Printf("Getting message status: %s", id)
		return nil
	},
}
