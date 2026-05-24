## ADDED Requirements

### Requirement: OMLX dashboard displays only OMLX metrics
The system SHALL provision a Grafana dashboard named "OMLX Exporter" that displays only metrics with the `omlx_` prefix, excluding `go_*`, `process_*`, `promhttp_*`, and other runtime metrics. Each panel SHALL use Prometheus queries with metric name selectors (e.g., `{__name__=~"omlx_.*"}`) to filter the metric namespace.

#### Scenario: Dashboard loads with OMLX-only panels
- **WHEN** a user navigates to the Grafana UI and opens the "OMLX Exporter" dashboard
- **THEN** all panels display data from `omlx_*` metrics only, with no Go runtime, process, or Prometheus HTTP handler metrics visible

#### Scenario: Memory pressure panel shows correct data
- **WHEN** the dashboard loads and the OMLX server is active
- **THEN** the memory pressure panel displays `omlx_active_models_memory_pressure_level`, `omlx_active_models_memory_current_bytes`, `omlx_active_models_memory_soft_bytes`, and `omlx_active_models_memory_hard_bytes` gauges

#### Scenario: Model request panels show per-model data
- **WHEN** the dashboard loads and models are processing requests
- **THEN** panels display per-model metrics (`omlx_model_active_requests`, `omlx_model_prefilling_requests`, `omlx_model_generating_requests`, `omlx_model_waiting_requests`) with the `model` label

#### Scenario: Server stats panel shows cumulative metrics
- **WHEN** the dashboard loads
- **THEN** the server stats panel displays `omlx_stats_total_prompt_tokens`, `omlx_stats_total_completion_tokens`, `omlx_stats_total_cached_tokens`, `omlx_stats_total_requests`, `omlx_stats_cache_efficiency`, `omlx_stats_avg_prefill_tps`, `omlx_stats_avg_generation_tps`, and `omlx_stats_uptime_seconds`

#### Scenario: Scrape health panel shows exporter status
- **WHEN** the dashboard loads
- **THEN** the scrape health panel displays `omlx_scrape_duration_seconds` histogram and `omlx_scrape_failures_total` counter

### Requirement: Dashboard is auto-provisioned in Grafana
The system SHALL provision the OMLX Exporter dashboard automatically via Grafana's dashboard provisioning mechanism. The dashboard JSON file SHALL be placed in a directory mounted to `/etc/grafana/provisioning/dashboards/` inside the Grafana container, with a provisioning configuration file that points to the dashboard JSON.

#### Scenario: Dashboard appears after Grafana restart
- **WHEN** the Grafana container is restarted with the dashboards provisioning directory mounted
- **THEN** the "OMLX Exporter" dashboard is automatically loaded and available in the Grafana UI

#### Scenario: Dashboard provisioning configuration is valid
- **WHEN** the dashboards provisioning YAML file is parsed by Grafana
- **THEN** it references the correct JSON file path, uses `file` as the `type`, and sets `overwrite: true`

### Requirement: Dashboard includes model filter variable
The system SHALL provision a Grafana template variable named `model` that allows users to filter per-model metrics. The variable SHALL use the Prometheus query `label_values(omlx_model_active_requests, model)` to populate available model options. Per-model panels SHALL include `{model="$model"}` in their Prometheus queries. Non-model metrics (memory pressure, server stats, scrape health) SHALL not be affected by the filter. The default value SHALL be `ALL` to show all models when no selection is made.

#### Scenario: Model filter dropdown is available
- **WHEN** a user opens the "OMLX Exporter" dashboard
- **THEN** a `model` filter dropdown is visible at the top of the dashboard with options listing all discovered models and an `ALL` option

#### Scenario: Selecting a model filters per-model panels
- **WHEN** a user selects a specific model from the filter dropdown
- **THEN** all panels displaying per-model metrics (`omlx_model_active_requests`, `omlx_model_prefilling_requests`, `omlx_model_generating_requests`, `omlx_model_waiting_requests`) show data only for the selected model

#### Scenario: Selecting ALL shows all models
- **WHEN** a user selects `ALL` from the model filter
- **THEN** per-model panels display aggregated data across all models (Prometheus queries use `{model="$model"}` which matches all models when `$model` is `ALL`)

#### Scenario: Non-model panels are unaffected by filter
- **WHEN** a user changes the model filter selection
- **THEN** panels for memory pressure, server stats, and scrape health continue to display their data without being filtered by model

### Requirement: Dashboard is organized into logical sections
The system SHALL organize the dashboard panels into four sections: Memory Pressure, Model Requests, Server Stats, and Scrape Health. Each section SHALL be a row in the Grafana dashboard layout.

#### Scenario: Dashboard rows are correctly labeled
- **WHEN** a user views the "OMLX Exporter" dashboard
- **THEN** four rows are visible: "Memory Pressure", "Model Requests", "Server Stats", and "Scrape Health"

#### Scenario: Each row contains the correct panels
- **WHEN** each row is expanded in the dashboard
- **THEN** the Memory Pressure row contains 3 panels (pressure level, current vs soft/hard limits, total active/waiting requests), the Model Requests row contains 4 panels (active, prefilling, generating, waiting per model), the Server Stats row contains 4 panels (tokens throughput, cache efficiency, TPS, uptime), and the Scrape Health row contains 2 panels (scrape duration, failures)
