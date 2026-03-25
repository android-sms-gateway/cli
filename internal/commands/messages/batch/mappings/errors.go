package mappings

import "errors"

var (
	ErrMappingParseFailed = errors.New("mapping parse failed")
	ErrValidationFailed   = errors.New("validation failed")
)
