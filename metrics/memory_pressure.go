package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/txix-open/isp-kit/metrics"
)

var (
	memoryPressureLevel = metrics.GetOrRegister(metrics.DefaultRegistry,
		prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "omlx_active_models_memory_pressure_level",
			Help: "Memory pressure level: 0=ok, 1=warning, 2=critical",
		}))
	memoryCurrentBytes = metrics.GetOrRegister(metrics.DefaultRegistry,
		prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "omlx_active_models_memory_current_bytes",
			Help: "Current memory usage in bytes",
		}))
	memorySoftBytes = metrics.GetOrRegister(metrics.DefaultRegistry,
		prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "omlx_active_models_memory_soft_bytes",
			Help: "Soft memory limit in bytes (warning threshold)",
		}))
	memoryHardBytes = metrics.GetOrRegister(metrics.DefaultRegistry,
		prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "omlx_active_models_memory_hard_bytes",
			Help: "Hard memory limit in bytes (eviction threshold)",
		}))
)

func UpdateMemoryPressure(level string, current, soft, hard float64) {
	var levelFloat float64
	switch level {
	case "ok":
		levelFloat = 0
	case "warning":
		levelFloat = 1
	case "critical":
		levelFloat = 2
	default:
		levelFloat = -1
	}
	memoryPressureLevel.Set(levelFloat)
	memoryCurrentBytes.Set(current)
	memorySoftBytes.Set(soft)
	memoryHardBytes.Set(hard)
}
