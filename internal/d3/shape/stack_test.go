package d3shape

import (
	"math"
	"testing"
)

type stackDatum struct {
	Apples, Oranges, Pears float64
}

func stackValFn(d stackDatum, key string, _ int, _ []stackDatum) float64 {
	switch key {
	case "apples":
		return d.Apples
	case "oranges":
		return d.Oranges
	case "pears":
		return d.Pears
	}
	return 0
}

func stackValABC(d stackDatum, key string, _ int, _ []stackDatum) float64 {
	switch key {
	case "a":
		return d.Apples
	case "b":
		return d.Oranges
	case "c":
		return d.Pears
	}
	return 0
}

var stackData = []stackDatum{{10, 5, 7}, {15, 8, 3}, {12, 6, 9}}

func TestStackOffsetNone(t *testing.T) {
	s := NewStack[stackDatum]().Keys([]string{"apples", "oranges", "pears"}).Value(stackValFn)
	r := s.Call(stackData)
	want := map[string][][2]float64{
		"apples":  {{0, 10}, {0, 15}, {0, 12}},
		"oranges": {{10, 15}, {15, 23}, {12, 18}},
		"pears":   {{15, 22}, {23, 26}, {18, 27}},
	}
	for _, series := range r {
		exp, ok := want[series.Key]
		if !ok {
			t.Errorf("unexpected key %s", series.Key)
			continue
		}
		for j, p := range series.Stats {
			if p.Lo != exp[j][0] || p.Hi != exp[j][1] {
				t.Errorf("series %s point %d = [%.0f, %.0f] want [%.0f, %.0f]",
					series.Key, j, p.Lo, p.Hi, exp[j][0], exp[j][1])
			}
		}
	}
}

func TestStackDivergingMixed(t *testing.T) {
	data := []stackDatum{{10, -5, 3}, {-2, 8, -4}}
	s := NewStack[stackDatum]().Keys([]string{"a", "b", "c"}).Value(stackValABC).Offset(StackOffsetDiverging)
	r := s.Call(data)
	want := map[string][][2]float64{
		"a": {{0, 10}, {-2, 0}},
		"b": {{-5, 0}, {0, 8}},
		"c": {{10, 13}, {-6, -2}},
	}
	for _, series := range r {
		exp, ok := want[series.Key]
		if !ok {
			t.Errorf("unexpected key %s", series.Key)
			continue
		}
		for j, p := range series.Stats {
			if p.Lo != exp[j][0] || p.Hi != exp[j][1] {
				t.Errorf("series %s point %d = [%g, %g] want [%g, %g]",
					series.Key, j, p.Lo, p.Hi, exp[j][0], exp[j][1])
			}
		}
	}
}

func TestStackExpand(t *testing.T) {
	s := NewStack[stackDatum]().Keys([]string{"apples", "oranges", "pears"}).Value(stackValFn).Offset(StackOffsetExpand)
	r := s.Call(stackData)
	// Check the last series reaches 1.0 at each point
	pears := r[2]
	for j, p := range pears.Stats {
		if math.Abs(p.Hi-1.0) > 1e-9 {
			t.Errorf("expand pears[%d].hi = %g want 1.0", j, p.Hi)
		}
	}
	// Check first series starts at 0
	apples := r[0]
	for j, p := range apples.Stats {
		if p.Lo != 0 {
			t.Errorf("expand apples[%d].lo = %g want 0", j, p.Lo)
		}
	}
}

func TestStackOrderAscending(t *testing.T) {
	s := NewStack[stackDatum]().Keys([]string{"apples", "oranges", "pears"}).Value(stackValFn).Order(StackOrderAscending)
	r := s.Call(stackData)
	want := map[string]int{"apples": 2, "oranges": 0, "pears": 1}
	for _, series := range r {
		if series.Index != want[series.Key] {
			t.Errorf("ascending %s.index = %d want %d", series.Key, series.Index, want[series.Key])
		}
	}
}

func TestStackOrderDescending(t *testing.T) {
	s := NewStack[stackDatum]().Keys([]string{"apples", "oranges", "pears"}).Value(stackValFn).Order(StackOrderDescending)
	r := s.Call(stackData)
	want := map[string]int{"apples": 0, "oranges": 2, "pears": 1}
	for _, series := range r {
		if series.Index != want[series.Key] {
			t.Errorf("descending %s.index = %d want %d", series.Key, series.Index, want[series.Key])
		}
	}
}

func TestStackOrderReverse(t *testing.T) {
	s := NewStack[stackDatum]().Keys([]string{"apples", "oranges", "pears"}).Value(stackValFn).Order(StackOrderReverse)
	r := s.Call(stackData)
	want := map[string]int{"apples": 2, "oranges": 1, "pears": 0}
	for _, series := range r {
		if series.Index != want[series.Key] {
			t.Errorf("reverse %s.index = %d want %d", series.Key, series.Index, want[series.Key])
		}
	}
}

func TestStackOrderAppearance(t *testing.T) {
	s := NewStack[stackDatum]().Keys([]string{"apples", "oranges", "pears"}).Value(stackValFn).Order(StackOrderAppearance)
	r := s.Call(stackData)
	// Just verify it produces a valid permutation
	seen := map[int]bool{}
	for _, series := range r {
		if seen[series.Index] {
			t.Errorf("appearance: duplicate index %d", series.Index)
		}
		seen[series.Index] = true
	}
	if len(seen) != 3 {
		t.Errorf("appearance: expected 3 unique indices, got %d", len(seen))
	}
}

func TestStackOrderInsideOut(t *testing.T) {
	s := NewStack[stackDatum]().Keys([]string{"apples", "oranges", "pears"}).Value(stackValFn).Order(StackOrderInsideOut)
	r := s.Call(stackData)
	seen := map[int]bool{}
	for _, series := range r {
		if seen[series.Index] {
			t.Errorf("insideOut: duplicate index %d", series.Index)
		}
		seen[series.Index] = true
	}
	if len(seen) != 3 {
		t.Errorf("insideOut: expected 3 unique indices, got %d", len(seen))
	}
}

func TestStackOffsetSilhouette(t *testing.T) {
	s := NewStack[stackDatum]().Keys([]string{"apples", "oranges", "pears"}).Value(stackValFn).Offset(StackOffsetSilhouette)
	r := s.Call(stackData)
	// The first series (in order) should be centered around 0.
	// Sum at each point: 22, 26, 27 → centered at -11, -13, -13.5
	s0 := r[0] // apples is first in order (orderNone)
	for j, p := range s0.Stats {
		sum := stackData[j].Apples + stackData[j].Oranges + stackData[j].Pears
		if math.Abs(p.Lo-(-sum/2)) > 1e-9 {
			t.Errorf("silhouette s0[%d].lo = %g want %g", j, p.Lo, -sum/2)
		}
	}
}

func TestStackOffsetWiggle(t *testing.T) {
	s := NewStack[stackDatum]().Keys([]string{"apples", "oranges", "pears"}).Value(stackValFn).Offset(StackOffsetWiggle)
	r := s.Call(stackData)
	// Just verify it runs and produces valid (non-NaN) values
	for _, series := range r {
		for j, p := range series.Stats {
			if math.IsNaN(p.Lo) || math.IsNaN(p.Hi) {
				t.Errorf("wiggle %s[%d] has NaN: [%g, %g]", series.Key, j, p.Lo, p.Hi)
			}
		}
	}
}
