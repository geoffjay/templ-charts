package d3scale

import (
	"math"
	"testing"
)

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

// approxEps allows a custom tolerance for tests where d3's float behavior
// introduces tiny drift.
func approxEps(a, b, eps float64) bool { return math.Abs(a-b) < eps }

func TestLinearDefaults(t *testing.T) {
	s := NewLinear()
	if !approx(s.Call(0), 0) {
		t.Errorf("default Call(0) = %v want 0", s.Call(0))
	}
	if !approx(s.Call(1), 1) {
		t.Errorf("default Call(1) = %v want 1", s.Call(1))
	}
	if !approx(s.Call(0.5), 0.5) {
		t.Errorf("default Call(0.5) = %v want 0.5", s.Call(0.5))
	}
}

// d3: scaleLinear().domain([10, 130]).range([0, 960])(50) === 320
func TestLinearD3DocExample(t *testing.T) {
	s := NewLinear().SetDomain(10, 130).SetRange(0, 960)
	got := s.Call(50)
	if !approx(got, 320) {
		t.Errorf("Call(50) = %v want 320 (d3 doc example)", got)
	}
}

// d3: scaleLinear().domain([10, 130]).range([0, 960]).invert(320) === 50
func TestLinearInvert(t *testing.T) {
	s := NewLinear().SetDomain(10, 130).SetRange(0, 960)
	got := s.Invert(320)
	if !approx(got, 50) {
		t.Errorf("Invert(320) = %v want 50", got)
	}
}

// d3: scaleLinear().domain([10, 130]).range([0, 960]).clamp(true)(-100) === 0
func TestLinearClamp(t *testing.T) {
	s := NewLinear().SetDomain(10, 130).SetRange(0, 960).SetClamp(true)
	if !approx(s.Call(-100), 0) {
		t.Errorf("clamped Call(-100) = %v want 0", s.Call(-100))
	}
	if !approx(s.Call(500), 960) {
		t.Errorf("clamped Call(500) = %v want 960", s.Call(500))
	}
}

// d3: scaleLinear().domain([0, 1]).rangeRound([0, 960])(0.499) === 479
// interpolateRound = Math.round(0 + 960*0.499) = Math.round(479.04) = 479
func TestLinearRangeRound(t *testing.T) {
	s := NewLinear().SetDomain(0, 1).SetRangeRound(0, 960)
	got := s.Call(0.499)
	if got != 479 { // exact integer comparison — rounding produces 479
		t.Errorf("rangeRound Call(0.499) = %v want 479", got)
	}
}

// d3: scaleLinear().domain([0.201, 0.979]).nice(10).domain() === [0.2, 1.0]
// (d3 extends outward to nearest step)
func TestLinearNice(t *testing.T) {
	s := NewLinear().SetDomain(0.201, 0.979)
	s.Nice(10)
	d := s.Domain()
	if !approxEps(d[0], 0.2, 1e-9) || !approxEps(d[1], 1.0, 1e-9) {
		t.Errorf("nice(10) domain = %v want [0.2, 1.0]", d)
	}
}

// d3: scaleLinear().domain([1, 10]).range([0, 100])(5) === 44.444...
// (linear interpolation: (5-1)/(10-1) * 100 = 400/9 ≈ 44.4444)
func TestLinearMidpoint(t *testing.T) {
	s := NewLinear().SetDomain(1, 10).SetRange(0, 100)
	got := s.Call(5)
	if !approxEps(got, 400.0/9.0, 1e-9) {
		t.Errorf("Call(5) = %v want 44.4444", got)
	}
}

// d3: scaleLinear().domain([0, 100]).range([0, 960]).ticks(5) === [0, 20, 40, 60, 80, 100]
func TestLinearTicks(t *testing.T) {
	s := NewLinear().SetDomain(0, 100).SetRange(0, 960)
	got := s.Ticks(5)
	want := []float64{0, 20, 40, 60, 80, 100}
	if len(got) != len(want) {
		t.Fatalf("Ticks(5) len = %d want %d (%v)", len(got), len(want), got)
	}
	for i, g := range got {
		if !approxEps(g, want[i], 1e-9) {
			t.Errorf("Ticks(5)[%d] = %v want %v", i, g, want[i])
		}
	}
}

// d3: scaleLinear().domain([0, 10]).range([0, 960]).ticks(10) === [0,1,2,...,10]
func TestLinearTicksCount10(t *testing.T) {
	s := NewLinear().SetDomain(0, 10).SetRange(0, 960)
	got := s.Ticks(10)
	want := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	if len(got) != len(want) {
		t.Fatalf("Ticks(10) len = %d want %d (%v)", len(got), len(want), got)
	}
	for i, g := range got {
		if !approxEps(g, want[i], 1e-9) {
			t.Errorf("Ticks(10)[%d] = %v want %v", i, g, want[i])
		}
	}
}

// d3: scaleLinear().domain([130, 10]).range([0, 960])(50) === 640 (reversed)
func TestLinearReversedDomain(t *testing.T) {
	s := NewLinear().SetDomain(130, 10).SetRange(0, 960)
	got := s.Call(50)
	if !approx(got, 640) {
		t.Errorf("reversed Call(50) = %v want 640", got)
	}
}

// Copy is independent of the original.
func TestLinearCopy(t *testing.T) {
	s := NewLinear().SetDomain(0, 10).SetRange(0, 100)
	c := s.Copy()
	c.SetDomain(0, 100)
	if !approx(s.Call(5), 50) {
		t.Errorf("original after copy-mutate Call(5) = %v want 50", s.Call(5))
	}
	if !approx(c.Call(50), 50) {
		t.Errorf("copy Call(50) = %v want 50", c.Call(50))
	}
}

func TestLinearType(t *testing.T) {
	if NewLinear().Type() != "linear" {
		t.Errorf("Type() != linear")
	}
}
