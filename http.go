package main

import (
	"fmt"
	"net/http"
)

type DummyMetricHandler struct {
	Metrics []Metric
}

func (h DummyMetricHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	for _, metric := range h.Metrics {
		// If no name is given, skip!
		if metric.Name == "" {
			continue
		}

		fmt.Fprintln(w, metric.String())
	}
}
