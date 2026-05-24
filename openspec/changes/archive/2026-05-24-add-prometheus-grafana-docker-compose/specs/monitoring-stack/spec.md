## ADDED Requirements

### Requirement: Prometheus service scrapes omlx-exporter
The system SHALL run a Prometheus container that periodically scrapes the omlx-exporter metrics endpoint at `http://omlx-exporter:8001/metrics` with the job name `omlx-exporter` and a scrape interval of 15 seconds.

#### Scenario: Prometheus scrapes omlx-exporter metrics
- **WHEN** Prometheus is running and the docker-compose stack is up
- **THEN** Prometheus discovers the `omlx-exporter` job and scrapes `/metrics` every 15 seconds

#### Scenario: Prometheus stores scraped metrics
- **WHEN** Prometheus has been running for at least one scrape cycle
- **THEN** the scraped metric time series are persisted in the Prometheus data volume and queryable via the Prometheus API

### Requirement: Grafana service runs with defaults
The system SHALL run a Grafana container with default configuration, accessible on host port 3000, using the default admin/admin credentials.

#### Scenario: Grafana starts successfully
- **WHEN** `docker-compose up` is executed with the Grafana service defined
- **THEN** Grafana starts and becomes accessible at `http://localhost:3000`

#### Scenario: Grafana default credentials work
- **WHEN** a user navigates to Grafana login page
- **THEN** the user can log in with username `admin` and password `admin`

### Requirement: Grafana datasource pre-configured via provisioning
The system SHALL provision Prometheus as a Grafana datasource using a YAML configuration file mounted at `/etc/grafana/provisioning/datasources/` inside the Grafana container. The datasource SHALL be named `Prometheus`, type `prometheus`, and URL `http://prometheus:9090`.

#### Scenario: Grafana discovers Prometheus datasource on startup
- **WHEN** the Grafana container starts with the provisioning config mounted
- **THEN** Grafana automatically loads the Prometheus datasource configuration and makes it available

#### Scenario: Prometheus datasource is accessible in Grafana
- **WHEN** a user logs into Grafana and navigates to Data Sources
- **THEN** a `Prometheus` datasource is listed with URL `http://prometheus:9090` and status `Healthy`

#### Scenario: Datasource provisioning file format
- **WHEN** the provisioning YAML file is parsed by Grafana
- **THEN** it contains a datasource with `orgId: 1`, `uid: prometheus`, `type: prometheus`, `url: http://prometheus:9090`, and `access: proxy`

### Requirement: Prometheus data persists across restarts
The system SHALL use a named Docker volume (`prometheus-data`) to store Prometheus TSDB data so that metrics survive container restarts.

#### Scenario: Data persists after restart
- **WHEN** the Prometheus container is restarted (e.g., `docker-compose restart prometheus`)
- **THEN** previously scraped metric data remains available in the Prometheus UI and API

### Requirement: Prometheus and Grafana are configured via compose and mounted config
The system SHALL define both services in `docker-compose.yaml` and mount the Prometheus configuration file (`prometheus/prometheus.yml`) as a read-only volume.

#### Scenario: Prometheus uses mounted config
- **WHEN** the Prometheus container starts
- **THEN** it reads its configuration from the mounted `prometheus.yml` file instead of the default config

#### Scenario: Services are defined in docker-compose
- **WHEN** `docker-compose config` is run
- **THEN** the output includes `prometheus` and `grafana` services with correct image names, port mappings, and volume mounts
