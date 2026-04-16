package tracing

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

type correlationContextKey string

const key correlationContextKey = "correlation-id"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := r.Header.Get(CorrelationIDHeader)
		if correlationID == "" {
			correlationID = newCorrelationID()
		}
		w.Header().Set(CorrelationIDHeader, correlationID)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), key, correlationID)))
	})
}

func FromContext(ctx context.Context) string {
	value, _ := ctx.Value(key).(string)
	return value
}

func newCorrelationID() string {
	buffer := make([]byte, 16)
	if _, err := rand.Read(buffer); err != nil {
		return "generated-correlation-id"
	}
	return hex.EncodeToString(buffer)
}
