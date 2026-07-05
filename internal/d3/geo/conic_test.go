package geo

import (
	"testing"

	"github.com/geoffjay/templ-charts/internal/golden"
)

// Conic projection tests. Projected points were produced by d3-geo@3's
// geoConicConformal/geoConicEqualArea/geoConicEquidistant with parallels [20,50].

func TestConicProjections_Project(t *testing.T) {
	cases := []struct {
		name string
		p    *Projection
		// proj(10,20) and proj(-30,45)
		ax, ay, bx, by float64
	}{
		{"conicConformal", GeoConicConformal(), 504.559, 192.9498, 426.4147, 122.2004},
		{"conicEqualArea", GeoConicEqualArea(), 504.5628, 199.0771, 426.4309, 125.3052},
		{"conicEquidistant", GeoConicEquidistant(), 504.561, 196.4238, 426.4244, 124.1779},
	}
	for _, c := range cases {
		p := c.p.Parallels(20, 50).Scale(150).Translate(480, 250)
		ax, ay := p.Project(10, 20)
		bx, by := p.Project(-30, 45)
		approx(t, c.name+".a.x", ax, c.ax)
		approx(t, c.name+".a.y", ay, c.ay)
		approx(t, c.name+".b.x", bx, c.bx)
		approx(t, c.name+".b.y", by, c.by)
	}
}

func TestConicProjections_ByType(t *testing.T) {
	for _, name := range []string{"conicConformal", "conicEqualArea", "conicEquidistant"} {
		p := ProjectionByType(name)
		if p.rawFactory == nil {
			t.Errorf("ProjectionByType(%q) is not a conic projection", name)
		}
	}
}

func TestConic_ParallelsRoundTrip(t *testing.T) {
	p := GeoConicConformal().Parallels(15, 45)
	p0, p1 := p.ParallelsValue()
	approx(t, "parallels.0", p0, 15)
	approx(t, "parallels.1", p1, 45)
}

// TestConicPathGolden pins the rendered path strings for the conic family (drift
// guard); the numbers themselves are validated against d3 in the Project test.
func TestConicPathGolden(t *testing.T) {
	for _, name := range []string{"conicConformal", "conicEqualArea", "conicEquidistant"} {
		p := ProjectionByType(name).Parallels(20, 50).Scale(150).Translate(480, 250).Precision(0)
		golden.Assert(t, "path-"+name, NewPath(p).Geometry(diamond()))
	}
}
