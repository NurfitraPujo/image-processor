package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.Register(TotalRequests)
	prometheus.Register(ResponseStatus)
	prometheus.Register(HttpDuration)
}
