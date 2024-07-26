package middleware

import (
	"encoding/base64"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/monishth/snail/internal/auth"
)

type Credentials struct {
	user string
	pass string
}

func BasicAuthMiddleware(next http.Handler, authProvider auth.AuthProvider) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		logger := req.Context().Value(requestLoggerKey).(*log.Logger)
		if logger == nil {
			logger = log.New(os.Stderr)
		}

		proxyAuthHeader := req.Header.Get("Proxy-Authorization")
		if proxyAuthHeader == "" {
			logger.Error("Auth Error", "message", "No auth header provided")
			http.Error(w, "Proxy Authorization Required", http.StatusProxyAuthRequired)
			return
		}

		// Split Basic and base64 encoded credentials
		parts := strings.SplitN(proxyAuthHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "basic" {
			logger.Error("Auth Error", "message", "No basic auth credentials provided")
			http.Error(w, "Proxy Authorization Required", http.StatusProxyAuthRequired)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(parts[1])
		// Get user + pass
		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			logger.Error("Auth Error", "message", "No basic auth credentials provided")
			http.Error(w, "Proxy Authorization Required", http.StatusProxyAuthRequired)
			return
		}

		if !authProvider.Verify(pair[0], pair[1]) {
			logger.Error("Auth Error", "message", "credentials incorrect")
			http.Error(w, "Proxy Authorization Required", http.StatusProxyAuthRequired)
			return
		}

		next.ServeHTTP(w, req)
	})
}
