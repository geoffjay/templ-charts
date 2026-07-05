// Package network mirrors @nivo/network: a force-directed node/link graph laid
// out server-side with internal/d3/force (link + many-body + center forces run
// for a fixed iteration count). It reuses charts/core (SvgWrapper), charts/interact
// (optional per-node hover tooltips) and charts/theming.
//
// SVG only, static render. The layout is deterministic (phyllotaxis
// seeding + d3's LCG, a fixed Iterations count), so goldens are byte-stable.
// Annotations are deferred. Interactive adds per-node hover tooltips and (unless
// UseMesh routes hover through the voronoi overlay) a chord/sankey-style
// hover-highlight via scoped CSS :has(); Animate adds a radius enter
// transition.
package network

import (
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/theming"
)

// NetworkInputNode is one input node. Size and Color are optional per-node
// overrides; when zero/empty the chart's NodeSize/NodeColor apply. Mirrors
// @nivo/network InputNode (with the size/color accessors folded in).
type NetworkInputNode struct {
	ID    string
	Size  float64
	Color string
}

// NetworkInputLink is one input link between two node ids. Distance is an
// optional per-link target length; when zero the chart's LinkDistance applies.
// Mirrors @nivo/network InputLink.
type NetworkInputLink struct {
	Source   string
	Target   string
	Distance float64
}

// ComputedNode is one positioned, styled node ready to render.
type ComputedNode struct {
	ID          string
	Index       int
	X, Y        float64
	Size        float64
	Color       string
	BorderWidth float64
	BorderColor string
}

// ComputedLink is one positioned, styled link (a straight segment between the
// centers of its endpoints).
type ComputedLink struct {
	ID        string
	Source    string
	Target    string
	X1, Y1    float64
	X2, Y2    float64
	Thickness float64
	Color     string
}

// NetworkLayerId enumerates the render layers. Mirrors @nivo/network LayerId.
type NetworkLayerId string

const (
	NetworkLayerLinks       NetworkLayerId = "links"
	NetworkLayerNodes       NetworkLayerId = "nodes"
	NetworkLayerMesh        NetworkLayerId = "mesh"
	NetworkLayerAnnotations NetworkLayerId = "annotations"
)

// DefaultLayers mirrors @nivo/network svgDefaultProps.layers.
var DefaultLayers = []NetworkLayerId{NetworkLayerLinks, NetworkLayerNodes, NetworkLayerMesh, NetworkLayerAnnotations}

// NetworkProps mirrors @nivo/network NetworkSvgProps (the supported subset).
// Fields left zero fall back to Defaults via applyDefaults.
type NetworkProps struct {
	Nodes []NetworkInputNode
	Links []NetworkInputLink

	Width  float64
	Height float64
	Margin core.Margin
	// Responsive makes the svg scale fluidly to its container.
	Responsive bool

	// Force-simulation controls (see internal/d3/force). LinkDistance is the
	// spring rest length; CenteringStrength is the link-force strength (nivo's
	// naming); Repulsivity is the positive charge magnitude (applied negatively).
	LinkDistance      float64
	CenteringStrength float64
	Repulsivity       float64
	DistanceMin       float64
	DistanceMax       float64 // 0 → +Inf
	Iterations        int

	NodeSize        float64
	NodeColor       string // default node fill; per-node Color overrides
	NodeBorderWidth float64
	NodeBorderColor string // "" → inherit node color

	LinkThickness float64
	LinkColor     string // "" → inherit source node color

	// FitView rescales the settled layout (uniformly, preserving aspect ratio)
	// to fill the inner area, insetting by the largest node radius so nothing
	// clips. Off by default — the raw force layout, like nivo, sizes itself by
	// its physics (LinkDistance/Repulsivity) and often occupies only the center.
	FitView bool

	// Interactive enables per-node client-side hover tooltips (charts/interact)
	// AND — unless UseMesh routes hover through the voronoi overlay — a
	// chord/sankey-style hover-highlight: hovering a node dims the rest and
	// re-lights that node, its links and its neighbours; hovering a link
	// re-lights it and its two endpoints. Driven by a scoped CSS :has() <style>
	// block (no JS/server round-trip). The *HoverOpacity / *HoverOthersOpacity
	// fields set the highlighted / dimmed opacities; zero falls back to Defaults.
	Interactive bool

	// Hover-highlight opacities (used when Interactive && !UseMesh).
	// NodeHoverOpacity / LinkHoverOpacity apply to the hovered element and its
	// connected elements; NodeHoverOthersOpacity / LinkHoverOthersOpacity dim
	// everything else.
	NodeHoverOpacity       float64
	NodeHoverOthersOpacity float64
	LinkHoverOpacity       float64
	LinkHoverOthersOpacity float64

	// UseMesh routes hover through an accurate voronoi mesh (charts/interact,
	// backed by internal/d3/delaunay) built over the node centers rather than
	// per-node tooltips: hovering anywhere resolves to the nearest node. Active
	// only when Interactive.
	UseMesh bool
	// DebugMesh draws the voronoi cells as a faint guide when UseMesh is on.
	DebugMesh bool
	// DetectionRadius, when > 0, bounds mesh hit-testing to this pixel distance.
	DetectionRadius float64

	Theme  *theming.Theme
	Layers []NetworkLayerId

	Role            string
	AriaLabel       string
	AriaLabelledBy  string
	AriaDescribedBy string
	Title           string
	Desc            string
	IsFocusable     bool

	// Animate, when true, emits a SMIL enter animation on each node scaling its
	// radius from 0 to its final value (600ms). MotionStagger delays successive
	// nodes by that many seconds (0 ⇒ all nodes enter together). Defaults off,
	// so static output is unchanged.
	Animate       bool
	MotionStagger float64
}

// NetworkResult is the computed model produced by UseNetwork.
type NetworkResult struct {
	Nodes []ComputedNode
	Links []ComputedLink
}
