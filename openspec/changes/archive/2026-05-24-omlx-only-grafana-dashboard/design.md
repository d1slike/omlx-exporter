## Context

The OMLX exporter exposes ~50 Prometheus metrics including OMLX-specific metrics and Go runtime/process metrics. Grafana is already configured with a Prometheus datasource via provisioning, but no dashboard exists. Users need a focused view of OMLX performance without noise from Go runtime metrics.

## Goals / Non-Goals

**Goals:**
- Create a Grafana dashboard JSON file that filters to only OMLX metrics using Prometheus label matchers
- Provision the dashboard automatically via Grafana's dashboard provisioning mechanism
- Organize dashboard panels into logical sections: memory pressure, model requests, server stats, scrape health
- Add a model filter dropdown to scope per-model metrics to a specific model instance

**Non-Goals:**
- Modifying the exporter or metrics output
- Adding alerts or alerting rules
- Supporting multiple OMLX instances in a single dashboard
- Customizing Grafana theme or layout beyond the dashboard JSON

## Decisions

1. **Dashboard provisioning via JSON file** - Use Grafana's built-in JSON dashboard provisioning (mounted to `/etc/grafana/provisioning/dashboards/`) rather than provisioning via Grafana API or external tooling. This is the simplest approach that requires no additional dependencies.

2. **Metric filtering via Prometheus selectors** - Each panel uses Prometheus queries with `{__name__=~"omlx_.*"}` or specific metric name matchers to exclude Go/runtime/process metrics. This keeps filtering declarative in the dashboard rather than requiring exporter-side metric naming changes.

3. **Single dashboard file** - One comprehensive dashboard JSON file rather than multiple smaller dashboards. OMLX metrics are cohesive and a single view is more useful for operators.

4. **Directory structure** - Place dashboard JSON at `grafana/provisioning/dashboards/omlx-exporter.json` alongside a `dashboards.yml` provisioning config, mirroring the existing `datasources/` structure.

5. **Model filter variable** - Add a Grafana template variable `$model` with query `label_values(omlx_model_active_requests, model)` to populate available models. Per-model panels SHALL use `{model="$model"}` in their Prometheus queries. Non-model metrics (memory pressure, server stats, scrape health) SHALL remain unaffected by the filter.

## Risks / Trade-offs

[Risk] Dashboard JSON is large and hard to maintain → Use a programmatic approach (Go script or Python) to generate the JSON from metric definitions in the future if needed.
[Risk] Grafana provisioning requires directory mount in docker-compose → Add the mount to `docker-compose.yaml` Grafana service volumes.
[Trade-off] Hardcoded metric names in panels → Panels reference specific metric names; adding new OMLX metrics requires dashboard updates. This is acceptable since new metrics are infrequent.

## Migration Plan

1. Create dashboard JSON and provisioning config
2. Add volume mount in docker-compose.yaml for dashboards directory
3. Restart Grafana container - dashboard appears automatically
4. No rollback needed - removing the volume mount and files reverts the change

## Open Questions

None identified.
