package output

import (
	"errors"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

type Format string

const (
	Text  Format = "text"
	JSON  Format = "json"
	RAW   Format = "raw"
	Table Format = "table"
)

type Renderer interface {
	MessageState(src smsgateway.MessageState) (string, error)
	Logs(src []smsgateway.LogEntry) (string, error)
	Webhook(src smsgateway.Webhook) (string, error)
	Webhooks(src []smsgateway.Webhook) (string, error)
	Success() (string, error)
}

const EmptyResult = "Empty result"

var ErrUnsupportedFormat = errors.New("unsupported format")

func New(format Format) (Renderer, error) {
	switch format {
	case Text:
		return NewTextOutput(), nil
	case JSON:
		return NewJSONOutput(), nil
	case RAW:
		return NewRawOutput(), nil
	case Table:
		return NewTableOutput(), nil
	default:
		return nil, ErrUnsupportedFormat
	}
}
