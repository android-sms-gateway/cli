package output

import (
	"encoding/json"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

type JSONOutput struct {
	marshaler func(any) (string, error)
}

func NewJSONOutput() *JSONOutput {
	return &JSONOutput{
		marshaler: func(a any) (string, error) {
			b, err := json.MarshalIndent(a, "", "  ")
			if err != nil {
				return "", err
			}
			return string(b), nil
		},
	}
}

func (o *JSONOutput) MessageState(src smsgateway.MessageState) (string, error) {
	return o.marshaler(src)
}
