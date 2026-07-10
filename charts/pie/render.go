package pie

import (
	"context"
	"strings"

	"github.com/a-h/templ"

	"github.com/geoffjay/templ-charts/charts/arcs"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/legends"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// renderComponent renders a templ.Component to a string.
func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}

// renderArcs renders the Arcs component to a string.
func renderArcs(props ArcsProps) string {
	return renderComponent(Arcs(props))
}

// renderArcLinkLabels renders ArcLinkLabels to a string.
func renderArcLinkLabels(props ArcLinkLabelsProps) string {
	return renderComponent(ArcLinkLabels(props))
}

// renderArcLabels renders ArcLabels to a string.
func renderArcLabels(props ArcLabelsProps) string {
	return renderComponent(ArcLabels(props))
}

// renderBoxLegend renders a legends.BoxLegendSvg to a string.
func renderBoxLegend(props legends.BoxLegendSvgProps) string {
	return renderComponent(legends.BoxLegendSvg(props))
}

// resolveTheme returns props.Theme or the default theme when nil.
func resolveTheme(t *theming.Theme) *theming.Theme {
	if t == nil {
		return &theming.DefaultTheme
	}
	return t
}

// themeBackground returns the theme background color.
func themeBackground(t *theming.Theme) string {
	if t == nil {
		return ""
	}
	return t.Background
}

// themeLabelFontSize extracts the labels.text.fontSize as float64.
func themeLabelFontSize(t *theming.Theme) float64 {
	if t == nil {
		return 0
	}
	return toFloatAny(t.Labels.Text.FontSize)
}

func toFloatAny(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	case int64:
		return float64(x)
	}
	return 0
}

// renderPieLayers renders the enabled layers as an SVG string.
func renderPieLayers(layers []PieLayerId, props PieProps, result PieResult, dims core.Dimensions, theme *theming.Theme, bound core.SvgDefsAndFill) string {
	var b strings.Builder
	for _, layer := range layers {
		switch layer {
		case PieLayerArcs:
			b.WriteString(renderArcsLayer(props, result, theme, bound))
		case PieLayerArcLinkLabels:
			if props.ArcLinkLabelsEnabled() {
				b.WriteString(renderArcLinkLabelsLayer(props, result, theme))
			}
		case PieLayerArcLabels:
			if props.ArcLabelsEnabled() {
				b.WriteString(renderArcLabelsLayer(props, result, theme))
			}
		case PieLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, dims))
		}
	}
	return b.String()
}

func renderArcsLayer(props PieProps, result PieResult, theme *theming.Theme, bound core.SvgDefsAndFill) string {
	borderColorFn := resolveBorderColor(props, theme)
	// Apply fill overrides from bound defs.
	data := result.DataWithArc
	if len(bound.FillByNodeIndex) > 0 {
		for i := range data {
			if fill, ok := bound.FillByNodeIndex[i]; ok {
				data[i].Fill = fill
			}
		}
	}
	return renderArcs(ArcsProps{
		CenterX:        result.CenterX,
		CenterY:        result.CenterY,
		Data:           data,
		ArcGenerator:   result.ArcGenerator,
		BorderWidth:    props.BorderWidth,
		BorderColor:    borderColorFn,
		Interactive:    props.Interactive,
		ChartID:        props.ChartID,
		Animate:        props.Animate,
		TransitionMode: props.TransitionMode,
	})
}

