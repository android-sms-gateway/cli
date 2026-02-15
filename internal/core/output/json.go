package output

import (
	"encoding/json"
	"fmt"

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
				return "", fmt.Errorf("failed to marshal: %w", err)
			}
			return string(b), nil
		},
	}
}

func (o *JSONOutput) MessageState(src smsgateway.MessageState) (string, error) {
	return o.marshaler(src)
}

func (o *JSONOutput) Webhook(src smsgateway.Webhook) (string, error) {
	return o.marshaler(src)
}

func (o *JSONOutput) Webhooks(src []smsgateway.Webhook) (string, error) {
	return o.marshaler(src)
}

func (o *JSONOutput) Success() (string, error) {
	return "", nil
}

var _ Renderer = (*JSONOutput)(nil)
