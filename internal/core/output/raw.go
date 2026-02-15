package output

import (
	"encoding/json"
	"fmt"
)

func NewRawOutput() *JSONOutput {
	return &JSONOutput{
		marshaler: func(a any) (string, error) {
			b, err := json.Marshal(a)
			if err != nil {
				return "", fmt.Errorf("failed to marshal: %w", err)
			}
			return string(b), nil
		},
	}
}
