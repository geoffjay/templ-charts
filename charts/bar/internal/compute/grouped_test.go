package compute_test

import (
	"math"
	"testing"

	"github.com/geoffjay/templ-charts/charts/bar/internal/compute"
	"github.com/geoffjay/templ-charts/charts/scales"
)

// groupedParams builds the base GroupedParams used across the grouped tests:
// two indexes (A, B), two keys (k1, k2), a 100x100 plot area, no padding.
// Value scale is linear [0, auto]. Hand-checked geometry: index band scale
// over 100 with 2 bands → A@0, B@50, bandwidth 50 → per-key bar size 25.
func groupedParams() compute.GroupedParams {
	return compute.GroupedParams{
		Data: []map[string]any{
			{"idx": "A", "k1": 10.0, "k2": 20.0},
			{"idx": "B", "k1": 30.0, "k2": 5.0},
		},
		Keys:   []string{"k1", "k2"},
		Layout: "vertical",
		Width:  100, Height: 100,
		Padding: 0, InnerPadding: 0,
		ValueScale:      scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat()},
		IndexScale:      scales.ScaleBandSpec{},
		GetIndex:        getIndex,
		GetColor:        getColor,
		GetTooltipLabel: getTooltipLabel,
		FormatValue:     formatValue,
		Margin:          compute.MarginLike{Top: 10, Right: 0, Bottom: 0, Left: 20},
	}
}

func TestGenerateGroupedBars_Vertical(t *testing.T) {
	res := compute.GenerateGroupedBars(groupedParams())
	if len(res.Bars) != 4 {
		t.Fatalf("bars = %d, want 4 (2 keys x 2 indexes)", len(res.Bars))
	}
	// Value domain [0, 30] over height 100 → yScale(v) = 100 - v*100/30.
	yOf := func(v float64) float64 { return 100 - v*100.0/30.0 }

	cases := []struct {
		key        string
		x, y, w, h float64
	}{
		{"k1.A", 0, yOf(10), 25, 100 - yOf(10)},
		{"k1.B", 50, yOf(30), 25, 100 - yOf(30)},
		{"k2.A", 25, yOf(20), 25, 100 - yOf(20)},
		{"k2.B", 75, yOf(5), 25, 100 - yOf(5)},
	}
	for _, c := range cases {
		b := barByKey(t, res.Bars, c.key)
		approx(t, c.key+".X", b.X, c.x)
		approx(t, c.key+".Y", b.Y, c.y)
		approx(t, c.key+".Width", b.Width, c.w)
		approx(t, c.key+".Height", b.Height, c.h)
		// Abs positions include the margin.
		approx(t, c.key+".AbsX", b.AbsX, 20+c.x)
		approx(t, c.key+".AbsY", b.AbsY, 10+c.y)
	}

	// Metadata comes from the accessors.
	b := barByKey(t, res.Bars, "k1.B")
	if b.Data.ID != "k1" || b.Data.IndexValue != "B" || b.Data.Value != 30.0 {
		t.Errorf("datum = %+v, want ID=k1 IndexValue=B Value=30", b.Data)
	}
	if b.Data.FormattedValue != "30" {
		t.Errorf("FormattedValue = %q, want 30", b.Data.FormattedValue)
	}
	if b.Color != "c-k1" || b.Label != "k1 - B" {
		t.Errorf("color/label = %q/%q, want c-k1 / k1 - B", b.Color, b.Label)
	}

	// Scales: vertical → x is the band scale, y the value scale.
	if res.XScale.Type() != scales.ScaleTypeBand || res.YScale.Type() != scales.ScaleTypeLinear {
		t.Errorf("scale types = %v/%v, want band/linear", res.XScale.Type(), res.YScale.Type())
	}
	approx(t, "YScale(0)", res.YScale.Call(0.0), 100)
	approx(t, "YScale(30)", res.YScale.Call(30.0), 0)
}

