package middleware

import (
	"4_3/internal/myLogger"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(next http.Handler, logger *myLogger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		next.ServeHTTP(lw, r)

		duration := time.Since(start)
		logger.LogHTTP("request handled", r.Method, r.URL.Path, lw.statusCode, duration)
	})
}
