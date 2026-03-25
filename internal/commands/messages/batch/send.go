package batch

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/android-sms-gateway/cli/internal/commands/flags"
	"github.com/android-sms-gateway/cli/internal/commands/messages/batch/mappings"
	"github.com/android-sms-gateway/cli/internal/core/codes"
	"github.com/android-sms-gateway/cli/internal/utils/metadata"
	"github.com/android-sms-gateway/cli/pkg/io/tabular"
	"github.com/android-sms-gateway/client-go/smsgateway"
	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

func batchSendCmd() *cli.Command {
	fl := []cli.Flag{
		&cli.StringFlag{
			Name:     "sheet",
			Category: "Excel",
			Usage:    "XLSX sheet name (defaults to the first sheet)",
			Value:    "",
		},
		&cli.StringFlag{
			Name:     "delimiter",
			Category: "CSV",
			Usage:    "CSV delimiter (single character)",
			Value:    ",",
		},

		&cli.BoolFlag{
			Name:     "header",
			Category: "Input",
			Usage:    "Treat first row as header",
			Value:    true,
		},
		&cli.StringFlag{
			Name:     "map",
			Category: "Input",
			Usage:    "Column mapping, e.g. phone=Phone,text=Message,id=ID,device_id=Device,sim_number=SIM,priority=Priority",
			Required: true,
		},

		&cli.BoolFlag{
			Name:     "dry-run",
			Category: "Preview",
			Usage:    "Validate and print normalized rows without sending",
			Value:    false,
		},
		&cli.BoolFlag{
			Name:     "validate-only",
			Category: "Preview",
			Usage:    "Validate input only (no preview, no sending)",
			Value:    false,
		},

		&cli.IntFlag{
			Name:     "concurrency",
			Category: "Send",
			Usage:    "Number of concurrent send workers",
			Value:    runtime.NumCPU(),
		},
		&cli.BoolFlag{
			Name:     "continue-on-error",
			Category: "Send",
			Usage:    "Continue sending after per-row failures",
			Value:    false,
		},
	}
	fl = append(fl, flags.Send()...)

	return &cli.Command{
		Name:      "send",
		Usage:     "Validate, preview, and send bulk messages from CSV/XLSX",
		ArgsUsage: "filename.[csv|xlsx]",
		Args:      true,
		Flags:     fl,
		Before:    batchSendBefore,
		Action:    batchSendAction,
	}
}

func batchSendBefore(c *cli.Context) error {
	if c.Args().Len() != 1 {
		return cli.Exit("filename is required", codes.ParamsError)
	}

	delimiter := c.String("delimiter")
	if len([]rune(delimiter)) != 1 {
		return cli.Exit("delimiter must be exactly one character", codes.ParamsError)
	}

	if c.Int("concurrency") < 1 {
		return cli.Exit("concurrency must be at least 1", codes.ParamsError)
	}

	mapping, err := mappings.ParseColumnMapping(c.String("map"))
	if err != nil {
		return cli.Exit(err.Error(), codes.ParamsError)
	}

	if _, ok := mapping["phone"]; !ok {
		return cli.Exit("map must include phone=<column>", codes.ParamsError)
	}
	if _, ok := mapping["text"]; !ok {
		return cli.Exit("map must include text=<column>", codes.ParamsError)
	}

	return nil
}

