package compute_test

import (
	"math"
	"strconv"
	"testing"

	"github.com/geoffjay/templ-charts/charts/bar/compute"
	"github.com/geoffjay/templ-charts/charts/scales"
)

// --- shared test helpers -----------------------------------------------------

func getIndex(d map[string]any) string {
	if s, ok := d["idx"].(string); ok {
		return s
	}
	return ""
}

func getColor(d compute.ComputedDatum) string { return "c-" + d.ID }

func getTooltipLabel(d compute.ComputedDatum) string { return d.ID + " - " + d.IndexValue }

func formatValue(v float64) string { return strconv.FormatFloat(v, 'f', -1, 64) }

func approx(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

// barByKey finds a computed bar by its Key ("<key>.<indexValue>").
func barByKey(t *testing.T, bars []compute.ComputedBarDatum, key string) compute.ComputedBarDatum {
	t.Helper()
	for _, b := range bars {
		if b.Key == key {
			return b
		}
	}
	t.Fatalf("bar %q not found (have %d bars)", key, len(bars))
	return compute.ComputedBarDatum{}
}

// --- NormalizeData -----------------------------------------------------------

func TestNormalizeData(t *testing.T) {
	data := []map[string]any{
		{"idx": "A", "k1": 10.0},
		{"idx": "B", "k2": 20.0},
	}
	out := compute.NormalizeData(data, []string{"k1", "k2"})
	if len(out) != 2 {
		t.Fatalf("len = %d, want 2", len(out))
	}
	if out[0]["k1"] != 10.0 || out[0]["k2"] != nil {
		t.Errorf("datum 0 = %v, want k1=10 k2=nil", out[0])
	}
	if out[1]["k1"] != nil || out[1]["k2"] != 20.0 {
		t.Errorf("datum 1 = %v, want k1=nil k2=20", out[1])
	}
	// Non-key fields are preserved.
	if out[0]["idx"] != "A" {
		t.Errorf("idx not preserved: %v", out[0])
	}
	// Input is not mutated.
	if _, ok := data[0]["k2"]; ok {
		t.Errorf("NormalizeData mutated its input")
	}
}

// --- FilterNullValues --------------------------------------------------------

func TestFilterNullValues(t *testing.T) {
	in := map[string]any{
		"nil":    nil,
		"empty":  "",
		"zero":   0.0,
		"keep-s": "x",
		"keep-f": 1.5,
		"keep-i": 0, // int zero is NOT filtered (only float64 zero is)
	}
	out := compute.FilterNullValues(in)
	for _, k := range []string{"nil", "empty", "zero"} {
		if _, ok := out[k]; ok {
			t.Errorf("key %q should have been filtered", k)
		}
	}
	if out["keep-s"] != "x" || out["keep-f"] != 1.5 || out["keep-i"] != 0 {
		t.Errorf("kept values wrong: %v", out)
	}
}

// --- CoerceValue -------------------------------------------------------------

func TestCoerceValue(t *testing.T) {
	cases := []struct {
		name    string
		in      any
		wantRaw any
		wantF   float64
	}{
		{"nil", nil, nil, 0},
		{"float64", 3.5, 3.5, 3.5},
		{"float32", float32(2), float32(2), 2},
		{"int", 7, 7, 7},
		{"int32", int32(4), int32(4), 4},
		{"int64", int64(9), int64(9), 9},
		{"bool-true", true, true, 1},
		{"bool-false", false, false, 0},
		{"numeric-string", "12.5", "12.5", 12.5},
		{"non-numeric-string", "abc", "abc", 0},
		{"struct", struct{}{}, struct{}{}, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			raw, f := compute.CoerceValue(c.in)
			if raw != c.wantRaw {
				t.Errorf("raw = %v, want %v", raw, c.wantRaw)
			}
			if f != c.wantF {
				t.Errorf("f = %v, want %v", f, c.wantF)
			}
		})
	}
}

// --- GetIndexScale -----------------------------------------------------------

func TestGetIndexScale_XAxis(t *testing.T) {
	data := []map[string]any{{"idx": "A"}, {"idx": "B"}}
	s := compute.GetIndexScale(data, getIndex, 0, scales.ScaleBandSpec{}, 100, scales.ScaleAxisX)
	if s.Type() != scales.ScaleTypeBand {
		t.Fatalf("type = %v, want band", s.Type())
	}
	approx(t, `Call("A")`, s.Call("A"), 0)
	approx(t, `Call("B")`, s.Call("B"), 50)
	bw := s.(scales.ScaleWithBandwidth)
	approx(t, "Bandwidth", bw.Bandwidth(), 50)
	approx(t, "Step", bw.Step(), 50)
}

func TestGetIndexScale_YAxisReversesRange(t *testing.T) {
	data := []map[string]any{{"idx": "A"}, {"idx": "B"}}
	s := compute.GetIndexScale(data, getIndex, 0, scales.ScaleBandSpec{}, 100, scales.ScaleAxisY)
	// Y axis uses range [size, 0]: first domain value sits at the top band.
	approx(t, `Call("A")`, s.Call("A"), 50)
	approx(t, `Call("B")`, s.Call("B"), 0)
}

