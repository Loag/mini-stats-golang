package main

import (
	"time"

	"github.com/Loag/mini-stats-golang/pkg/client"
)

func main() {
	counter := client.NewCounter("new_test_counter")
	gauge := client.NewGauge("new_test_gauge")

	ministatsopts := client.MiniStatsClientOptions{
		Debug:    true,
		ApiKey:   "key_here",
		Endpoint: "host.docker.internal:5000",
		Interval: 15,
	}
	ministats := client.New(ministatsopts)

	go ministats.
		AddMetric(counter).
		AddMetric(gauge).
		Start()

	go func() {
		t := time.NewTicker(5 * time.Second)

		for range t.C {
			counter.Inc()
			gauge.Set(50)
		}
	}()

	select {}
}
