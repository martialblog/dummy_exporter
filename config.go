package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Metrics []Metric
}

// TODO LabelName regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
type Metric struct {
	Name   string              `json:"metric"`
	Labels map[string][]string `json:"labels"`
	Type   string              `json:"type"`
}

// Renders the Metric in the Exposition Format with the given value.
// # HELP http_requests_total The total number of HTTP requests.
// # TYPE http_requests_total counter
// http_requests_total{method="post",code="200"}
func (m *Metric) Render(value int) string {
	renderedLabels := m.RenderLabels()
	labels := permute(renderedLabels...)

	var sb strings.Builder

	for _, lbs := range labels {
		s := fmt.Sprintf("%s\n%s\n%s\n",
			fmt.Sprintf("# HELP %s", m.Name),
			fmt.Sprintf("# TYPE %s %s", m.Name, m.Type),
			fmt.Sprintf("dummy_%s{%s} %d", m.Name, strings.TrimRight(lbs, ","), value))

		sb.WriteString(s)
	}

	return sb.String()
}

func (c *Config) LoadConfig(confFile string) (err error) {
	conf, err := os.ReadFile(confFile)

	if err != nil {
		return fmt.Errorf("Error while opening configuration: %v", err)
	}

	var m []Metric
	err = json.Unmarshal(conf, &m)

	if err != nil {
		return fmt.Errorf("Error during Unmarshal: %v", err)
	}

	c.Metrics = m

	return nil
}

func (m *Metric) RenderLabels() [][]string {
	result := make([][]string, 0, len(m.Labels))

	for name, values := range m.Labels {
		var at = make([]string, 0, len(values))

		for _, value := range values {
			at = append(at, fmt.Sprintf("%s=\"%s\",", name, value))
		}

		result = append(result, at)
	}

	return result
}

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
