package main

import (
	"log/slog"
	"os"
	"strings"
	"tvheadend-prometheus-exporter/prometheus_exporter"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	tvheadendHost := os.Getenv("TVHEADEND_HOST")
	tvheadendUsername := os.Getenv("TVHEADEND_USERNAME")
	tvheadendPassword := os.Getenv("TVHEADEND_PASSWORD")

	if tvheadendHost == "" || tvheadendUsername == "" || tvheadendPassword == "" {
		logger.Error("one or more environment variables are not set")
		return
	}

	if !strings.HasPrefix(tvheadendHost, "http://") && strings.HasPrefix(tvheadendHost, "https://") {
		logger.Error("tvheadend domain should have http:// or https:// prefix")
		return
	}

	config := prometheus_exporter.TvheadendConfig{
		Host:     tvheadendHost,
		Username: tvheadendUsername,
		Password: tvheadendPassword,
	}
	prometheus_exporter.RunHTTPServer(logger, config)
}
