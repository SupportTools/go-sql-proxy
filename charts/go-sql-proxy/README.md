# go-sql-proxy Helm Chart

This Helm chart deploys the go-sql-proxy MySQL proxy server on Kubernetes.

## Installation

```bash
helm install my-proxy ./charts/go-sql-proxy
```

## Configuration

The following table lists the configurable parameters of the go-sql-proxy chart and their default values.

| Parameter | Description | Default |
| --------- | ----------- | ------- |
| `settings.bind.host` | Bind address for proxy | `0.0.0.0` |
| `settings.bind.port` | Bind port for proxy | `3306` |
| `settings.debug` | Enable debug logging | `false` |
| `settings.metrics.enabled` | Enable metrics endpoint | `true` |
| `settings.metrics.port` | Metrics port | `9090` |
| `settings.source.host` | Target MySQL server hostname | `example.db.ondigitalocean.com` |
| `settings.source.port` | Target MySQL server port | `25060` |
| `settings.source.user` | MySQL username | `doadmin` |
| `settings.source.password` | MySQL password | `password` |
| `settings.source.database` | Default database name | `defaultdb` |
| `settings.ssl.enabled` | Enable SSL/TLS connection | `false` |
| `settings.ssl.skipVerify` | Skip SSL certificate verification | `false` |
| `settings.ssl.caFile` | Path to CA certificate file | `""` |
| `settings.ssl.certFile` | Path to client certificate file | `""` |
| `settings.ssl.keyFile` | Path to client key file | `""` |

## SSL/TLS Configuration

To connect to SSL-enabled MySQL servers (like PlanetScale), enable SSL:

```yaml
settings:
  source:
    host: your-database.planetscale.com
    port: 3306
  ssl:
    enabled: true
    skipVerify: true  # For self-signed certificates
```

For proper certificate verification, provide CA certificate:

```yaml
settings:
  ssl:
    enabled: true
    skipVerify: false
    caFile: /path/to/ca.pem
```

For mutual TLS authentication:

```yaml
settings:
  ssl:
    enabled: true
    certFile: /path/to/client-cert.pem
    keyFile: /path/to/client-key.pem
```

## Monitoring

The proxy exposes Prometheus metrics on the configured metrics port:
- `/metrics` - Prometheus metrics
- `/healthz` - Liveness probe
- `/readyz` - Readiness probe
- `/version` - Version information