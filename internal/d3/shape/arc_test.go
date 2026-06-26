package d3shape

import (
	"math"
	"testing"
)

// arc0 is a shorthand for an arc with innerRadius=0, outerRadius=50.
func arc0() *Arc {
	return NewArc().InnerRadius(func(d ArcDatum) float64 { return d.InnerRadius }).OuterRadius(func(d ArcDatum) float64 { return d.OuterRadius })
}

func TestArcFullCircle(t *testing.T) {
	got := arc0().Call(ArcDatum{InnerRadius: 0, OuterRadius: 50, StartAngle: 0, EndAngle: 2 * math.Pi})
	want := "M0,-50A50,50,0,1,1,0,50A50,50,0,1,1,0,-50Z"
	if got != want {
		t.Errorf("full circle:\n got=%q\nwant=%q", got, want)
	}
}

func TestArcHalfCircle(t *testing.T) {
	got := arc0().Call(ArcDatum{InnerRadius: 0, OuterRadius: 50, StartAngle: 0, EndAngle: math.Pi})
	want := "M0,-50A50,50,0,1,1,0,50L0,0Z"
	if got != want {
		t.Errorf("half circle:\n got=%q\nwant=%q", got, want)
	}
}

func TestArcQuarter(t *testing.T) {
	got := arc0().Call(ArcDatum{InnerRadius: 0, OuterRadius: 50, StartAngle: 0, EndAngle: math.Pi / 2})
	want := "M0,-50A50,50,0,0,1,50,0L0,0Z"
	if got != want {
		t.Errorf("quarter:\n got=%q\nwant=%q", got, want)
	}
}

func TestArcDonutHalf(t *testing.T) {
	got := arc0().Call(ArcDatum{InnerRadius: 20, OuterRadius: 50, StartAngle: 0, EndAngle: math.Pi})
	want := "M0,-50A50,50,0,1,1,0,50L0,20A20,20,0,1,0,0,-20Z"
	if got != want {
		t.Errorf("donut half:\n got=%q\nwant=%q", got, want)
	}
}

func TestArcDonutQuarter(t *testing.T) {
	got := arc0().Call(ArcDatum{InnerRadius: 20, OuterRadius: 50, StartAngle: 0, EndAngle: math.Pi / 2})
	want := "M0,-50A50,50,0,0,1,50,0L20,0A20,20,0,0,0,0,-20Z"
	if got != want {
		t.Errorf("donut quarter:\n got=%q\nwant=%q", got, want)
	}
}

func TestArcCornerRadius(t *testing.T) {
	got := arc0().CornerRadius(5).Call(ArcDatum{InnerRadius: 20, OuterRadius: 50, StartAngle: 0, EndAngle: math.Pi / 2})
	want := "M0,-44.721A5,5,0,0,1,5.556,-49.69A50,50,0,0,1,49.69,-5.556A5,5,0,0,1,44.721,0L24.495,0A5,5,0,0,1,19.596,-4A20,20,0,0,0,4,-19.596A5,5,0,0,1,0,-24.495Z"
	if got != want {
		t.Errorf("corner:\n got=%q\nwant=%q", got, want)
	}
}

func TestArcPadAngle(t *testing.T) {
	got := arc0().PadAngle(0.1).Call(ArcDatum{InnerRadius: 20, OuterRadius: 50, StartAngle: 0, EndAngle: math.Pi / 2})
	want := "M2.691,-49.928A50,50,0,0,1,49.928,-2.691L19.818,-2.691A20,20,0,0,0,2.691,-19.818Z"
	if got != want {
		t.Errorf("pad:\n got=%q\nwant=%q", got, want)
	}
}

func TestArcCornerAndPad(t *testing.T) {
	got := arc0().CornerRadius(5).PadAngle(0.1).Call(ArcDatum{InnerRadius: 20, OuterRadius: 50, StartAngle: 0, EndAngle: math.Pi / 2})
	want := "M2.691,-44.338A5,5,0,0,1,8.546,-49.264A50,50,0,0,1,49.264,-8.546A5,5,0,0,1,44.338,-2.691L23.787,-2.691A5,5,0,0,1,19.03,-6.153A20,20,0,0,0,6.153,-19.03A5,5,0,0,1,2.691,-23.787Z"
	if got != want {
		t.Errorf("corner+pad:\n got=%q\nwant=%q", got, want)
	}
}

func TestArcPoint(t *testing.T) {
	got := arc0().Call(ArcDatum{InnerRadius: 0, OuterRadius: 50, StartAngle: 0, EndAngle: 0})
	want := "M0,-50L0,0Z"
	if got != want {
		t.Errorf("point:\n got=%q\nwant=%q", got, want)
	}
}

func TestArcAnnulus(t *testing.T) {
	got := arc0().Call(ArcDatum{InnerRadius: 20, OuterRadius: 50, StartAngle: 0, EndAngle: 2 * math.Pi})
	want := "M0,-50A50,50,0,1,1,0,50A50,50,0,1,1,0,-50M0,-20A20,20,0,1,0,0,20A20,20,0,1,0,0,-20Z"
	if got != want {
		t.Errorf("annulus:\n got=%q\nwant=%q", got, want)
	}
}

func TestArcCentroid(t *testing.T) {
	a := arc0()
	c := a.Centroid(ArcDatum{InnerRadius: 0, OuterRadius: 50, StartAngle: 0, EndAngle: math.Pi / 2})
	wantX := 17.67766952966369
	wantY := -17.677669529663685
	if math.Abs(c[0]-wantX) > 1e-9 || math.Abs(c[1]-wantY) > 1e-9 {
		t.Errorf("centroid = [%g, %g] want [%g, %g]", c[0], c[1], wantX, wantY)
	}
}

func TestArcOffsetStart(t *testing.T) {
	got := arc0().Call(ArcDatum{InnerRadius: 20, OuterRadius: 50, StartAngle: math.Pi / 4, EndAngle: 3 * math.Pi / 4})
	want := "M35.355,-35.355A50,50,0,0,1,35.355,35.355L14.142,14.142A20,20,0,0,0,14.142,-14.142Z"
	if got != want {
		t.Errorf("offset start:\n got=%q\nwant=%q", got, want)
	}
}

func TestArcSmallCorner(t *testing.T) {
	got := arc0().CornerRadius(10).Call(ArcDatum{InnerRadius: 0, OuterRadius: 50, StartAngle: 0, EndAngle: 0.3})
	want := "M0,-43.011A6.5,6.5,0,1,1,12.711,-41.09L0,0Z"
	if got != want {
		t.Errorf("small corner:\n got=%q\nwant=%q", got, want)
	}
}

// TestArcPieIntegration verifies arc + pie work together as nivo's pie uses them.
func TestArcPieIntegration(t *testing.T) {
	data := []float64{1, 1, 1}
	pie := NewPieNumeric().SortValues(nil)
	arcs := pie.Call(data)
	gen := NewArc().InnerRadius(func(d ArcDatum) float64 { return d.InnerRadius }).OuterRadius(func(d ArcDatum) float64 { return d.OuterRadius })
	for _, a := range arcs {
		path := gen.Call(ArcDatum{
			StartAngle:  a.StartAngle,
			EndAngle:    a.EndAngle,
			InnerRadius: 0,
			OuterRadius: 50,
		})
		if len(path) == 0 {
			t.Errorf("empty arc path for pie slice start=%g end=%g", a.StartAngle, a.EndAngle)
		}
		// Each slice is 120° → should start with M and contain an arc
		if path[0] != 'M' {
			t.Errorf("arc path should start with M, got %c", path[0])
		}
	}
}
