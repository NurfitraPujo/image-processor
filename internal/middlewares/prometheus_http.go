package middlewares

import (
	"net/http"
	"strconv"

	"github.com/NurfitraPujo/image-processor/internal/metrics"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusHttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		timer := prometheus.NewTimer(metrics.HttpDuration.WithLabelValues(path))
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)

		metrics.ResponseStatus.WithLabelValues(strconv.Itoa(ww.Status())).Inc()
		metrics.TotalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})

}
