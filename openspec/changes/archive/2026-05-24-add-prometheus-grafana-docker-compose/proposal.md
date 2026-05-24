## Why

The omlx-exporter exposes Prometheus metrics but there is no service to scrape, store, or visualize them. Adding Prometheus and Grafana to docker-compose provides observability out of the box — metric collection via Prometheus and dashboards via Grafana — so operators can monitor the oMLX server's health, performance, and model status without external tooling.

## What Changes

- Add a `prometheus` service to `docker-compose.yaml` using the latest official image
- Add a `grafana` service to `docker-compose.yaml` using the latest official image
- Create a `prometheus/prometheus.yml` configuration file that scrapes `omlx-exporter` at `:8001/metrics` with the job name `omlx-exporter`
- Create a `grafana/provisioning/datasources/datasources.yml` configuration file that auto-configures Prometheus as a datasource
- Add a `prometheus-data` named volume for persistent storage of Prometheus data
- Expose Prometheus on host port `9090` and Grafana on host port `3000`

## Capabilities

### New Capabilities
- `monitoring-stack`: Infrastructure capability for deploying Prometheus and Grafana alongside omlx-exporter for metric scraping and visualization

### Modified Capabilities
<!-- No existing spec requirements are changing — this is purely deployment configuration -->

## Impact

- `docker-compose.yaml` — adds two new services and a volume
- `prometheus/prometheus.yml` — new file with scrape configuration for omlx-exporter
- `grafana/provisioning/datasources/datasources.yml` — new file with Prometheus datasource configuration
- No changes to Go code, API, or existing exporter behavior
- New external dependencies: `prom/prometheus` and `grafana/grafana` Docker images
