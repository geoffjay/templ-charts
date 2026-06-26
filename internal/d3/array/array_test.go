package d3array

import (
	"math"
	"testing"
)

func approx(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func TestExtentEmpty(t *testing.T) {
	e := Extent(nil)
	if !math.IsNaN(e[0]) || !math.IsNaN(e[1]) {
		t.Errorf("Extent(nil) = %v want [NaN NaN]", e)
	}
}

func TestExtentBasic(t *testing.T) {
	e := Extent([]float64{3, 1, 4, 1, 5, 9, 2, 6, 5})
	if !approx(e[0], 1) || !approx(e[1], 9) {
		t.Errorf("Extent = %v want [1 9]", e)
	}
}

func TestExtentIgnoresNaN(t *testing.T) {
	e := Extent([]float64{math.NaN(), 2, math.NaN(), 8})
	if !approx(e[0], 2) || !approx(e[1], 8) {
		t.Errorf("Extent with NaNs = %v want [2 8]", e)
	}
}

func TestMinMaxSum(t *testing.T) {
	vs := []float64{1, 2, 3, 4, 5}
	if !approx(Min(vs), 1) {
		t.Errorf("Min = %v want 1", Min(vs))
	}
	if !approx(Max(vs), 5) {
		t.Errorf("Max = %v want 5", Max(vs))
	}
	if !approx(Sum(vs), 15) {
		t.Errorf("Sum = %v want 15", Sum(vs))
	}
}

func TestMeanMedian(t *testing.T) {
	if !approx(Mean([]float64{1, 2, 3, 4}), 2.5) {
		t.Errorf("Mean = %v want 2.5", Mean([]float64{1, 2, 3, 4}))
	}
	if !approx(Median([]float64{1, 2, 3, 4}), 2.5) {
		t.Errorf("Median even = %v want 2.5", Median([]float64{1, 2, 3, 4}))
	}
	if !approx(Median([]float64{1, 2, 3}), 2) {
		t.Errorf("Median odd = %v want 2", Median([]float64{1, 2, 3}))
	}
	if !math.IsNaN(Mean(nil)) {
		t.Errorf("Mean(nil) should be NaN, got %v", Mean(nil))
	}
}

func TestRange(t *testing.T) {
	cases := []struct {
		args []float64
		want []float64
	}{
		{[]float64{5}, []float64{0, 1, 2, 3, 4}},
		{[]float64{2, 5}, []float64{2, 3, 4}},
		{[]float64{0, 10, 2}, []float64{0, 2, 4, 6, 8}},
		{[]float64{0, 1, 0.25}, []float64{0, 0.25, 0.5, 0.75}},
		{[]float64{0, 0, 1}, []float64{}},
	}
	for _, c := range cases {
		got := Range(c.args...)
		if len(got) != len(c.want) {
			t.Errorf("Range(%v) = %v want %v", c.args, got, c.want)
			continue
		}
		for i := range got {
			if !approx(got[i], c.want[i]) {
				t.Errorf("Range(%v)[%d] = %v want %v", c.args, i, got[i], c.want[i])
			}
		}
	}
}

func TestAscendingDescending(t *testing.T) {
	if Ascending(1, 2) != -1 {
		t.Errorf("Ascending(1,2) should be -1")
	}
	if Ascending(2, 1) != 1 {
		t.Errorf("Ascending(2,1) should be 1")
	}
	if Ascending(1, 1) != 0 {
		t.Errorf("Ascending(1,1) should be 0")
	}
	if Ascending(math.NaN(), 1) != 1 {
		t.Errorf("NaN should sort last (asc)")
	}
	if Descending(2, 1) != -1 {
		t.Errorf("Descending(2,1) should be -1")
	}
}

func TestBisect(t *testing.T) {
	vs := []float64{1, 2, 3, 4, 5}
	// left bisect: equal values go before
	if i := BisectLeft(vs, 3); i != 2 {
		t.Errorf("BisectLeft(3) = %d want 2", i)
	}
	// right bisect: equal values go after
	if i := BisectRight(vs, 3); i != 3 {
		t.Errorf("BisectRight(3) = %d want 3", i)
	}
	// out-of-range
	if i := BisectLeft(vs, 0); i != 0 {
		t.Errorf("BisectLeft(0) = %d want 0", i)
	}
	if i := BisectRight(vs, 6); i != 5 {
		t.Errorf("BisectRight(6) = %d want 5", i)
	}
}

func TestTicks(t *testing.T) {
	cases := []struct {
		start, stop float64
		count       int
		wantMin     float64
		wantMax     float64
	}{
		{0, 10, 5, 0, 10},
		{0, 100, 5, 0, 100},
		{1, 9, 5, 2, 8}, // d3 returns [2,4,6,8] (step 2)
	}
	for _, c := range cases {
		got := Ticks(c.start, c.stop, c.count)
		if len(got) == 0 {
			t.Errorf("Ticks(%v,%v,%d) = [] want non-empty", c.start, c.stop, c.count)
			continue
		}
		if !approx(got[0], c.wantMin) {
			t.Errorf("Ticks(%v,%v,%d)[0] = %v want %v", c.start, c.stop, c.count, got[0], c.wantMin)
		}
		if !approx(got[len(got)-1], c.wantMax) {
			t.Errorf("Ticks(%v,%v,%d)[-1] = %v want %v", c.start, c.stop, c.count, got[len(got)-1], c.wantMax)
		}
	}
}

func TestTicksEdgeCases(t *testing.T) {
	if got := Ticks(math.NaN(), 10, 5); len(got) != 0 {
		t.Errorf("Ticks(NaN, 10, 5) = %v want []", got)
	}
	if got := Ticks(0, 10, 0); len(got) != 0 {
		t.Errorf("Ticks(0, 10, 0) = %v want []", got)
	}
	// zero span -> just [start]
	got := Ticks(5, 5, 5)
	if len(got) != 1 || !approx(got[0], 5) {
		t.Errorf("Ticks(5, 5, 5) = %v want [5]", got)
	}
}

func TestNumberable(t *testing.T) {
	if !Numberable(5) {
		t.Errorf("5 should be Numberable")
	}
	if Numberable(math.NaN()) {
		t.Errorf("NaN should not be Numberable")
	}
	if Numberable(math.Inf(1)) {
		t.Errorf("+Inf should not be Numberable")
	}
}
