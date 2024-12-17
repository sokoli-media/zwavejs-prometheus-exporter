package prometheus_exporter

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"sync"
)

func RunHTTPServer(logger *slog.Logger, config MosquittoConfig) {
	var wg sync.WaitGroup
	quitZWaveMetrics := make(chan bool, 1)

	wg.Add(1)
	go CollectZWaveMetrics(logger, config, &wg, quitZWaveMetrics)

	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		logger.Error("failed to run http server", "error", err)
		return
	}
}