func TestGenerateGroupedBars_Horizontal(t *testing.T) {
	p := groupedParams()
	p.Layout = "horizontal"
	res := compute.GenerateGroupedBars(p)
	if len(res.Bars) != 4 {
		t.Fatalf("bars = %d, want 4", len(res.Bars))
	}
	// Horizontal: value axis is x with range [0, 100], domain [0, 30] →
	// xScale(v) = v*100/30. Index band scale over height 100 on the y axis
	// (reversed range) → A@50, B@0, band 50, per-key bar height 25.
	xOf := func(v float64) float64 { return v * 100.0 / 30.0 }
	cases := []struct {
		key        string
		x, y, w, h float64
	}{
		{"k1.A", 0, 50, xOf(10), 25},
		{"k1.B", 0, 0, xOf(30), 25},
		{"k2.A", 0, 75, xOf(20), 25},
		{"k2.B", 0, 25, xOf(5), 25},
	}
	for _, c := range cases {
		b := barByKey(t, res.Bars, c.key)
		approx(t, c.key+".X", b.X, c.x)
		approx(t, c.key+".Y", b.Y, c.y)
		approx(t, c.key+".Width", b.Width, c.w)
		approx(t, c.key+".Height", b.Height, c.h)
	}
	if res.XScale.Type() != scales.ScaleTypeLinear || res.YScale.Type() != scales.ScaleTypeBand {
		t.Errorf("scale types = %v/%v, want linear/band", res.XScale.Type(), res.YScale.Type())
	}
}

func TestGenerateGroupedBars_NegativeValues(t *testing.T) {
	p := groupedParams()
	p.Data = []map[string]any{
		{"idx": "A", "k1": -20.0},
		{"idx": "B", "k1": 20.0},
	}
	p.Keys = []string{"k1"}
	// Min auto → clamped-to-zero min = -20; max = 20. Domain [-20, 20] over
	// height 100 → yScale(v) = 50 - 2.5v; zero line at y=50.
	p.ValueScale = scales.ScaleLinearSpec{Min: scales.AutoFloat(), Max: scales.AutoFloat()}
	res := compute.GenerateGroupedBars(p)
	if len(res.Bars) != 2 {
		t.Fatalf("bars = %d, want 2", len(res.Bars))
	}
	neg := barByKey(t, res.Bars, "k1.A")
	approx(t, "neg.Y", neg.Y, 50) // negative bar hangs from the zero line
	approx(t, "neg.Height", neg.Height, 50)
	pos := barByKey(t, res.Bars, "k1.B")
	approx(t, "pos.Y", pos.Y, 0)
	approx(t, "pos.Height", pos.Height, 50)
}

func TestGenerateGroupedBars_NegativeValuesHorizontal(t *testing.T) {
	p := groupedParams()
	p.Layout = "horizontal"
	p.Data = []map[string]any{
		{"idx": "A", "k1": -20.0},
		{"idx": "B", "k1": 20.0},
	}
	p.Keys = []string{"k1"}
	p.ValueScale = scales.ScaleLinearSpec{Min: scales.AutoFloat(), Max: scales.AutoFloat()}
	// Domain [-20, 20] over width 100 → xScale(v) = 50 + 2.5v; zero at x=50.
	res := compute.GenerateGroupedBars(p)
	neg := barByKey(t, res.Bars, "k1.A")
	approx(t, "neg.X", neg.X, 0) // extends left of the zero line
	approx(t, "neg.Width", neg.Width, 50)
	pos := barByKey(t, res.Bars, "k1.B")
	approx(t, "pos.X", pos.X, 50)
	approx(t, "pos.Width", pos.Width, 50)
}

func TestGenerateGroupedBars_Reverse(t *testing.T) {
	p := groupedParams()
	p.Data = []map[string]any{{"idx": "A", "k1": 10.0}}
	p.Keys = []string{"k1"}
	p.ValueScale = scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat(), Reverse: true}
	res := compute.GenerateGroupedBars(p)
	b := barByKey(t, res.Bars, "k1.A")
	// Reversed vertical scale: domain [10, 0] → yScale(0)=0 (the reference),
	// yScale(10)=100 → bar grows downward from the top.
	approx(t, "Y", b.Y, 0)
	approx(t, "Height", b.Height, 100)
}

func TestGenerateGroupedBars_HiddenIDs(t *testing.T) {
	p := groupedParams()
	p.HiddenIDs = []string{"k2"}
	res := compute.GenerateGroupedBars(p)
	if len(res.Bars) != 2 {
		t.Fatalf("bars = %d, want 2 (k2 hidden)", len(res.Bars))
	}
	for _, b := range res.Bars {
		if b.Data.ID != "k1" {
			t.Errorf("unexpected bar for hidden key: %q", b.Key)
		}
	}
	// With only k1 visible the full bandwidth (50) goes to k1.
	approx(t, "Width", res.Bars[0].Width, 50)
}

