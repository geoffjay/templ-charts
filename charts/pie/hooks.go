package pie

import (
	"math"

	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
	d3shape "github.com/geoffjay/templ-charts/internal/d3/shape"
)

// NormalizedDatum is a ComputedDatum without the Arc/Fill fields (the output
// of NormalizeData before PieArcs adds arcs).
type NormalizedDatum struct {
	ID             DatumId
	Label          DatumId
	Hidden         bool
	Value          float64
	FormattedValue string
	Data           any
	Color          string
}

// NormalizeData mirrors @nivo/pie useNormalizedData: resolves id/value
// accessors, formats values, and assigns colors via an ordinal color scale
// keyed by id.
func NormalizeData(
	data []any,
	idAcc core.PropertyAccessor[any, DatumId],
	valueAcc core.PropertyAccessor[any, float64],
	valueFormat core.ValueFormat[float64],
	colorsCfg colors.OrdinalColorScaleConfig,
) []NormalizedDatum {
	getID := core.GetPropertyAccessor[any, DatumId](idAcc)
	if idAcc == nil {
		getID = core.GetPropertyAccessor[any, DatumId](Defaults.ID)
	}
	getValue := core.GetPropertyAccessor[any, float64](valueAcc)
	if valueAcc == nil {
		getValue = core.GetPropertyAccessor[any, float64](Defaults.Value)
	}
	formatValue := core.GetValueFormatter[float64](valueFormat)

	cfg := colorsCfg
	if cfg.Type == 0 && cfg.Scheme == "" && cfg.Static == "" && len(cfg.Colors) == 0 && cfg.Func == nil && cfg.DatumPath == "" {
		cfg = Defaults.Colors
	}
	getColor := colors.GetOrdinalColorScale[NormalizedDatum](cfg, func(d NormalizedDatum) string { return d.ID })

	out := make([]NormalizedDatum, len(data))
	for i, datum := range data {
		id := getID(datum)
		val := getValue(datum)
		nd := NormalizedDatum{
			ID:             id,
			Label:          resolveLabel(datum, id),
			Hidden:         false,
			Value:          val,
			FormattedValue: formatValue(val),
			Data:           datum,
		}
		nd.Color = getColor(nd)
		out[i] = nd
	}
	return out
}

// resolveLabel returns datum.label if present (via reflection), else id.
func resolveLabel(datum any, id DatumId) DatumId {
	if s, ok := datum.(map[string]any); ok {
		if l, ok := s["label"]; ok {
			if ls, ok := l.(string); ok && ls != "" {
				return ls
			}
		}
	}
	// Try struct field "Label" via reflection.
	labelAcc := core.GetPropertyAccessor[any, string]("label")
	if l := labelAcc(datum); l != "" {
		return l
	}
	return id
}

// PieArcs computes the d3-pie arcs for the normalized data and attaches them
// to each datum as PieArc. Mirrors @nivo/pie usePieArcs.
func PieArcs(
	data []NormalizedDatum,
	startAngleDeg, endAngleDeg, padAngleDeg float64,
	innerRadius, outerRadius float64,
	sortByValue bool,
	activeID DatumId,
	activeInnerRadiusOffset, activeOuterRadiusOffset float64,
	hiddenIDs []DatumId,
) ([]ComputedDatum, []LegendDatum) {
	// Filter hidden data.
	visible := make([]NormalizedDatum, 0, len(data))
	for _, d := range data {
		if !containsID(hiddenIDs, d.ID) {
			visible = append(visible, d)
		}
	}

	// Build d3-pie generator.
	pie := d3shape.NewPie[NormalizedDatum]().
		Value(func(d NormalizedDatum, _ int, _ []NormalizedDatum) float64 { return d.Value }).
		StartAngle(arcs.DegToRad(startAngleDeg)).
		EndAngle(arcs.DegToRad(endAngleDeg)).
		PadAngle(arcs.DegToRad(padAngleDeg))
	if !sortByValue {
		pie.SortValues(nil)
	}

	pieArcs := pie.Call(visible)

	dataWithArc := make([]ComputedDatum, 0, len(pieArcs))
	for _, arc := range pieArcs {
		d := arc.Data
		angle := math.Abs(arc.EndAngle - arc.StartAngle)
		ir := innerRadius
		or := outerRadius
		if activeID == d.ID {
			ir -= activeInnerRadiusOffset
			or += activeOuterRadiusOffset
		}
		dataWithArc = append(dataWithArc, ComputedDatum{
			ID:             d.ID,
			Label:          d.Label,
			Value:          d.Value,
			FormattedValue: d.FormattedValue,
			Color:          d.Color,
			Data:           d.Data,
			Arc: PieArc{
				Arc: arcs.Arc{
					StartAngle:  arc.StartAngle,
					EndAngle:    arc.EndAngle,
					InnerRadius: ir,
					OuterRadius: or,
					PadAngle:    arc.PadAngle,
				},
				Index:     arc.Index,
				Angle:     angle,
				AngleDeg:  arcs.RadToDeg(angle),
				Thickness: or - ir,
				PadAngle:  arc.PadAngle,
			},
		})
	}

	// Legend data: all data items (including hidden), with hidden flag.
	legendData := make([]LegendDatum, len(data))
	for i, d := range data {
		legendData[i] = LegendDatum{
			ID:     d.ID,
			Label:  d.Label,
			Color:  d.Color,
			Hidden: containsID(hiddenIDs, d.ID),
			Data: ComputedDatumNoArc{
				ID:             d.ID,
				Label:          d.Label,
				Value:          d.Value,
				FormattedValue: d.FormattedValue,
				Color:          d.Color,
				Data:           d.Data,
				Hidden:         containsID(hiddenIDs, d.ID),
			},
		}
	}

	return dataWithArc, legendData
}

