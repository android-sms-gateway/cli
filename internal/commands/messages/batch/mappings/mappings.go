package mappings

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/android-sms-gateway/cli/pkg/io/tabular"
	"github.com/samber/lo"
)

func ParseColumnMapping(raw string) (map[string]string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, fmt.Errorf("%w: empty map", ErrMappingParseFailed)
	}

	const piecesCount = 2
	allowed := map[string]struct{}{
		"id":         {},
		"phone":      {},
		"text":       {},
		"device_id":  {},
		"sim_number": {},
		"priority":   {},
	}
	result := map[string]string{}
	for part := range strings.SplitSeq(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		pieces := strings.SplitN(part, "=", piecesCount)
		if len(pieces) != piecesCount {
			return nil, fmt.Errorf("%w: %q, expected key=value", ErrMappingParseFailed, part)
		}

		key := strings.TrimSpace(strings.ToLower(pieces[0]))
		column := strings.TrimSpace(pieces[1])
		if key == "" || column == "" {
			return nil, fmt.Errorf("%w: %q, key and value must be non-empty", ErrMappingParseFailed, part)
		}
		if _, ok := allowed[key]; !ok {
			return nil, fmt.Errorf("%w: unsupported mapping field %q", ErrMappingParseFailed, key)
		}
		if _, exists := result[key]; exists {
			return nil, fmt.Errorf("%w: duplicate mapping field %q", ErrMappingParseFailed, key)
		}

		result[key] = column
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("%w: empty map", ErrMappingParseFailed)
	}

	return result, nil
}

func MapAndValidateRows(records []tabular.Record, mapping map[string]string) ([]SendRow, []string) {
	rows := make([]SendRow, 0, len(records))
	errs := make([]string, 0)

	for _, record := range records {
		row, err := mapRow(record, mapping)
		if err != nil {
			errs = append(errs, fmt.Sprintf("row %d: %v", record.RowNumber, err))
			continue
		}
		rows = append(rows, row)
	}

	return rows, errs
}

func parsePriority(record tabular.Record, mapping map[string]string) (*int8, error) {
	column, ok := mapping["priority"]
	if !ok {
		return nil, nil //nolint:nilnil // value is not provided
	}

	value, exists := record.Values[column]
	if !exists {
		return nil, fmt.Errorf("%w: invalid priority mapping: column %q not found", ErrMappingParseFailed, column)
	}
	raw := strings.TrimSpace(value)
	if raw == "" {
		return nil, nil //nolint:nilnil // value is not provided
	}

	priority, err := strconv.ParseInt(raw, 10, 8)
	if err != nil {
		return nil, fmt.Errorf("invalid priority: %w", err)
	}

	return lo.ToPtr(int8(priority)), nil
}

func mapRow(record tabular.Record, mapping map[string]string) (SendRow, error) {
	row := SendRow{
		RowNumber: record.RowNumber,
		ID:        strings.TrimSpace(record.Values[mapping["id"]]),
		Phone:     strings.TrimSpace(record.Values[mapping["phone"]]),
		Text:      strings.TrimSpace(record.Values[mapping["text"]]),
		DeviceID:  strings.TrimSpace(record.Values[mapping["device_id"]]),
		SimNumber: nil,
		Priority:  nil,
	}

	if row.Phone == "" {
		return SendRow{}, fmt.Errorf("%w: phone is empty", ErrValidationFailed)
	}
	if row.Text == "" {
		return SendRow{}, fmt.Errorf("%w: text is empty", ErrValidationFailed)
	}

	if column, ok := mapping["sim_number"]; ok {
		raw := strings.TrimSpace(record.Values[column])
		if raw != "" {
			sim, err := strconv.ParseUint(raw, 10, 8)
			if err != nil {
				return SendRow{}, fmt.Errorf("invalid sim_number: %w", err)
			}
			sim8 := uint8(sim)
			row.SimNumber = &sim8
		}
	}

	if priority, err := parsePriority(record, mapping); err != nil {
		return SendRow{}, err
	} else if priority != nil {
		row.Priority = priority
	}

	return row, nil
}
