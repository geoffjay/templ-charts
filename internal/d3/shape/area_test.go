package d3shape

import (
	"math"
	"testing"
)

var areaTestData = []Point2D{{0, 0}, {10, 20}, {20, 10}, {30, 40}, {40, 30}}

func TestAreaLinear(t *testing.T) {
	got := NewArea().Call(areaTestData)
	want := "M0,0L10,20L20,10L30,40L40,30L40,0L30,0L20,0L10,0L0,0Z"
	if got != want {
		t.Errorf("linear:\n got=%q\nwant=%q", got, want)
	}
}

func TestAreaStepY0(t *testing.T) {
	got := NewArea().Y0(func(_ Point2D, _ int, _ []Point2D) float64 { return 0 }).Curve(CurveStep).Call(areaTestData)
	want := "M0,0L5,0L5,20L15,20L15,10L25,10L25,40L35,40L35,30L40,30L40,0L35,0L35,0L25,0L25,0L15,0L15,0L5,0L5,0L0,0Z"
	if got != want {
		t.Errorf("step y0=0:\n got=%q\nwant=%q", got, want)
	}
}

func TestAreaMonotoneXY0_5(t *testing.T) {
	got := NewArea().Y0(func(_ Point2D, _ int, _ []Point2D) float64 { return 5 }).Curve(CurveMonotoneX).Call(areaTestData)
	want := "M0,0C3.333,10,6.667,20,10,20C13.333,20,16.667,10,20,10C23.333,10,26.667,40,30,40C33.333,40,36.667,35,40,30L40,5C36.667,5,33.333,5,30,5C26.667,5,23.333,5,20,5C16.667,5,13.333,5,10,5C6.667,5,3.333,5,0,5Z"
	if got != want {
		t.Errorf("monotoneX y0=5:\n got=%q\nwant=%q", got, want)
	}
}

func TestAreaBasis(t *testing.T) {
	got := NewArea().Curve(CurveBasis).Call(areaTestData)
	want := "M0,0L1.667,3.333C3.333,6.667,6.667,13.333,10,15C13.333,16.667,16.667,13.333,20,16.667C23.333,20,26.667,30,30,33.333C33.333,36.667,36.667,33.333,38.333,31.667L40,30L40,0L38.333,0C36.667,0,33.333,0,30,0C26.667,0,23.333,0,20,0C16.667,0,13.333,0,10,0C6.667,0,3.333,0,1.667,0L0,0Z"
	if got != want {
		t.Errorf("basis:\n got=%q\nwant=%q", got, want)
	}
}

func TestAreaDefaultY0(t *testing.T) {
	got := NewArea().Call(areaTestData)
	want := "M0,0L10,20L20,10L30,40L40,30L40,0L30,0L20,0L10,0L0,0Z"
	if got != want {
		t.Errorf("default y0:\n got=%q\nwant=%q", got, want)
	}
}

func TestAreaDefinedGap(t *testing.T) {
	data := []Point2D{{0, 0}, {10, 20}, {20, math.NaN()}, {30, 40}, {40, 30}}
	a := NewArea().Defined(func(d Point2D, _ int, _ []Point2D) bool { return !math.IsNaN(d[1]) })
	got := a.Call(data)
	want := "M0,0L10,20L10,0L0,0ZM30,40L40,30L40,0L30,0Z"
	if got != want {
		t.Errorf("defined-gap:\n got=%q\nwant=%q", got, want)
	}
}

func TestAreaSinglePoint(t *testing.T) {
	got := NewArea().Call([]Point2D{{5, 5}})
	want := "M5,5L5,0Z"
	if got != want {
		t.Errorf("single: got=%q want=%q", got, want)
	}
}

func TestAreaEmpty(t *testing.T) {
	got := NewArea().Call([]Point2D{})
	if got != "" {
		t.Errorf("empty: got=%q want %q", got, "")
	}
}

func TestAreaX0X1Y0Y1(t *testing.T) {
	a := NewArea().
		X0(func(d Point2D, _ int, _ []Point2D) float64 { return d[0] }).
		X1(func(d Point2D, _ int, _ []Point2D) float64 { return d[0] + 5 }).
		Y0(func(d Point2D, _ int, _ []Point2D) float64 { return d[1] }).
		Y1(func(d Point2D, _ int, _ []Point2D) float64 { return d[1] + 10 })
	got := a.Call(areaTestData)
	want := "M5,10L15,30L25,20L35,50L45,40L40,30L30,40L20,10L10,20L0,0Z"
	if got != want {
		t.Errorf("x0x1y0y1:\n got=%q\nwant=%q", got, want)
	}
}
