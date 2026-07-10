package compute_test

import (
	"math"
	"testing"

	"github.com/geoffjay/templ-charts/charts/bar/internal/compute"
	"github.com/geoffjay/templ-charts/charts/scales"
)

// TestGenerateGroupedBars_LogScale exercises the log value scale through the
// grouped path: geometry stays finite, heights stay non-negative, and the log
// mapping stays monotonic (a larger value yields a taller bar). Min is pinned
// to 1 so the value 10 sits above the baseline.
func TestGenerateGroupedBars_LogScale(t *testing.T) {
	p := groupedParams()
	p.Data = []map[string]any{
		{"idx": "A", "k1": 10.0, "k2": 1000.0},
		{"idx": "B", "k1": 100.0, "k2": 5.0},
	}
	p.ValueScale = scales.ScaleLogSpec{Base: 10, Min: scales.FloatVal(1), Max: scales.AutoFloat()}

	res := compute.GenerateGroupedBars(p)
	if len(res.Bars) != 4 {
		t.Fatalf("bars = %d, want 4", len(res.Bars))
	}
	for _, b := range res.Bars {
		if math.IsNaN(b.Y) || math.IsInf(b.Y, 0) || math.IsNaN(b.Height) || math.IsInf(b.Height, 0) {
			t.Errorf("bar %s non-finite: Y=%v H=%v", b.Key, b.Y, b.Height)
		}
		if b.Height < 0 {
			t.Errorf("bar %s height negative: %v", b.Key, b.Height)
		}
	}
	// Log is monotonic: 1000 must map to a taller bar than 10.
	small := barByKey(t, res.Bars, "k1.A") // 10
	large := barByKey(t, res.Bars, "k2.A") // 1000
	if large.Height <= small.Height {
		t.Errorf("log: value 1000 (h=%v) should be taller than 10 (h=%v)", large.Height, small.Height)
	}
}

// TestGenerateStackedBars_LogFiltersZeros drives the stacked path with a log
// value scale so filterZerosIfLog runs. A zero-valued key would break the log
// domain; the result must stay finite regardless.
func TestGenerateStackedBars_LogFiltersZeros(t *testing.T) {
	p := stackedParams()
	p.Data = []map[string]any{
		{"idx": "A", "k1": 0.0, "k2": 100.0},
		{"idx": "B", "k1": 50.0, "k2": 0.0},
	}
	p.ValueScale = scales.ScaleLogSpec{Base: 10, Min: scales.FloatVal(1), Max: scales.AutoFloat()}

	res := compute.GenerateStackedBars(p)
	for _, b := range res.Bars {
		if math.IsNaN(b.Y) || math.IsInf(b.Y, 0) || math.IsNaN(b.Height) || math.IsInf(b.Height, 0) {
			t.Errorf("bar %s non-finite: Y=%v H=%v", b.Key, b.Y, b.Height)
		}
	}
}
