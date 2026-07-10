package compute_test

import (
	"testing"

	"github.com/geoffjay/templ-charts/charts/bar/internal/compute"
	"github.com/geoffjay/templ-charts/charts/scales"
)

// stackedParams: two indexes (A, B) and two keys stacking to 40 at each index.
// Value domain [0, 40] over a 100px value axis → 2.5 px per unit.
// Index band scale over 100 with 2 bands → bandwidth 50.
func stackedParams() compute.StackedParams {
	return compute.StackedParams{
		Data: []map[string]any{
			{"idx": "A", "k1": 10.0, "k2": 30.0},
			{"idx": "B", "k1": 20.0, "k2": 20.0},
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

func TestGenerateStackedBars_Vertical(t *testing.T) {
	res := compute.GenerateStackedBars(stackedParams())
	if len(res.Bars) != 4 {
		t.Fatalf("bars = %d, want 4", len(res.Bars))
	}
	// yScale(v) = 100 - 2.5v. Stack (diverging, all positive):
	//   A: k1 [0,10], k2 [10,40];  B: k1 [0,20], k2 [20,40].
	cases := []struct {
		key        string
		x, y, w, h float64
	}{
		{"k1.A", 0, 75, 50, 25},  // y=yScale(10)=75, h=yScale(0)-75=25
		{"k1.B", 50, 50, 50, 50}, // y=yScale(20)=50, h=50
		{"k2.A", 0, 0, 50, 75},   // y=yScale(40)=0, h=yScale(10)-0=75
		{"k2.B", 50, 0, 50, 50},  // y=yScale(40)=0, h=yScale(20)-0=50
	}
	for _, c := range cases {
		b := barByKey(t, res.Bars, c.key)
		approx(t, c.key+".X", b.X, c.x)
		approx(t, c.key+".Y", b.Y, c.y)
		approx(t, c.key+".Width", b.Width, c.w)
		approx(t, c.key+".Height", b.Height, c.h)
		approx(t, c.key+".AbsX", b.AbsX, 20+c.x)
		approx(t, c.key+".AbsY", b.AbsY, 10+c.y)
	}
	// The stacked segments of an index tile exactly: k1.A bottom edge = 100,
	// k2.A bottom edge = k1.A top edge.
	k1a := barByKey(t, res.Bars, "k1.A")
	k2a := barByKey(t, res.Bars, "k2.A")
	approx(t, "k2.A bottom == k1.A top", k2a.Y+k2a.Height, k1a.Y)
	approx(t, "k1.A bottom == baseline", k1a.Y+k1a.Height, 100)

	// Datum metadata.
	if k1a.Data.ID != "k1" || k1a.Data.IndexValue != "A" || k1a.Data.Value != 10.0 {
		t.Errorf("datum = %+v, want ID=k1 IndexValue=A Value=10", k1a.Data)
	}
	if k1a.Color != "c-k1" || k1a.Label != "k1 - A" || k1a.Data.FormattedValue != "10" {
		t.Errorf("color/label/fv = %q/%q/%q", k1a.Color, k1a.Label, k1a.Data.FormattedValue)
	}
	if res.XScale.Type() != scales.ScaleTypeBand || res.YScale.Type() != scales.ScaleTypeLinear {
		t.Errorf("scale types = %v/%v, want band/linear", res.XScale.Type(), res.YScale.Type())
	}
}

func TestGenerateStackedBars_Horizontal(t *testing.T) {
	p := stackedParams()
	p.Layout = "horizontal"
	res := compute.GenerateStackedBars(p)
	if len(res.Bars) != 4 {
		t.Fatalf("bars = %d, want 4", len(res.Bars))
	}
	// xScale(v) = 2.5v; index band on y (reversed range): A@50, B@0, band 50.
	cases := []struct {
		key        string
		x, y, w, h float64
	}{
		{"k1.A", 0, 50, 25, 50},  // x=xScale(0)=0, w=xScale(10)-0=25
		{"k1.B", 0, 0, 50, 50},   // w=xScale(20)=50
		{"k2.A", 25, 50, 75, 50}, // x=xScale(10)=25, w=xScale(40)-25=75
		{"k2.B", 50, 0, 50, 50},  // x=xScale(20)=50, w=xScale(40)-50=50
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

func TestGenerateStackedBars_InnerPadding(t *testing.T) {
	p := stackedParams()
	p.InnerPadding = 4
	res := compute.GenerateStackedBars(p)
	// Vertical: y shifts down by innerPadding/2 and the height shrinks by
	// 1.5*innerPadding (half from the shifted top, a full pad off the bottom).
	b := barByKey(t, res.Bars, "k1.A")
	approx(t, "Y", b.Y, 75+2)
	approx(t, "Height", b.Height, 25-6)
}

func TestGenerateStackedBars_InnerPaddingHorizontal(t *testing.T) {
	p := stackedParams()
	p.Layout = "horizontal"
	p.InnerPadding = 4
	res := compute.GenerateStackedBars(p)
	b := barByKey(t, res.Bars, "k2.A")
	// Unpadded: x=25, w=75 → padded: x=25+2=27, w=75-1.5*4=69.
	approx(t, "X", b.X, 27)
	approx(t, "Width", b.Width, 69)
}

func TestGenerateStackedBars_Reverse(t *testing.T) {
	p := stackedParams()
	p.Data = []map[string]any{
		{"idx": "A", "k1": 10.0, "k2": 30.0},
	}
	p.ValueScale = scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat(), Reverse: true}
	res := compute.GenerateStackedBars(p)
	// Reversed domain [40, 0] over range [100, 0] → yScale(v) = 2.5v.
	// k1 [0,10]: y = yScale(lo) = 0, h = yScale(hi) - y = 25.
	b := barByKey(t, res.Bars, "k1.A")
	approx(t, "Y", b.Y, 0)
	approx(t, "Height", b.Height, 25)
	// k2 [10,40]: y = yScale(10) = 25, h = yScale(40) - 25 = 75.
	b2 := barByKey(t, res.Bars, "k2.A")
	approx(t, "k2.Y", b2.Y, 25)
	approx(t, "k2.Height", b2.Height, 75)
}

func TestGenerateStackedBars_ReverseHorizontal(t *testing.T) {
	p := stackedParams()
	p.Layout = "horizontal"
	p.Data = []map[string]any{
		{"idx": "A", "k1": 10.0, "k2": 30.0},
	}
	p.ValueScale = scales.ScaleLinearSpec{Min: scales.FloatVal(0), Max: scales.AutoFloat(), Reverse: true}
	res := compute.GenerateStackedBars(p)
	// Reversed domain [40, 0] over range [0, 100] → xScale(v) = 100 - 2.5v.
	// k1 [0,10]: x = xScale(hi) = 75, w = xScale(lo) - x = 25.
	b := barByKey(t, res.Bars, "k1.A")
	approx(t, "X", b.X, 75)
	approx(t, "Width", b.Width, 25)
}

func TestGenerateStackedBars_DivergingNegative(t *testing.T) {
	p := stackedParams()
	p.Data = []map[string]any{
		{"idx": "A", "k1": -20.0, "k2": 20.0},
	}
	p.ValueScale = scales.ScaleLinearSpec{Min: scales.AutoFloat(), Max: scales.AutoFloat()}
	res := compute.GenerateStackedBars(p)
	// Diverging stack: k1 [-20,0] below zero, k2 [0,20] above. Domain
	// [-20,20] over [100,0] → yScale(v) = 50 - 2.5v; zero line at 50.
	neg := barByKey(t, res.Bars, "k1.A")
	pos := barByKey(t, res.Bars, "k2.A")
	approx(t, "neg.Y", neg.Y, 50) // starts at the zero line, extends down
	approx(t, "neg.Height", neg.Height, 50)
	approx(t, "pos.Y", pos.Y, 0)
	approx(t, "pos.Height", pos.Height, 50)
	if neg.Data.Value != -20.0 || pos.Data.Value != 20.0 {
		t.Errorf("raw values = %v/%v, want -20/20", neg.Data.Value, pos.Data.Value)
	}
}

func TestGenerateStackedBars_HiddenIDs(t *testing.T) {
	p := stackedParams()
	p.HiddenIDs = []string{"k1"}
	res := compute.GenerateStackedBars(p)
	if len(res.Bars) != 2 {
		t.Fatalf("bars = %d, want 2 (k1 hidden)", len(res.Bars))
	}
	// With k1 hidden, k2 stacks from zero: A k2 [0,30], domain [0,30].
	b := barByKey(t, res.Bars, "k2.A")
	approx(t, "Y", b.Y, 0) // yScale(30) with domain [0,30] = 0
	approx(t, "Height", b.Height, 100)
}

func TestGenerateStackedBars_NilValueTreatedAsZero(t *testing.T) {
	p := stackedParams()
	p.Data = []map[string]any{
		{"idx": "A", "k1": 10.0}, // k2 missing
		{"idx": "B", "k1": 20.0, "k2": 20.0},
	}
	res := compute.GenerateStackedBars(p)
	b := barByKey(t, res.Bars, "k2.A")
	if b.Data.Value != nil {
		t.Errorf("missing value should stay nil, got %v", b.Data.Value)
	}
	approx(t, "zero segment height", b.Height, 0)
}
