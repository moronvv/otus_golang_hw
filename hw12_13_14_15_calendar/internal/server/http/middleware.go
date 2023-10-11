package internalhttp

import (
	"log/slog"
	"net/http"
	"time"
)

type httpStatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *httpStatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

type loggerMiddleware struct {
	logger *slog.Logger
}

func (m *loggerMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		recorder := &httpStatusRecorder{
			ResponseWriter: w,
			Status:         http.StatusOK,
		}
		next.ServeHTTP(recorder, r)

		m.logger.Info("",
			slog.Group("http",
				"client_ip", r.RemoteAddr,
				"user_agent", r.Header.Get("User-Agent"),
				"protocol", r.Proto,
				"status", recorder.Status,
				"method", r.Method,
				"path", r.URL.Path,
			),
			"latency", time.Since(start),
		)
	})
}

func newLoggerMiddleware(logger *slog.Logger) *loggerMiddleware {
	return &loggerMiddleware{
		logger: logger,
	}
}
