// Package curves resolves the public curve and stack enum ids declared in
// charts/core to the concrete d3-shape factory functions in internal/d3/shape.
//
// It lives under internal so those d3-shape types never appear in the frozen
// public API of charts/core: consumers select a strategy with the public
// core.CurveFactoryId / core.StackOrder / core.StackOffset enums, and chart
// packages call these resolvers to obtain the internal generators.
package curves

import (
	"github.com/geoffjay/templ-charts/charts/core"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// CurvePropMapping maps core.CurveFactoryId → d3-shape CurveFactory. Mirrors
// nivo's curvePropMapping.
var CurvePropMapping = map[core.CurveFactoryId]d3shape.CurveFactory{
	core.CurveBasis:        d3shape.CurveBasis,
	core.CurveCardinal:     d3shape.CurveCardinal,
	core.CurveCatmullRom:   d3shape.CurveCatmullRom,
	core.CurveLinear:       d3shape.CurveLinear,
	core.CurveLinearClosed: d3shape.CurveLinearClosed,
	core.CurveMonotoneX:    d3shape.CurveMonotoneX,
	core.CurveMonotoneY:    d3shape.CurveMonotoneY,
	core.CurveNatural:      d3shape.CurveNatural,
	core.CurveStep:         d3shape.CurveStep,
	core.CurveStepAfter:    d3shape.CurveStepAfter,
	core.CurveStepBefore:   d3shape.CurveStepBefore,
	// Closed/open variants not in the current d3 port fall back to the nearest
	// equivalent so charts still render.
	core.CurveBasisClosed:      d3shape.CurveLinearClosed,
	core.CurveBasisOpen:        d3shape.CurveBasis,
	core.CurveBundle:           d3shape.CurveBasis,
	core.CurveCardinalClosed:   d3shape.CurveLinearClosed,
	core.CurveCardinalOpen:     d3shape.CurveCardinal,
	core.CurveCatmullRomClosed: d3shape.CurveLinearClosed,
	core.CurveCatmullRomOpen:   d3shape.CurveCatmullRom,
}

// CurveFromProp resolves a core.CurveFactoryId to a d3-shape CurveFactory.
// Returns CurveLinear if the id is unknown (matches nivo's throw fallback
// behavior without panicking — chart packages opt into validation).
func CurveFromProp(id core.CurveFactoryId) d3shape.CurveFactory {
	if f, ok := CurvePropMapping[id]; ok {
		return f
	}
	return d3shape.CurveLinear
}

// StackOrderPropMapping maps core.StackOrder → d3-shape StackOrderFunc.
var StackOrderPropMapping = map[core.StackOrder]d3shape.StackOrderFunc{
	core.StackOrderAscending:  d3shape.StackOrderAscending,
	core.StackOrderDescending: d3shape.StackOrderDescending,
	core.StackOrderInsideOut:  d3shape.StackOrderInsideOut,
	core.StackOrderNone:       d3shape.StackOrderNone,
	core.StackOrderReverse:    d3shape.StackOrderReverse,
}

// StackOrderFromProp resolves a core.StackOrder. Defaults to StackOrderNone.
func StackOrderFromProp(p core.StackOrder) d3shape.StackOrderFunc {
	if f, ok := StackOrderPropMapping[p]; ok {
		return f
	}
	return d3shape.StackOrderNone
}

// StackOffsetPropMapping maps core.StackOffset → d3-shape StackOffsetFunc.
var StackOffsetPropMapping = map[core.StackOffset]d3shape.StackOffsetFunc{
	core.StackOffsetExpand:     d3shape.StackOffsetExpand,
	core.StackOffsetDiverging:  d3shape.StackOffsetDiverging,
	core.StackOffsetNone:       d3shape.StackOffsetNone,
	core.StackOffsetSilhouette: d3shape.StackOffsetSilhouette,
	core.StackOffsetWiggle:     d3shape.StackOffsetWiggle,
}

// StackOffsetFromProp resolves a core.StackOffset. Defaults to StackOffsetNone.
func StackOffsetFromProp(p core.StackOffset) d3shape.StackOffsetFunc {
	if f, ok := StackOffsetPropMapping[p]; ok {
		return f
	}
	return d3shape.StackOffsetNone
}
