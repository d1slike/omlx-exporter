## Context

The omlx-exporter project runs as a single service in docker-compose, exposing Prometheus metrics on port 8001. There is currently no metric storage or visualization layer. Operators have no way to query historical metrics or build dashboards without external infrastructure.

## Goals / Non-Goals

**Goals:**
- Add Prometheus to docker-compose with a scrape config targeting omlx-exporter
- Add Grafana to docker-compose for dashboard visualization
- Provide a working observability stack with zero additional configuration beyond `docker-compose up`
- Use latest official images for both services

**Non-Goals:**
- Pre-configured Grafana dashboards (operators add these separately)
- Alerting rules or Alertmanager integration
- Authentication or TLS for Prometheus/Grafana (internal deployment only)

## Decisions

1. **Latest official images** (`prom/prometheus:latest`, `grafana/grafana:latest`)
   - Rationale: keeps the stack up-to-date without pinning versions; the project already uses `alpine:3.23` and `golang:1.26-alpine` without explicit version pinning in docker-compose
   - Alternatives considered: pinned versions for reproducibility — rejected as unnecessary for this internal tooling

2. **External Prometheus config file mounted as volume**
   - Rationale: clean separation of config from image; operators can edit `prometheus.yml` without rebuilding
   - Alternatives considered: inline config in docker-compose command — rejected as harder to maintain

3. **Named volume for Prometheus data** (`prometheus-data`)
   - Rationale: persists scraped metrics across container restarts
   - Alternatives considered: bind mount to host directory — rejected for portability

4. **Service dependency ordering with `depends_on`**
   - Rationale: ensures proper startup order — omlx-exporter before Prometheus before Grafana
   - Implementation: Prometheus depends on omlx-exporter, Grafana depends on Prometheus
   - Alternatives considered: relying on Prometheus retry logic — rejected as adds unnecessary startup delay

5. **Grafana datasource pre-configured via provisioning**
   - Rationale: eliminates manual setup step; operators get a working stack with `docker-compose up`
   - Implementation: mount `grafana/provisioning/datasources/datasources.yml` into the Grafana container at `/etc/grafana/provisioning/datasources/`
   - Datasource points to `http://prometheus:9090` (Prometheus service name)
   - Alternatives considered: manual datasource configuration after first login — rejected as adds friction

## Risks / Trade-offs

[Risk] Latest images may introduce breaking changes
→ Mitigation: pin versions in docker-compose if stability issues arise

[Risk] Prometheus data volume grows over time
→ Mitigation: add retention config (`--storage.tsdb.retention.time`) in a follow-up change

[Risk] Grafana port 3000 may conflict with existing services
→ Mitigation: change host port mapping if needed; no code changes required

## Migration Plan

1. Add services and volume to `docker-compose.yaml`
2. Create `prometheus/prometheus.yml` with scrape config
3. Create `grafana/provisioning/datasources/datasources.yml` with Prometheus datasource pointing to `http://prometheus:9090`
4. Run `docker-compose up -d` to verify all services start and Prometheus scrapes omlx-exporter
5. Access Grafana at `http://localhost:3000` (default admin/admin) — Prometheus datasource is pre-configured

**Rollback:** Remove the added services and volume from docker-compose.yaml and delete the prometheus directory.

## Open Questions

- Should alerting rules be added for key omlx metrics? (deferred to follow-up)
