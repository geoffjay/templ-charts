package core_test

import (
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
)

func TestCurveFromProp(t *testing.T) {
	// Every documented curve id must resolve, including the closed/open
	// variants that fall back to a nearest equivalent.
	ids := []core.CurveFactoryId{
		core.CurveBasis, core.CurveBasisClosed, core.CurveBasisOpen, core.CurveBundle,
		core.CurveCardinal, core.CurveCardinalClosed, core.CurveCardinalOpen,
		core.CurveCatmullRom, core.CurveCatmullRomClosed, core.CurveCatmullRomOpen,
		core.CurveLinear, core.CurveLinearClosed, core.CurveMonotoneX, core.CurveMonotoneY,
		core.CurveNatural, core.CurveStep, core.CurveStepAfter, core.CurveStepBefore,
	}
	for _, id := range ids {
		if core.CurveFromProp(id) == nil {
			t.Errorf("CurveFromProp(%q) = nil", id)
		}
		if _, ok := core.CurvePropMapping[id]; !ok {
			t.Errorf("CurvePropMapping missing %q", id)
		}
	}
	if core.CurveFromProp("bogus") == nil {
		t.Errorf("unknown curve should fall back to linear, got nil")
	}
}

func TestClosedAndAreaCurveKeys(t *testing.T) {
	for _, id := range core.ClosedCurvePropKeys {
		if !strings.HasSuffix(string(id), "Closed") {
			t.Errorf("ClosedCurvePropKeys contains non-closed id %q", id)
		}
	}
	for _, id := range core.AreaCurvePropKeys {
		if _, ok := core.CurvePropMapping[id]; !ok {
			t.Errorf("AreaCurvePropKeys id %q missing from mapping", id)
		}
	}
}

func TestStackOrderFromProp(t *testing.T) {
	orders := []core.StackOrder{
		core.StackOrderAscending, core.StackOrderDescending, core.StackOrderInsideOut,
		core.StackOrderNone, core.StackOrderReverse,
	}
	for _, o := range orders {
		if core.StackOrderFromProp(o) == nil {
			t.Errorf("StackOrderFromProp(%q) = nil", o)
		}
	}
	if core.StackOrderFromProp("bogus") == nil {
		t.Errorf("unknown order should fall back to none, got nil")
	}
}

func TestStackOffsetFromProp(t *testing.T) {
	offsets := []core.StackOffset{
		core.StackOffsetExpand, core.StackOffsetDiverging, core.StackOffsetNone,
		core.StackOffsetSilhouette, core.StackOffsetWiggle,
	}
	for _, o := range offsets {
		if core.StackOffsetFromProp(o) == nil {
			t.Errorf("StackOffsetFromProp(%q) = nil", o)
		}
	}
	if core.StackOffsetFromProp("bogus") == nil {
		t.Errorf("unknown offset should fall back to none, got nil")
	}
}

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
