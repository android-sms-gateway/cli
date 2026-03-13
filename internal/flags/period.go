package flags

import (
	"time"

	"github.com/urfave/cli/v2"
)

func Period() []cli.Flag {
	return []cli.Flag{
		&cli.TimestampFlag{
			Name:     "from",
			Usage:    "Start of time range (RFC3339 format)",
			Layout:   time.RFC3339,
			Timezone: time.Local,
			Value:    cli.NewTimestamp(time.Now().Add(-24 * time.Hour)),
		},
		&cli.TimestampFlag{
			Name:     "to",
			Usage:    "End of time range (RFC3339 format)",
			Layout:   time.RFC3339,
			Timezone: time.Local,
			Value:    cli.NewTimestamp(time.Now()),
		},
	}
}

type PeriodFlags struct {
	From *time.Time
	To   *time.Time
}

func ParsePeriodFlags(c *cli.Context) PeriodFlags {
	return PeriodFlags{
		From: c.Timestamp("from"),
		To:   c.Timestamp("to"),
	}
}
