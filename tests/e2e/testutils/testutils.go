package testutils

import (
	"net/http"
	"net/http/httptest"
)

func CreateMockServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}
