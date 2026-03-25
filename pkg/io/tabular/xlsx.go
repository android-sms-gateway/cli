package tabular

import (
	"context"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

type XLSXConfig struct {
	Path      string
	Sheet     string
	HasHeader bool
}

type XLSXReader struct {
	cfg XLSXConfig
}

func NewXLSXReader(cfg XLSXConfig) *XLSXReader {
	return &XLSXReader{cfg: cfg}
}

func (r *XLSXReader) Read(ctx context.Context) ([]Record, error) {
	if err := ctx.Err(); err != nil {
		return nil, fmt.Errorf("read xlsx rows: %w", err)
	}

	f, err := excelize.OpenFile(r.cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("open xlsx file: %w", err)
	}
	defer f.Close()

	sheet := strings.TrimSpace(r.cfg.Sheet)
	if sheet == "" {
		sheetList := f.GetSheetList()
		if len(sheetList) == 0 {
			return []Record{}, nil
		}
		sheet = sheetList[0]
	}

	rows, err := f.Rows(sheet)
	if err != nil {
		return nil, fmt.Errorf("read xlsx rows: %w", err)
	}
	defer rows.Close()

	start := 1
	var headers []string
	if r.cfg.HasHeader {
		if !rows.Next() {
			return []Record{}, nil
		}

		line, colErr := rows.Columns()
		if colErr != nil {
			return nil, fmt.Errorf("read xlsx header: %w", colErr)
		}

		headers = normalizeHeaders(line)
		start = 2
	}

	return r.readRecords(ctx, headers, start, rows)
}

func (r *XLSXReader) readRecords(
	ctx context.Context,
	headers []string,
	start int,
	rows *excelize.Rows,
) ([]Record, error) {
	records := make([]Record, 0)
	for rows.Next() {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, fmt.Errorf("read xlsx rows: %w", ctxErr)
		}

		line, rdErr := rows.Columns()
		if rdErr != nil {
			return nil, fmt.Errorf("read xlsx row: %w", rdErr)
		}

		values := map[string]string{}
		for col, v := range line {
			key := fmt.Sprintf("col_%d", col+1)
			if len(headers) > col {
				key = headers[col]
			}
			values[key] = strings.TrimSpace(v)
		}

		records = append(records, Record{RowNumber: len(records) + start, Values: values})
	}

	return records, nil
}

var _ Reader = (*XLSXReader)(nil)
