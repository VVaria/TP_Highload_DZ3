package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Middleware struct {
	histogramQ *prometheus.HistogramVec
}

func NewMiddleware(hist *prometheus.HistogramVec) *Middleware {
	return &Middleware{
		histogramQ: hist,
	}
}

func (mw *Middleware) AppJSONMiddleware(next http.Handler) http.Handler {
	fmt.Println("received connection")
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		start := time.Now()
		next.ServeHTTP(w, req)

		var status string
		mw.histogramQ.WithLabelValues(status).Observe(float64(start.Second()))
	})
}
