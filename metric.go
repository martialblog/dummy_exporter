package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

var counters map[string]int

func init() {
	counters = make(map[string]int)
}

type Config struct {
	Metrics []Metric
}

// TODO LabelName regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
type Metric struct {
	Name   string              `json:"metric"`
	Labels map[string][]string `json:"labels"`
	Type   string              `json:"type"`
	// Min    int                 `json:"min"`
	// Max    int                 `json:"max"`
}

// Renders the Metric in the Exposition Format with the given value.
// # HELP http_requests_total The total number of HTTP requests.
// # TYPE http_requests_total counter
// http_requests_total{method="post",code="200"}
func (m *Metric) String() string {
	renderedLabels := m.RenderLabels()
	labels := permute(renderedLabels...)

	var sb strings.Builder
	var value int

	for _, lbs := range labels {
		// If Gauge, use a random value
		if m.Type == "gauge" {
			value = rand.Intn(10-1) + 1
		}

		// If Counter, increase counter and append
		if m.Type == "counter" {
			counters[lbs] = counters[lbs] + rand.Intn(10-1) + 1
			value = counters[lbs]
		}

		s := fmt.Sprintf("%s\n%s\n%s\n",
			fmt.Sprintf("# HELP %s", m.Name),
			fmt.Sprintf("# TYPE %s %s", m.Name, m.Type),
			fmt.Sprintf("dummy_%s{%s} %d", m.Name, strings.TrimRight(lbs, ","), value))

		sb.WriteString(s)
	}

	return sb.String()
}

// Renders Labels are a Slice of Slices
func (m *Metric) RenderLabels() [][]string {
	result := make([][]string, 0, len(m.Labels))

	for name, values := range m.Labels {
		var at = make([]string, 0, len(values))

		for _, value := range values {
			at = append(at, fmt.Sprintf("%s=\"%s\",", name, value))
		}

		result = append(result, at)
	}

	// Sort before return for easier testing
	sort.Slice(result, func(i, j int) bool {
		return len(result[i]) < len(result[j])
	})

	return result
}

// Permutes Slices of Strings
func permute(input ...[]string) (result []string) {
	{
		var n = 1
		for _, elems := range input {
			n *= len(elems)
		}
		// Calculate length of return slice
		result = make([]string, 0, n)
	}

	var at = make([]int, len(input))

	var buf bytes.Buffer

loop:
	for {
		// Position Counter
		for i := len(input) - 1; i >= 0; i-- {
			if at[i] > 0 && at[i] >= len(input[i]) {
				if i == 0 || (i == 1 && at[i-1] == len(input[0])-1) {
					break loop
				}
				at[i] = 0
				at[i-1]++
			}
		}
		// Construct string
		buf.Reset()
		for i, ar := range input {
			var p = at[i]
			if p >= 0 && p < len(ar) {
				buf.WriteString(ar[p])
			}
		}
		result = append(result, buf.String())
		at[len(input)-1]++
	}

	return result
}