func TestGetIndexScale_PaddingShrinksBandwidth(t *testing.T) {
	data := []map[string]any{{"idx": "A"}, {"idx": "B"}}
	none := compute.GetIndexScale(data, getIndex, 0, scales.ScaleBandSpec{}, 100, scales.ScaleAxisX).(scales.ScaleWithBandwidth)
	padded := compute.GetIndexScale(data, getIndex, 0.2, scales.ScaleBandSpec{}, 100, scales.ScaleAxisX).(scales.ScaleWithBandwidth)
	if padded.Bandwidth() >= none.Bandwidth() {
		t.Errorf("padding should shrink bandwidth: %v >= %v", padded.Bandwidth(), none.Bandwidth())
	}
	// d3 band math: step = 100 / (2 - 0.2 + 2*0.2) = 100/2.2; bw = step*0.8.
	approx(t, "padded bandwidth", padded.Bandwidth(), 100/2.2*0.8)
}

func TestGetIndexScale_NonStringDomainValues(t *testing.T) {
	// Indexes coming from numeric fields exercise the toString coercion in
	// bandScale.Call.
	data := []map[string]any{{"idx": "1"}, {"idx": "2.5"}}
	s := compute.GetIndexScale(data, getIndex, 0, scales.ScaleBandSpec{}, 100, scales.ScaleAxisX)
	approx(t, "Call(float64 1)", s.Call(float64(1)), 0) // "1"
	approx(t, "Call(float64 2.5)", s.Call(2.5), 50)     // "2.5"
	approx(t, "Call(int 1)", s.Call(1), 0)              // "1"
	approx(t, "Call(int64 1)", s.Call(int64(1)), 0)     // "1"
	if v := s.Call(nil); !math.IsNaN(v) {               // "" is not in the domain
		t.Errorf("Call(nil) = %v, want NaN", v)
	}
	if v := s.Call(struct{}{}); !math.IsNaN(v) {
		t.Errorf("Call(struct{}) = %v, want NaN", v)
	}
}

// --- ComputeLabelLayout ------------------------------------------------------

func TestComputeLabelLayout(t *testing.T) {
	const w, h, off = 40.0, 30.0, 5.0
	cases := []struct {
		name       string
		layout     string
		reverse    bool
		position   string
		wantX      float64
		wantY      float64
		wantAnchor string
	}{
		{"vertical-middle", "vertical", false, "middle", w / 2, h/2 - off, "middle"},
		{"vertical-start", "vertical", false, "start", w / 2, h - off, "middle"},
		{"vertical-end", "vertical", false, "end", w / 2, 0 - off, "middle"},
		{"vertical-middle-reverse", "vertical", true, "middle", w / 2, h/2 + off, "middle"},
		{"vertical-start-reverse", "vertical", true, "start", w / 2, 0 + off, "middle"},
		{"vertical-end-reverse", "vertical", true, "end", w / 2, h + off, "middle"},
		{"horizontal-middle", "horizontal", false, "middle", w/2 + off, h / 2, "middle"},
		{"horizontal-start", "horizontal", false, "start", 0 + off, h / 2, "start"},
		{"horizontal-end", "horizontal", false, "end", w + off, h / 2, "start"},
		{"horizontal-middle-reverse", "horizontal", true, "middle", w/2 - off, h / 2, "middle"},
		{"horizontal-start-reverse", "horizontal", true, "start", w - off, h / 2, "end"},
		{"horizontal-end-reverse", "horizontal", true, "end", 0 - off, h / 2, "end"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := compute.ComputeLabelLayout(c.layout, c.reverse, c.position, off)
			got := f(w, h)
			approx(t, "LabelX", got.LabelX, c.wantX)
			approx(t, "LabelY", got.LabelY, c.wantY)
			if got.TextAnchor != c.wantAnchor {
				t.Errorf("TextAnchor = %q, want %q", got.TextAnchor, c.wantAnchor)
			}
		})
	}
}

func TestGetIndexScale_Round(t *testing.T) {
	// Three bands over 100px do not divide evenly; Round snaps band starts
	// and widths to integers.
	data := []map[string]any{{"idx": "A"}, {"idx": "B"}, {"idx": "C"}}
	s := compute.GetIndexScale(data, getIndex, 0, scales.ScaleBandSpec{Round: true}, 100, scales.ScaleAxisX)
	r, ok := s.(interface{ Round() bool })
	if !ok {
		t.Fatalf("index scale should expose Round()")
	}
	if !r.Round() {
		t.Errorf("Round() = false, want true")
	}
	bw := s.(scales.ScaleWithBandwidth).Bandwidth()
	if bw != math.Trunc(bw) {
		t.Errorf("rounded bandwidth = %v, want an integer", bw)
	}
	for _, idx := range []string{"A", "B", "C"} {
		v := s.Call(idx)
		if v != math.Trunc(v) {
			t.Errorf("Call(%s) = %v, want an integer", idx, v)
		}
	}
}
