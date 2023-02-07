package main

import (
	"reflect"
	"strings"
	"testing"
)

func BenchmarkProduct(b *testing.B) {
	b.ReportAllocs()
	x := [][]string{
		{"1"},
		{"2"},
	}

	for n := 0; n < b.N; n++ {
		Product(x...)
	}
}

func BenchmarkRenderLabels(b *testing.B) {
	m := Metric{
		Name: "foo",
		Min:  1,
		Max:  1,
		Labels: map[string][]string{
			"x": []string{"1", "2"},
			"y": []string{"1", "2", "3"},
			"z": []string{"1", "2", "3"},
		},
		Type: "gauge",
	}

	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		m.RenderLabels()
	}
}

func TestRenderLabels(t *testing.T) {
	testcases := map[string]struct {
		metric   Metric
		expected [][]string
	}{
		"simple": {
			metric: Metric{
				Name: "foo",
				Min:  1,
				Max:  1,
				Labels: map[string][]string{
					"m": []string{"1"},
				},
			},
			expected: [][]string{
				[]string{"m=\"1\","},
			},
		},
		"sorted": {
			metric: Metric{
				Name: "foo",
				Min:  1,
				Max:  1,
				Labels: map[string][]string{
					"z": []string{"1", "2"},
					"a": []string{"1", "2"},
					"m": []string{"1", "2", "3"},
				},
			},
			expected: [][]string{
				[]string{"a=\"1\",", "a=\"2\","},
				[]string{"z=\"1\",", "z=\"2\","},
				[]string{"m=\"1\",", "m=\"2\",", "m=\"3\","},
			},
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actual := tc.metric.RenderLabels()
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Error("Actual", actual, "Expected", tc.expected)
			}
		})
	}
}

func TestMetrics(t *testing.T) {
	testcases := map[string]struct {
		metric   Metric
		value    int
		expected string
	}{
		"no-type": {
			metric: Metric{
				Name: "foo",
				Min:  1,
				Max:  1,
				Labels: map[string][]string{
					"in": []string{"bar"},
				},
			},
			expected: "# HELP foo\n# TYPE foo \ndummy_foo{in=\"bar\"}",
		},
		"wrong-type": {
			metric: Metric{
				Name: "foo",
				Min:  1,
				Max:  1,
				Labels: map[string][]string{
					"in": []string{"bar"},
				},
				Type: "unknown",
			},
			expected: "# HELP foo\n# TYPE foo unknown\ndummy_foo{in=\"bar\"}",
		},
		"gauge": {
			metric: Metric{
				Name: "foogauge",
				Min:  1,
				Max:  1,
				Labels: map[string][]string{
					"in": []string{"bar"},
				},
				Type: "gauge",
			},
			expected: "# HELP foogauge\n# TYPE foogauge gauge\ndummy_foogauge{in=\"bar\"}",
		},
		"counter": {
			metric: Metric{
				Name: "foocounter",
				Min:  1,
				Max:  1,
				Labels: map[string][]string{
					"job": []string{"foo"},
				},
				Type: "counter",
			},
			expected: "# HELP foocounter\n# TYPE foocounter counter\ndummy_foocounter{job=\"foo\"}",
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			actual := tc.metric.String()
			if !strings.Contains(actual, tc.expected) {
				t.Error("\nActual: ", actual, "\nExpected: ", tc.expected)
			}
		})
	}
}
