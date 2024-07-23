package middleware

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
)

type contextKey string

var requestIDKey contextKey = contextKey("requestID")

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		requestID := strconv.FormatInt(rand.Int63(), 16)
		ctx := context.WithValue(req.Context(), requestIDKey, requestID)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}
