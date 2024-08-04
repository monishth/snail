package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestLoggerMiddleware(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Could not create HTTP request: %v", err)
	}

	req = req.WithContext(context.WithValue(req.Context(), requestIDKey, "123"))

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := GetLogger(r.Context())
		assert.NotNil(t, logger)
	})

	RequestLoggerMiddleware(handler).ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

}
