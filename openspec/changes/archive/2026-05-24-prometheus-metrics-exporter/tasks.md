## 1. Configuration

- [x] 1.1 Extend `config/config.go` with OmlxExporter config struct (target_url, api_key, scrape_interval)
- [x] 1.2 Extend `main.go` Config struct to include exporter settings
- [x] 1.3 Update `config.yaml` with omlx exporter configuration section

## 2. Metrics Registration

### 2.1 Server Stats (`/api/stats`)
- [x] 2.1.1 Create `metrics/stats.go` - register counters: total_prompt_tokens, total_completion_tokens, total_cached_tokens, total_requests
- [x] 2.1.2 Register gauges: cache_efficiency, avg_prefill_tps, avg_generation_tps, uptime_seconds

### 2.2 Active Models (`/api/stats` → `active_models`)
- [x] 2.2.1 Register gauge: total_active_requests, total_waiting_requests
- [x] 2.2.2 Register gauges: memory_pressure_level, memory_current_bytes, memory_soft_bytes, memory_hard_bytes

### 2.3 Per-Model Request Detail (`/api/stats` → `active_models.models[]`)
- [x] 2.3.1 Register gauges: model_waiting_requests, model_prefilling_requests, model_generating_requests, model_active_requests (labels: model)
- [x] 2.3.2 Register gauge: model_generating_tokens_per_second, model_generating_generated_tokens, model_generating_elapsed_seconds (labels: model)
- [x] 2.3.3 Register gauge: model_waiting_queue_position, model_waiting_elapsed_seconds (labels: model)

### 2.4 Health / Error Metrics
- [x] 2.4.1 Create `metrics/health.go` - register histogram `scrape_duration_seconds`, counter `scrape_failures_total`

## 3. API Client and Response Models

- [x] 3.1 Create `client/client.go` - HTTP client using `github.com/txix-open/isp-kit/http/httpcli` with:
  - `New(baseURL, apiKey)` - creates client with API key
  - `RefreshSession(ctx)` - POST /admin/api/login with API key, extracts `omlx_admin_session` cookie
  - `FetchStats(ctx)` - GET /admin/api/stats with session cookie
  - `NeedsSessionRefresh(err)` - detects 401 errors to trigger session refresh
- [x] 3.2 Create `client/models.go` - Go structs for oMLX API responses:
  - `StatsResponse` (top-level: total_prompt_tokens, total_completion_tokens, total_cached_tokens, total_requests, cache_efficiency, avg_prefill_tps, avg_generation_tps, uptime_seconds, active_models)
  - `ActiveModelsResponse` (total_active_requests, total_waiting_requests, memory_pressure, models)
  - `MemoryPressureResponse` (enabled, current_bytes, soft_bytes, hard_bytes, pressure_level as string)
  - `ActiveModel` (model_id, active_requests, waiting_requests, prefilling, generating, waiting, idle_seconds, model_memory_used, model_memory_max)
  - `WaitingRequest` (request_id, queue_position, elapsed_seconds, prompt_tokens)
  - `GeneratingRequest` (request_id, elapsed_seconds, generated_tokens, tokens_per_second, last_activity_age_seconds, prompt_tokens, max_tokens)
- [x] 3.3 Implement `client.Client.FetchStats()` - GET /admin/api/stats, scope=session using httpcli

## 4. Scraper

- [x] 4.1 Create `exporter/scraper.go` - Scraper struct holding HTTP client (httpcli), prometheus registry, worker, and error counter
- [x] 4.2 Implement `scraper.Run()` - periodic scrape loop using `github.com/txix-open/isp-kit/worker`
- [x] 4.3 Implement `scraper.updateStats()` - parse StatsResponse, update stats counters/gauges
- [x] 4.4 Implement `scraper.updateActiveModels()` - parse active_models, update request/memory gauges and per-model request detail (waiting, prefilling, generating counts, generating tokens_per_second with request_id label)
- [x] 4.5 Implement `scraper.recordScrapeDuration()` - record histogram for scrape time
- [x] 4.6 Implement `scraper.recordFailure()` - increment scrape_failures_total on error
- [x] 4.7 Implement graceful stop via context cancellation

## 5. Integration

- [x] 5.1 Wire scraper into `main.go` as an `app.Runner`
- [x] 5.2 Add scraper stop as an `app.Closer`
- [x] 5.3 Add startup and shutdown logging
- [x] 5.4 Ensure `/metrics` endpoint is served by the existing infra.Server

## 6. Build and Verification

- [x] 6.1 Run `go mod tidy` to update dependencies
- [x] 6.2 Verify build succeeds with `go build ./...`
- [x] 6.3 Test exporter starts and scrapes metrics correctly
- [x] 6.4 Verify `/metrics` endpoint returns valid Prometheus format
