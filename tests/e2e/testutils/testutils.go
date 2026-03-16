package testutils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func CreateMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

func RequireBinPath(t *testing.T) string {
	t.Helper()
	binPath := os.Getenv("SMSGATE_BIN")
	if binPath == "" {
		t.Fatal("SMSGATE_BIN environment variable is not set")
	}
	return binPath
}
