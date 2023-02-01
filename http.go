package main

import (
	"fmt"
	"net/http"
)

type DummyMetricHandler struct{}

func (h DummyMetricHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Some text")
}
