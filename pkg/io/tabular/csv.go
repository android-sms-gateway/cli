package tabular

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

const defaultDelimiter = ','

type CSVConfig struct {
	Path      string
	Delimiter rune
	HasHeader bool
}

type CSVReader struct {
	cfg CSVConfig
}

func NewCSVReader(cfg CSVConfig) *CSVReader {
	if cfg.Delimiter == 0 {
		cfg.Delimiter = defaultDelimiter
	}

	return &CSVReader{cfg: cfg}
}

func (r *CSVReader) Read(ctx context.Context) ([]Record, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("read csv rows: %w", err)
	}

	f, err := os.Open(r.cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("open csv file: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.Comma = r.cfg.Delimiter

	start := 1
	var headers []string
	if r.cfg.HasHeader {
		row, rdErr := reader.Read()
		if rdErr != nil {
			return nil, fmt.Errorf("read csv header: %w", rdErr)
		}
		headers = normalizeHeaders(row)
		start = 2
	}

	records := make([]Record, 0)
	for {
		line, rdErr := reader.Read()
		if rdErr != nil {
			if errors.Is(rdErr, io.EOF) {
				break
			}
			return nil, fmt.Errorf("read csv rows: %w", rdErr)
		}

		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("read csv rows: %w", ctxErr)
		}

		values := map[string]string{}

		for col, v := range line {
			key := fmt.Sprintf("col_%d", col+1)
			if len(headers) > col {
				key = headers[col]
			}
			values[key] = strings.TrimSpace(v)
		}

		records = append(records, Record{
			RowNumber: len(records) + start,
			Values:    values,
		})
	}

	return records, nil
}

func normalizeHeaders(src []string) []string {
	headers := make([]string, 0, len(src))
	used := map[string]struct{}{}

	for i, h := range src {
		base := strings.TrimSpace(h)
		if base == "" {
			base = fmt.Sprintf("col_%d", i+1)
		}

		key := base
		suffix := 1
		for {
			if _, exists := used[key]; !exists {
				break
			}
			suffix++
			key = fmt.Sprintf("%s_%d", base, suffix)
		}
		used[key] = struct{}{}

		headers = append(headers, key)
	}

	return headers
}

var _ Reader = (*CSVReader)(nil)
