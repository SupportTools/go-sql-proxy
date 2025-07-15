# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

go-sql-proxy is a MySQL proxy server that acts as an intermediary between MySQL clients and servers. It provides transparent traffic forwarding with protocol decoding capabilities, metrics collection, health checks, and connection management.

## Development Commands

### Build and Run
```bash
# Run locally
go run main.go

# Build binary
go build -o go-sql-proxy

# Build Docker image
make build TAG=1.0.0

# Push to registry
make push TAG=1.0.0

# Build and push
make all TAG=1.0.0
```

### Testing
Currently, there are no unit tests in the codebase. When adding new functionality, consider creating appropriate test files.

## Architecture

### Core Components

1. **main.go**: Entry point that initializes configuration, starts metrics server, creates proxy instance, and handles graceful shutdown
2. **pkg/proxy**: Core proxy logic - accepts client connections, establishes upstream connections, and manages bidirectional data transfer
3. **pkg/protocol**: MySQL protocol decoding/encoding for packet inspection
4. **pkg/metrics**: Prometheus metrics collection and HTTP endpoints
5. **pkg/health**: Health check endpoints that verify proxy and upstream connectivity
6. **pkg/config**: Environment-based configuration management

### Connection Flow

1. Client connects to proxy on `BIND_PORT` (default: 3306)
2. Proxy establishes connection to `SOURCE_DATABASE_SERVER:SOURCE_DATABASE_PORT`
3. Data is transferred bidirectionally using `io.Copy` with optional protocol decoding
4. Metrics are collected for bytes transferred and connection counts

### Key Design Patterns

- **Context-based lifecycle management**: Uses Go contexts for graceful shutdown
- **Concurrent connection handling**: Each client connection runs in its own goroutine
- **Centralized logging**: All packages use the shared Logrus logger from pkg/logging
- **Environment configuration**: All settings come from environment variables

### Environment Variables

- `DEBUG`: Enable debug logging
- `METRICS_PORT`: Port for metrics/health endpoints (default: 9090)
- `SOURCE_DATABASE_SERVER`: Target MySQL server hostname
- `SOURCE_DATABASE_PORT`: Target MySQL server port (default: 25060)
- `SOURCE_DATABASE_USER`: MySQL username
- `SOURCE_DATABASE_PASSWORD`: MySQL password
- `SOURCE_DATABASE_NAME`: Default database name
- `BIND_ADDRESS`: Proxy bind address (default: 0.0.0.0)
- `BIND_PORT`: Proxy listening port (default: 3306)
- `USE_SSL`: Enable SSL/TLS connection to upstream MySQL (default: false)
- `SSL_SKIP_VERIFY`: Skip SSL certificate verification (default: false)
- `SSL_CA_FILE`: Path to CA certificate file for SSL verification
- `SSL_CERT_FILE`: Path to client certificate file for mutual TLS
- `SSL_KEY_FILE`: Path to client key file for mutual TLS

### Metrics and Health Endpoints

Available on `METRICS_PORT`:
- `/metrics`: Prometheus metrics
- `/healthz`: Liveness check (tests proxy connectivity)
- `/readyz`: Readiness check (tests upstream MySQL connectivity)
- `/version`: Version information

### Important Implementation Details

- The proxy uses `io.Copy` for efficient data transfer between connections
- Protocol decoding is optional and controlled by configuration
- Connection errors are logged but don't crash the proxy
- Each connection tracks bytes transferred in both directions
- Version information is injected at build time using LDFLAGS
- SSL/TLS support is controlled by the `USE_SSL` flag instead of port-based logic
- Health checks also respect SSL settings when connecting to the database