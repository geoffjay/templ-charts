package voronoi

import (
	"context"
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"github.com/geoffjay/templ-charts/charts/interact"
)

// renderComponent renders a templ.Component to a string.
func renderComponent(c templ.Component) string {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return ""
	}
	return b.String()
}

// applyDefaults fills zero-valued VoronoiProps fields from Defaults.
func applyDefaults(p VoronoiProps) VoronoiProps {
	if p.XDomain == [2]float64{} {
		p.XDomain = Defaults.XDomain
	}
	if p.YDomain == [2]float64{} {
		p.YDomain = Defaults.YDomain
	}
	if len(p.Layers) == 0 {
		p.Layers = Defaults.Layers
	}
	if p.EnableLinks == nil {
		p.EnableLinks = Defaults.EnableLinks
	}
	if p.LinkLineWidth == 0 {
		p.LinkLineWidth = Defaults.LinkLineWidth
	}
	if p.LinkLineColor == "" {
		p.LinkLineColor = Defaults.LinkLineColor
	}
	if p.EnableCells == nil {
		p.EnableCells = Defaults.EnableCells
	}
	if p.CellLineWidth == 0 {
		p.CellLineWidth = Defaults.CellLineWidth
	}
	if p.CellLineColor == "" {
		p.CellLineColor = Defaults.CellLineColor
	}
	if p.EnablePoints == nil {
		p.EnablePoints = Defaults.EnablePoints
	}
	if p.PointSize == 0 {
		p.PointSize = Defaults.PointSize
	}
	if p.PointColor == "" {
		p.PointColor = Defaults.PointColor
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	return p
}

// renderLayers renders the enabled layers as an inner SVG string, in the
// configured layer order.
func renderLayers(props VoronoiProps, result VoronoiResult) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case LayerLinks:
			if props.LinksEnabled() {
				b.WriteString(renderLinksLayer(props, result))
			}
		case LayerCells:
			if props.CellsEnabled() {
				b.WriteString(renderCellsLayer(props, result))
			}
		case LayerPoints:
			if props.PointsEnabled() {
				b.WriteString(renderPointsLayer(props, result))
			}
		case LayerBounds:
			b.WriteString(renderBoundsLayer(props, result))
		}
	}
	return b.String()
}

// renderLinksLayer draws the Delaunay edges as a single path.
func renderLinksLayer(props VoronoiProps, result VoronoiResult) string {
	d := result.Delaunay.Render()
	if d == "" {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<path d="`)
	b.WriteString(d)
	b.WriteString(`" fill="none" stroke="`)
	b.WriteString(props.LinkLineColor)
	b.WriteString(`" stroke-width="`)
	b.WriteString(fmtF(props.LinkLineWidth))
	b.WriteString(`"></path>`)
	return b.String()
}

// renderCellsLayer draws the Voronoi cells. When Interactive, each cell is a
// separate path carrying its datum's tooltip; otherwise all cells are a single
// path for a compact static render.
func renderCellsLayer(props VoronoiProps, result VoronoiResult) string {
	var b strings.Builder
	if props.Interactive {
		for i, p := range result.Points {
			cell := result.Voronoi.RenderCell(i)
			if cell == "" {
				continue
			}
			b.WriteString(`<path d="`)
			b.WriteString(cell)
			b.WriteString(`" fill="none" stroke="`)
			b.WriteString(props.CellLineColor)
			b.WriteString(`" stroke-width="`)
			b.WriteString(fmtF(props.CellLineWidth))
			// pointer-events:all makes the whole (unfilled) cell area hoverable,
			// not just its stroke, so the tooltip shows anywhere inside the cell.
			b.WriteString(`" pointer-events="all" data-tc-tooltip="`)
			b.WriteString(interact.EscapeAttr(interact.TooltipHTML(props.PointColor, p.ID, pointValue(p))))
			b.WriteString(`"></path>`)
		}
		return b.String()
	}
	d := result.Voronoi.Render()
	if d == "" {
		return ""
	}
	b.WriteString(`<path d="`)
	b.WriteString(d)
	b.WriteString(`" fill="none" stroke="`)
	b.WriteString(props.CellLineColor)
	b.WriteString(`" stroke-width="`)
	b.WriteString(fmtF(props.CellLineWidth))
	b.WriteString(`"></path>`)
	return b.String()
}

// renderPointsLayer draws each input point as a small circle.
func renderPointsLayer(props VoronoiProps, result VoronoiResult) string {
	var b strings.Builder
	r := props.PointSize / 2
	for _, p := range result.Points {
		b.WriteString(`<circle cx="`)
		b.WriteString(fmtF(p.X))
		b.WriteString(`" cy="`)
		b.WriteString(fmtF(p.Y))
		b.WriteString(`" r="`)
		b.WriteString(fmtF(r))
		b.WriteString(`" fill="`)
		b.WriteString(props.PointColor)
		b.WriteString(`"></circle>`)
	}
	return b.String()
}

// renderBoundsLayer draws the clip rectangle outline.
func renderBoundsLayer(props VoronoiProps, result VoronoiResult) string {
	var b strings.Builder
	b.WriteString(`<path d="`)
	b.WriteString(result.Voronoi.RenderBounds())
	b.WriteString(`" fill="none" stroke="`)
	b.WriteString(props.CellLineColor)
	b.WriteString(`" stroke-width="`)
	b.WriteString(fmtF(props.CellLineWidth))
	b.WriteString(`" stroke-opacity="0.35"></path>`)
	return b.String()
}

// pointValue formats a datum's coordinates for the tooltip.
func pointValue(p ComputedPoint) string {
	return "x " + trimNum(p.Data.X) + ", y " + trimNum(p.Data.Y)
}

func trimNum(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

// fmtF formats a float rounded to 3 decimals (matching the d3-path serializer).
func fmtF(v float64) string {
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
}
