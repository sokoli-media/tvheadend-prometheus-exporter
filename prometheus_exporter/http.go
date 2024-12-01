package prometheus_exporter

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

func RunHTTPServer(logger *slog.Logger, config TvheadendConfig) {
	var wg sync.WaitGroup
	quitTvMetrics := make(chan bool, 1)

	wg.Add(1)
	go CollectTvheadendMetrics(logger, config, &wg, quitTvMetrics)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/dashboard.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/dashboards/dashboard.json")
	})

	sigIntChannel := make(chan os.Signal, 1)
	signal.Notify(sigIntChannel, os.Interrupt)
	go func() {
		<-sigIntChannel

		quitTvMetrics <- true

		wg.Wait()
		os.Exit(0)
	}()

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		logger.Error("failed to run http server", "error", err)
		return
	}
}
