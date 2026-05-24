package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/txix-open/isp-kit/metrics"
)

var (
	modelGeneratingTPS     *prometheus.GaugeVec
	modelGeneratingTokens  *prometheus.GaugeVec
	modelGeneratingElapsed *prometheus.GaugeVec
)

func init() {
	modelGeneratingTPS = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_generating_tokens_per_second",
		Help: "Tokens per second for generating requests (sum across all generating requests per model)",
	}, []string{"model"})
	modelGeneratingTokens = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_generating_generated_tokens",
		Help: "Total generated tokens across all generating requests per model",
	}, []string{"model"})
	modelGeneratingElapsed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_generating_elapsed_seconds",
		Help: "Total elapsed time across all generating requests per model",
	}, []string{"model"})

	metrics.GetOrRegister(metrics.DefaultRegistry, modelGeneratingTPS)
	metrics.GetOrRegister(metrics.DefaultRegistry, modelGeneratingTokens)
	metrics.GetOrRegister(metrics.DefaultRegistry, modelGeneratingElapsed)
}

func ResetRequestGauges() {
	modelGeneratingTPS.Reset()
	modelGeneratingTokens.Reset()
	modelGeneratingElapsed.Reset()
}

func SetModelGeneratingTPS(model string, value float64) {
	modelGeneratingTPS.WithLabelValues(model).Set(value)
}

func SetModelGeneratingTokens(model string, value float64) {
	modelGeneratingTokens.WithLabelValues(model).Set(value)
}

func SetModelGeneratingElapsed(model string, value float64) {
	modelGeneratingElapsed.WithLabelValues(model).Set(value)
}
