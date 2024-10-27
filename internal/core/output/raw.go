package output

import "encoding/json"

func NewRawOutput() *JSONOutput {
	return &JSONOutput{
		marshaler: func(a any) (string, error) {
			b, err := json.Marshal(a)
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
	}
}