func batchSendAction(c *cli.Context) error {
	renderer := metadata.GetRenderer(c.App.Metadata)
	sendFlags, err := flags.NewSendFlags(c)
	if err != nil {
		return cli.Exit(err.Error(), codes.ParamsError)
	}

	reader, err := newTabularReader(c)
	if err != nil {
		return cli.Exit(err.Error(), codes.ParamsError)
	}

	records, err := reader.Read(c.Context)
	if err != nil {
		return cli.Exit(err.Error(), codes.ParamsError)
	}

	mapping, _ := mappings.ParseColumnMapping(c.String("map"))
	rows, errs := mappings.MapAndValidateRows(records, mapping)
	if len(errs) > 0 {
		return cli.Exit(strings.Join(errs, "\n"), codes.ParamsError)
	}

	if c.Bool("validate-only") {
		fmt.Fprintf(os.Stderr, "Validation passed: %d rows\n", len(rows))
		return nil
	}

	if c.Bool("dry-run") {
		fmt.Fprintf(os.Stderr, "Dry run successful: %d rows\n", len(rows))
		for _, row := range rows {
			fmt.Fprintf(
				os.Stderr,
				"row=%d id=%q phone=%q text=%q device_id=%q\n",
				row.RowNumber,
				row.ID,
				row.Phone,
				row.Text,
				row.DeviceID,
			)
		}
		return nil
	}

	client := metadata.GetClient(c.App.Metadata)
	results := runBatchSend(
		c.Context,
		client,
		rows,
		sendFlags,
		c.Int("concurrency"),
		c.Bool("continue-on-error"),
	)

	failed := 0
	for _, result := range results {
		if result.Error != nil {
			failed++
		}
	}

	skipped := len(rows) - len(results)
	sent := len(rows) - failed - skipped
	fmt.Fprintf(
		os.Stderr,
		"Batch send summary: total=%d enqueued=%d failed=%d skipped=%d\n",
		len(rows),
		sent,
		failed,
		skipped,
	)

	for _, result := range results {
		if result.Error != nil {
			continue
		}

		line, renderErr := renderer.MessageState(result.State)
		if renderErr != nil {
			return cli.Exit(renderErr.Error(), codes.OutputError)
		}

		if _, outErr := fmt.Fprintln(
			os.Stdout,
			line,
		); outErr != nil {
			return cli.Exit(outErr.Error(), codes.OutputError)
		}
	}

	if failed > 0 {
		return cli.Exit("one or more rows failed", codes.ClientError)
	}

	return nil
}

func newTabularReader(c *cli.Context) (tabular.Reader, error) {
	path := c.Args().Get(0)
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".csv":
		delimiter := []rune(c.String("delimiter"))[0]
		return tabular.NewCSVReader(tabular.CSVConfig{
			Path:      path,
			Delimiter: delimiter,
			HasHeader: c.Bool("header"),
		}), nil
	case ".xlsx", ".xlsm":
		return tabular.NewXLSXReader(tabular.XLSXConfig{
			Path:      path,
			Sheet:     c.String("sheet"),
			HasHeader: c.Bool("header"),
		}), nil
	default:
		return nil, fmt.Errorf("%w: unsupported file extension %q; use .csv or .xlsx", ErrValidationFailed, ext)
	}
}

func runBatchSend(
	ctx context.Context,
	client *smsgateway.Client,
	rows []mappings.SendRow,
	sendFlags *flags.SendFlags,
	concurrency int,
	continueOnError bool,
) []batchRowResult {
	workerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	jobs := make(chan mappings.SendRow)
	results := make(chan batchRowResult, len(rows))

	var wg sync.WaitGroup
	for range concurrency {
		wg.Go(func() {
			sendWorker(workerCtx, client, jobs, results, sendFlags, continueOnError, cancel)
		})
	}

	go func() {
		defer close(jobs)
		for _, row := range rows {
			if workerCtx.Err() != nil {
				return
			}
			jobs <- row
		}
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]batchRowResult, 0, len(rows))
	for result := range results {
		var state string
		if result.Error != nil {
			state = fmt.Sprintf("failed: %s", result.Error.Error())
		} else {
			state = string(result.State.State)
		}
		fmt.Fprintf(os.Stderr, "[%d] %s: %s\n", result.RowNumber, result.Identifier, state)
		out = append(out, result)
	}

	return out
}

func sendWorker(
	ctx context.Context,
	client *smsgateway.Client,

	jobs <-chan mappings.SendRow,
	results chan<- batchRowResult,

	sendFlags *flags.SendFlags,
	continueOnError bool,
	cancel context.CancelFunc,
) {
	for row := range jobs {
		identifier := row.ID
		if identifier == "" {
			identifier = uuid.Must(uuid.NewV7()).String()
		}

		req := smsgateway.Message{
			ID:       identifier,
			DeviceID: row.DeviceID,
			Message:  "",
			TextMessage: &smsgateway.TextMessage{
				Text: row.Text,
			},
			DataMessage:        nil,
			PhoneNumbers:       []string{row.Phone},
			IsEncrypted:        false,
			SimNumber:          row.SimNumber,
			WithDeliveryReport: nil,
			Priority:           0,
			TTL:                nil,
			ValidUntil:         nil,
		}

		req = sendFlags.Merge(req)

		if row.Priority != nil {
			req.Priority = smsgateway.MessagePriority(*row.Priority)
		}

		options := sendFlags.Option()

		state, err := client.Send(ctx, req, options...)
		results <- batchRowResult{RowNumber: row.RowNumber, Identifier: identifier, Error: err, State: state}
		if err != nil && !continueOnError {
			cancel()
		}
	}
}
