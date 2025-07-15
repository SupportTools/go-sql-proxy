package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/supporttools/go-sql-proxy/pkg/config"
	"github.com/supporttools/go-sql-proxy/pkg/logging"
	"github.com/supporttools/go-sql-proxy/pkg/metrics"
	"github.com/supporttools/go-sql-proxy/pkg/proxy"
)

var logger = logging.SetupLogging()

func main() {
	logger.Println("Starting go-sql-proxy server...")
	config.LoadConfiguration()
	if config.CFG.Debug {
		logger.Println("Debug mode enabled")
		logger.Println("Configuration:")
		logger.Printf("Debug: %t", config.CFG.Debug)
		logger.Printf("Metrics Port: %d", config.CFG.MetricsPort)
		logger.Printf("Source Database Server: %s", config.CFG.SourceDatabaseServer)
		logger.Printf("Source Database Port: %d", config.CFG.SourceDatabasePort)
		logger.Printf("Source Database User: %s", config.CFG.SourceDatabaseUser)
		//logger.Printf("Source Database Password: %s", config.CFG.SourceDatabasePassword)
		logger.Printf("Bind Address: %s", config.CFG.BindAddress)
		logger.Printf("Bind Port: %d", config.CFG.BindPort)
	}

	go func() {
		logger.Println("Starting metrics server...")
		metrics.StartMetricsServer()
	}()

	// Create a context to manage the proxy server
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called to release resources if main exits before signal

	p := proxy.NewProxy(ctx, config.CFG.SourceDatabaseServer, config.CFG.SourceDatabasePort, config.CFG.UseSSL)
	p.EnableDecoding = true

	var wg sync.WaitGroup

	// Setup signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Signal received, stopping and exiting...")
		cancel()   // Notify all operations to start shutting down
		wg.Wait()  // Wait for all goroutines to finish
		os.Exit(0) // Ensure the program exits
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Start the proxy server in a goroutine to allow shutdown process to proceed
		err := proxy.StartProxy(p, config.CFG.BindPort)
		if err != nil {
			log.Fatalf("Failed to start proxy server: %v", err)
		}
	}()

	// Wait here until signal is received and handled
	<-ctx.Done()
}
