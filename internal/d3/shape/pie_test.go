package d3shape

import (
	"math"
	"testing"
)

func approxPie(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestPieDefault(t *testing.T) {
	data := []float64{1, 2, 3, 4}
	r := NewPieNumeric().Call(data)
	if len(r) != 4 {
		t.Fatalf("len = %d want 4", len(r))
	}
	// d3 default: sortValues descending → value 4 gets index 0, value 1 gets index 3
	want := []struct {
		index      int
		value      float64
		start, end float64
	}{
		{3, 1, 5.654866776461628, 6.283185307179586},
		{2, 2, 4.39822971502571, 5.654866776461628},
		{1, 3, 2.5132741228718345, 4.39822971502571},
		{0, 4, 0, 2.5132741228718345},
	}
	for i, w := range want {
		if r[i].Index != w.index || !approxPie(r[i].Value, w.value) ||
			!approxPie(r[i].StartAngle, w.start) || !approxPie(r[i].EndAngle, w.end) {
			t.Errorf("arc[%d] = {idx:%d val:%g start:%g end:%g} want {idx:%d val:%g start:%g end:%g}",
				i, r[i].Index, r[i].Value, r[i].StartAngle, r[i].EndAngle,
				w.index, w.value, w.start, w.end)
		}
	}
}

func TestPieSortValuesNull(t *testing.T) {
	data := []float64{1, 2, 3, 4}
	r := NewPieNumeric().SortValues(nil).Call(data)
	// Without sorting, arcs stay in data order with index = position
	for i, a := range r {
		if a.Index != i {
			t.Errorf("arc[%d].Index = %d want %d", i, a.Index, i)
		}
		if !approxPie(a.Value, data[i]) {
			t.Errorf("arc[%d].Value = %g want %g", i, a.Value, data[i])
		}
	}
	// First arc starts at 0
	if !approxPie(r[0].StartAngle, 0) {
		t.Errorf("first arc start = %g want 0", r[0].StartAngle)
	}
}

func TestPiePadAngle(t *testing.T) {
	data := []float64{1, 2, 3, 4}
	r := NewPieNumeric().PadAngle(0.1).Call(data)
	for _, a := range r {
		if !approxPie(a.PadAngle, 0.1) {
			t.Errorf("arc padAngle = %g want 0.1", a.PadAngle)
		}
	}
}

func TestPieStartEndAngle(t *testing.T) {
	data := []float64{1, 2, 3, 4}
	r := NewPieNumeric().StartAngle(math.Pi / 4).EndAngle(math.Pi).Call(data)
	// First arc (by index 0 = value 4) should start at π/4
	for _, a := range r {
		if a.Index == 0 {
			if !approxPie(a.StartAngle, math.Pi/4) {
				t.Errorf("first arc start = %g want π/4", a.StartAngle)
			}
		}
		if a.Index == 3 {
			if !approxPie(a.EndAngle, math.Pi) {
				t.Errorf("last arc end = %g want π", a.EndAngle)
			}
		}
	}
}

func TestPieCustomValue(t *testing.T) {
	type datum struct{ N float64 }
	data := []datum{{10}, {20}, {30}}
	p := NewPie[datum]().Value(func(d datum, _ int, _ []datum) float64 { return d.N })
	r := p.Call(data)
	// Descending: 30 gets index 0, 10 gets index 2
	if r[2].Index != 0 || !approxPie(r[2].Value, 30) {
		t.Errorf("expected value 30 at index 0, got idx=%d val=%g", r[2].Index, r[2].Value)
	}
	if r[0].Index != 2 || !approxPie(r[0].Value, 10) {
		t.Errorf("expected value 10 at index 2, got idx=%d val=%g", r[0].Index, r[0].Value)
	}
}

func TestPieEmpty(t *testing.T) {
	r := NewPieNumeric().Call([]float64{})
	if r != nil {
		t.Errorf("empty pie should return nil, got %v", r)
	}
}

func TestPieFullCircle(t *testing.T) {
	data := []float64{1, 1, 1}
	r := NewPieNumeric().SortValues(nil).Call(data)
	// Three equal slices: each gets 2π/3
	last := r[len(r)-1]
	if !approxPie(last.EndAngle, tau) {
		t.Errorf("last arc end = %g want 2π (%g)", last.EndAngle, tau)
	}
}
