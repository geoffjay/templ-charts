package compute_test

import (
	"testing"

	"github.com/geoffjay/templ-charts/charts/bar/internal/compute"
	"github.com/geoffjay/templ-charts/charts/scales"
)

// stubScale is a hand-checkable scales.Scale: floats map to themselves
// (identity) and strings map through a fixed position table.
type stubScale struct{ pos map[string]float64 }

func (s stubScale) Type() scales.ScaleType { return scales.ScaleTypeLinear }
func (s stubScale) Call(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case string:
		return s.pos[x]
	}
	return 0
}

func totalsBars() []compute.ComputedBarDatum {
	mk := func(id, iv string, v any, w, h float64) compute.ComputedBarDatum {
		return compute.ComputedBarDatum{
			Key: id + "." + iv, Width: w, Height: h,
			Data: compute.ComputedDatum{ID: id, IndexValue: iv, Value: v},
		}
	}
	return []compute.ComputedBarDatum{
		mk("k1", "A", 10.0, 20, 30),
		mk("k1", "B", -5.0, 20, 30),
		mk("k2", "A", 20.0, 20, 30),
		mk("k2", "B", 15.0, 20, 30),
	}
}

func totalByKey(t *testing.T, totals []compute.BarTotalData, key string) compute.BarTotalData {
	t.Helper()
	for _, d := range totals {
		if d.Key == key {
			return d
		}
	}
	t.Fatalf("total %q not found (have %d)", key, len(totals))
	return compute.BarTotalData{}
}

var totalsIndexPos = map[string]float64{"A": 0, "B": 50}

func TestComputeBarTotals_StackedVertical(t *testing.T) {
	x := stubScale{pos: totalsIndexPos}
	y := stubScale{}
	got := compute.ComputeBarTotals(totalsBars(), x, y, "vertical", "stacked", 10, formatValue)
	if len(got) != 2 {
		t.Fatalf("totals = %d, want 2", len(got))
	}
	// Index A: total 30, positives 30. x = xScale("A") + barWidth/2 = 10;
	// y = yScale(30) - offset = 20; animOff = yScale(15) = 15.
	a := totalByKey(t, got, "total_A")
	approx(t, "A.X", a.X, 10)
	approx(t, "A.Y", a.Y, 20)
	approx(t, "A.Value", a.Value, 30)
	approx(t, "A.AnimationOffset", a.AnimationOffset, 15)
	if a.FormattedValue != "30" {
		t.Errorf("A.FormattedValue = %q, want 30", a.FormattedValue)
	}
	// Index B: total 10 (includes the -5), positives 15. y = 15 - 10 = 5.
	b := totalByKey(t, got, "total_B")
	approx(t, "B.X", b.X, 60)
	approx(t, "B.Y", b.Y, 5)
	approx(t, "B.Value", b.Value, 10)
	approx(t, "B.AnimationOffset", b.AnimationOffset, 7.5)
}

func TestComputeBarTotals_StackedHorizontal(t *testing.T) {
	x := stubScale{}
	y := stubScale{pos: totalsIndexPos}
	got := compute.ComputeBarTotals(totalsBars(), x, y, "horizontal", "stacked", 10, formatValue)
	// Index A: positives 30 → x = xScale(30) + offset = 40;
	// y = yScale("A") + barHeight/2 = 15.
	a := totalByKey(t, got, "total_A")
	approx(t, "A.X", a.X, 40)
	approx(t, "A.Y", a.Y, 15)
	approx(t, "A.AnimationOffset", a.AnimationOffset, 15)
	b := totalByKey(t, got, "total_B")
	approx(t, "B.X", b.X, 25) // xScale(15) + 10
	approx(t, "B.Y", b.Y, 65) // 50 + 30/2
	approx(t, "B.Value", b.Value, 10)
}

