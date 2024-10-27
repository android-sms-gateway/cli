package client

import "github.com/android-sms-gateway/client-go/smsgateway"

func New(username, password, endpoint string) *smsgateway.Client {
	return smsgateway.NewClient(smsgateway.Config{
		BaseURL:  endpoint,
		User:     username,
		Password: password,
	})
}
