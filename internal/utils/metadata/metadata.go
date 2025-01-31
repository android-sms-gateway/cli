package metadata

import (
	"github.com/android-sms-gateway/cli/internal/core/output"
	"github.com/android-sms-gateway/client-go/ca"
	"github.com/android-sms-gateway/client-go/smsgateway"
)

const (
	ClientKey   = "client"
	CAClientKey = "caclient"
	RendererKey = "renderer"
)

func GetClient(metadata map[string]any) *smsgateway.Client {
	return metadata[ClientKey].(*smsgateway.Client)
}

func GetCAClient(metadata map[string]any) *ca.Client {
	return metadata[CAClientKey].(*ca.Client)
}

func GetRenderer(metadata map[string]any) output.Renderer {
	return metadata[RendererKey].(output.Renderer)
}
