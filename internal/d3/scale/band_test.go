package d3scale

import (
	"math"
	"testing"
)

func approxStr(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	return math.Abs(a-b) < 1e-9
}

// d3: scaleBand().domain(["a","b","c"]).range([0,960])
//
//	→ a=0, b=320, c=640, bandwidth=320, step=320
func TestBandBasic(t *testing.T) {
	s := NewBand().SetDomain([]string{"a", "b", "c"}).SetRange(0, 960)
	cases := []struct {
		x    string
		want float64
	}{
		{"a", 0},
		{"b", 320},
		{"c", 640},
	}
	for _, c := range cases {
		if got := s.Call(c.x); !approxStr(got, c.want) {
			t.Errorf("Call(%q) = %v want %v", c.x, got, c.want)
		}
	}
	if !approxStr(s.Bandwidth(), 320) {
		t.Errorf("Bandwidth = %v want 320", s.Bandwidth())
	}
	if !approxStr(s.Step(), 320) {
		t.Errorf("Step = %v want 320", s.Step())
	}
}

// d3: scaleBand().domain(["a","b","c"]).range([0,960]).padding(0.1)
//
//	→ a=30.9677..., b=340.645..., c=650.322..., bw=278.709..., step=309.677...
func TestBandPadding(t *testing.T) {
	s := NewBand().SetDomain([]string{"a", "b", "c"}).SetRange(0, 960).SetPadding(0.1)
	want := map[string]float64{
		"a": 30.967741935483843,
		"b": 340.64516129032256,
		"c": 650.3225806451612,
	}
	for x, w := range want {
		if got := s.Call(x); !approxEps(got, w, 1e-6) {
			t.Errorf("Call(%q) = %v want %v", x, got, w)
		}
	}
	if !approxEps(s.Bandwidth(), 278.7096774193549, 1e-6) {
		t.Errorf("Bandwidth = %v want 278.7097", s.Bandwidth())
	}
	if !approxEps(s.Step(), 309.6774193548387, 1e-6) {
		t.Errorf("Step = %v want 309.6774", s.Step())
	}
}

// d3: scaleBand().domain(["a","b","c"]).range([0,960])
//
//	.paddingInner(0.2).paddingOuter(0.3)
//	→ a=84.705..., b=367.058..., c=649.411..., bw=225.882..., step=282.352...
func TestBandPaddingInnerOuter(t *testing.T) {
	s := NewBand().
		SetDomain([]string{"a", "b", "c"}).
		SetRange(0, 960).
		SetPaddingInner(0.2).
		SetPaddingOuter(0.3)
	want := map[string]float64{
		"a": 84.70588235294116,
		"b": 367.05882352941177,
		"c": 649.4117647058824,
	}
	for x, w := range want {
		if got := s.Call(x); !approxEps(got, w, 1e-6) {
			t.Errorf("Call(%q) = %v want %v", x, got, w)
		}
	}
	if !approxEps(s.Bandwidth(), 225.8823529411765, 1e-6) {
		t.Errorf("Bandwidth = %v want 225.8824", s.Bandwidth())
	}
}

// d3: scaleBand()...rangeRound([0,960]).padding(0.1) → rounded values
//
//	a=32, b=341, c=650, bw=278, step=309
func TestBandRound(t *testing.T) {
	s := NewBand().SetDomain([]string{"a", "b", "c"}).SetRangeRound(0, 960).SetPadding(0.1)
	want := map[string]float64{
		"a": 32,
		"b": 341,
		"c": 650,
	}
	for x, w := range want {
		if got := s.Call(x); got != w {
			t.Errorf("round Call(%q) = %v want %v", x, got, w)
		}
	}
	if s.Bandwidth() != 278 {
		t.Errorf("round Bandwidth = %v want 278", s.Bandwidth())
	}
	if s.Step() != 309 {
		t.Errorf("round Step = %v want 309", s.Step())
	}
}

// d3: align=0,0.5,1 with padding 0.2
//
//	align0:  a=0   b=300 c=600  bw=240 step=300
//	align1:  a=120 b=420 c=720  bw=240 step=300
//	default: a=60  b=360 c=660  bw=240 step=300
func TestBandAlign(t *testing.T) {
	check := func(align, a, b, c float64) {
		s := NewBand().SetDomain([]string{"a", "b", "c"}).SetRange(0, 960).SetPadding(0.2).SetAlign(align)
		if got := s.Call("a"); !approxStr(got, a) {
			t.Errorf("align=%v Call(a) = %v want %v", align, got, a)
		}
		if got := s.Call("b"); !approxStr(got, b) {
			t.Errorf("align=%v Call(b) = %v want %v", align, got, b)
		}
		if got := s.Call("c"); !approxStr(got, c) {
			t.Errorf("align=%v Call(c) = %v want %v", align, got, c)
		}
	}
	check(0, 0, 300, 600)
	check(1, 120, 420, 720)
	check(0.5, 60, 360, 660)
}

