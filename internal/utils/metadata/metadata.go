package metadata

import (
	"github.com/android-sms-gateway/cli/internal/core/output"
	"github.com/android-sms-gateway/client-go/smsgateway"
)

const (
	ClientKey   = "client"
	RendererKey = "renderer"
)

func GetClient(metadata map[string]any) *smsgateway.Client {
	return metadata[ClientKey].(*smsgateway.Client)
}

func GetRenderer(metadata map[string]any) output.Renderer {
	return metadata[RendererKey].(output.Renderer)
}
