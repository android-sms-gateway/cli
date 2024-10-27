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
	builder.WriteString("\nState: ")
	builder.WriteString(string(src.State))
	builder.WriteString("\nIsHashed: ")
	builder.WriteString(boolToString(src.IsHashed))
	builder.WriteString("\nIsEncrypted: ")
	builder.WriteString(boolToString(src.IsEncrypted))

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
			builder.WriteString(v.Local().Format(time.DateTime))
		}
	}

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

	return builder.String(), nil
}
