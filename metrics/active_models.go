package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/txix-open/isp-kit/metrics"
)

var (
	totalActiveRequests = metrics.GetOrRegister(metrics.DefaultRegistry,
		prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "omlx_active_models_total_active_requests",
			Help: "Total number of active requests across all models",
		}))
	totalWaitingRequests = metrics.GetOrRegister(metrics.DefaultRegistry,
		prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "omlx_active_models_total_waiting_requests",
			Help: "Total number of requests waiting in queue across all models",
		}))
)

func UpdateActiveRequests(active, waiting float64) {
	totalActiveRequests.Set(active)
	totalWaitingRequests.Set(waiting)
}
