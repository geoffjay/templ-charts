package waffle

import (
	"context"
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/colors"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/grid"
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
		case WaffleLayerAreas:
			b.WriteString(renderAreasLayer(props, result))
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
	for i, cell := range result.Cells {
		cp := WaffleCellShapeProps{
			Cell:         cell,
			BorderWidth:  props.BorderWidth,
			BorderRadius: props.BorderRadius,
			Animate:      props.Animate,
		}
		if props.Animate {
			cp.AnimateBegin = core.StaggerBegin(i, props.MotionStagger)
		}
		if props.Interactive && cell.HasData {
			cp.Tooltip = interact.TooltipHTML(cell.Color, cell.Label, "")
		}
		inner.WriteString(renderComponent(WaffleCellShape(cp)))
	}
	return fmt.Sprintf(`<g transform="translate(%s,%s)">%s</g>`, fmtW(result.GridX), fmtW(result.GridY), inner.String())
}

// renderAreasLayer draws one union-outline <path> per datum (in computed-data
// order for determinism) instead of per-cell rects. Each datum's data cells are
// merged into minimal polygons via grid.GetCellsPolygons; each polygon becomes a
// closed subpath (M … L … Z). Fill/stroke mirror how renderCellsLayer styles a
// data cell (datum Color at full opacity, BorderColor at BorderWidth).
func renderAreasLayer(props WaffleProps, result WaffleResult) string {
	// Group data cells by DatumID.
	byDatum := map[string][]WaffleCell{}
	for _, cell := range result.Cells {
		if cell.HasData {
			byDatum[cell.DatumID] = append(byDatum[cell.DatumID], cell)
		}
	}

	var inner strings.Builder
	// Iterate in computed-data order (not map order) for deterministic output.
	for _, cd := range result.ComputedData {
		cells := byDatum[cd.ID]
		if len(cells) == 0 {
			continue
		}
		polygons := grid.GetCellsPolygons(cells)
		var d strings.Builder
		for _, poly := range polygons {
			for j, v := range poly {
				if j == 0 {
					d.WriteString("M")
				} else {
					d.WriteString("L")
				}
				d.WriteString(fmtW(v[0]))
				d.WriteString(",")
				d.WriteString(fmtW(v[1]))
				if j < len(poly)-1 {
					d.WriteString(" ")
				}
			}
			d.WriteString("Z")
		}
		if d.Len() == 0 {
			continue
		}
		inner.WriteString(fmt.Sprintf(`<path d="%s" fill="%s" fill-opacity="%s"`, d.String(), cd.Color, fmtW(1)))
		if props.BorderWidth > 0 && cd.BorderColor != "" {
			inner.WriteString(fmt.Sprintf(` stroke="%s" stroke-width="%s"`, cd.BorderColor, fmtW(props.BorderWidth)))
		}
		inner.WriteString(`></path>`)
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
