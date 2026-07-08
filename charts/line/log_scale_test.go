package line_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/scales"
)

var nonFinite = regexp.MustCompile(`NaN|Inf`)

// TestLine_LogYScale confirms the widened YScale field accepts a log spec end
// to end: an exponentially growing series renders a finite SVG that differs
// from the linear default.
func TestLine_LogYScale(t *testing.T) {
	data := []line.LineSeries{{ID: "s", Data: []line.LinePointData{
		{X: "a", Y: 10.0}, {X: "b", Y: 1000.0}, {X: "c", Y: 100000.0},
	}}}
	base := line.LineProps{Width: 700, Height: 400, Data: data}
	lin := renderLineComp(t, line.Line(base))

	logged := base
	logged.YScale = scales.ScaleLogSpec{Base: 10, Min: scales.AutoFloat(), Max: scales.AutoFloat()}
	out := renderLineComp(t, line.Line(logged))

	if nonFinite.MatchString(out) {
		t.Errorf("log render contains NaN/Inf: %q", out)
	}
	if out == lin {
		t.Error("log render should differ from linear render")
	}
}

// TestLine_NilYScaleDefaultsLinear guards the nil-check that replaced the old
// struct-equality default.
func TestLine_NilYScaleDefaultsLinear(t *testing.T) {
	out := renderLineComp(t, line.Line(line.LineProps{
		Width: 400, Height: 300,
		Data: []line.LineSeries{{ID: "s", Data: []line.LinePointData{{X: "a", Y: 1.0}, {X: "b", Y: 2.0}}}},
	}))
	if !strings.Contains(out, "<svg") {
		t.Errorf("nil YScale should render a chart, got %q", out)
	}
}
