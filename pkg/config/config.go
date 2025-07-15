package config

import (
	"log"
	"os"
	"strconv"
)

// AppConfig structure for environment-based configurations.
type AppConfig struct {
	Debug                  bool   `json:"debug"`
	MetricsPort            int    `json:"metricsPort"`
	SourceDatabaseServer   string `json:"sourceDatabaseServer"`
	SourceDatabasePort     int    `json:"sourceDatabasePort"`
	SourceDatabaseUser     string `json:"sourceDatabaseUser"`
	SourceDatabasePassword string `json:"sourceDatabasePassword"`
	SourceDatabaseName     string `json:"sourceDatabaseName"`
	BindAddress            string `json:"bindAddress"`
	BindPort               int    `json:"bindPort"`
	UseSSL                 bool   `json:"useSSL"`
	SSLSkipVerify          bool   `json:"sslSkipVerify"`
	SSLCAFile              string `json:"sslCAFile"`
	SSLCertFile            string `json:"sslCertFile"`
	SSLKeyFile             string `json:"sslKeyFile"`
}

var CFG AppConfig

// LoadConfiguration loads configuration from environment variables.
func LoadConfiguration() {
	CFG.Debug = parseEnvBool("DEBUG", false)            // Assuming false as the default value
	CFG.MetricsPort = parseEnvInt("METRICS_PORT", 9090) // Assuming 9090 as the default port
	CFG.SourceDatabaseServer = getEnvOrDefault("SOURCE_DATABASE_SERVER", "example-db-do-user-1022819-0.c.db.ondigitalocean.com")
	CFG.SourceDatabasePort = parseEnvInt("SOURCE_DATABASE_PORT", 25060)
	CFG.SourceDatabaseUser = getEnvOrDefault("SOURCE_DATABASE_USER", "doadmin")
	CFG.SourceDatabasePassword = getEnvOrDefault("SOURCE_DATABASE_PASSWORD", "")
	CFG.SourceDatabaseName = getEnvOrDefault("SOURCE_DATABASE_NAME", "defaultdb")
	CFG.BindAddress = getEnvOrDefault("BIND_ADDRESS", "0.0.0.0")
	CFG.BindPort = parseEnvInt("BIND_PORT", 3306)
	CFG.UseSSL = parseEnvBool("USE_SSL", false)
	CFG.SSLSkipVerify = parseEnvBool("SSL_SKIP_VERIFY", false)
	CFG.SSLCAFile = getEnvOrDefault("SSL_CA_FILE", "")
	CFG.SSLCertFile = getEnvOrDefault("SSL_CERT_FILE", "")
	CFG.SSLKeyFile = getEnvOrDefault("SSL_KEY_FILE", "")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseEnvInt(key string, defaultValue int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Error parsing %s as int: %v. Using default value: %d", key, err, defaultValue)
		return defaultValue
	}
	return intValue
}

func parseEnvBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("Error parsing %s as bool: %v. Using default value: %t", key, err, defaultValue)
		return defaultValue
	}
	return boolValue
}
