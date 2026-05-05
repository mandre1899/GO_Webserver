package middleware

import (
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	FileserverHits	atomic.Int32
	JWTSecret       string
}

func (a *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.FileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

