package tree

import (
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/interact"
	"github.com/geoffjay/templ-charts/charts/theming"
)

func applyDefaults(p TreeProps) TreeProps {
	if p.Mode == "" {
		p.Mode = Defaults.Mode
	}
	if p.Layout == "" {
		p.Layout = Defaults.Layout
	}
	if p.NodeSize == 0 {
		p.NodeSize = Defaults.NodeSize
	}
	if p.LinkThickness == 0 {
		p.LinkThickness = Defaults.LinkThickness
	}
	if p.LinkOpacity == 0 {
		p.LinkOpacity = Defaults.LinkOpacity
	}
	if p.EnableLabel == nil {
		p.EnableLabel = Defaults.EnableLabel
	}
	if p.LabelOffset == 0 {
		p.LabelOffset = Defaults.LabelOffset
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	return p
}

// renderTree draws the links, then the nodes, then optional labels, then —
// when Interactive && UseMesh — an accurate voronoi-mesh hover overlay
// (charts/interact, backed by internal/d3/delaunay). When Interactive without
// UseMesh, each node carries its own hover tooltip instead.
func renderTree(props TreeProps, result TreeResult, theme *theming.Theme) string {
	var b strings.Builder
	perNodeTooltip := props.Interactive && !props.UseMesh

	// Links (behind nodes). Fade in on enter when animating.
	for i, l := range result.Links {
		if l.Path == "" {
			continue
		}
		b.WriteString(`<path d="`)
		b.WriteString(l.Path)
		b.WriteString(`" fill="none" stroke="`)
		b.WriteString(l.Color)
		b.WriteString(`" stroke-width="`)
		b.WriteString(fmtF(props.LinkThickness))
		b.WriteString(`" stroke-opacity="`)
		b.WriteString(fmtF(props.LinkOpacity))
		b.WriteString(`">`)
		if props.Animate {
			b.WriteString(core.SMILFadeIn(core.StaggerBegin(i, props.MotionStagger)))
		}
		b.WriteString(`</path>`)
	}

	// Nodes. Scale radius from 0 on enter when animating.
	r := props.NodeSize / 2
	for i, n := range result.Nodes {
		b.WriteString(`<circle cx="`)
		b.WriteString(fmtF(n.X))
		b.WriteString(`" cy="`)
		b.WriteString(fmtF(n.Y))
		b.WriteString(`" r="`)
		b.WriteString(fmtF(r))
		b.WriteString(`" fill="`)
		b.WriteString(n.Color)
		b.WriteString(`"`)
		if perNodeTooltip {
			b.WriteString(` `)
			b.WriteString(interact.TooltipAttrName)
			b.WriteString(`="`)
			b.WriteString(templ.EscapeString(interact.TooltipHTML(n.Color, n.ID, "")))
			b.WriteString(`" style="pointer-events:auto"`)
		}
		b.WriteString(`>`)
		if props.Animate {
			b.WriteString(core.SMILAnimate("r", "0", fmtF(r), core.StaggerBegin(i, props.MotionStagger)))
		}
		b.WriteString(`</circle>`)
	}

	// Labels.
	if props.LabelEnabled() {
		fill, fontSize, fontFamily := labelsTextStyle(theme)
		horizontal := props.Layout == LayoutLeftToRight || props.Layout == LayoutRightToLeft
		for _, n := range result.Nodes {
			lx, ly, anchor := n.X, n.Y, "middle"
			if horizontal {
				lx = n.X + r + props.LabelOffset
				anchor = "start"
			} else {
				ly = n.Y + r + props.LabelOffset
			}
			b.WriteString(text(lx, ly, anchor, n.ID, fill, fontSize, fontFamily))
		}
	}

	// Voronoi-mesh hover overlay (on top, so it captures the whole area).
	if props.Interactive && props.UseMesh {
		pts := make([]interact.MeshPoint, 0, len(result.Nodes))
		for _, n := range result.Nodes {
			pts = append(pts, interact.MeshPoint{
				X:    n.X,
				Y:    n.Y,
				HTML: interact.TooltipHTML(n.Color, n.ID, ""),
			})
		}
		b.WriteString(interact.MeshOverlay(pts, props.Width, props.Height, props.DebugMesh, props.DetectionRadius))
	}
	return b.String()
}

func text(x, y float64, anchor, s, fill string, fontSize float64, fontFamily string) string {
	var b strings.Builder
	b.WriteString(`<text x="`)
	b.WriteString(fmtF(x))
	b.WriteString(`" y="`)
	b.WriteString(fmtF(y))
	b.WriteString(`" text-anchor="`)
	b.WriteString(anchor)
	b.WriteString(`" dominant-baseline="central" style="pointer-events:none;fill:`)
	b.WriteString(fill)
	if fontSize > 0 {
		b.WriteString(`;font-size:`)
		b.WriteString(fmtF(fontSize))
		b.WriteString(`px`)
	}
	if fontFamily != "" {
		b.WriteString(`;font-family:`)
		b.WriteString(fontFamily)
	}
	b.WriteString(`">`)
	b.WriteString(templ.EscapeString(s))
	b.WriteString(`</text>`)
	return b.String()
}

func labelsTextStyle(theme *theming.Theme) (fill string, fontSize float64, fontFamily string) {
	if theme == nil {
		return "", 0, ""
	}
	t := theme.Labels.Text
	fill = t.Fill
	if fill == "" {
		fill = theme.Text.Fill
	}
	fs := t.FontSize
	if fs == nil {
		fs = theme.Text.FontSize
	}
	fontFamily = t.FontFamily
	if fontFamily == "" {
		fontFamily = theme.Text.FontFamily
	}
	return fill, toNum(fs), fontFamily
}

func toNum(v any) float64 {
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	}
	return 0
}

func fmtF(v float64) string {
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
}
