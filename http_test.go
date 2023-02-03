package main

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServiceHTTP(t *testing.T) {
	testcases := map[string]struct {
		handler  DummyMetricHandler
		expected string
	}{
		"http-gauge-metric": {
			expected: "dummy_foogauge{foo=\"bar\"}",
			handler: DummyMetricHandler{
				Metrics: []Metric{
					Metric{
						Name: "foogauge",
						Labels: map[string][]string{
							"foo": []string{"bar"},
						},
						Type: "gauge",
					},
				},
			},
		},
		"http-counter-type": {
			expected: "# TYPE foocounter counter",
			handler: DummyMetricHandler{
				Metrics: []Metric{
					Metric{
						Name: "foocounter",
						Labels: map[string][]string{
							"foo": []string{"bar"},
						},
						Type: "counter",
					},
				},
			},
		},
		"http-no-name": {
			expected: "",
			handler: DummyMetricHandler{
				Metrics: []Metric{
					Metric{
						Labels: map[string][]string{
							"foo": []string{"bar"},
						},
						Type: "counter",
					},
				},
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {

			req := httptest.NewRequest("GET", "http://localhost", nil)
			w := httptest.NewRecorder()

			tc.handler.ServeHTTP(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			actual := string(body)

			if !strings.Contains(actual, tc.expected) {
				t.Error("\nActual: ", actual, "\nExpected: ", tc.expected)
			}
		})
	}
}
