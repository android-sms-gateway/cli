package output

import (
	"errors"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

type Format string

const (
	Text Format = "text"
	JSON Format = "json"
	RAW  Format = "raw"
)

type Renderer interface {
	MessageState(src smsgateway.MessageState) (string, error)
}

var ErrUnsupportedFormat = errors.New("unsupported format")

func New(format Format) (Renderer, error) {
	switch format {
	case Text:
		return NewTextOutput(), nil
	case JSON:
		return NewJSONOutput(), nil
	case RAW:
		return NewRawOutput(), nil
	default:
		return nil, ErrUnsupportedFormat
	}
}
