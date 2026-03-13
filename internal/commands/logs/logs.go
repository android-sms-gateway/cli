package logs

import (
	"fmt"
	"os"

	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/flags"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/urfave/cli/v2"
)

func Commands() []*cli.Command {
	return []*cli.Command{logsCmd()}
}

func logsCmd() *cli.Command {
	f := flags.Period()

	return &cli.Command{
		Name:     "logs",
		Aliases:  []string{"log"},
		Usage:    "Get logs for a specific time range",
		Category: "Logs",
		Flags:    f,
		Before:   logsBefore,
		Action:   logsAction,
	}
}

func logsBefore(c *cli.Context) error {
	period := flags.ParsePeriodFlags(c)

	if period.From == nil || period.To == nil {
		return cli.Exit("From and To dates are required", codes.ParamsError)
	}

	if period.From.After(*period.To) {
		return cli.Exit("From date must be less than or equal to To date", codes.ParamsError)
	}

	return nil
}

func logsAction(c *cli.Context) error {
	period := flags.ParsePeriodFlags(c)

	from := period.From
	to := period.To

	client := metadata.GetClient(c.App.Metadata)
	renderer := metadata.GetRenderer(c.App.Metadata)

	res, err := client.GetLogs(c.Context, *from, *to)
	if err != nil {
		return cli.Exit(err.Error(), codes.ClientError)
	}

	output, err := renderer.Logs(res)
	if err != nil {
		return cli.Exit(err.Error(), codes.OutputError)
	}
	fmt.Fprintln(os.Stdout, output)

	return nil
}
