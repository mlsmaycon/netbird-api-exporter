package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"netbird-api-exporter/pkg/exporters"
	"netbird-api-exporter/pkg/utils"
)

func main() {
	// Configuration from environment variables
	netbirdURL := utils.GetEnvWithDefault("NETBIRD_API_URL", "https://api.netbird.io")
	netbirdToken := os.Getenv("NETBIRD_API_TOKEN")
	listenAddr := utils.GetEnvWithDefault("LISTEN_ADDRESS", ":8080")
	metricsPath := utils.GetEnvWithDefault("METRICS_PATH", "/metrics")
	logLevel := utils.GetEnvWithDefault("LOG_LEVEL", "info")

	// Set log level
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// Validate required configuration
	if netbirdToken == "" {
		logrus.Fatal("NETBIRD_API_TOKEN environment variable is required")
	}

	logrus.WithFields(logrus.Fields{
		"netbird_url":  netbirdURL,
		"listen_addr":  listenAddr,
		"metrics_path": metricsPath,
		"log_level":    logLevel,
	}).Info("Starting NetBird API Exporter")

	// Create exporter
	exporter := exporters.NewNetBirdExporter(netbirdURL, netbirdToken)

	// Register exporter
	prometheus.MustRegister(exporter)

	// Create HTTP server
	mux := http.NewServeMux()

	// Metrics endpoint
	mux.Handle(metricsPath, promhttp.Handler())

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339)); err != nil {
			logrus.WithError(err).Error("Failed to write health check response")
		}
	})

	// Root endpoint with information
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		if _, err := fmt.Fprintf(w, `
		<html>
		<head><title>NetBird API Exporter</title></head>
		<body>
		<h1>NetBird API Exporter</h1>
		<p>This is a Prometheus exporter for NetBird API metrics.</p>
		<ul>
		<li><a href="%s">Metrics</a></li>
		<li><a href="/health">Health Check</a></li>
		</ul>
		<h2>Available Metrics</h2>
		<ul>
		<li><strong>Peers API:</strong> Connection status, OS distribution, geographic distribution, SSH status, login status</li>
		<li><strong>Groups API:</strong> Group counts and membership</li>
		<li><strong>Users API:</strong> User counts, roles, status, permissions</li>
		<li><strong>DNS API:</strong> Nameserver groups, DNS settings, nameserver configurations</li>
		<li><strong>Networks API:</strong> Network counts, routers, resources, policies, routing peers</li>
		</ul>
		</body>
		</html>
		`, metricsPath); err != nil {
			logrus.WithError(err).Error("Failed to write root page response")
		}
	})

	server := &http.Server{
		Addr:              listenAddr,
		Handler:           mux,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logrus.Info("Shutting down server...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logrus.WithError(err).Error("Error during server shutdown")
		}
		cancel()
	}()

	// Start server
	logrus.WithField("address", listenAddr).Info("Starting HTTP server")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("HTTP server error")
	}

	<-ctx.Done()
	logrus.Info("Server stopped")
}
