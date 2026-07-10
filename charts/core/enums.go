package core

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

// StackOrder enumerates d3-shape stack order strategies.
type StackOrder string

const (
	StackOrderAscending  StackOrder = "ascending"
	StackOrderDescending StackOrder = "descending"
	StackOrderInsideOut  StackOrder = "insideOut"
	StackOrderNone       StackOrder = "none"
	StackOrderReverse    StackOrder = "reverse"
)

// StackOffset enumerates d3-shape stack offset strategies.
type StackOffset string

const (
	StackOffsetExpand     StackOffset = "expand"
	StackOffsetDiverging  StackOffset = "diverging"
	StackOffsetNone       StackOffset = "none"
	StackOffsetSilhouette StackOffset = "silhouette"
	StackOffsetWiggle     StackOffset = "wiggle"
)

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
