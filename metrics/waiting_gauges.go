package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/txix-open/isp-kit/metrics"
)

var (
	modelWaitingQueuePos *prometheus.GaugeVec
	modelWaitingElapsed  *prometheus.GaugeVec
)

func init() {
	modelWaitingQueuePos = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_waiting_queue_position",
		Help: "Queue position for waiting requests (sum across all waiting requests per model)",
	}, []string{"model"})
	modelWaitingElapsed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_waiting_elapsed_seconds",
		Help: "Total elapsed time across all waiting requests per model",
	}, []string{"model"})

	metrics.GetOrRegister(metrics.DefaultRegistry, modelWaitingQueuePos)
	metrics.GetOrRegister(metrics.DefaultRegistry, modelWaitingElapsed)
}

func ResetWaitingGauges() {
	modelWaitingQueuePos.Reset()
	modelWaitingElapsed.Reset()
}

func SetModelWaitingQueuePos(model string, value float64) {
	modelWaitingQueuePos.WithLabelValues(model).Set(value)
}

func SetModelWaitingElapsed(model string, value float64) {
	modelWaitingElapsed.WithLabelValues(model).Set(value)
}
