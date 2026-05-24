## ADDED Requirements

### Requirement: Exporter scrapes oMLX server stats
The system SHALL scrape the `/api/stats` endpoint of the oMLX server and expose the following metrics as Prometheus counters and gauges:
- `omlx_server_stats_total_prompt_tokens` (counter)
- `omlx_server_stats_total_completion_tokens` (counter)
- `omlx_server_stats_total_cached_tokens` (counter)
- `omlx_server_stats_total_requests` (counter)
- `omlx_server_stats_cache_efficiency` (gauge, percentage)
- `omlx_server_stats_avg_prefill_tps` (gauge, tokens per second)
- `omlx_server_stats_avg_generation_tps` (gauge, tokens per second)
- `omlx_server_stats_uptime_seconds` (gauge)

#### Scenario: Scrape stats endpoint successfully
- **WHEN** the oMLX server is running and the stats endpoint returns a valid JSON response
- **THEN** the exporter registers all stats metrics with the prometheus registry and exposes them on `/metrics`

#### Scenario: Scrape stats endpoint returns zero values
- **WHEN** the oMLX server has not processed any requests yet
- **THEN** all counter metrics are 0 and gauge metrics reflect the current uptime

#### Scenario: Scrape stats endpoint fails
- **WHEN** the oMLX server is unreachable or returns an error
- **THEN** the exporter logs the error and retains the last known metric values without updating

### Requirement: Exporter scrapes oMLX model status
The system SHALL scrape the `/api/models` endpoint and expose model status metrics:
- `omlx_model_loaded` (gauge, value 1=loaded, 0=not loaded, label `model`)
- `omlx_model_count` (gauge, total number of discovered models)
- `omlx_model_count_loaded` (gauge, number of loaded models)
- `omlx_model_count_unloaded` (gauge, number of unloaded models)

#### Scenario: Scrape models endpoint successfully
- **WHEN** the oMLX server returns a list of models
- **THEN** the exporter sets `omlx_model_loaded` with a `model` label for each model and updates count gauges

#### Scenario: Models endpoint returns empty list
- **WHEN** no models are discovered by the oMLX server
- **THEN** `omlx_model_count` is 0 and no `omlx_model_loaded` time series exist

### Requirement: Exporter scrapes oMLX device info
The system SHALL scrape the `/api/device-info` endpoint and expose hardware metrics:
- `omlx_device_info_memory_gb` (gauge)
- `omlx_device_info_gpu_cores` (gauge)
- `omlx_device_info_chip_name` (gauge, label `chip`)

#### Scenario: Scrape device info successfully
- **WHEN** the oMLX server returns device information
- **THEN** the exporter exposes memory, GPU core count, and chip name as gauge metrics

### Requirement: Exporter scrapes oMLX cache state
The system SHALL scrape cache-related metrics from the oMLX server and expose:
- `omlx_cache_ssd_total_bytes` (gauge)
- `omlx_cache_ssd_used_bytes` (gauge)
- `omlx_cache_hot_max_bytes` (gauge)

#### Scenario: Scrape cache state successfully
- **WHEN** the oMLX server returns cache information
- **THEN** the exporter exposes SSD and hot cache metrics as gauge values

### Requirement: Exporter exposes metrics endpoint
The system SHALL expose a `/metrics` HTTP endpoint that returns Prometheus-formatted metrics from the default registry.

#### Scenario: Metrics endpoint is accessible
- **WHEN** the exporter is running
- **THEN** an HTTP GET to `/metrics` returns `text/plain; version=0.0.4` content type with Prometheus exposition format data

#### Scenario: Metrics endpoint serves registered metrics
- **WHEN** the exporter has scraped the oMLX server at least once
- **THEN** the `/metrics` endpoint includes all scraped metrics with their current values

### Requirement: Exporter is configurable
The system SHALL accept the following configuration via `config.yaml`:
- `omlx.target_url`: URL of the oMLX server (required)
- `omlx.api_key`: API key for authentication (optional)
- `omlx.scrape_interval`: Scrape interval in seconds, default 15 (optional)

#### Scenario: Configuration with all options
- **WHEN** config.yaml contains `omlx.target_url`, `omlx.api_key`, and `omlx.scrape_interval`
- **THEN** the exporter connects to the specified URL, authenticates with the API key, and scrapes at the configured interval

#### Scenario: Configuration with minimal options
- **WHEN** config.yaml contains only `omlx.target_url`
- **THEN** the exporter connects without authentication and uses the default scrape interval of 15 seconds

### Requirement: Exporter lifecycle management
The system SHALL integrate with the `isp-kit/app` lifecycle:
- Scraper goroutine starts as an `app.Runner`
- Scraper goroutine stops gracefully as an `app.Closer`
- Server startup/shutdown is logged via the application logger

#### Scenario: Exporter starts with application
- **WHEN** the application runs
- **THEN** the scraper begins fetching metrics from the oMLX server

#### Scenario: Exporter stops with application
- **WHEN** the application receives a shutdown signal
- **THEN** the scraper goroutine stops cleanly and the server shuts down gracefully
