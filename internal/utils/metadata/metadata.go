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
	v, ok := metadata[ClientKey].(*smsgateway.Client)
	if !ok {
		return nil
	}
	return v
}

func GetCAClient(metadata map[string]any) *ca.Client {
	v, ok := metadata[CAClientKey].(*ca.Client)
	if !ok {
		return nil
	}
	return v
}

func GetRenderer(metadata map[string]any) output.Renderer {
	v, ok := metadata[RendererKey].(output.Renderer)
	if !ok {
		return nil
	}
	return v
}
