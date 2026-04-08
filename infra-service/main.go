package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type InfraStatus struct {
	Nginx   string `json:"nginx"`
	Redis   string `json:"redis"`
	Kubernetes string `json:"kubernetes"`
	Status  string `json:"status"`
}

type ServerTime struct {
	ServerTime string `json:"server_time"`
}

func infraHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(InfraStatus{
		Nginx:      "healthy",
		Redis:      "healthy",
		Kubernetes: "healthy",
		Status:     "all systems operational",
	})
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ServerTime{
		ServerTime: time.Now().UTC().Format(time.RFC3339),
	})
}

func main() {
	http.HandleFunc("/", infraHandler)
	http.HandleFunc("/time", timeHandler)
	http.ListenAndServe(":8080", nil)
}