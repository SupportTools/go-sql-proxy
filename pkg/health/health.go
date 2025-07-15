package health

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	// Import MySQL driver for database connectivity
	_ "github.com/go-sql-driver/mysql"
	"github.com/supporttools/go-sql-proxy/pkg/config"
	"github.com/supporttools/go-sql-proxy/pkg/logging"
)

// VersionInfo represents the structure of version information.
type VersionInfo struct {
	Version   string `json:"version"`
	GitCommit string `json:"gitCommit"`
	BuildTime string `json:"buildTime"`
}

var logger = logging.SetupLogging()

// version holds the application version. It's set during the build process.
var version = "MISSING VERSION INFO"

// GitCommit holds the Git commit hash of the build. It's set during the build process.
var GitCommit = "MISSING GIT COMMIT"

// BuildTime holds the timestamp of when the build was created. It's set during the build process.
var BuildTime = "MISSING BUILD TIME"

// HealthzHandler returns an HTTP handler function that checks database connectivity.
func HealthzHandler(username, password, host string, port int, database string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("HealthzHandler")

		// Construct the DSN (Data Source Name) string
		dsn := buildDSN(username, password, host, port, database)

		// Open a new database connection
		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			logger.Error("HealthzHandler: Failed to connect to database via the proxy", err)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
		defer conn.Close()

		// Ping the database to check connectivity
		if err := conn.Ping(); err != nil {
			logger.Error("HealthzHandler: Database ping failed via the proxy", err)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}

		logger.Info("HealthzHandler: Database is reachable via the proxy")
		// If the database is reachable, return "ok"
		fmt.Fprintf(w, "ok")
	}
}

// ReadyzHandler returns an HTTP handler function that checks database connectivity.
func ReadyzHandler(username, password, host string, port int, database string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("ReadyzHandler")

		// Construct the DSN (Data Source Name) string
		dsn := buildDSN(username, password, host, port, database)

		// Open a new database connection
		conn, err := sql.Open("mysql", dsn)
		if err != nil {
			logger.Error("ReadyzHandler: Failed to connect to database directly", err)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}
		defer conn.Close()

		// Ping the database to check connectivity
		if err := conn.Ping(); err != nil {
			logger.Error("ReadyzHandler: Database ping failed directly", err)
			http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
			return
		}

		logger.Info("ReadyzHandler: Database is reachable directly")
		// If the database is reachable, return "ok"
		fmt.Fprintf(w, "ok")
	}
}

// VersionHandler returns version information as JSON.
func VersionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("VersionHandler")

		versionInfo := VersionInfo{
			Version:   version,
			GitCommit: GitCommit,
			BuildTime: BuildTime,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(versionInfo); err != nil {
			logger.Error("Failed to encode version info to JSON", err)
			http.Error(w, "Failed to encode version info", http.StatusInternalServerError)
		}
	}
}

// buildDSN constructs a MySQL DSN with optional TLS parameters
func buildDSN(username, password, host string, port int, database string) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, database)
	
	if config.CFG.UseSSL {
		// Add TLS parameters to the DSN
		tlsConfig := "?tls=true"
		
		// For custom CA verification, we'd need to register a custom TLS config
		// with the MySQL driver, but for basic SSL with skip-verify, this works
		if config.CFG.SSLSkipVerify {
			tlsConfig = "?tls=skip-verify"
		}
		
		dsn += tlsConfig
	}
	
	return dsn
}
