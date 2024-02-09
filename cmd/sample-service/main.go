package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"prometheus-grafana/internal/gometrics"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func recordMetrics() {
	go func() {
		i := 0
		for {
			fmt.Printf("hi..")
			opsProcessed.Inc()
			time.Sleep(1 * time.Second)
			if i%2 == 0 {
				tempCelsius.Inc()
			} else {
				tempCelsius.Dec()
			}
			i = i + 1
		}
	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})

	tempCelsius = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "current_temperature_api_celsius",
			Help: "Current temperature",
		},
	)
)

func main() {
	recordMetrics()
	gometrics.StartManagementServer(":8090", "/metrics", nil)

	startServer()
}

func startServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", GetStatus)
	httpServer := http.Server{
		Addr:    "0.0.0.0:8091",
		Handler: mux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// GetStatus implements Kubernetes readiness and liveness probe endpoint (GET /status).
func GetStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
