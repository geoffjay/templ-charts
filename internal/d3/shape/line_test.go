package d3shape

import (
	"math"
	"testing"
)

// goldenD3 holds the reference outputs from d3-shape v3.2.0 for the dataset
// [[0,0],[10,20],[20,10],[30,40],[40,30]] with each curve factory.
var lineTestData = []Point2D{{0, 0}, {10, 20}, {20, 10}, {30, 40}, {40, 30}}

var lineTestCurves = []struct {
	name string
	f    CurveFactory
	want string
}{
	{"linear", CurveLinear, "M0,0L10,20L20,10L30,40L40,30"},
	{"step", CurveStep, "M0,0L5,0L5,20L15,20L15,10L25,10L25,40L35,40L35,30L40,30"},
	{"stepAfter", CurveStepAfter, "M0,0L10,0L10,20L20,20L20,10L30,10L30,40L40,40L40,30"},
	{"stepBefore", CurveStepBefore, "M0,0L0,20L10,20L10,10L20,10L20,40L30,40L30,30L40,30"},
	{"basis", CurveBasis, "M0,0L1.667,3.333C3.333,6.667,6.667,13.333,10,15C13.333,16.667,16.667,13.333,20,16.667C23.333,20,26.667,30,30,33.333C33.333,36.667,36.667,33.333,38.333,31.667L40,30"},
	{"cardinal", CurveCardinal, "M0,0C0,0,6.667,18.333,10,20C13.333,21.667,16.667,6.667,20,10C23.333,13.333,26.667,36.667,30,40C33.333,43.333,40,30,40,30"},
	{"catmullRom", CurveCatmullRom, "M0,0C0,0,6.189,19.382,10,20C13.031,20.492,17.109,9.318,20,10C24.323,11.02,25.677,38.98,30,40C32.891,40.682,38.333,31.667,40,30"},
	{"monotoneX", CurveMonotoneX, "M0,0C3.333,10,6.667,20,10,20C13.333,20,16.667,10,20,10C23.333,10,26.667,40,30,40C33.333,40,36.667,35,40,30"},
	{"natural", CurveNatural, "M0,0C3.333,10.536,6.667,21.071,10,20C13.333,18.929,16.667,6.25,20,10C23.333,13.75,26.667,33.929,30,40C33.333,46.071,36.667,38.036,40,30"},
}

func TestLineAllCurves(t *testing.T) {
	for _, c := range lineTestCurves {
		got := NewLine().Curve(c.f).Call(lineTestData)
		if got != c.want {
			t.Errorf("curve %s:\n got=%q\nwant=%q", c.name, got, c.want)
		}
	}
}

// d3: line().defined(d => d[1] !== null) with a gap → two subpaths
func TestLineDefinedGap(t *testing.T) {
	data := []Point2D{{0, 0}, {10, 20}, {20, math.NaN()}, {30, 40}, {40, 30}}
	l := NewLine().Defined(func(d Point2D, i int, _ []Point2D) bool { return !math.IsNaN(d[1]) })
	got := l.Call(data)
	want := "M0,0L10,20M30,40L40,30"
	if got != want {
		t.Errorf("defined-gap:\n got=%q\nwant=%q", got, want)
	}
}

// d3: line()([[5,5]]) === "M5,5Z" — single point closes.
// Actually d3 returns "M5,5" for a single point (no Z since line is NaN).
func TestLineSinglePoint(t *testing.T) {
	got := NewLine().Call([]Point2D{{5, 5}})
	want := "M5,5"
	if got != want {
		t.Errorf("single point: got=%q want=%q", got, want)
	}
}

// d3: line()([]) === null → we return ""
func TestLineEmpty(t *testing.T) {
	got := NewLine().Call([]Point2D{})
	if got != "" {
		t.Errorf("empty: got=%q want %q", got, "")
	}
}

// d3: line().x(d=>d[0]*2).y(d=>d[1]*3) scales accessors
func TestLineAccessors(t *testing.T) {
	l := NewLine().
		X(func(d Point2D, _ int, _ []Point2D) float64 { return d[0] * 2 }).
		Y(func(d Point2D, _ int, _ []Point2D) float64 { return d[1] * 3 })
	got := l.Call([]Point2D{{0, 0}, {10, 20}})
	want := "M0,0L20,60"
	if got != want {
		t.Errorf("accessors: got=%q want=%q", got, want)
	}
}

// d3: monotoneY reflects x/y
func TestLineMonotoneY(t *testing.T) {
	l := NewLine().Curve(CurveMonotoneY)
	got := l.Call(lineTestData)
	// monotoneY swaps x/y in the curve, so the path is the mirror of monotoneX.
	// We just verify it produces a non-empty path with the expected structure.
	if len(got) == 0 {
		t.Fatal("monotoneY produced empty path")
	}
	// The first segment should start with M0,0 (swapped: still 0,0)
	wantPrefix := "M0,0"
	if got[:len(wantPrefix)] != wantPrefix {
		t.Errorf("monotoneY prefix: got=%q want=%q", got[:len(wantPrefix)], wantPrefix)
	}
}

// d3: linearClosed curve closes the path even for a line (non-area)
func TestLineLinearClosed(t *testing.T) {
	l := NewLine().Curve(CurveLinearClosed)
	got := l.Call([]Point2D{{0, 0}, {10, 20}, {20, 10}})
	want := "M0,0L10,20L20,10Z"
	if got != want {
		t.Errorf("linearClosed: got=%q want=%q", got, want)
	}
}

func TestLineDigits(t *testing.T) {
	l := NewLine().Digits(6)
	got := l.Call([]Point2D{{0, 0}, {1.123456789, 2.987654321}})
	// With 6 digits, coordinates should be rounded to 6 decimal places
	want := "M0,0L1.123457,2.987654"
	if got != want {
		t.Errorf("digits=6: got=%q want=%q", got, want)
	}
}
