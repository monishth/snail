package middleware

import (
	"encoding/base64"
	"github.com/monishth/snail/internal/auth"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/charmbracelet/log"
)

func TestBasicAuthMiddleware(t *testing.T) {
	tests := []struct {
		name               string
		proxyAuthHeader    string
		expectErrorMessage string
	}{
		{
			name:               "NoAuthHeader",
			proxyAuthHeader:    "",
			expectErrorMessage: "Proxy Authorization Required",
		},
		{
			name:               "NonBasicAuth",
			proxyAuthHeader:    "Bearer token",
			expectErrorMessage: "Proxy Authorization Required",
		},
		{
			name:               "NotEnoughCredentials",
			proxyAuthHeader:    "Basic " + base64.StdEncoding.EncodeToString([]byte("username")),
			expectErrorMessage: "Proxy Authorization Required",
		},
		{
			name:               "InvalidCredentials",
			proxyAuthHeader:    "Basic " + base64.StdEncoding.EncodeToString([]byte("username:wrongpass")),
			expectErrorMessage: "Proxy Authorization Required",
		},
		{
			name:            "ValidCredentials",
			proxyAuthHeader: "Basic " + base64.StdEncoding.EncodeToString([]byte("username:password")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Proxy-Authorization", tt.proxyAuthHeader)
			req = req.WithContext(log.WithContext(req.Context(), log.New(os.Stdout)))

			rr := httptest.NewRecorder()

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			authProvider := &auth.SimpleAuthProvider{User: "username", Pass: "password"}

			middleware := RequestLoggerMiddleware(BasicAuthMiddleware(handler, authProvider))
			middleware.ServeHTTP(rr, req)

			if rr.Code == http.StatusProxyAuthRequired {
				if !strings.Contains(rr.Body.String(), tt.expectErrorMessage) {
					t.Errorf("got %q, want %q", rr.Body.String(), tt.expectErrorMessage)
				}
			} else if rr.Code != http.StatusOK {
				t.Errorf("got status %v, want %v", rr.Code, http.StatusOK)
			}
		})
	}
}
