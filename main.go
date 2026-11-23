package main

import (
	"math/big"
	"uuid"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func sum(a, b int) int {
	return a + b
}

func greet(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hello World! %s", time.Now())
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
}

func main() {
	reg := prometheus.NewRegistry()
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	fmt.Println("Server started at localhost:8080")
	http.HandleFunc("/", greet)
	go generateLog()
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error ListenAndServe %v\n", err)
	}
}

func generateLog() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		for i := 0; i < 30; i++ {
			id := uuid.New()
			num := new(big.Int).SetBytes(id[:])
			log.Printf("Log from application %d\n", num)
		}
	}
}
