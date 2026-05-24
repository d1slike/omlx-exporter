## Why

The monitoring stack runs Grafana with a pre-configured Prometheus datasource, but no dashboard exists. The `/metrics` endpoint exposes both OMLX-specific metrics and Go runtime/process metrics, making it difficult to focus on OMLX performance in Grafana without noise from unrelated metrics.

## What Changes

- Create a Grafana dashboard JSON file that displays only OMLX metrics (excluding `go_*`, `process_*`, `promhttp_*`, and other runtime metrics)
- Provision the dashboard in Grafana using Grafana's dashboard provisioning via a JSON file mounted to `/etc/grafana/provisioning/dashboards/`
- Add a `dashboards` provisioning datasource configuration alongside the existing `datasources` provisioning

## Capabilities

### New Capabilities
- `omlx-dashboard`: Grafana dashboard that visualizes only OMLX metrics (memory pressure, model request stats, server stats, scrape health) with Go/runtime/process metrics excluded, including a model filter variable to scope per-model metrics

### Modified Capabilities
- `monitoring-stack`: Adds dashboard provisioning requirement to the existing Grafana provisioning capability

## Impact

- `grafana/provisioning/dashboards/` - New directory for dashboard JSON and provisioning config
- `grafana/provisioning/datasources/datasources.yml` - Existing datasource config unchanged
- `docker-compose.yaml` - May need volume mount for dashboards provisioning directory
- Grafana UI - New "OMLX Exporter" dashboard with model filter variable for selecting specific models
