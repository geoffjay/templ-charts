package core

import (
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// CurveFactoryId enumerates the d3-shape curve interpolators nivo exposes
// via the `curve` prop. Mirrors @nivo/core props/curve.js.
type CurveFactoryId string

const (
	CurveBasis            CurveFactoryId = "basis"
	CurveBasisClosed      CurveFactoryId = "basisClosed"
	CurveBasisOpen        CurveFactoryId = "basisOpen"
	CurveBundle           CurveFactoryId = "bundle"
	CurveCardinal         CurveFactoryId = "cardinal"
	CurveCardinalClosed   CurveFactoryId = "cardinalClosed"
	CurveCardinalOpen     CurveFactoryId = "cardinalOpen"
	CurveCatmullRom       CurveFactoryId = "catmullRom"
	CurveCatmullRomClosed CurveFactoryId = "catmullRomClosed"
	CurveCatmullRomOpen   CurveFactoryId = "catmullRomOpen"
	CurveLinear           CurveFactoryId = "linear"
	CurveLinearClosed     CurveFactoryId = "linearClosed"
	CurveMonotoneX        CurveFactoryId = "monotoneX"
	CurveMonotoneY        CurveFactoryId = "monotoneY"
	CurveNatural          CurveFactoryId = "natural"
	CurveStep             CurveFactoryId = "step"
	CurveStepAfter        CurveFactoryId = "stepAfter"
	CurveStepBefore       CurveFactoryId = "stepBefore"
)

// CurvePropMapping maps CurveFactoryId → d3-shape CurveFactory. Mirrors nivo's
// curvePropMapping.
var CurvePropMapping = map[CurveFactoryId]d3shape.CurveFactory{
	CurveBasis:        d3shape.CurveBasis,
	CurveCardinal:     d3shape.CurveCardinal,
	CurveCatmullRom:   d3shape.CurveCatmullRom,
	CurveLinear:       d3shape.CurveLinear,
	CurveLinearClosed: d3shape.CurveLinearClosed,
	CurveMonotoneX:    d3shape.CurveMonotoneX,
	CurveMonotoneY:    d3shape.CurveMonotoneY,
	CurveNatural:      d3shape.CurveNatural,
	CurveStep:         d3shape.CurveStep,
	CurveStepAfter:    d3shape.CurveStepAfter,
	CurveStepBefore:   d3shape.CurveStepBefore,
	// Closed/open variants not in the current d3 port fall back to the
	// nearest equivalent so charts still render.
	CurveBasisClosed:      d3shape.CurveLinearClosed,
	CurveBasisOpen:        d3shape.CurveBasis,
	CurveBundle:           d3shape.CurveBasis,
	CurveCardinalClosed:   d3shape.CurveLinearClosed,
	CurveCardinalOpen:     d3shape.CurveCardinal,
	CurveCatmullRomClosed: d3shape.CurveLinearClosed,
	CurveCatmullRomOpen:   d3shape.CurveCatmullRom,
}

// ClosedCurvePropKeys lists the curve ids ending in "Closed".
var ClosedCurvePropKeys = []CurveFactoryId{
	CurveBasisClosed, CurveCardinalClosed, CurveCatmullRomClosed, CurveLinearClosed,
}

// AreaCurvePropKeys are safe curves for d3's area generator (excludes
// closed/open variants unsupported by area).
var AreaCurvePropKeys = []CurveFactoryId{
	CurveBasis, CurveBundle, CurveCardinal, CurveCatmullRom, CurveLinear,
	CurveMonotoneX, CurveMonotoneY, CurveNatural, CurveStep, CurveStepAfter, CurveStepBefore,
}

// LineCurvePropKeys are safe curves for d3's line generator.
var LineCurvePropKeys = AreaCurvePropKeys

// CurveFromProp resolves a CurveFactoryId to a d3-shape CurveFactory.
// Returns CurveLinear if the id is unknown (matches nivo's throw fallback
// behavior without panicking — chart packages opt into validation).
func CurveFromProp(id CurveFactoryId) d3shape.CurveFactory {
	if f, ok := CurvePropMapping[id]; ok {
		return f
	}
	return d3shape.CurveLinear
}

// StackOrder enumerates d3-shape stack order strategies.
type StackOrder string

const (
	StackOrderAscending  StackOrder = "ascending"
	StackOrderDescending StackOrder = "descending"
	StackOrderInsideOut  StackOrder = "insideOut"
	StackOrderNone       StackOrder = "none"
	StackOrderReverse    StackOrder = "reverse"
)

// StackOrderPropMapping maps StackOrder → d3-shape StackOrderFunc.
var StackOrderPropMapping = map[StackOrder]d3shape.StackOrderFunc{
	StackOrderAscending:  d3shape.StackOrderAscending,
	StackOrderDescending: d3shape.StackOrderDescending,
	StackOrderInsideOut:  d3shape.StackOrderInsideOut,
	StackOrderNone:       d3shape.StackOrderNone,
	StackOrderReverse:    d3shape.StackOrderReverse,
}

// StackOrderFromProp resolves a StackOrder. Defaults to StackOrderNone.
func StackOrderFromProp(p StackOrder) d3shape.StackOrderFunc {
	if f, ok := StackOrderPropMapping[p]; ok {
		return f
	}
	return d3shape.StackOrderNone
}

// StackOffset enumerates d3-shape stack offset strategies.
type StackOffset string

const (
	StackOffsetExpand     StackOffset = "expand"
	StackOffsetDiverging  StackOffset = "diverging"
	StackOffsetNone       StackOffset = "none"
	StackOffsetSilhouette StackOffset = "silhouette"
	StackOffsetWiggle     StackOffset = "wiggle"
)

// StackOffsetPropMapping maps StackOffset → d3-shape StackOffsetFunc.
var StackOffsetPropMapping = map[StackOffset]d3shape.StackOffsetFunc{
	StackOffsetExpand:     d3shape.StackOffsetExpand,
	StackOffsetDiverging:  d3shape.StackOffsetDiverging,
	StackOffsetNone:       d3shape.StackOffsetNone,
	StackOffsetSilhouette: d3shape.StackOffsetSilhouette,
	StackOffsetWiggle:     d3shape.StackOffsetWiggle,
}

// StackOffsetFromProp resolves a StackOffset. Defaults to StackOffsetNone.
func StackOffsetFromProp(p StackOffset) d3shape.StackOffsetFunc {
	if f, ok := StackOffsetPropMapping[p]; ok {
		return f
	}
	return d3shape.StackOffsetNone
}

// CssMixBlendMode enumerates CSS mix-blend-mode values nivo exposes.
type CssMixBlendMode string

const (
	MixBlendNormal     CssMixBlendMode = "normal"
	MixBlendMultiply   CssMixBlendMode = "multiply"
	MixBlendScreen     CssMixBlendMode = "screen"
	MixBlendOverlay    CssMixBlendMode = "overlay"
	MixBlendDarken     CssMixBlendMode = "darken"
	MixBlendLighten    CssMixBlendMode = "lighten"
	MixBlendColorDodge CssMixBlendMode = "color-dodge"
	MixBlendColorBurn  CssMixBlendMode = "color-burn"
	MixBlendHardLight  CssMixBlendMode = "hard-light"
	MixBlendSoftLight  CssMixBlendMode = "soft-light"
	MixBlendDifference CssMixBlendMode = "difference"
	MixBlendExclusion  CssMixBlendMode = "exclusion"
	MixBlendHue        CssMixBlendMode = "hue"
	MixBlendSaturation CssMixBlendMode = "saturation"
	MixBlendColor      CssMixBlendMode = "color"
	MixBlendLuminosity CssMixBlendMode = "luminosity"
)
