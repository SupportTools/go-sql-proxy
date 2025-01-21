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

You can customize the server configurations by setting the environment variables as defined in `pkg/config/config.go`.

## Note

- This implementation serves as a basic demonstration of a SQL proxy server and may require additional features for production use.
- Ensure proper configuration of database access and proxy settings before running the server in a production environment.

Feel free to explore and extend this project for your specific use case or requirements!
