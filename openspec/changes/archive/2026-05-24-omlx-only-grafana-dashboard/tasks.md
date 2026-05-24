## 1. Create dashboard provisioning structure

- [x] 1.1 Create `grafana/provisioning/dashboards/` directory
- [x] 1.2 Create `grafana/provisioning/dashboards/dashboards.yml` provisioning config

## 2. Create OMLX dashboard JSON

- [x] 2.1 Create dashboard JSON with "OMLX Exporter" title, Prometheus datasource, and model template variable using `label_values(omlx_model_active_requests, model)`
- [x] 2.2 Add Memory Pressure row: (a) pressure level gauge with thresholds 0/1/2, (b) current vs soft/hard memory bytes bar gauge, (c) total active/waiting requests time series
- [x] 2.3 Add Model Requests row: (a) active requests, (b) prefilling requests, (c) generating requests, (d) waiting requests — each using `{model="$model"}` filter as time series
- [x] 2.4 Add Server Stats row: (a) prompt/completion/cached tokens throughput as time series, (b) cache efficiency percentage gauge, (c) prefill/generation TPS as time series, (d) uptime as singlestat
- [x] 2.5 Add Scrape Health row: (a) scrape duration histogram/SLO panel, (b) scrape failures counter

## 3. Update docker-compose

- [x] 3.1 Add volume mount for `grafana/provisioning/dashboards/` directory to Grafana service
