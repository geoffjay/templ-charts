package geo

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/interact"
	"github.com/geoffjay/templ-charts/charts/legends"
)

func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}

// applyGeoBase fills zero-valued shared props from commonDefaults.
func applyGeoBase(b GeoBase, defaultLayers []GeoLayerId) GeoBase {
	if b.ProjectionType == "" {
		b.ProjectionType = commonDefaults.ProjectionType
	}
	if b.ProjectionScale == 0 {
		b.ProjectionScale = commonDefaults.ProjectionScale
	}
	if b.ProjectionTranslation == [2]float64{0, 0} {
		b.ProjectionTranslation = commonDefaults.ProjectionTranslation
	}
	if b.GraticuleLineWidth == 0 {
		b.GraticuleLineWidth = commonDefaults.GraticuleLineWidth
	}
	if b.GraticuleLineColor == "" {
		b.GraticuleLineColor = commonDefaults.GraticuleLineColor
	}
	if b.BorderColor == "" {
		b.BorderColor = commonDefaults.BorderColor
	}
	if b.Role == "" {
		b.Role = commonDefaults.Role
	}
	if len(b.Layers) == 0 {
		b.Layers = defaultLayers
	}
	return b
}

func applyGeoMapDefaults(p GeoMapProps) GeoMapProps {
	p.GeoBase = applyGeoBase(p.GeoBase, DefaultGeoMapLayers)
	if p.FillColor == "" {
		p.FillColor = GeoMapDefaults.FillColor
	}
	return p
}

func applyChoroplethDefaults(p ChoroplethProps) ChoroplethProps {
	p.GeoBase = applyGeoBase(p.GeoBase, DefaultChoroplethLayers)
	if p.Colors == "" {
		p.Colors = ChoroplethDefaults.Colors
	}
	if p.Steps == 0 {
		p.Steps = ChoroplethDefaults.Steps
	}
	if p.UnknownColor == "" {
		p.UnknownColor = ChoroplethDefaults.UnknownColor
	}
	return p
}

// renderGeoMapLayers renders the enabled GeoMap layers as an inner SVG string.
func renderGeoMapLayers(props GeoMapProps, result GeoResult) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case GeoLayerGraticule:
			b.WriteString(renderGraticule(props.GeoBase, result))
		case GeoLayerFeatures:
			b.WriteString(renderFeatures(result.Features, props.Interactive, false, props.Animate, props.MotionStagger))
		}
	}
	return b.String()
}

// renderChoroplethLayers renders the enabled Choropleth layers, including the
// continuous-color legend.
func renderChoroplethLayers(props ChoroplethProps, result GeoResult, width, height float64) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case GeoLayerGraticule:
			b.WriteString(renderGraticule(props.GeoBase, result))
		case GeoLayerFeatures:
			b.WriteString(renderFeatures(result.Features, props.Interactive, true, props.Animate, props.MotionStagger))
		case GeoLayerLegends:
			b.WriteString(renderChoroplethLegends(props, result, width, height))
		}
	}
	return b.String()
}

func renderGraticule(b GeoBase, result GeoResult) string {
	if !b.EnableGraticule || result.GraticulePath == "" {
		return ""
	}
	var s strings.Builder
	s.WriteString(`<path fill="none" stroke-width="`)
	s.WriteString(fmtF(b.GraticuleLineWidth))
	s.WriteString(`" stroke="`)
	s.WriteString(b.GraticuleLineColor)
	s.WriteString(`" d="`)
	s.WriteString(result.GraticulePath)
	s.WriteString(`"></path>`)
	return s.String()
}

// renderFeatures emits one <path> per feature. withValue controls the tooltip
// content (Choropleth shows the value; GeoMap shows only the id).
func renderFeatures(features []ComputedFeature, interactive, withValue, animate bool, stagger float64) string {
	var b strings.Builder
	for i, f := range features {
		if f.Path == "" {
			continue
		}
		b.WriteString(`<path fill="`)
		b.WriteString(f.FillColor)
		b.WriteString(`" stroke-width="`)
		b.WriteString(fmtF(f.BorderWidth))
		b.WriteString(`" stroke="`)
		b.WriteString(f.BorderColor)
		b.WriteString(`" stroke-linejoin="bevel" d="`)
		b.WriteString(f.Path)
		b.WriteString(`"`)
		if interactive {
			label := f.ID
			value := ""
			if withValue && f.HasValue {
				value = f.FormattedValue
			}
			b.WriteString(` `)
			b.WriteString(interact.TooltipAttrName)
			b.WriteString(`="`)
			b.WriteString(templ.EscapeString(interact.TooltipHTML(f.FillColor, label, value)))
			b.WriteString(`" style="pointer-events:auto"`)
		} else {
			b.WriteString(` pointer-events="none"`)
		}
		b.WriteString(`>`)
		if animate {
			b.WriteString(core.SMILFadeIn(core.StaggerBegin(i, stagger)))
		}
		b.WriteString(`</path>`)
	}
	return b.String()
}

func renderChoroplethLegends(props ChoroplethProps, result GeoResult, width, height float64) string {
	if len(props.Legends) == 0 || result.ColorScale == nil {
		return ""
	}
	var b strings.Builder
	for _, lg := range props.Legends {
		b.WriteString(renderComponent(legends.ContinuousColorsLegendSvg(legends.ContinuousColorsLegendProps{
			Scale:       result.ColorScale,
			Min:         result.ValueMin,
			Max:         result.ValueMax,
			Width:       legendDim(lg.ItemWidth, 200),
			Height:      legendDim(lg.ItemHeight, 20),
			Anchor:      lg.Anchor,
			TranslateX:  lg.TranslateX,
			TranslateY:  lg.TranslateY,
			Samples:     props.Steps,
			ChartWidth:  width,
			ChartHeight: height,
		})))
	}
	return b.String()
}

func legendDim(v, fallback float64) float64 {
	if v == 0 {
		return fallback
	}
	return v
}

// fmtF formats a float rounded to 3 decimals (matching the d3-path serializer).
func fmtF(v float64) string {
	r := math.Round(v*1000) / 1000
	if r == 0 {
		r = 0 // normalize -0 to +0 for cross-platform-stable output
	}
	return strconv.FormatFloat(r, 'g', -1, 64)
}
