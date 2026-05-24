## Why

The oMLX server exposes rich operational metrics through its REST API (`/api/stats`, `/api/models`, `/api/device-info`, etc.) but these metrics are only available in JSON format for the admin dashboard. There is no Prometheus-compatible endpoint for monitoring via external tooling (Grafana, alerting systems, etc.). A dedicated Prometheus metrics exporter is needed to scrape these metrics and expose them in Prometheus format.

## What Changes

- Create a Go-based Prometheus metrics exporter that scrapes the oMLX server API
- Register all scraped metrics using `metrics.GetOrRegister(metrics.DefaultRegistry, ...)` from `isp-kit/metrics`
- Expose a `/metrics` endpoint serving Prometheus-formatted metrics
- Export metrics for: server stats (tokens, requests, cache efficiency, throughput), model status (loaded/unloaded, model info), device/hardware info, cache state (SSD/hot cache)
- Add configuration for oMLX server target URL and scrape interval

## Capabilities

### New Capabilities
- `prometheus-exporter`: Scrape oMLX server API endpoints and expose metrics in Prometheus format

## Impact

- New Go module: `github.com/d1slike/omlx-exporter` (already initialized)
- New packages: `exporter/` (scraping logic), `config/` (configuration)
- Dependency: `github.com/prometheus/client_golang` (already transitive via isp-kit)
- Configuration: `config.yaml` schema extended with exporter settings
- `main.go` extended with exporter initialization and metrics endpoint
