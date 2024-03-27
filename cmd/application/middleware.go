package application

import (
	"net/http"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
)

func (a *App) withMiddleware() {
	a.handler = a.logging(a.handler)
	a.handler = a.panicRecovery(a.handler)
	a.handler = a.contentType(a.handler)
}

func (a *App) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		recorder := &StatusRecorder{
			ResponseWriter: w,
			Status:         200,
		}
		next.ServeHTTP(recorder, req)
		a.logger.Info("request",
			zap.String("method", req.Method),
			zap.String("url", req.RequestURI),
			zap.Int("response code", recorder.Status),
			zap.Duration("duration", time.Since(start)))
	})
}

func (a *App) panicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				a.logger.Error("panic recover", zap.String("stack", string(debug.Stack())))
			}
		}()
		next.ServeHTTP(w, req)
	})
}

func (a *App) contentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}
