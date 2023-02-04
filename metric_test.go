package main

import (
	"reflect"
	"strings"
	"testing"
)

func TestPermute(t *testing.T) {
	var actual = permute(
		[]string{"job=foo", "job=bar"},
		[]string{},
		[]string{"in=1", "in=2", "in=3"},
		[]string{"code=200"},
	)

	expected := []string{"job=fooin=1code=200", "job=fooin=2code=200", "job=fooin=3code=200", "job=barin=1code=200", "job=barin=2code=200", "job=barin=3code=200"}

	if !reflect.DeepEqual(actual, expected) {
		t.Error("Actual", actual, "Expected", expected)
	}
}

func BenchmarkPermute(b *testing.B) {
	for n := 0; n < b.N; n++ {
		permute(
			[]string{"job=foo", "job=bar"},
			[]string{},
			[]string{"in=1", "in=2", "in=3"},
			[]string{"code=200"},
		)
	}
}

func TestRenderLabels(t *testing.T) {
	m := Metric{
		Name: "foo",
		Labels: map[string][]string{
			"job": []string{"foo", "bar"},
			"in":  []string{"1", "2", "3"},
		},
		Type: "gauge",
	}

	expected := [][]string{
		[]string{"job=\"foo\",", "job=\"bar\","},
		[]string{"in=\"1\",", "in=\"2\",", "in=\"3\","},
	}

	actual := m.RenderLabels()
	if !reflect.DeepEqual(actual, expected) {
		t.Error("Actual", actual, "Expected", expected)
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
				Labels: map[string][]string{
					"in": []string{"bar"},
				},
			},
			expected: "# HELP foo\n# TYPE foo \ndummy_foo{in=\"bar\"}",
		},
		"wrong-type": {
			metric: Metric{
				Name: "foo",
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
