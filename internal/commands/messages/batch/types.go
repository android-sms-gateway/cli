package batch

import "github.com/android-sms-gateway/client-go/smsgateway"

type batchRowResult struct {
	RowNumber  int
	Identifier string
	State      smsgateway.MessageState
	Error      error
}
