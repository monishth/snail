package middleware

import (
	"context"
	"net/http"

	"github.com/charmbracelet/log"
)

var requestLoggerKey contextKey = contextKey("requestLogger")

func RequestLoggerMiddleware(next http.Handler) http.Handler {
	return RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		requestID := GetRequestID(ctx)
		ctx = context.WithValue(ctx, requestLoggerKey, log.With("ID", requestID))
		next.ServeHTTP(w, req.WithContext(ctx))
	}))
}

func GetLogger(ctx context.Context) *log.Logger {
	if requestLogger, ok := ctx.Value(requestLoggerKey).(*log.Logger); ok {
		return requestLogger
	}
	return nil
}
