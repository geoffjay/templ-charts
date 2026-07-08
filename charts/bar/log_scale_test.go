package bar_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/scales"
)

var nonFinite = regexp.MustCompile(`NaN|Inf`)

// TestBar_LogValueScale confirms the widened ValueScale field accepts a log
// spec end to end: the chart renders a finite SVG that differs from the linear
// default, and the small values produce visible (nonzero-height) bars where a
// linear axis would flatten them.
func TestBar_LogValueScale(t *testing.T) {
	data := []bar.BarDatum{
		{"tier": "cache", "requests": 1200000.0},
		{"tier": "db", "requests": 3800.0},
		{"tier": "audit", "requests": 12.0},
	}
	base := bar.BarProps{
		Width: 700, Height: 400, IndexBy: "tier",
		Keys: []string{"requests"}, Data: data,
		GroupMode: bar.GroupModeGrouped,
	}
	lin := renderComp(t, bar.Bar(base))

	logged := base
	logged.ValueScale = scales.ScaleLogSpec{Base: 10, Min: scales.FloatVal(1), Max: scales.AutoFloat()}
	out := renderComp(t, bar.Bar(logged))

	if nonFinite.MatchString(out) {
		t.Errorf("log render contains NaN/Inf: %q", out)
	}
	if out == lin {
		t.Error("log render should differ from linear render")
	}

	// Every data bar must have a positive height (the log baseline keeps the
	// small tiers visible instead of collapsing them to the zero line).
	heights := regexp.MustCompile(`<rect[^>]*\sheight="([\d.]+)"`).FindAllStringSubmatch(out, -1)
	positive := 0
	for _, m := range heights {
		if m[1] != "0" && !strings.HasPrefix(m[1], "0.0") {
			positive++
		}
	}
	if positive < 3 {
		t.Errorf("expected >=3 visible bars under log, got %d (heights %v)", positive, heights)
	}
}

// TestBar_NilValueScaleDefaultsLinear guards the nil-check that replaced the
// old struct-equality default: a zero-value (nil) ValueScale still renders.
func TestBar_NilValueScaleDefaultsLinear(t *testing.T) {
	out := renderComp(t, bar.Bar(bar.BarProps{
		Width: 400, Height: 300, IndexBy: "tier",
		Keys: []string{"requests"},
		Data: []bar.BarDatum{{"tier": "a", "requests": 10.0}, {"tier": "b", "requests": 20.0}},
	}))
	if !strings.Contains(out, "<svg") {
		t.Errorf("nil ValueScale should render a chart, got %q", out)
	}
}
