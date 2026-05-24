# AGENTS.md - Project Guidelines for OMLX Exporter

## Project Overview

Go 1.26 Prometheus exporter for OMLX monitoring. Exposes metrics about OMLX models, memory pressure, health, and request stats via `/metrics` endpoint.

## Architecture

- **`main.go`** - Application entry point, wires together config, client, scraper, and HTTP server
- **`config.go`** - Root config struct
- **`config/`** - Config loading logic
- **`client/`** - OMLX API client and response models
- **`exporter/`** - Scraper that periodically fetches OMLX metrics
- **`metrics/`** - Prometheus metric definitions (gauges, health, stats)
- **`grafana/`** - Grafana dashboard configurations
- **`prometheus/`** - Prometheus scrape configurations

## Dependencies

- `github.com/txix-open/isp-kit` - Application framework (config, logging, lifecycle, metrics)
- `github.com/prometheus/client_golang` - Prometheus client
- `github.com/pkg/errors` - Error wrapping

## Code Conventions

- Module: `github.com/d1slike/omlx-exporter`
- No comments in code unless explicitly requested
- Use `isp-kit` patterns for config, logging, and application lifecycle
- Use `pkg/errors` for error wrapping
- Metrics are defined in `metrics/` package and exported via `prometheus/client_golang`
- HTTP server routes are registered in `main.go` via `infra.NewServer()`

## Building and Running

```bash
go run .
```

## Docker

```bash
docker compose up -d
```

## Adding New Metrics

1. Define the metric in the `metrics/` package using `prometheus.NewGaugeVec` or similar
2. Register it with `metrics.DefaultRegistry`
3. Update the scraper in `exporter/` to collect and set metric values
4. If needed, add corresponding API types in `client/models.go`

## Config

Configuration is loaded from `config.yaml` via `isp-kit`'s YAML config source. Any variable from `config.yaml` can be overwritten by an environment variable with the same name. Extend the `Config` struct in `config.go` and update `config.yaml` for new options.
