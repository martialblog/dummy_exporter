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
	Le               []float64           `json:"le"`
	Quantile         []float64           `json:"quantile"`
	permutatedLabels []string
}

// Returns a random int between min and max
func randInt(min int, max int) int {
	// nolint: gosec
	return rand.Intn(max-min+1) + min
}

func (m *Metric) UnmarshalJSON(b []byte) error {
	// Default values before unmarshaling
	m.Type = "gauge"
	m.Min = 1
	m.Max = 10
	m.Le = []float64{}
	m.Quantile = []float64{}

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

	sort.Float64s(m.Le)
	sort.Float64s(m.Quantile)
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

	sb.WriteString(fmt.Sprintf("# HELP %s\n", m.Name))
	sb.WriteString(fmt.Sprintf("# TYPE %s %s\n", m.Name, m.Type))

	for _, lbs := range m.permutatedLabels {
		switch m.Type {
		case "gauge":
			// Gauge, return a random value
			value = randInt(m.Min, m.Max)
			sb.WriteString(fmt.Sprintf("dummy_%s{%s} %d\n", m.Name, strings.TrimRight(lbs, ","), value))

		case "counter":
			// Counter, Set the counter and increment it with a random value
			counters[lbs] = counters[lbs] + randInt(m.Min, m.Max)
			value = counters[lbs]
			sb.WriteString(fmt.Sprintf("dummy_%s{%s} %d\n", m.Name, strings.TrimRight(lbs, ","), value))

		case "summary":
			// Summary, Set an initial counter for the lowest qu and increment it with a random value
			// TODO Basically works, but still needs some polishing
			counters[lbs] = counters[lbs] + randInt(m.Min, m.Max)
			value = counters[lbs]

			for _, qu := range m.Quantile {
				value += randInt(m.Min, m.Max)
				sb.WriteString(fmt.Sprintf("dummy_%s{%s quantile=\"%g\"} %d\n", m.Name, lbs, qu, value))
			}
			// Add count and sum
			sb.WriteString(fmt.Sprintf("dummy_%s_count{%s} %d\n", m.Name, lbs, value))
			sb.WriteString(fmt.Sprintf("dummy_%s_sum{%s} %d\n", m.Name, lbs, value*2))

		case "histogram":
			// Histogram, Set an initial counter for the lowest bucket and increment it with a random value
			counters[lbs] = counters[lbs] + randInt(m.Min, m.Max)
			value = counters[lbs]

			for _, le := range m.Le {
				value += randInt(m.Min, m.Max)
				sb.WriteString(fmt.Sprintf("dummy_%s_bucket{%s le=\"%g\"} %d\n", m.Name, lbs, le, value))
			}
			// Add final +Inf bucket
			sb.WriteString(fmt.Sprintf("dummy_%s_bucket{%s le=\"+Inf\"} %d\n", m.Name, lbs, value))
			// Add count and sum
			sb.WriteString(fmt.Sprintf("dummy_%s_count{%s} %d\n", m.Name, lbs, value))
			sb.WriteString(fmt.Sprintf("dummy_%s_sum{%s} %d\n", m.Name, lbs, value*2))

		default:
			sb.WriteString(fmt.Sprintf("dummy_%s{%s} %d\n", m.Name, strings.TrimRight(lbs, ","), value))
		}
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

	sort.Strings(lbs)
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
