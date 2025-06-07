# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an OpenTelemetry observability stack setup using Docker Compose. The architecture implements a complete telemetry pipeline for collecting, processing, and visualizing metrics, logs, and traces from microservices.

## Architecture

**Core Components:**
- **OpenTelemetry Collector** (port 4317/4318) - Central telemetry data processor and forwarder
- **Prometheus** (port 9090) - Metrics collection and storage
- **Grafana Loki** (port 3100) - Log aggregation and storage
- **Grafana Tempo** (port 3200) - Distributed tracing backend
- **Jaeger** (port 16686) - Distributed tracing UI and storage
- **Grafana** (port 3000) - Unified visualization dashboard (admin/admin)

**Data Flow:**
Applications → OpenTelemetry Collector → Prometheus/Loki/Tempo/Jaeger → Grafana

**Service Dependencies:**
The OpenTelemetry Collector depends on all backend services (Loki, Tempo, Prometheus, Jaeger) and must start after them. Applications depend on the collector being available.

## Common Commands

**Start the entire stack:**
```bash
docker compose up -d
```

**View logs for specific services:**
```bash
docker compose logs -f otel-collector
docker compose logs -f grafana
```

**Stop and remove everything:**
```bash
docker compose down -v
```

**Rebuild and restart after config changes:**
```bash
docker compose down && docker compose up -d
```

## Configuration

Each service has its configuration in dedicated directories:
- `otel-collector/config.yaml` - Collector receivers, processors, and exporters
- `prometheus/config.yaml` - Scrape targets and retention policies
- `loki/config.yaml` - Log ingestion and storage configuration
- `grafana/provisioning/` - Data sources and dashboard provisioning

## Application Integration

Applications should send telemetry to the OpenTelemetry Collector:
- **OTLP gRPC:** `http://otel-collector:4317`
- **OTLP HTTP:** `http://otel-collector:4318`

Set environment variables:
```bash
OTEL_EXPORTER_OTLP_ENDPOINT=http://otel-collector:4317
SERVICE_NAME=your-service-name
```

## Access Points

- Grafana: http://localhost:3000 (admin/admin)
- Prometheus: http://localhost:9090
- Loki: http://localhost:3100
- Tempo: http://localhost:3200
- Jaeger UI: http://localhost:16686