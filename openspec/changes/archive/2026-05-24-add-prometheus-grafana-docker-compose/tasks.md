## 1. Prometheus Configuration

- [x] 1.1 Create `prometheus/` directory
- [x] 1.2 Create `prometheus/prometheus.yml` with scrape config for omlx-exporter job at `http://omlx-exporter:8001/metrics`, interval 15s

## 2. Grafana Configuration

- [x] 2.1 Create `grafana/provisioning/datasources/` directory structure
- [x] 2.2 Create `grafana/provisioning/datasources/datasources.yml` with Prometheus datasource (name: Prometheus, type: prometheus, url: http://prometheus:9090, access: proxy, uid: prometheus)

## 3. Docker-Compose Updates

- [x] 3.1 Add `prometheus-data` named volume to top-level `volumes:` section
- [x] 3.2 Add `prometheus` service using `prom/prometheus:latest` image with mounted config volume, port 9090 mapping, and `depends_on: omlx-exporter`
- [x] 3.3 Add `grafana` service using `grafana/grafana:latest` image with port 3000 mapping, mounted provisioning volume, and `depends_on: prometheus`

## 4. Verification

- [x] 4.1 Run `docker-compose up -d` and verify all three services start without errors
- [x] 4.2 Confirm Prometheus targets page shows omlx-exporter as up
- [x] 4.3 Confirm Grafana is accessible at http://localhost:3000 with admin/admin login
- [x] 4.4 Confirm Prometheus datasource is pre-configured and healthy in Grafana Data Sources page
