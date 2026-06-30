package waffle

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/interact"
	"github.com/geoffjay/templ-charts/charts/legends"
)

// renderComponent renders a templ.Component to a string.
func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}

// applyDefaults fills zero-valued WaffleProps fields from Defaults.
func applyDefaults(p WaffleProps) WaffleProps {
	if p.FillDirection == "" {
		p.FillDirection = Defaults.FillDirection
	}
	if p.Padding == 0 {
		p.Padding = Defaults.Padding
	}
	if p.Colors.Type == 0 && p.Colors.Scheme == "" && p.Colors.Static == "" && len(p.Colors.Colors) == 0 && p.Colors.Func == nil && p.Colors.DatumPath == "" {
		p.Colors = Defaults.Colors
	}
	if p.EmptyColor == "" {
		p.EmptyColor = Defaults.EmptyColor
	}
	if p.EmptyOpacity == 0 {
		p.EmptyOpacity = Defaults.EmptyOpacity
	}
	if isZeroColor(p.BorderColor) {
		p.BorderColor = Defaults.BorderColor
	}
	if len(p.Layers) == 0 {
		p.Layers = DefaultLayers
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	if p.Rows <= 0 {
		p.Rows = 10
	}
	if p.Columns <= 0 {
		p.Columns = 10
	}
	if p.Total == 0 {
		p.Total = sumValues(p.Data)
	}
	return p
}

func sumValues(data []WaffleDatum) float64 {
	s := 0.0
	for _, d := range data {
		s += d.Value
	}
	return s
}

// isZeroColor reports whether an InheritedColorConfig is unset.
func isZeroColor(c colors.InheritedColorConfig) bool {
	return c.Type == 0 && c.Static == "" && c.ThemePath == "" && c.FromPath == "" && c.Func == nil
}

// renderLayers renders the enabled layers as an inner SVG string.
func renderLayers(props WaffleProps, result WaffleResult, dims core.Dimensions) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case WaffleLayerCells:
			b.WriteString(renderCellsLayer(props, result))
		case WaffleLayerLegends:
			b.WriteString(renderLegendsLayer(props, result, dims))
		}
	}
	return b.String()
}

func renderCellsLayer(props WaffleProps, result WaffleResult) string {
	// Cells are positioned relative to the grid origin; wrap them in a single
	// translate so each cell's X/Y stays in grid-local coordinates.
	var inner strings.Builder
	for _, cell := range result.Cells {
		cp := WaffleCellShapeProps{
			Cell:         cell,
			BorderWidth:  props.BorderWidth,
			BorderRadius: props.BorderRadius,
		}
		if props.Interactive && cell.HasData {
			cp.Tooltip = interact.TooltipHTML(cell.Color, cell.Label, "")
		}
		inner.WriteString(renderComponent(WaffleCellShape(cp)))
	}
	return fmt.Sprintf(`<g transform="translate(%s,%s)">%s</g>`, fmtW(result.GridX), fmtW(result.GridY), inner.String())
}

func renderLegendsLayer(props WaffleProps, result WaffleResult, dims core.Dimensions) string {
	if len(props.Legends) == 0 {
		return ""
	}
	items := make([]legends.Datum, 0, len(result.ComputedData))
	for _, cd := range result.ComputedData {
		items = append(items, legends.Datum{ID: cd.ID, Label: cd.Label, Color: cd.Color, Hidden: cd.Hidden})
	}
	var s strings.Builder
	for _, lg := range props.Legends {
		s.WriteString(renderComponent(legends.BoxLegendSvg(legends.BoxLegendSvgProps{
			Props: legends.LegendProps{
				Anchor:      lg.Anchor,
				Direction:   orDirection(lg.Direction, legends.LegendDirectionColumn),
				Items:       items,
				TranslateX:  lg.TranslateX,
				TranslateY:  lg.TranslateY,
				ItemWidth:   orFloat(lg.ItemWidth, 100),
				ItemHeight:  orFloat(lg.ItemHeight, 20),
				SymbolShape: orSymbol(lg.SymbolShape, legends.SymbolShapeSquare),
				SymbolSize:  14,
			},
			ChartWidth:  dims.InnerWidth,
			ChartHeight: dims.InnerHeight,
		})))
	}
	return s.String()
}

func orDirection(v, def legends.LegendDirection) legends.LegendDirection {
	if v == "" {
		return def
	}
	return v
}

func orSymbol(v, def legends.SymbolShape) legends.SymbolShape {
	if v == "" {
		return def
	}
	return v
}

func orFloat(v, def float64) float64 {
	if v == 0 {
		return def
	}
	return v
}
