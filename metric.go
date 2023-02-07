package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

// Counters for the Counter metrics
var counters map[string]int

func init() {
	counters = make(map[string]int)
}

type Metric struct {
	Name             string              `json:"metric"`
	Type             string              `json:"type"`
	Labels           map[string][]string `json:"labels"`
	Min              int                 `json:"min"`
	Max              int                 `json:"max"`
	permutatedLabels []string
}

func (m *Metric) UnmarshalJSON(b []byte) error {
	// Default values before unmarshaling
	m.Type = "gauge"
	m.Min = 1
	m.Max = 10

	// Create temporary struct
	type Temp Metric
	t := (*Temp)(m)

	err := json.Unmarshal(b, t)

	if err != nil {
		return err
	}

	if m.Name == "" {
		return fmt.Errorf("Metric Name cannot be empty")
	}

	if len(m.Labels) == 0 {
		return fmt.Errorf("Metric Labels cannot be empty")
	}

	m.PermuteLabels()

	return nil
}

// Renders the Metric in the Exposition Format with the given value.
// # HELP http_requests_total The total number of HTTP requests.
// # TYPE http_requests_total counter
// http_requests_total{method="post",code="200"}
func (m *Metric) String() string {

	// Might be too obscure?
	if len(m.permutatedLabels) == 0 {
		m.PermuteLabels()
	}

	var (
		sb    strings.Builder
		value int
	)

	for _, lbs := range m.permutatedLabels {
		// If Gauge, use a random value
		if m.Type == "gauge" {
			value = rand.Intn(m.Max-m.Min+1) + m.Min
		}

		// If Counter, increase counter and append
		if m.Type == "counter" {
			counters[lbs] = counters[lbs] + rand.Intn(m.Max-m.Min+1) + m.Min
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

// We precalculate these, since they do not change and are
// computationally expensive
func (m *Metric) PermuteLabels() {
	renderedLabels := m.RenderLabels()
	m.permutatedLabels = permute(renderedLabels...)
}

// Renders Labels (map of []string)
// into a [][]string containing
// the labels ready to print. E.g.
// [{x="foo",}, {x="bar",}]
func (m *Metric) RenderLabels() [][]string {
	result := make([][]string, 0, len(m.Labels))

	//Sort Labels
	sortedLabels := make([]string, 0, len(m.Labels))
	for k, _ := range m.Labels {
		sortedLabels = append(sortedLabels, k)
	}
	sort.Strings(sortedLabels)

	for _, name := range sortedLabels {
		var at = make([]string, 0, len(m.Labels[name]))

		for _, value := range m.Labels[name] {
			at = append(at, name+"=\""+value+"\",")
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
