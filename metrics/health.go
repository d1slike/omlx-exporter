package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/txix-open/isp-kit/metrics"
)

var (
	scrapeDuration = metrics.GetOrRegister(metrics.DefaultRegistry,
		prometheus.NewHistogram(prometheus.HistogramOpts{
			Name: "omlx_scrape_duration_seconds",
			Help: "Duration of each scrape from oMLX API in seconds",
		}))
	scrapeFailures = metrics.GetOrRegister(metrics.DefaultRegistry,
		prometheus.NewCounter(prometheus.CounterOpts{
			Name: "omlx_scrape_failures_total",
			Help: "Total number of failed scrapes from oMLX API",
		}))
)

func RecordScrapeDuration(duration float64) {
	scrapeDuration.Observe(duration)
}

func RecordFailure() {
	scrapeFailures.Inc()
}
