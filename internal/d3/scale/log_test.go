package d3scale

import (
	"math"
	"testing"
)

// d3: scaleLog().domain([1, 1000]).range([0, 960])
//
//	(1) = 0, (10) = 320, (100) = 640, (1000) = 960
func TestLogCall(t *testing.T) {
	s := NewLog().SetDomain(1, 1000).SetRange(0, 960)
	cases := []struct {
		x, want float64
	}{
		{1, 0},
		{10, 320},
		{100, 640},
		{1000, 960},
	}
	for _, c := range cases {
		if got := s.Call(c.x); !approxEps(got, c.want, 1e-6) {
			t.Errorf("Call(%v) = %v want %v", c.x, got, c.want)
		}
	}
}

// d3: scaleLog().domain([1, 1000]).range([0, 960])(50) === 543.6704...
func TestLogCallMid(t *testing.T) {
	s := NewLog().SetDomain(1, 1000).SetRange(0, 960)
	got := s.Call(50)
	if !approxEps(got, 543.670401387526, 1e-6) {
		t.Errorf("Call(50) = %v want 543.6704", got)
	}
}

// d3: scaleLog().domain([1, 1000]).range([0, 960]).invert(480) === 31.6228...
func TestLogInvert(t *testing.T) {
	s := NewLog().SetDomain(1, 1000).SetRange(0, 960)
	got := s.Invert(480)
	if !approxEps(got, 31.62277660168379, 1e-6) {
		t.Errorf("Invert(480) = %v want 31.6228", got)
	}
}

// d3: log ticks(5) and ticks(10) over [1,1000] → the decade-multiples list
func TestLogTicks(t *testing.T) {
	s := NewLog().SetDomain(1, 1000).SetRange(0, 960)
	want := []float64{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		20, 30, 40, 50, 60, 70, 80, 90, 100,
		200, 300, 400, 500, 600, 700, 800, 900, 1000,
	}
	got := s.Ticks(5)
	if len(got) != len(want) {
		t.Fatalf("ticks(5) len = %d want %d (%v)", len(got), len(want), got)
	}
	for i, g := range got {
		if !approxEps(g, want[i], 1e-9) {
			t.Errorf("ticks(5)[%d] = %v want %v", i, g, want[i])
		}
	}
	// ticks(10) yields the same set for this domain
	got10 := s.Ticks(10)
	if len(got10) != len(want) {
		t.Fatalf("ticks(10) len = %d want %d", len(got10), len(want))
	}
}

// d3: scaleLog().base(2).domain([1, 32]).range([0, 100]).ticks() ===
//
//	[1, 2, 4, 8, 16, 32]
func TestLogBase2Ticks(t *testing.T) {
	s := NewLog().SetBase(2).SetDomain(1, 32).SetRange(0, 100)
	got := s.Ticks(5)
	want := []float64{1, 2, 4, 8, 16, 32}
	if len(got) != len(want) {
		t.Fatalf("base2 ticks len = %d want %d (%v)", len(got), len(want), got)
	}
	for i, g := range got {
		if !approxEps(g, want[i], 1e-9) {
			t.Errorf("base2 ticks[%d] = %v want %v", i, g, want[i])
		}
	}
}

// d3: scaleLog().base(2).domain([1, 32]).range([0, 100])(8) === 60
func TestLogBase2Call(t *testing.T) {
	s := NewLog().SetBase(2).SetDomain(1, 32).SetRange(0, 100)
	cases := []struct {
		x, want float64
	}{
		{1, 0}, {2, 20}, {4, 40}, {8, 60}, {16, 80}, {32, 100},
	}
	for _, c := range cases {
		if got := s.Call(c.x); !approxEps(got, c.want, 1e-6) {
			t.Errorf("base2 Call(%v) = %v want %v", c.x, got, c.want)
		}
	}
}

// d3: scaleLog().domain([1, 1000]).nice().domain() === [1, 1000] (already nice)
func TestLogNice(t *testing.T) {
	s := NewLog().SetDomain(1, 1000)
	s.Nice(0)
	d := s.Domain()
	if !approxEps(d[0], 1, 1e-9) || !approxEps(d[1], 1000, 1e-9) {
		t.Errorf("nice [1,1000] = %v want [1, 1000]", d)
	}
	// A non-nice domain should be extended outward.
	s2 := NewLog().SetDomain(3, 997)
	s2.Nice(0)
	d2 := s2.Domain()
	if !approxEps(d2[0], 1, 1e-9) || !approxEps(d2[1], 1000, 1e-9) {
		t.Errorf("nice [3,997] = %v want [1, 1000]", d2)
	}
}

// d3: negative domain reflected
//
//	scaleLog().domain([-1000, -1]).range([0, 100])
//	(-1000) = 0, (-1) = 100, (-100) = 33.333...
func TestLogNegativeDomain(t *testing.T) {
	s := NewLog().SetDomain(-1000, -1).SetRange(0, 100)
	if got := s.Call(-1000); !approxEps(got, 0, 1e-6) {
		t.Errorf("neg Call(-1000) = %v want 0", got)
	}
	if got := s.Call(-1); !approxEps(got, 100, 1e-6) {
		t.Errorf("neg Call(-1) = %v want 100", got)
	}
	if got := s.Call(-100); !approxEps(got, 33.33333333333333, 1e-6) {
		t.Errorf("neg Call(-100) = %v want 33.3333", got)
	}
}

func TestLogClamp(t *testing.T) {
	s := NewLog().SetDomain(1, 1000).SetRange(0, 960).SetClamp(true)
	if got := s.Call(0.5); got != 0 {
		t.Errorf("clamped Call(0.5) = %v want 0", got)
	}
	if got := s.Call(2000); got != 960 {
		t.Errorf("clamped Call(2000) = %v want 960", got)
	}
}

func TestLogType(t *testing.T) {
	if NewLog().Type() != "log" {
		t.Error("Type() != log")
	}
}

func TestLogCopy(t *testing.T) {
	s := NewLog().SetDomain(1, 100).SetRange(0, 100).SetBase(2)
	c := s.Copy()
	c.SetDomain(1, 1000).SetBase(10)
	if s.Base() != 2 {
		t.Errorf("original Base after copy mutate = %v want 2", s.Base())
	}
	if !approxEps(s.Domain()[1], 100, 1e-9) {
		t.Errorf("original Domain after copy mutate = %v want [1,100]", s.Domain())
	}
}

func TestLogZeroDomain(t *testing.T) {
	// d3 throws; we just produce NaN rather than panic, since the d3 transform
	// of 0 is -Inf. This documents the behavior.
	s := NewLog().SetDomain(0, 100).SetRange(0, 100)
	if !math.IsNaN(s.Call(5)) {
		t.Errorf("log scale with 0 domain: Call(5) = %v want NaN (degenerate)", s.Call(5))
	}
}
