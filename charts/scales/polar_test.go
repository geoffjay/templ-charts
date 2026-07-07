package scales

import "testing"

func TestNewLinearScaleWithRange(t *testing.T) {
	s := NewLinearScaleWithRange(0, 10, 100, 200)
	if s.Type() != ScaleTypeLinear {
		t.Fatalf("type = %v", s.Type())
	}
	if got := s.Call(0.0); got != 100 {
		t.Errorf("f(0) = %v, want 100", got)
	}
	if got := s.Call(5.0); got != 150 {
		t.Errorf("f(5) = %v, want 150", got)
	}
	if got := s.Call(10.0); got != 200 {
		t.Errorf("f(10) = %v, want 200", got)
	}
	// The result is a full scaleImpl, so tick machinery works on it.
	if got := GetScaleTicks(s, TicksSpec{HasCount: true, Count: 5}); len(got) == 0 {
		t.Error("linear-with-range should produce ticks")
	}
}

func TestNewBandScaleWithRange(t *testing.T) {
	s := NewBandScaleWithRange([]string{"a", "b"}, 0, 100, 0, false)
	if s.Type() != ScaleTypeBand {
		t.Fatalf("type = %v", s.Type())
	}
	bw := s.(ScaleWithBandwidth).Bandwidth()
	if bw != 50 {
		t.Errorf("bandwidth = %v, want 50", bw)
	}
	if got := s.Call("a"); got != 0 {
		t.Errorf("f(a) = %v, want 0", got)
	}
	if got := s.Call("b"); got != 50 {
		t.Errorf("f(b) = %v, want 50", got)
	}
	// Rounded variant reports Round.
	sr := NewBandScaleWithRange([]string{"a", "b", "c"}, 0, 100, 0.1, true)
	if !sr.(ScaleWithBandwidth).Round() {
		t.Error("round variant should report Round() = true")
	}
	if sr.(ScaleWithBandwidth).Bandwidth() <= 0 {
		t.Error("padded band should still have positive bandwidth")
	}
}

func TestNewPointScaleWithRange(t *testing.T) {
	s := NewPointScaleWithRange([]string{"a", "b", "c"}, 0, 100, 0, false)
	if s.Type() != ScaleTypePoint {
		t.Fatalf("type = %v", s.Type())
	}
	if got := s.Call("a"); got != 0 {
		t.Errorf("f(a) = %v, want 0", got)
	}
	if got := s.Call("b"); got != 50 {
		t.Errorf("f(b) = %v, want 50", got)
	}
	if got := s.Call("c"); got != 100 {
		t.Errorf("f(c) = %v, want 100", got)
	}
	if got := s.(ScaleWithBandwidth).Bandwidth(); got != 0 {
		t.Errorf("point bandwidth = %v, want 0", got)
	}
	// Rounded + padded variant.
	sr := NewPointScaleWithRange([]string{"a", "b"}, 0, 101, 0.5, true)
	if got := sr.Call("a"); got != float64(int64(got)) {
		t.Errorf("rounded point f(a) = %v, want integral", got)
	}
}
