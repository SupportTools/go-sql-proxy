<p align="center">
  <img src="https://cdn.support.tools/go-sql-proxy/go-sql-proxy-no-bg.png">
</p>

# Go-SQL-Proxy

This project implements a SQL proxy server in Golang that acts as an intermediary between a client and a MySQL server. The proxy server can handle protocol decoding and data transfer while measuring and updating metrics for monitoring purposes.

[![Go Report Card](https://goreportcard.com/badge/github.com/SupportTools/go-sql-proxy)](https://goreportcard.com/report/github.com/SupportTools/go-sql-proxy)
[![Go Reference](https://pkg.go.dev/badge/github.com/SupportTools/go-sql-proxy.svg)](https://pkg.go.dev/github.com/SupportTools/go-sql-proxy)

## Project Structure

The project is structured as follows:

- **pkg/**
  - **config/**
    - `config.go`: Contains the configuration settings and loads them from environment variables.
  - **logging/**
    - `logging.go`: Handles setting up and configuring the logger using Logrus.
  - **metrics/**
    - `metrics.go`: Implements metrics collection and exposes Prometheus metrics endpoints.
  - **health/**
    - `health.go`: Handles health check endpoints for database connectivity.
  - **protocol/**
    - `flags.go`: Defines MySQL capability flags used in the protocol.
    - `math.go`: Contains utility functions for mathematical operations.
    - `protocol.go`: Implements decoding and encoding of MySQL protocol packets.
  - **proxy/**
    - `NewConnection.go`: Creates a new proxy connection to the target MySQL server.
    - `NewProxy.go`: Creates a new instance of the Proxy server.
    - `HandleConnection.go`: Manages data transfer and protocol decoding for a connection.
    - `proxyHandle.go`: Handles a new connection request in a goroutine.
    - `transferData.go`: Handles data transfer between client and server while measuring latency.
    - `StartProxy.go`: Starts the proxy server and accepts incoming connections.
    - `EnableDecoding.go`: Enables protocol decoding for the proxy.
    - `handleProtocolDecoding.go`: Decodes the MySQL protocol handshake.
  - **models/**
    - `Proxy.go`: Defines the structure for the proxy server configuration and state.
    - `Connection.go`: Represents a connection to a MySQL server.
- `main.go`: Main entry point of the proxy server application.

## Running the Server

To start the Go-SQL-Proxy server, the main entry point is `main.go`. The server can be started by running this file with the necessary environment variables configured to set up the proxy settings.

To run the server, execute the following command:

```bash
go run main.go
```

## Configuration

The proxy is configured through environment variables:

### Basic Configuration
- `DEBUG`: Enable debug logging (default: false)
- `METRICS_PORT`: Port for metrics/health endpoints (default: 9090)
- `BIND_ADDRESS`: Proxy bind address (default: 0.0.0.0)
- `BIND_PORT`: Proxy listening port (default: 3306)

### Database Connection
- `SOURCE_DATABASE_SERVER`: Target MySQL server hostname
- `SOURCE_DATABASE_PORT`: Target MySQL server port (default: 25060)
- `SOURCE_DATABASE_USER`: MySQL username
- `SOURCE_DATABASE_PASSWORD`: MySQL password
- `SOURCE_DATABASE_NAME`: Default database name

### SSL/TLS Configuration
- `USE_SSL`: Enable SSL/TLS connection to upstream MySQL (default: false)
- `SSL_SKIP_VERIFY`: Skip SSL certificate verification (default: false)
- `SSL_CA_FILE`: Path to CA certificate file for SSL verification
- `SSL_CERT_FILE`: Path to client certificate file for mutual TLS
- `SSL_KEY_FILE`: Path to client key file for mutual TLS

### Example: Connecting to PlanetScale

```bash
export SOURCE_DATABASE_SERVER=your-database.planetscale.com
export SOURCE_DATABASE_PORT=3306
export SOURCE_DATABASE_USER=your-username
export SOURCE_DATABASE_PASSWORD=your-password
export SOURCE_DATABASE_NAME=your-database
export USE_SSL=true
export SSL_SKIP_VERIFY=true

go run main.go
```

## Note

- This implementation serves as a basic demonstration of a SQL proxy server and may require additional features for production use.
- Ensure proper configuration of database access and proxy settings before running the server in a production environment.

Feel free to explore and extend this project for your specific use case or requirements!
