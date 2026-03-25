package tabular

import "context"

// Record is a normalized row from a tabular input source.
type Record struct {
	RowNumber int
	Values    map[string]string
}

// Reader defines the contract for reading tabular records.
type Reader interface {
	Read(ctx context.Context) ([]Record, error)
}
