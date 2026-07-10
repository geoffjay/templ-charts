package curves_test

import (
	"strings"
	"testing"

	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/internal/curves"
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
		if curves.CurveFromProp(id) == nil {
			t.Errorf("CurveFromProp(%q) = nil", id)
		}
		if _, ok := curves.CurvePropMapping[id]; !ok {
			t.Errorf("CurvePropMapping missing %q", id)
		}
	}
	if curves.CurveFromProp("bogus") == nil {
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
		if _, ok := curves.CurvePropMapping[id]; !ok {
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
		if curves.StackOrderFromProp(o) == nil {
			t.Errorf("StackOrderFromProp(%q) = nil", o)
		}
	}
	if curves.StackOrderFromProp("bogus") == nil {
		t.Errorf("unknown order should fall back to none, got nil")
	}
}

func TestStackOffsetFromProp(t *testing.T) {
	offsets := []core.StackOffset{
		core.StackOffsetExpand, core.StackOffsetDiverging, core.StackOffsetNone,
		core.StackOffsetSilhouette, core.StackOffsetWiggle,
	}
	for _, o := range offsets {
		if curves.StackOffsetFromProp(o) == nil {
			t.Errorf("StackOffsetFromProp(%q) = nil", o)
		}
	}
	if curves.StackOffsetFromProp("bogus") == nil {
		t.Errorf("unknown offset should fall back to none, got nil")
	}
}
