package compute

import (
	"math"
	"testing"

	"github.com/geoffjay/templ-charts/charts/scales"
)

func TestFilterZerosIfLog(t *testing.T) {
	in := []float64{0, 1, 0, 5, -3}
	// Non-log scale types pass the slice through untouched.
	if got := filterZerosIfLog(in, "linear"); len(got) != len(in) {
		t.Errorf("linear: len = %d, want %d (no filtering)", len(got), len(in))
	}
	// Log drops zeros (a log domain can't include 0), keeping other values.
	got := filterZerosIfLog(in, "log")
	if len(got) != 3 {
		t.Fatalf("log: len = %d, want 3", len(got))
	}
	for _, v := range got {
		if v == 0 {
			t.Errorf("log: zero survived filtering: %v", got)
		}
	}
}

func TestLogScaleBase(t *testing.T) {
	if got := logScaleBase(scales.ScaleLogSpec{Base: 0}); got != 10 {
		t.Errorf("zero base defaults to 10, got %v", got)
	}
	if got := logScaleBase(scales.ScaleLogSpec{Base: 2}); got != 2 {
		t.Errorf("base 2, got %v", got)
	}
}

func TestValueBaseline(t *testing.T) {
	logScale := scales.ComputeScale(
		scales.ScaleLogSpec{Base: 10, Min: scales.FloatVal(1), Max: scales.AutoFloat()},
		scales.ComputedSerieAxis{All: []any{1.0, 1000.0}, Min: 1.0, Max: 1000.0},
		100, scales.ScaleAxisY,
	)
	// Log scale(0) is non-finite, so the baseline falls back to the axis origin:
	// the bottom (size) of a Y range, the left (0) of an X range.
	if got := valueBaseline(logScale, scales.ScaleAxisY, 100); got != 100 {
		t.Errorf("log Y baseline = %v, want 100", got)
	}
	if got := valueBaseline(logScale, scales.ScaleAxisX, 100); got != 0 {
		t.Errorf("log X baseline = %v, want 0", got)
	}

	// A linear scale whose domain includes 0 keeps the finite scale(0) baseline.
	linScale := scales.ComputeScale(
		scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat()},
		scales.ComputedSerieAxis{All: []any{0.0, 100.0}, Min: 0.0, Max: 100.0},
		100, scales.ScaleAxisY,
	)
	if got := valueBaseline(linScale, scales.ScaleAxisY, 100); math.IsNaN(got) || math.IsInf(got, 0) {
		t.Errorf("linear baseline non-finite: %v", got)
	}
}
