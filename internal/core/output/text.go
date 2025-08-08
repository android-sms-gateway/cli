package output

import (
	"strings"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

var messageStates = []string{
	string(smsgateway.ProcessingStatePending),
	string(smsgateway.ProcessingStateProcessed),
	string(smsgateway.ProcessingStateSent),
	string(smsgateway.ProcessingStateDelivered),
	string(smsgateway.ProcessingStateFailed),
}

type TextOutput struct {
}

func NewTextOutput() *TextOutput {
	return &TextOutput{}
}

func (*TextOutput) MessageState(src smsgateway.MessageState) (string, error) {
	builder := strings.Builder{}
	builder.WriteString("ID: ")
	builder.WriteString(src.ID)
	builder.WriteString("\nDevice ID: ")
	builder.WriteString(src.DeviceID)
	builder.WriteString("\nState: ")
	builder.WriteString(string(src.State))
	builder.WriteString("\nIsHashed: ")
	builder.WriteString(boolToString(src.IsHashed))
	builder.WriteString("\nIsEncrypted: ")
	builder.WriteString(boolToString(src.IsEncrypted))

	builder.WriteString("\nRecipients: ")

	for _, r := range src.Recipients {
		builder.WriteString("\n\t")
		builder.WriteString(r.PhoneNumber)
		builder.WriteString("\t")
		builder.WriteString(string(r.State))
		builder.WriteString("\t")
		if r.Error != nil {
			builder.WriteString(*r.Error)
		} else {
			builder.WriteString("")
		}
	}

	if len(src.States) > 0 {
		builder.WriteString("\nStates: ")

		for _, k := range messageStates {
			v, ok := src.States[k]
			if !ok {
				continue
			}

			builder.WriteString("\n\t")
			builder.WriteString(k)
			builder.WriteString("\t")
			builder.WriteString(v.Local().Format(time.RFC3339))
		}
	}

	return builder.String(), nil
}

// Webhook formats a single smsgateway.Webhook into a string representation.
// The output includes the ID, Event, and URL of the webhook.
func (*TextOutput) Webhook(src smsgateway.Webhook) (string, error) {
	builder := strings.Builder{}
	builder.WriteString("ID: ")
	builder.WriteString(src.ID)
	builder.WriteString("\nEvent: ")
	builder.WriteString(src.Event)
	builder.WriteString("\nURL: ")
	builder.WriteString(src.URL)

	return builder.String(), nil
}

// Webhooks formats a slice of smsgateway.Webhook into a single string representation.
// Each webhook's string representation is separated by "---".
// Returns the formatted string and any error encountered during formatting.
func (o *TextOutput) Webhooks(src []smsgateway.Webhook) (string, error) {
	if len(src) == 0 {
		return "Empty result", nil
	}

	builder := strings.Builder{}

	for i, w := range src {
		b, err := o.Webhook(w)
		if err != nil {
			return "", err
		}
		builder.WriteString(b)

		if i == len(src)-1 {
			continue
		}

		builder.WriteString("\n---\n")
	}

	return builder.String(), nil
}

// Success returns a string indicating success.
func (*TextOutput) Success() (string, error) {
	return "Success", nil
}

var _ Renderer = (*TextOutput)(nil)
