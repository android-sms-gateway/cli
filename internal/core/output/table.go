package output

import (
	"fmt"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

const tabwriterPadding = 2

type TableOutput struct{}

func NewTableOutput() *TableOutput {
	return &TableOutput{}
}

func (*TableOutput) MessageState(src smsgateway.MessageState) (string, error) {
	var b strings.Builder

	fmt.Fprintf(&b, "ID:\t%s\n", src.ID)
	fmt.Fprintf(&b, "Device ID:\t%s\n", src.DeviceID)
	fmt.Fprintf(&b, "State:\t%s\n", src.State)
	fmt.Fprintf(&b, "IsHashed:\t%s\n", boolToString(src.IsHashed))
	fmt.Fprintf(&b, "IsEncrypted:\t%s\n", boolToString(src.IsEncrypted))

	if len(src.Recipients) > 0 {
		b.WriteString("\nRecipients:\n")
		tw := tabwriter.NewWriter(&b, 0, 0, tabwriterPadding, ' ', 0)
		fmt.Fprintln(tw, "PHONE\tSTATE\tERROR")
		for _, r := range src.Recipients {
			errStr := ""
			if r.Error != nil {
				errStr = *r.Error
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\n", r.PhoneNumber, r.State, errStr)
		}
		if err := tw.Flush(); err != nil {
			return "", fmt.Errorf("flush tabwriter: %w", err)
		}
	}

	if len(src.States) > 0 {
		messageStates := []string{
			string(smsgateway.ProcessingStatePending),
			string(smsgateway.ProcessingStateProcessed),
			string(smsgateway.ProcessingStateSent),
			string(smsgateway.ProcessingStateDelivered),
			string(smsgateway.ProcessingStateFailed),
		}

		b.WriteString("\nStates:\n")
		tw := tabwriter.NewWriter(&b, 0, 0, tabwriterPadding, ' ', 0)
		fmt.Fprintln(tw, "STATE\tTIME")
		for _, k := range messageStates {
			v, ok := src.States[k]
			if !ok {
				continue
			}
			fmt.Fprintf(tw, "%s\t%s\n", k, v.Local().Format(time.RFC3339))
		}
		if err := tw.Flush(); err != nil {
			return "", fmt.Errorf("flush tabwriter: %w", err)
		}
	}

	return strings.TrimRight(b.String(), "\n"), nil
}

func (*TableOutput) Logs(src []smsgateway.LogEntry) (string, error) {
	if len(src) == 0 {
		return EmptyResult, nil
	}

	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, tabwriterPadding, ' ', 0)

	fmt.Fprintln(tw, "ID\tPRIORITY\tMODULE\tMESSAGE\tCREATED AT")
	for _, entry := range src {
		fmt.Fprintf(
			tw,
			"%d\t%s\t%s\t%s\t%s\n",
			entry.ID,
			entry.Priority,
			entry.Module,
			entry.Message,
			entry.CreatedAt.Local().Format(time.RFC3339),
		)
	}

	if err := tw.Flush(); err != nil {
		return "", fmt.Errorf("flush tabwriter: %w", err)
	}

	return strings.TrimRight(b.String(), "\n"), nil
}

func (*TableOutput) Webhook(src smsgateway.Webhook) (string, error) {
	var b strings.Builder

	fmt.Fprintf(&b, "ID:\t%s\n", src.ID)
	fmt.Fprintf(&b, "Event:\t%s\n", src.Event)
	fmt.Fprintf(&b, "URL:\t%s\n", src.URL)

	if src.DeviceID != nil {
		fmt.Fprintf(&b, "Device ID:\t%s\n", *src.DeviceID)
	}

	return b.String(), nil
}

func (o *TableOutput) Webhooks(src []smsgateway.Webhook) (string, error) {
	if len(src) == 0 {
		return EmptyResult, nil
	}

	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, tabwriterPadding, ' ', 0)

	fmt.Fprintln(tw, "ID\tEVENT\tURL\tDEVICE ID")
	for _, w := range src {
		deviceID := ""
		if w.DeviceID != nil {
			deviceID = *w.DeviceID
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", w.ID, w.Event, w.URL, deviceID)
	}

	if err := tw.Flush(); err != nil {
		return "", fmt.Errorf("flush tabwriter: %w", err)
	}

	return strings.TrimRight(b.String(), "\n"), nil
}

func (*TableOutput) Success() (string, error) {
	return "Success", nil
}

var _ Renderer = (*TableOutput)(nil)
