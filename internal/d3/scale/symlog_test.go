package d3scale

import "testing"

// d3: scaleSymlog().domain([0, 100]).range([0, 960])
//
//	(0) = 0, (100) = 960, (50) = 817.8665..., (10) = 498.7908..., (1) = 144.1829...
func TestSymlogCall(t *testing.T) {
	s := NewSymlog().SetDomain(0, 100).SetRange(0, 960)
	cases := []struct {
		x, want float64
	}{
		{0, 0},
		{100, 960},
		{50, 817.8665310345526},
		{10, 498.7907582231431},
		{1, 144.18286389474042},
	}
	for _, c := range cases {
		if got := s.Call(c.x); !approxEps(got, c.want, 1e-6) {
			t.Errorf("Call(%v) = %v want %v", c.x, got, c.want)
		}
	}
}

// d3: scaleSymlog().constant(10).domain([0, 100]).range([0, 960])
//
//	(50) = 717.3329..., (10) = 277.5022...
func TestSymlogConstant(t *testing.T) {
	s := NewSymlog().SetConstant(10).SetDomain(0, 100).SetRange(0, 960)
	if got := s.Call(50); !approxEps(got, 717.3328668568455, 1e-6) {
		t.Errorf("c=10 Call(50) = %v want 717.3329", got)
	}
	if got := s.Call(10); !approxEps(got, 277.5022332651723, 1e-6) {
		t.Errorf("c=10 Call(10) = %v want 277.5022", got)
	}
}

// d3: scaleSymlog().domain([-100, 100]).range([0, 100])
//
//	(-100) = 0, (0) = 50, (100) = 100, (50) = 92.5972..., (-50) = 7.4028...
func TestSymlogNegativeDomain(t *testing.T) {
	s := NewSymlog().SetDomain(-100, 100).SetRange(0, 100)
	cases := []struct {
		x, want float64
	}{
		{-100, 0},
		{0, 50},
		{100, 100},
		{50, 92.59721515804961},
		{-50, 7.4027848419503846},
	}
	for _, c := range cases {
		if got := s.Call(c.x); !approxEps(got, c.want, 1e-6) {
			t.Errorf("neg Call(%v) = %v want %v", c.x, got, c.want)
		}
	}
}

// d3: symlog uses linearish ticks (same as linear)
func TestSymlogTicks(t *testing.T) {
	s := NewSymlog().SetDomain(0, 100).SetRange(0, 960)
	got := s.Ticks(5)
	want := []float64{0, 20, 40, 60, 80, 100}
	if len(got) != len(want) {
		t.Fatalf("ticks(5) len = %d want %d (%v)", len(got), len(want), got)
	}
	for i, g := range got {
		if !approxEps(g, want[i], 1e-9) {
			t.Errorf("ticks(5)[%d] = %v want %v", i, g, want[i])
		}
	}
}

// d3: symlog uses linearish nice
func TestSymlogNice(t *testing.T) {
	s := NewSymlog().SetDomain(0.201, 0.979)
	s.Nice(10)
	d := s.Domain()
	if !approxEps(d[0], 0.2, 1e-9) || !approxEps(d[1], 1.0, 1e-9) {
		t.Errorf("nice(10) domain = %v want [0.2, 1.0]", d)
	}
}

func TestSymlogInvert(t *testing.T) {
	s := NewSymlog().SetDomain(0, 100).SetRange(0, 960)
	// invert(Call(x)) ≈ x
	for _, x := range []float64{0, 1, 10, 50, 100} {
		y := s.Call(x)
		got := s.Invert(y)
		if !approxEps(got, x, 1e-6) {
			t.Errorf("invert(Call(%v)) = %v want %v", x, got, x)
		}
	}
}

func TestSymlogClamp(t *testing.T) {
	s := NewSymlog().SetDomain(0, 100).SetRange(0, 960).SetClamp(true)
	if got := s.Call(-50); got != 0 {
		t.Errorf("clamped Call(-50) = %v want 0", got)
	}
	if got := s.Call(200); got != 960 {
		t.Errorf("clamped Call(200) = %v want 960", got)
	}
}

func TestSymlogType(t *testing.T) {
	if NewSymlog().Type() != "symlog" {
		t.Error("Type() != symlog")
	}
}

func TestSymlogCopy(t *testing.T) {
	s := NewSymlog().SetDomain(0, 100).SetRange(0, 960).SetConstant(10)
	c := s.Copy()
	c.SetConstant(1).SetDomain(0, 1000)
	if s.Constant() != 10 {
		t.Errorf("original Constant after copy mutate = %v want 10", s.Constant())
	}
	if !approxEps(s.Domain()[1], 100, 1e-9) {
		t.Errorf("original Domain after copy mutate = %v want [0,100]", s.Domain())
	}
}
