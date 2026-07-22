package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type contextKey string

const requestIDKey contextKey = "requestID"

// RequestID generates a unique ID for each request and injects it into both
// the request context and the response headers (X-Request-ID).
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)
		
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger injects a request-scoped zerolog.Logger into the context. It adds
// fields like method, path, remote IP, and request ID. It also logs the
// request completion with status code and duration.
func Logger(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			reqID, _ := r.Context().Value(requestIDKey).(string)

			// Create a request-scoped logger
			reqLog := log.With().
				Str("req_id", reqID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("remote_ip", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Logger()

			// Inject logger into context
			ctx := reqLog.WithContext(r.Context())
			
			// Use response recorder to capture status code
			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			next.ServeHTTP(ww, r.WithContext(ctx))

			// Log completion
			reqLog.Info().
				Int("status", ww.status).
				Dur("duration", time.Since(start)).
				Msg("request completed")
		})
	}
}

// responseWriter is a wrapper around http.ResponseWriter that tracks the status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
