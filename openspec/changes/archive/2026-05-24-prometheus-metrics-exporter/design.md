## Context

The oMLX server (Python/FastAPI) exposes operational metrics through REST API endpoints that return JSON data. The existing Go project (`d1slike/omlx-exporter`) is a standalone process that already uses `isp-kit` infrastructure (app, config, logger, shutdown, server). The `isp-kit/metrics` package provides `metrics.DefaultRegistry` and `metrics.GetOrRegister` for registering Prometheus metrics. The exporter needs to periodically scrape the oMLX server API and convert JSON responses into Prometheus counters, gauges, and histograms.

## Goals / Non-Goals

**Goals:**
- Scrape oMLX server API endpoints at configurable intervals
- Convert API responses to Prometheus metrics using `isp-kit/metrics`
- Expose `/metrics` endpoint via the existing `infra.Server`
- Support authentication (API key) for oMLX server endpoints
- Export server stats, model status, device info, and cache metrics

**Non-Goals:**
- Modifying the oMLX server (it's an external dependency)
- Implementing a full scrape-and-store pipeline (no local storage)
- Alerting logic (only metric exposure)
- Support for non-HTTP scrape targets

## Decisions

### 1. Pull-based scraping via HTTP client
The exporter pulls metrics from the oMLX server API using `github.com/txix-open/isp-kit/http/httpcli`. This matches the existing `isp-kit` ecosystem and allows the exporter to be deployed independently.

**Rationale**: The oMLX server has no native Prometheus exposition. Pull-based approach means the exporter controls its own scrape cadence and can handle authentication. `httpcli` provides a typed, reusable HTTP client with built-in retry and error handling.

### 2. Registry-driven metrics registration
All metrics are registered with `metrics.GetOrRegister(metrics.DefaultRegistry, ...)` from `isp-kit/metrics`. This ensures compatibility with the existing `metrics.DefaultRegistry.MetricsHandler()` already wired in `main.go`.

**Rationale**: No need to add prometheus/client_golang as a direct dependency. The existing infrastructure handles registry and exposition.

### 3. Structured scraping with typed response models
Each oMLX API endpoint gets its own Go struct for JSON deserialization and a dedicated scraper function. This provides type safety and makes it easy to add new endpoints.

**Rationale**: The oMLX API responses are well-structured JSON. Go structs map cleanly to these shapes and provide compile-time validation.

### 4. Metrics naming convention
Metrics follow Prometheus naming: `omlx_<endpoint>_<metric>`. Examples:
- `omlx_server_stats_total_prompt_tokens` (counter)
- `omlx_server_stats_cache_efficiency` (gauge)
- `omlx_model_loaded` (gauge with `model` label)
- `omlx_device_info_memory_gb` (gauge)

### 5. Periodic scraping with worker
A goroutine runs the scraper at a configurable interval (default: 15s) using `github.com/txix-open/isp-kit/worker`. The worker is started as an `app.Runner` and stopped via `app.Closer`.

**Rationale**: Uses the existing `isp-kit/app` lifecycle management and `isp-kit/worker` for reliable periodic task execution with clean start/stop semantics.

## Risks / Trade-offs

[Risk] oMLX server API schema changes between versions → [Mitigation] Scraper returns partial metrics on parse errors rather than failing entirely. Log warnings for unexpected fields.

[Risk] Authentication credentials in config.yaml → [Mitigation] Document security best practices. Consider future support for env var injection via isp-kit config.

[Risk] High scrape frequency causing load on oMLX server → [Mitigation] Configurable interval with sensible default (15s).

## Migration Plan

No migration needed. This is a new standalone exporter deployed alongside the oMLX server.

1. Add exporter code to existing `d1slike/omlx-exporter` module
2. Extend `config.yaml` with exporter settings (`target_url`, `api_key`, `scrape_interval`)
3. Wire exporter into `main.go` lifecycle
4. Build and deploy as separate binary/container

## Metrics Specification

### Source: `GET /api/stats` (scope=session)

| Prometheus Metric | Type | oMLX JSON Field | Description |
|---|---|---|---|
| `omlx_stats_total_prompt_tokens` | counter | `total_prompt_tokens` | Total prompt tokens processed |
| `omlx_stats_total_completion_tokens` | counter | `total_completion_tokens` | Total completion tokens generated |
| `omlx_stats_total_cached_tokens` | counter | `total_cached_tokens` | Total tokens served from cache |
| `omlx_stats_total_requests` | counter | `total_requests` | Total requests processed |
| `omlx_stats_cache_efficiency` | gauge | `cache_efficiency` | Cache hit ratio percentage |
| `omlx_stats_avg_prefill_tps` | gauge | `avg_prefill_tps` | Average prefill tokens/sec |
| `omlx_stats_avg_generation_tps` | gauge | `avg_generation_tps` | Average generation tokens/sec |
| `omlx_stats_uptime_seconds` | gauge | `uptime_seconds` | Server uptime in seconds |

### Source: `GET /api/stats` — `active_models`

| Prometheus Metric | Type | oMLX JSON Field | Labels | Description |
|---|---|---|---|---|
| `omlx_active_models_model_memory_used_bytes` | gauge | `active_models.model_memory_used` | `model` | Per-model memory used |
| `omlx_active_models_total_active_requests` | gauge | `active_models.total_active_requests` | — | Total active requests across all models |
| `omlx_active_models_total_waiting_requests` | gauge | `active_models.total_waiting_requests` | — | Total waiting requests |
| `omlx_active_models_memory_pressure_level` | gauge | `active_models.memory_pressure.pressure_level` | — | Memory pressure: 0=ok, 1=warning, 2=critical |
| `omlx_active_models_memory_current_bytes` | gauge | `active_models.memory_pressure.current_bytes` | — | Current memory enforcement bytes |
| `omlx_active_models_memory_soft_bytes` | gauge | `active_models.memory_pressure.soft_bytes` | — | Soft memory limit |
| `omlx_active_models_memory_hard_bytes` | gauge | `active_models.memory_pressure.hard_bytes` | — | Hard memory limit |

### Source: `GET /api/stats` — `active_models` per-model request detail

| Prometheus Metric | Type | oMLX JSON Field | Labels | Description |
|---|---|---|---|---|
| `omlx_model_waiting_requests` | gauge | `active_models.models[].waiting` length | `model` | Number of requests waiting in queue |
| `omlx_model_prefilling_requests` | gauge | `active_models.models[].prefilling` length | `model` | Number of requests in prefill phase |
| `omlx_model_generating_requests` | gauge | `active_models.models[].generating` length | `model` | Number of requests in generation phase |
| `omlx_model_generating_tokens_per_second` | gauge | `active_models.models[].generating[].tokens_per_second` | `model`, `request_id` | Real-time throughput per generating request |
| `omlx_model_generating_generated_tokens` | gauge | `active_models.models[].generating[].generated_tokens` | `model`, `request_id` | Tokens generated so far for this request |
| `omlx_model_generating_elapsed_seconds` | gauge | `active_models.models[].generating[].elapsed_seconds` | `model`, `request_id` | How long this request has been generating |
| `omlx_model_waiting_queue_position` | gauge | `active_models.models[].waiting[].queue_position` | `model`, `request_id` | Position in waiting queue |
| `omlx_model_waiting_elapsed_seconds` | gauge | `active_models.models[].waiting[].elapsed_seconds` | `model`, `request_id` | How long this request has been waiting |
| `omlx_model_active_requests` | gauge | `active_models.models[].active_requests` | `model` | Active requests for this model |

### Source: `GET /api/stats` — `runtime_cache`

| Prometheus Metric | Type | oMLX JSON Field | Description |
|---|---|---|---|
| `omlx_cache_ssd_total_files` | gauge | `runtime_cache.total_num_files` | Total SSD cache files |
| `omlx_cache_ssd_total_size_bytes` | gauge | `runtime_cache.total_size_bytes` | Total SSD cache size |
| `omlx_cache_ssd_max_bytes` | gauge | `runtime_cache.disk_max_bytes` | SSD cache max configured size |
| `omlx_cache_hot_size_bytes` | gauge | `runtime_cache.hot_cache_size_bytes` | Hot cache size |
| `omlx_cache_hot_max_bytes` | gauge | `runtime_cache.hot_cache_max_bytes` | Hot cache max configured size |
| `omlx_cache_hot_entries` | gauge | `runtime_cache.hot_cache_entries` | Hot cache entry count |

### Error / Health Metrics

| Prometheus Metric | Type | Description |
|---|---|---|
| `omlx_scrape_duration_seconds` | histogram | Time to complete one full scrape cycle |
| `omlx_scrape_failures_total` | counter | Total number of scrape failures |

## Open Questions

- Should the exporter support scraping multiple oMLX servers simultaneously? (single target assumed for now)
