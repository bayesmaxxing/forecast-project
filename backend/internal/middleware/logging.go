package middleware

import (
	"backend/internal/logger"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"
)

func generateRequestID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		requestID := generateRequestID()

		log := slog.Default().With(
			slog.String("request_id", requestID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		// store logger in context
		ctx := logger.WithContext(r.Context(), log)

		r = r.WithContext(ctx)

		rw := &responseWriter{w, http.StatusOK}

		log.Info("request started")
		// Call the next handler
		next.ServeHTTP(rw, r)

		// Log the request details after it's processed
		log.Info("request completed",
			slog.Int("status", rw.status),
			slog.Duration("duration", time.Since(start)),
		)
	})
}

// responseWriter is a wrapper for http.ResponseWriter that captures the status code
type responseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader captures the status code before passing it to the underlying ResponseWriter
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