func renderArcLinkLabelsLayer(props PieProps, result PieResult, theme *theming.Theme) string {
	textColor := resolveInheritedColor(props.ArcLinkLabelsTextColor, Defaults.ArcLinkLabelsTextColor, theme)
	linkColor := resolveInheritedColor(props.ArcLinkLabelsColor, Defaults.ArcLinkLabelsColor, theme)
	getLabel := core.GetPropertyAccessor[ComputedDatum, string](props.ArcLinkLabel)
	if props.ArcLinkLabel == nil {
		getLabel = core.GetPropertyAccessor[ComputedDatum, string](Defaults.ArcLinkLabel)
	}
	return renderArcLinkLabels(ArcLinkLabelsProps{
		CenterX:        result.CenterX,
		CenterY:        result.CenterY,
		Data:           result.DataWithArc,
		GetLabel:       getLabel,
		SkipAngle:      props.ArcLinkLabelsSkipAngle,
		Offset:         props.ArcLinkLabelsOffset,
		DiagonalLength: props.ArcLinkLabelsDiagonalLength,
		StraightLength: props.ArcLinkLabelsStraightLength,
		Thickness:      props.ArcLinkLabelsThickness,
		TextOffset:     props.ArcLinkLabelsTextOffset,
		TextColor:      textColor,
		LinkColor:      linkColor,
		FontSize:       themeLabelFontSize(theme),
		FontFamily:     theme.Labels.Text.FontFamily,
	})
}

func renderArcLabelsLayer(props PieProps, result PieResult, theme *theming.Theme) string {
	textColor := resolveInheritedColor(props.ArcLabelsTextColor, Defaults.ArcLabelsTextColor, theme)
	getLabel := core.GetPropertyAccessor[ComputedDatum, string](props.ArcLabel)
	if props.ArcLabel == nil {
		getLabel = core.GetPropertyAccessor[ComputedDatum, string](Defaults.ArcLabel)
	}
	return renderArcLabels(ArcLabelsProps{
		CenterX:      result.CenterX,
		CenterY:      result.CenterY,
		Data:         result.DataWithArc,
		GetLabel:     getLabel,
		RadiusOffset: props.ArcLabelsRadiusOffset,
		SkipAngle:    props.ArcLabelsSkipAngle,
		SkipRadius:   props.ArcLabelsSkipRadius,
		TextColor:    textColor,
		FontSize:     themeLabelFontSize(theme),
		FontFamily:   theme.Labels.Text.FontFamily,
	})
}

func renderLegendsLayer(props PieProps, result PieResult, dims core.Dimensions) string {
	if len(props.Legends) == 0 {
		return ""
	}
	var b strings.Builder
	for _, legend := range props.Legends {
		lp := legend
		if len(lp.Items) == 0 {
			items := make([]legends.Datum, len(result.LegendData))
			for i, d := range result.LegendData {
				items[i] = legends.Datum{ID: d.ID, Label: d.Label, Color: d.Color, Hidden: d.Hidden}
			}
			lp.Items = items
		}
		if lp.ChartID == "" {
			lp.ChartID = props.ChartID
		}
		if lp.Toggle == "" && props.ChartID != "" {
			lp.Toggle = "1"
		}
		b.WriteString(renderBoxLegend(legends.BoxLegendSvgProps{
			Props:       lp,
			ChartWidth:  dims.InnerWidth,
			ChartHeight: dims.InnerHeight,
		}))
	}
	return b.String()
}

// resolveBorderColor builds the border-color function from the inherited
// color config + theme.
func resolveBorderColor(props PieProps, theme *theming.Theme) func(ComputedDatum) string {
	cfg := props.BorderColor
	if cfg.Type == 0 && cfg.Static == "" && cfg.ThemePath == "" && cfg.FromPath == "" && cfg.Func == nil {
		cfg = Defaults.BorderColor
	}
	gen := colors.GetInheritedColorGenerator(cfg, theme)
	return func(d ComputedDatum) string {
		return gen(map[string]any{
			"id":             d.ID,
			"label":          d.Label,
			"value":          d.Value,
			"formattedValue": d.FormattedValue,
			"color":          d.Color,
			"data":           d.Data,
		})
	}
}

