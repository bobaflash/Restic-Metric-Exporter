package main

import (
	"log/slog"
	"net/http"
	"os"
	"restic_metric_exporter/internal"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	config, err := internal.NewConfig()
	if err != nil {
		slog.Error("Could not read config", "error", err)
		os.Exit(1)
	}

	prometheusReg := prometheus.NewRegistry()
	metrics := internal.NewMetrics(prometheusReg)

	go startWebserver(prometheusReg, ":8000")

	var requesters []internal.ResticRequester

	wg := &sync.WaitGroup{}

	for _, repo := range config.Repositories {
		requesters = append(requesters, internal.NewResticRequester(repo, wg, *metrics))
	}

	for {
		for _, requester := range requesters {
			wg.Add(1)
			requester.Run()
			wg.Wait()
		}

		time.Sleep(30 * time.Second)
	}

}

func startWebserver(prometheusReg *prometheus.Registry, listenSocket string) {

	prometheusReg.MustRegister()

	slog.Info("Starting webserver...", "Socket", listenSocket)
	http.Handle("/metrics", promhttp.HandlerFor(prometheusReg, promhttp.HandlerOpts{}))
	if err := http.ListenAndServe(listenSocket, nil); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
