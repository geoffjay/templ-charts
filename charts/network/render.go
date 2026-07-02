package network

import (
	"math"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/geoffjay/templ-charts/charts/interact"
)

// applyDefaults fills zero-valued NetworkProps fields from Defaults.
func applyDefaults(p NetworkProps) NetworkProps {
	if p.LinkDistance == 0 {
		p.LinkDistance = Defaults.LinkDistance
	}
	if p.CenteringStrength == 0 {
		p.CenteringStrength = Defaults.CenteringStrength
	}
	if p.Repulsivity == 0 {
		p.Repulsivity = Defaults.Repulsivity
	}
	if p.DistanceMin == 0 {
		p.DistanceMin = Defaults.DistanceMin
	}
	if p.Iterations == 0 {
		p.Iterations = Defaults.Iterations
	}
	if p.NodeSize == 0 {
		p.NodeSize = Defaults.NodeSize
	}
	if p.NodeColor == "" {
		p.NodeColor = Defaults.NodeColor
	}
	if p.LinkThickness == 0 {
		p.LinkThickness = Defaults.LinkThickness
	}
	if len(p.Layers) == 0 {
		p.Layers = Defaults.Layers
	}
	if p.Role == "" {
		p.Role = Defaults.Role
	}
	return p
}

// renderLayers renders the enabled layers as an inner SVG string. Annotations
// are deferred (no network annotation specs in v3).
func renderLayers(props NetworkProps, result NetworkResult) string {
	var b strings.Builder
	for _, layer := range props.Layers {
		switch layer {
		case NetworkLayerLinks:
			b.WriteString(renderLinksLayer(result))
		case NetworkLayerNodes:
			b.WriteString(renderNodesLayer(props, result))
		case NetworkLayerAnnotations:
			// annotations: deferred in v3.
		}
	}
	return b.String()
}

func renderLinksLayer(result NetworkResult) string {
	var b strings.Builder
	for _, l := range result.Links {
		b.WriteString(`<line x1="`)
		b.WriteString(fmtF(l.X1))
		b.WriteString(`" y1="`)
		b.WriteString(fmtF(l.Y1))
		b.WriteString(`" x2="`)
		b.WriteString(fmtF(l.X2))
		b.WriteString(`" y2="`)
		b.WriteString(fmtF(l.Y2))
		b.WriteString(`" stroke="`)
		b.WriteString(l.Color)
		b.WriteString(`" stroke-width="`)
		b.WriteString(fmtF(l.Thickness))
		b.WriteString(`"></line>`)
	}
	return b.String()
}

func renderNodesLayer(props NetworkProps, result NetworkResult) string {
	var b strings.Builder
	for _, node := range result.Nodes {
		b.WriteString(`<circle cx="`)
		b.WriteString(fmtF(node.X))
		b.WriteString(`" cy="`)
		b.WriteString(fmtF(node.Y))
		b.WriteString(`" r="`)
		b.WriteString(fmtF(node.Size / 2))
		b.WriteString(`" fill="`)
		b.WriteString(node.Color)
		b.WriteString(`"`)
		if node.BorderWidth > 0 {
			b.WriteString(` stroke="`)
			b.WriteString(node.BorderColor)
			b.WriteString(`" stroke-width="`)
			b.WriteString(fmtF(node.BorderWidth))
			b.WriteString(`"`)
		}
		if props.Interactive {
			b.WriteString(` `)
			b.WriteString(interact.TooltipAttrName)
			b.WriteString(`="`)
			b.WriteString(templ.EscapeString(interact.TooltipHTML(node.Color, node.ID, "")))
			b.WriteString(`" style="pointer-events:auto"`)
		}
		b.WriteString(`></circle>`)
	}
	return b.String()
}

func fmtF(v float64) string {
	return strconv.FormatFloat(math.Round(v*1000)/1000, 'g', -1, 64)
}
