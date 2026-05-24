package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/txix-open/isp-kit/metrics"
)

var (
	modelGeneratingTPS     *prometheus.GaugeVec
	modelGeneratingTokens  *prometheus.GaugeVec
	modelGeneratingElapsed *prometheus.GaugeVec
	modelPrefillingTPS     *prometheus.GaugeVec
	modelPrefillingTokens  *prometheus.GaugeVec
	modelPrefillingElapsed *prometheus.GaugeVec
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
	modelPrefillingTPS = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_prefilling_tokens_per_second",
		Help: "Tokens per second for prefilling requests (sum across all prefilling requests per model)",
	}, []string{"model"})
	modelPrefillingTokens = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_prefilling_prompt_tokens",
		Help: "Total prompt tokens across all prefilling requests per model",
	}, []string{"model"})
	modelPrefillingElapsed = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "omlx_model_prefilling_elapsed_seconds",
		Help: "Total elapsed time across all prefilling requests per model",
	}, []string{"model"})

	metrics.GetOrRegister(metrics.DefaultRegistry, modelGeneratingTPS)
	metrics.GetOrRegister(metrics.DefaultRegistry, modelGeneratingTokens)
	metrics.GetOrRegister(metrics.DefaultRegistry, modelGeneratingElapsed)
	metrics.GetOrRegister(metrics.DefaultRegistry, modelPrefillingTPS)
	metrics.GetOrRegister(metrics.DefaultRegistry, modelPrefillingTokens)
	metrics.GetOrRegister(metrics.DefaultRegistry, modelPrefillingElapsed)
}

func ResetRequestGauges() {
	modelGeneratingTPS.Reset()
	modelGeneratingTokens.Reset()
	modelGeneratingElapsed.Reset()
	modelPrefillingTPS.Reset()
	modelPrefillingTokens.Reset()
	modelPrefillingElapsed.Reset()
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

func SetModelPrefillingTPS(model string, value float64) {
	modelPrefillingTPS.WithLabelValues(model).Set(value)
}

func SetModelPrefillingTokens(model string, value float64) {
	modelPrefillingTokens.WithLabelValues(model).Set(value)
}

func SetModelPrefillingElapsed(model string, value float64) {
	modelPrefillingElapsed.WithLabelValues(model).Set(value)
}
