package mappings

// SendRow is a normalized row for the batch send flow.
type SendRow struct {
	RowNumber int

	ID       string
	Phone    string
	Text     string
	DeviceID string

	SimNumber *uint8
	Priority  *int8
}
