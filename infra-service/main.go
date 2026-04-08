package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	rootCount int64
	timeCount int64
)

type InfraStatus struct {
	Nginx      string `json:"nginx"`
	Redis      string `json:"redis"`
	Kubernetes string `json:"kubernetes"`
	Status     string `json:"status"`
}

type ServerTime struct {
	ServerTime string `json:"server_time"`
}

func infraHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&rootCount, 1)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(InfraStatus{
		Nginx:      "healthy",
		Redis:      "healthy",
		Kubernetes: "healthy",
		Status:     "all systems operational",
	})
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(&timeCount, 1)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ServerTime{
		ServerTime: time.Now().UTC().Format(time.RFC3339),
	})
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "# HELP infra_requests_total Total de requests\n")
	fmt.Fprintf(w, "# TYPE infra_requests_total counter\n")
	fmt.Fprintf(w, "infra_requests_total{endpoint=\"/\"} %d\n", atomic.LoadInt64(&rootCount))
	fmt.Fprintf(w, "infra_requests_total{endpoint=\"/time\"} %d\n", atomic.LoadInt64(&timeCount))
}

func main() {
	http.HandleFunc("/", infraHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/metrics", metricsHandler)
	http.ListenAndServe(":8080", nil)
}