func TestGenerateGroupedBars_NilValue(t *testing.T) {
	p := groupedParams()
	p.Data = []map[string]any{
		{"idx": "A", "k1": 10.0}, // k2 missing → normalized to nil
		{"idx": "B", "k1": 30.0, "k2": 5.0},
	}
	res := compute.GenerateGroupedBars(p)
	b := barByKey(t, res.Bars, "k2.A")
	if b.Data.Value != nil {
		t.Errorf("missing value should stay nil, got %v", b.Data.Value)
	}
	approx(t, "nil bar height", b.Height, 0)
	// The nil-valued key is filtered from the datum map handed to tooltips.
	if _, ok := b.Data.Data["k2"]; ok {
		t.Errorf("nil key should be filtered from Data.Data: %v", b.Data.Data)
	}
}

func TestGenerateGroupedBars_InnerPadding(t *testing.T) {
	p := groupedParams()
	p.InnerPadding = 4
	res := compute.GenerateGroupedBars(p)
	// bandwidth = (50 - 4*(2-1)) / 2 = 23.
	b1 := barByKey(t, res.Bars, "k1.A")
	b2 := barByKey(t, res.Bars, "k2.A")
	approx(t, "Width", b1.Width, 23)
	// Second key is offset by bar width + inner padding.
	approx(t, "k2.X - k1.X", b2.X-b1.X, 27)
}

func TestGenerateGroupedBars_NoBarsWhenBandwidthNonPositive(t *testing.T) {
	p := groupedParams()
	p.InnerPadding = 60 // (50 - 60)/2 < 0
	res := compute.GenerateGroupedBars(p)
	if len(res.Bars) != 0 {
		t.Fatalf("bars = %d, want 0 when computed bandwidth <= 0", len(res.Bars))
	}
	if res.XScale == nil || res.YScale == nil {
		t.Errorf("scales should still be returned")
	}
}

func TestGenerateGroupedBars_ReverseHorizontal(t *testing.T) {
	p := groupedParams()
	p.Layout = "horizontal"
	p.Data = []map[string]any{{"idx": "A", "k1": 10.0}}
	p.Keys = []string{"k1"}
	p.ValueScale = scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat(), Reverse: true}
	res := compute.GenerateGroupedBars(p)
	b := barByKey(t, res.Bars, "k1.A")
	// Reversed domain [10, 0] over range [0, 100]: xScale(10)=0, zero ref at
	// x=100 → the bar spans the full width starting at the left edge.
	approx(t, "X", b.X, 0)
	approx(t, "Width", b.Width, 100)
}

func TestGenerateGroupedBars_AutoMinClampsToZero(t *testing.T) {
	p := groupedParams()
	// All-positive data with Min auto: nivo clamps the domain min to 0 so
	// bars grow from the axis baseline, not from the smallest value.
	p.ValueScale = scales.ScaleLinearSpec{Min: scales.AutoFloat(), Max: scales.AutoFloat()}
	res := compute.GenerateGroupedBars(p)
	approx(t, "YScale(0)", res.YScale.Call(0.0), 100)
	approx(t, "YScale(30)", res.YScale.Call(30.0), 0)
	b := barByKey(t, res.Bars, "k1.B")
	approx(t, "k1.B.Y", b.Y, 0)
	approx(t, "k1.B.Height", b.Height, 100)
}

func TestGenerateGroupedBars_AllZeroValues(t *testing.T) {
	p := groupedParams()
	p.Data = []map[string]any{
		{"idx": "A", "k1": 0.0, "k2": 0.0},
		{"idx": "B", "k1": 0.0, "k2": 0.0},
	}
	res := compute.GenerateGroupedBars(p)
	if len(res.Bars) != 4 {
		t.Fatalf("bars = %d, want 4", len(res.Bars))
	}
	for _, b := range res.Bars {
		approx(t, b.Key+".Height", b.Height, 0)
	}
}

func TestGenerateGroupedBars_InfiniteValueClampedMax(t *testing.T) {
	p := groupedParams()
	p.Data = []map[string]any{{"idx": "A", "k1": math.Inf(1)}}
	p.Keys = []string{"k1"}
	res := compute.GenerateGroupedBars(p)
	// Mirrors nivo: a non-finite max collapses to 0 rather than producing an
	// unbounded scale. The raw value survives on the datum.
	if len(res.Bars) != 1 {
		t.Fatalf("bars = %d, want 1", len(res.Bars))
	}
	if v, ok := res.Bars[0].Data.Value.(float64); !ok || !math.IsInf(v, 1) {
		t.Errorf("raw value = %v, want +Inf", res.Bars[0].Data.Value)
	}
	approx(t, "YScale(0) with clamped max", res.YScale.Call(0.0), 100)
}
