package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

type DummyMetricHandler struct {
	Metrics []Metric
}

var counters map[string]int

func init() {
	counters = make(map[string]int)
}

func (h DummyMetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, metric := range h.Metrics {
		// If no name is given, skip!
		if metric.Name == "" {
			continue
		}

		// If gauge, generate random value and append
		if metric.Type == "gauge" {
			fmt.Fprintln(w, metric.Render(rand.Intn(10-1)+1))
		}

		// If Counter, increase counter and append
		if metric.Type == "counter" {
			counters[metric.Name] = counters[metric.Name] + rand.Intn(10-1) + 1
			fmt.Fprintln(w, metric.Render(counters[metric.Name]))
		}
	}
}