// d3: reversed range → values reversed
//
//	scaleBand().domain(["a","b","c"]).range([960,0]) → a=640 b=320 c=0
func TestBandReversedRange(t *testing.T) {
	s := NewBand().SetDomain([]string{"a", "b", "c"}).SetRange(960, 0)
	if got := s.Call("a"); !approxStr(got, 640) {
		t.Errorf("reversed Call(a) = %v want 640", got)
	}
	if got := s.Call("b"); !approxStr(got, 320) {
		t.Errorf("reversed Call(b) = %v want 320", got)
	}
	if got := s.Call("c"); !approxStr(got, 0) {
		t.Errorf("reversed Call(c) = %v want 0", got)
	}
}

// Unknown domain value → NaN
func TestBandUnknown(t *testing.T) {
	s := NewBand().SetDomain([]string{"a", "b"}).SetRange(0, 100)
	if !math.IsNaN(s.Call("zzz")) {
		t.Errorf("Call(unknown) = %v want NaN", s.Call("zzz"))
	}
}

// Domain dedup
func TestBandDomainDedup(t *testing.T) {
	s := NewBand().SetDomain([]string{"a", "b", "a", "c"}).SetRange(0, 300)
	if len(s.Domain()) != 3 {
		t.Errorf("dedup Domain len = %d want 3 (%v)", len(s.Domain()), s.Domain())
	}
}

// Copy independence
func TestBandCopy(t *testing.T) {
	s := NewBand().SetDomain([]string{"a", "b"}).SetRange(0, 100).SetPadding(0.2)
	c := s.Copy()
	c.SetDomain([]string{"x", "y", "z"}).SetRange(0, 999)
	if len(s.Domain()) != 2 {
		t.Errorf("original Domain after copy mutate = %v want 2", s.Domain())
	}
	if !approxEps(s.Call("a"), 9.090909090909093, 1e-9) { // padding 0.2 shifts a
		t.Errorf("original Call(a) = %v want 9.0909", s.Call("a"))
	}
}

func TestBandType(t *testing.T) {
	if NewBand().Type() != "band" {
		t.Error("Type() != band")
	}
}

// --- Point scale ----------------------------------------------------------

// d3: scalePoint().domain(["a","b","c"]).range([0,960])
//
//	→ a=0, b=480, c=960, bandwidth=0, step=480
func TestPointBasic(t *testing.T) {
	s := NewPoint().SetDomain([]string{"a", "b", "c"}).SetRange(0, 960)
	cases := []struct {
		x    string
		want float64
	}{
		{"a", 0},
		{"b", 480},
		{"c", 960},
	}
	for _, c := range cases {
		if got := s.Call(c.x); !approxStr(got, c.want) {
			t.Errorf("point Call(%q) = %v want %v", c.x, got, c.want)
		}
	}
	if s.Bandwidth() != 0 {
		t.Errorf("point Bandwidth = %v want 0", s.Bandwidth())
	}
	if !approxStr(s.Step(), 480) {
		t.Errorf("point Step = %v want 480", s.Step())
	}
}

// d3: scalePoint().domain(["a","b","c"]).range([0,960]).padding(0.5)
//
//	→ a=160, b=480, c=800, step=320
func TestPointPadding(t *testing.T) {
	s := NewPoint().SetDomain([]string{"a", "b", "c"}).SetRange(0, 960).SetPadding(0.5)
	if got := s.Call("a"); !approxStr(got, 160) {
		t.Errorf("point pad0.5 Call(a) = %v want 160", got)
	}
	if got := s.Call("b"); !approxStr(got, 480) {
		t.Errorf("point pad0.5 Call(b) = %v want 480", got)
	}
	if got := s.Call("c"); !approxStr(got, 800) {
		t.Errorf("point pad0.5 Call(c) = %v want 800", got)
	}
	if !approxStr(s.Step(), 320) {
		t.Errorf("point pad0.5 Step = %v want 320", s.Step())
	}
}

func TestPointType(t *testing.T) {
	if NewPoint().Type() != "point" {
		t.Error("Type() != point")
	}
}

func TestPointCopy(t *testing.T) {
	s := NewPoint().SetDomain([]string{"a", "b"}).SetRange(0, 100).SetPadding(0.2)
	c := s.Copy()
	c.SetDomain([]string{"x", "y", "z"})
	if len(s.Domain()) != 2 {
		t.Errorf("original Domain after copy mutate = %v want 2", s.Domain())
	}
}
