package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"sync"
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
	// There might be a better way.
	if len(m.permutatedLabels) == 0 {
		m.PermuteLabels()
	}

	var (
		sb    strings.Builder
		value int
	)

	for _, lbs := range m.permutatedLabels {
		// If Gauge, use a random value
		// nolint: gosec
		if m.Type == "gauge" {
			value = rand.Intn(m.Max-m.Min+1) + m.Min
		}

		// If Counter, increase counter and append
		if m.Type == "counter" {
			// nolint: gosec
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
	l := Product(renderedLabels...)

	// Calculate capacity
	cap := 0
	for _, v := range m.Labels {
		cap = +len(v)
	}

	// for m.Lables =+ len(lab)
	lbs := make([]string, 0, cap)
	for p := range l {
		lbs = append(lbs, strings.Join(p, ""))
	}

	m.permutatedLabels = lbs
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

// Generates all combinations from a slice of slices
// A × B = {1,2} × {3,4} = {(1,3), (1,4), (2,3), (2,4)}
func Product(items ...[]string) chan []string {
	ch := make(chan []string)

	var wg sync.WaitGroup

	wg.Add(1)

	iterate(&wg, ch, []string{}, items...)

	go func() { wg.Wait(); close(ch) }()

	return ch
}

func iterate(wg *sync.WaitGroup, channel chan []string, result []string, items ...[]string) {
	defer wg.Done()

	// Return of no more elements left
	if len(items) == 0 {
		channel <- result
		return
	}

	// Shift Items
	item, items := items[0], items[1:]

	for i := 0; i < len(item); i++ {
		wg.Add(1)

		// Copy of the result
		copyOfResults := append([]string{}, result...)
		// Recursion with remaining items
		go iterate(wg, channel, append(copyOfResults, item[i]), items...)
	}
}