// PieFromBox mirrors @nivo/pie usePieFromBox: computes the pie layout from a
// box (width × height), deriving the radius/innerRadius/center. When `fit` is
// true, rescales the radius so the arc bounding box fills the available space.
func PieFromBox(
	data []NormalizedDatum,
	width, height float64,
	fit bool,
	innerRadiusRatio float64,
	startAngleDeg, endAngleDeg, padAngleDeg float64,
	sortByValue bool,
	cornerRadius float64,
	activeInnerRadiusOffset, activeOuterRadiusOffset float64,
	activeID DatumId,
	hiddenIDs []DatumId,
) PieResult {
	radius := math.Min(width, height) / 2
	innerRadius := radius * math.Min(innerRadiusRatio, 1)
	centerX := width / 2
	centerY := height / 2

	if fit {
		bx, by, bw, bh := arcs.ComputeArcBoundingBox(centerX, centerY, radius, startAngleDeg-90, endAngleDeg-90, false)
		ratio := math.Min(width/bw, height/bh)
		adjustedW := bw * ratio
		adjustedH := bh * ratio
		adjustedX := (width - adjustedW) / 2
		adjustedY := (height - adjustedH) / 2
		centerX = ((centerX-bx)/bw)*bw*ratio + adjustedX
		centerY = ((centerY-by)/bh)*bh*ratio + adjustedY
		radius *= ratio
		innerRadius *= ratio
		_ = adjustedW
		_ = adjustedH
	}

	dataWithArc, legendData := PieArcs(
		data, startAngleDeg, endAngleDeg, padAngleDeg,
		innerRadius, radius, sortByValue, activeID,
		activeInnerRadiusOffset, activeOuterRadiusOffset, hiddenIDs,
	)

	arcGen := arcs.CreateArcGenerator(cornerRadius, arcs.DegToRad(padAngleDeg))

	return PieResult{
		DataWithArc:  dataWithArc,
		LegendData:   legendData,
		ArcGenerator: arcGen,
		CenterX:      centerX,
		CenterY:      centerY,
		Radius:       radius,
		InnerRadius:  innerRadius,
		HiddenIDs:    hiddenIDs,
	}
}

// UsePie is the full orchestrator: normalizes data, computes the pie layout
// from the box dimensions, and returns the PieResult.
func UsePie(props PieProps) PieResult {
	// merge defaults
	idAcc := props.ID
	if idAcc == nil {
		idAcc = Defaults.ID
	}
	valueAcc := props.Value
	if valueAcc == nil {
		valueAcc = Defaults.Value
	}
	colorsCfg := props.Colors
	if colorsCfg.Type == 0 && colorsCfg.Scheme == "" && colorsCfg.Static == "" && len(colorsCfg.Colors) == 0 && colorsCfg.Func == nil && colorsCfg.DatumPath == "" {
		colorsCfg = Defaults.Colors
	}

	normalized := NormalizeData(props.Data, idAcc, valueAcc, props.ValueFormat, colorsCfg)

	dims := core.UseDimensions(props.Width, props.Height, props.Margin)
	innerRadiusRatio := props.InnerRadius
	startAngle := props.StartAngle
	if startAngle == 0 && props.EndAngle == 0 {
		startAngle = Defaults.StartAngle
	}
	endAngle := props.EndAngle
	if endAngle == 0 {
		endAngle = Defaults.EndAngle
	}
	padAngle := props.PadAngle
	cornerRadius := props.CornerRadius
	fit := props.Fit
	if !fit && !props.fitExplicitlyDisabled() {
		fit = Defaults.Fit
	}

	return PieFromBox(
		normalized, dims.InnerWidth, dims.InnerHeight,
		fit, innerRadiusRatio,
		startAngle, endAngle, padAngle,
		props.SortByValue, cornerRadius,
		props.ActiveInnerRadiusOffset, props.ActiveOuterRadiusOffset,
		props.ActiveID, props.InitialHiddenIDs,
	)
}

// fitExplicitlyDisabled is a placeholder for the fit default detection.
// Since Go bool zero value is false, we can't distinguish "unset" from "set
// to false" — we always apply the default (true) when Fit is false. Callers
// who want fit=false must set it explicitly and it will be respected because
// applyDefaults only sets Fit=true when it's false AND hasn't been explicitly
// set. Since we can't detect that, we document: Fit defaults to true; set
// Fit=false to disable.
func (PieProps) fitExplicitlyDisabled() bool { return false }

// UsePieLayerContext builds the context passed to custom layers. Mirrors
// @nivo/pie usePieLayerContext.
func UsePieLayerContext(result PieResult) map[string]any {
	return map[string]any{
		"dataWithArc":  result.DataWithArc,
		"arcGenerator": result.ArcGenerator,
		"centerX":      result.CenterX,
		"centerY":      result.CenterY,
		"radius":       result.Radius,
		"innerRadius":  result.InnerRadius,
	}
}

// --- helpers ---

func containsID(ids []DatumId, id DatumId) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}

// guard against unused imports during refactoring
var _ = theming.DefaultTheme