// resolveInheritedColor resolves an inherited color config to a string,
// falling back to the default config when the input is unset.
func resolveInheritedColor(cfg colors.InheritedColorConfig, def colors.InheritedColorConfig, theme *theming.Theme) string {
	if cfg.Type == 0 && cfg.Static == "" && cfg.ThemePath == "" && cfg.FromPath == "" && cfg.Func == nil {
		cfg = def
	}
	return colors.GetInheritedColorGenerator(cfg, theme)(nil)
}

// dataAsAny converts []ComputedDatum to []any for BindDefs.
func dataAsAny(ds []ComputedDatum) []any {
	out := make([]any, len(ds))
	for i, d := range ds {
		out[i] = d
	}
	return out
}

// applyDefaults overlays Defaults onto a partial PieProps.
func applyDefaults(p PieProps) PieProps {
	if p.ID == nil {
		p.ID = Defaults.ID
	}
	if p.Value == nil {
		p.Value = Defaults.Value
	}
	if p.Colors.Type == 0 && p.Colors.Scheme == "" && p.Colors.Static == "" && len(p.Colors.Colors) == 0 && p.Colors.Func == nil && p.Colors.DatumPath == "" {
		p.Colors = Defaults.Colors
	}
	if p.BorderColor.Type == 0 && p.BorderColor.Static == "" && p.BorderColor.ThemePath == "" && p.BorderColor.FromPath == "" && p.BorderColor.Func == nil {
		p.BorderColor = Defaults.BorderColor
	}
	if p.ArcLabelsTextColor.Type == 0 && p.ArcLabelsTextColor.Static == "" && p.ArcLabelsTextColor.ThemePath == "" && p.ArcLabelsTextColor.FromPath == "" && p.ArcLabelsTextColor.Func == nil {
		p.ArcLabelsTextColor = Defaults.ArcLabelsTextColor
	}
	if p.ArcLinkLabelsTextColor.Type == 0 && p.ArcLinkLabelsTextColor.Static == "" && p.ArcLinkLabelsTextColor.ThemePath == "" && p.ArcLinkLabelsTextColor.FromPath == "" && p.ArcLinkLabelsTextColor.Func == nil {
		p.ArcLinkLabelsTextColor = Defaults.ArcLinkLabelsTextColor
	}
	if p.ArcLinkLabelsColor.Type == 0 && p.ArcLinkLabelsColor.Static == "" && p.ArcLinkLabelsColor.ThemePath == "" && p.ArcLinkLabelsColor.FromPath == "" && p.ArcLinkLabelsColor.Func == nil {
		p.ArcLinkLabelsColor = Defaults.ArcLinkLabelsColor
	}
	if p.ArcLabel == nil {
		p.ArcLabel = Defaults.ArcLabel
	}
	if p.ArcLinkLabel == nil {
		p.ArcLinkLabel = Defaults.ArcLinkLabel
	}
	if p.ArcLabelsRadiusOffset == 0 {
		p.ArcLabelsRadiusOffset = Defaults.ArcLabelsRadiusOffset
	}
	if p.ArcLinkLabelsDiagonalLength == 0 {
		p.ArcLinkLabelsDiagonalLength = Defaults.ArcLinkLabelsDiagonalLength
	}
	if p.ArcLinkLabelsStraightLength == 0 {
		p.ArcLinkLabelsStraightLength = Defaults.ArcLinkLabelsStraightLength
	}
	if p.ArcLinkLabelsThickness == 0 {
		p.ArcLinkLabelsThickness = Defaults.ArcLinkLabelsThickness
	}
	if p.ArcLinkLabelsTextOffset == 0 {
		p.ArcLinkLabelsTextOffset = Defaults.ArcLinkLabelsTextOffset
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	if len(p.Layers) == 0 {
		p.Layers = DefaultLayers
	}
	if p.TransitionMode == "" {
		p.TransitionMode = Defaults.TransitionMode
	}
	if p.EndAngle == 0 {
		p.EndAngle = Defaults.EndAngle
	}
	return p
}

// guard against unused imports during refactoring
var _ = arcs.Arc{}