func TestComputeBarTotals_GroupedVertical(t *testing.T) {
	x := stubScale{pos: totalsIndexPos}
	y := stubScale{}
	got := compute.ComputeBarTotals(totalsBars(), x, y, "vertical", "grouped", 10, formatValue)
	// Index A: greatest 20, 2 bars. x = xScale("A") + 2*20/2 = 20;
	// y = yScale(20) - 10 = 10; animOff = yScale(10) = 10.
	a := totalByKey(t, got, "total_A")
	approx(t, "A.X", a.X, 20)
	approx(t, "A.Y", a.Y, 10)
	approx(t, "A.Value", a.Value, 30)
	approx(t, "A.AnimationOffset", a.AnimationOffset, 10)
	// Index B: greatest 15 (the -5 never becomes "greatest"), total 10.
	b := totalByKey(t, got, "total_B")
	approx(t, "B.X", b.X, 70) // 50 + 2*20/2
	approx(t, "B.Y", b.Y, 5)  // 15 - 10
	approx(t, "B.Value", b.Value, 10)
}

func TestComputeBarTotals_GroupedHorizontal(t *testing.T) {
	x := stubScale{}
	y := stubScale{pos: totalsIndexPos}
	got := compute.ComputeBarTotals(totalsBars(), x, y, "horizontal", "grouped", 10, formatValue)
	// Index A: greatest 20 → x = 20 + 10 = 30; y = yScale("A") + 2*30/2 = 30.
	a := totalByKey(t, got, "total_A")
	approx(t, "A.X", a.X, 30)
	approx(t, "A.Y", a.Y, 30)
	approx(t, "A.AnimationOffset", a.AnimationOffset, 10)
	b := totalByKey(t, got, "total_B")
	approx(t, "B.X", b.X, 25) // 15 + 10
	approx(t, "B.Y", b.Y, 80) // 50 + 30
}

func TestComputeBarTotals_EmptyBars(t *testing.T) {
	got := compute.ComputeBarTotals(nil, stubScale{}, stubScale{}, "vertical", "stacked", 10, formatValue)
	if got != nil {
		t.Errorf("empty bars should yield nil, got %v", got)
	}
}

func TestComputeBarTotals_NilAndNonNumericValues(t *testing.T) {
	bars := []compute.ComputedBarDatum{
		{Key: "k1.A", Width: 20, Height: 30, Data: compute.ComputedDatum{ID: "k1", IndexValue: "A", Value: nil}},
		{Key: "k2.A", Width: 20, Height: 30, Data: compute.ComputedDatum{ID: "k2", IndexValue: "A", Value: "n/a"}},
		{Key: "k3.A", Width: 20, Height: 30, Data: compute.ComputedDatum{ID: "k3", IndexValue: "A", Value: 5.0}},
	}
	got := compute.ComputeBarTotals(bars, stubScale{pos: totalsIndexPos}, stubScale{}, "vertical", "stacked", 0, formatValue)
	if len(got) != 1 {
		t.Fatalf("totals = %d, want 1", len(got))
	}
	approx(t, "Value", got[0].Value, 5) // nil and non-numeric count as 0
}

func TestComputeBarTotals_IntegrationWithGeneratedBars(t *testing.T) {
	// End-to-end: generate stacked bars and feed them plus their real scales
	// back into ComputeBarTotals, as charts/bar does.
	res := compute.GenerateStackedBars(stackedParams())
	got := compute.ComputeBarTotals(res.Bars, res.XScale, res.YScale, "vertical", "stacked", 10, formatValue)
	if len(got) != 2 {
		t.Fatalf("totals = %d, want 2", len(got))
	}
	// Both indexes total 40 → top of stack at yScale(40)=0 → label y = -10.
	for _, key := range []string{"total_A", "total_B"} {
		d := totalByKey(t, got, key)
		approx(t, key+".Value", d.Value, 40)
		approx(t, key+".Y", d.Y, -10)
	}
	// Label x is centered on the 50px band: A@0+25, B@50+25.
	approx(t, "A.X", totalByKey(t, got, "total_A").X, 25)
	approx(t, "B.X", totalByKey(t, got, "total_B").X, 75)
}
