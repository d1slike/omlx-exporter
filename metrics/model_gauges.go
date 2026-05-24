package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/txix-open/isp-kit/metrics"
)

var (
	modelWaitingRequests    *prometheus.GaugeVec
	modelPrefillingRequests *prometheus.GaugeVec
	modelGeneratingRequests *prometheus.GaugeVec
	modelActiveRequests     *prometheus.GaugeVec
)

func init() {
	modelWaitingRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_waiting_requests",
		Help: "Number of requests waiting in queue for the model",
	}, []string{"model"})
	modelPrefillingRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_prefilling_requests",
		Help: "Number of requests currently in prefill phase for the model",
	}, []string{"model"})
	modelGeneratingRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_generating_requests",
		Help: "Number of requests currently generating tokens for the model",
	}, []string{"model"})
	modelActiveRequests = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_active_requests",
		Help: "Total active requests (prefilling + generating) for the model",
	}, []string{"model"})

	metrics.GetOrRegister(metrics.DefaultRegistry, modelWaitingRequests)
	metrics.GetOrRegister(metrics.DefaultRegistry, modelPrefillingRequests)
	metrics.GetOrRegister(metrics.DefaultRegistry, modelGeneratingRequests)
	metrics.GetOrRegister(metrics.DefaultRegistry, modelActiveRequests)
}

func ResetModelGauges() {
	modelWaitingRequests.Reset()
	modelPrefillingRequests.Reset()
	modelGeneratingRequests.Reset()
	modelActiveRequests.Reset()
}

func SetModelWaitingRequests(model string, value float64) {
	modelWaitingRequests.WithLabelValues(model).Set(value)
}

func SetModelPrefillingRequests(model string, value float64) {
	modelPrefillingRequests.WithLabelValues(model).Set(value)
}

func SetModelGeneratingRequests(model string, value float64) {
	modelGeneratingRequests.WithLabelValues(model).Set(value)
}

func SetModelActiveRequests(model string, value float64) {
	modelActiveRequests.WithLabelValues(model).Set(value)
}
