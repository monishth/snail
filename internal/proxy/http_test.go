package proxy

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/log"
	"io"
	"net/http"

	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleHTTPRequest(t *testing.T) {
	tests := []struct {
		name         string
		request      *http.Request
		logOutput    string
		expectedRes  string
		expectedCode int
	}{
		{
			name:         "UnsupportedProtocol",
			request:      httptest.NewRequest("GET", "https://example.com", nil),
			logOutput:    "ERRO unsupported protocol scheme https",
			expectedRes:  "unsupported protocol scheme https",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:         "InvalidHost",
			request:      httptest.NewRequest("GET", "http://:80", nil),
			logOutput:    "Error message=\"Failed to connect to upstream server\" error=\"dial tcp :80: connect: connection refused\"",
			expectedRes:  "Failed to connect to upstream server",
			expectedCode: http.StatusBadGateway,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			var buf bytes.Buffer
			logger := log.New(&buf)
			proxy := &ForwardProxy{}

			proxy.handleHTTPRequest(recorder, tt.request, logger)
			res := recorder.Result()

			if res.StatusCode != tt.expectedCode {
				t.Fatalf("expected status %v; got %v", tt.expectedCode, res.StatusCode)
			}

			defer res.Body.Close()
			body, _ := io.ReadAll(res.Body)

			if strings.TrimSpace(string(body)) != tt.expectedRes {
				fmt.Print(string(body))
				fmt.Print(tt.expectedRes)
				t.Fatalf("expected body '%v'; got '%v'", tt.expectedRes, string(body))
			}

			logList := strings.Split(strings.TrimSpace(buf.String()), "\n")
			logline := logList[len(logList)-1]

			if strings.Contains(logline, tt.logOutput) == false {
				t.Fatalf("expected to find '%v' in logline '%v'", tt.logOutput, logline)
			}
		})
	}
}
