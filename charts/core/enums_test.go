package core_test

import (
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
)

func TestUseDimensions(t *testing.T) {
	d := core.UseDimensions(400, 300, core.Margin{Top: 10, Left: 20})
	if d.OuterWidth != 400 || d.OuterHeight != 300 {
		t.Errorf("outer = %v x %v", d.OuterWidth, d.OuterHeight)
	}
	if d.InnerWidth != 380 || d.InnerHeight != 290 {
		t.Errorf("inner = %v x %v, want 380 x 290", d.InnerWidth, d.InnerHeight)
	}
	if d.Margin != (core.Margin{Top: 10, Left: 20}) {
		t.Errorf("margin = %+v", d.Margin)
	}

	full := core.UseDimensions(100, 100, core.Margin{Top: 1, Right: 2, Bottom: 3, Left: 4})
	if full.InnerWidth != 94 || full.InnerHeight != 96 {
		t.Errorf("full margin inner = %v x %v, want 94 x 96", full.InnerWidth, full.InnerHeight)
	}
}
