package htmx

import (
	"context"
	"strings"

	"github.com/a-h/templ"

	"github.com/geoffjay/templ-charts/charts/bar"
	cp "github.com/geoffjay/templ-charts/charts/circlepacking"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/icicle"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/sunburst"
	"github.com/geoffjay/templ-charts/charts/treemap"
)

// renderFull renders the full SVG for the instance, applying the current
// state (hidden ids, pie active id) to a clone of the props snapshot. The
// returned string is a complete <svg>…</svg> document.
func renderFull(inst *ChartInstance) (string, error) {
	st := inst.State()
	switch inst.Kind {
	case KindBar:
		props := inst.Props.(bar.BarProps)
		applyBarState(&props, inst.ID, st)
		return renderComponent(bar.Bar(props))
	case KindLine:
		props := inst.Props.(line.LineProps)
		applyLineState(&props, inst.ID, st)
		return renderComponent(line.Line(props))
	case KindPie:
		props := inst.Props.(pie.PieProps)
		applyPieState(&props, inst.ID, st)
		return renderComponent(pie.Pie(props))
	case KindHeatmap:
		props := inst.Props.(heatmap.HeatMapProps)
		applyHeatmapState(&props, inst.ID, st)
		return renderComponent(heatmap.HeatMap(props))
	case KindIcicle:
		props := inst.Props.(icicle.IcicleProps)
		applyIcicleState(&props, inst.ID, st)
		return renderComponent(icicle.Icicle(props))
	case KindTreemap:
		props := inst.Props.(treemap.TreemapProps)
		applyTreemapState(&props, inst.ID, st)
		return renderComponent(treemap.Treemap(props))
	case KindCirclePack:
		props := inst.Props.(cp.CirclePackingProps)
		applyCirclePackingState(&props, inst.ID, st)
		return renderComponent(cp.CirclePacking(props))
	case KindSunburst:
		props := inst.Props.(sunburst.SunburstProps)
		applySunburstState(&props, inst.ID, st)
		return renderComponent(sunburst.Sunburst(props))
	}
	return "", errUnknownKind
}

// applyBarState overlays the per-instance state onto a bar.BarProps clone.
func applyBarState(props *bar.BarProps, id string, st State) {
	props.ChartID = id
	props.Interactive = true
	props.InitialHiddenIDs = st.HiddenIDs
	props.HoveredKey = st.HoveredKey
	if props.Legends != nil {
		for i := range props.Legends {
			props.Legends[i].ChartID = id
			if props.Legends[i].Toggle == "" {
				props.Legends[i].Toggle = "1"
			}
		}
	}
}

// applyLineState overlays the per-instance state onto a line.LineProps clone.
func applyLineState(props *line.LineProps, id string, st State) {
	props.ChartID = id
	props.Interactive = true
	// The htmx registry IS the server-round-trip path, so it opts into the
	// legacy per-mousemove server hover (mesh/slice hx-get endpoints); the
	// standalone default is now the client path (the fallback is retired).
	props.ServerHover = true
	props.InitialHiddenIDs = st.HiddenIDs
	props.HoverX = st.HoverX
	props.HoverY = st.HoverY
	props.HasHover = st.HasHover
	if props.Legends != nil {
		for i := range props.Legends {
			props.Legends[i].ChartID = id
			if props.Legends[i].Toggle == "" {
				props.Legends[i].Toggle = "1"
			}
		}
	}
}

// applyPieState overlays the per-instance state onto a pie.PieProps clone.
func applyPieState(props *pie.PieProps, id string, st State) {
	props.ChartID = id
	props.Interactive = true
	props.InitialHiddenIDs = st.HiddenIDs
	props.ActiveID = st.ActiveID
	if props.Legends != nil {
		for i := range props.Legends {
			props.Legends[i].ChartID = id
			if props.Legends[i].Toggle == "" {
				props.Legends[i].Toggle = "1"
			}
		}
	}
}

// applyHeatmapState overlays the per-instance state onto a heatmap.HeatMapProps
// clone (the hovered cell drives the active/dim opacity).
func applyHeatmapState(props *heatmap.HeatMapProps, id string, st State) {
	props.ChartID = id
	props.Interactive = true
	props.HoveredKey = st.HoveredKey
}

// applyIcicleState overlays the per-instance state onto an icicle.IcicleProps
// clone: the chart id (so hx-* zoom wiring is scoped) and the current zoom
// focus. EnableZooming stays whatever the caller registered.
func applyIcicleState(props *icicle.IcicleProps, id string, st State) {
	props.ChartID = id
	props.Interactive = true
	props.FocusID = st.FocusID
}

// applyTreemapState overlays the per-instance state onto a treemap.TreemapProps
// clone.
func applyTreemapState(props *treemap.TreemapProps, id string, st State) {
	props.ChartID = id
	props.Interactive = true
	props.FocusID = st.FocusID
}

// applyCirclePackingState overlays the per-instance state onto a
// cp.CirclePackingProps clone.
func applyCirclePackingState(props *cp.CirclePackingProps, id string, st State) {
	props.ChartID = id
	props.Interactive = true
	props.FocusID = st.FocusID
}

// applySunburstState overlays the per-instance state onto a
// sunburst.SunburstProps clone.
func applySunburstState(props *sunburst.SunburstProps, id string, st State) {
	props.ChartID = id
	props.Interactive = true
	props.FocusID = st.FocusID
}

// zoomParentID resolves the parent-node id of nodeID for a hierarchy instance,
// walking the chart's input data. It returns "" when nodeID is the root or a
// depth-1 node (whose parent is the root) — i.e. zooming out returns to the
// full view. Returns ("", false) for a non-hierarchy kind.
func zoomParentID(inst *ChartInstance, nodeID string) (string, bool) {
	switch inst.Kind {
	case KindIcicle:
		root := inst.Props.(icicle.IcicleProps).Data
		return icicleParent(root, "", nodeID), true
	case KindTreemap:
		root := inst.Props.(treemap.TreemapProps).Data
		return treemapParent(root, "", nodeID), true
	case KindCirclePack:
		root := inst.Props.(cp.CirclePackingProps).Data
		return cpParent(root, "", nodeID), true
	case KindSunburst:
		root := inst.Props.(sunburst.SunburstProps).Data
		return sunburstParent(root, "", nodeID), true
	}
	return "", false
}

// The parent resolvers return the id of parentID when a child matches target,
// or "" if the parent is the root (so zoom-out reaches the full view). They
// recurse depth-first over the input hierarchy.
func icicleParent(n icicle.IcicleNode, parentID, target string) string {
	for _, c := range n.Children {
		if c.ID == target {
			return parentID
		}
		if childHasIcicle(c, target) {
			return icicleParent(c, c.ID, target)
		}
	}
	return ""
}

func childHasIcicle(n icicle.IcicleNode, target string) bool {
	for _, c := range n.Children {
		if c.ID == target || childHasIcicle(c, target) {
			return true
		}
	}
	return false
}

func treemapParent(n treemap.TreemapNode, parentID, target string) string {
	for _, c := range n.Children {
		if c.ID == target {
			return parentID
		}
		if childHasTreemap(c, target) {
			return treemapParent(c, c.ID, target)
		}
	}
	return ""
}

func childHasTreemap(n treemap.TreemapNode, target string) bool {
	for _, c := range n.Children {
		if c.ID == target || childHasTreemap(c, target) {
			return true
		}
	}
	return false
}

func cpParent(n cp.CirclePackingNode, parentID, target string) string {
	for _, c := range n.Children {
		if c.ID == target {
			return parentID
		}
		if childHasCP(c, target) {
			return cpParent(c, c.ID, target)
		}
	}
	return ""
}

func childHasCP(n cp.CirclePackingNode, target string) bool {
	for _, c := range n.Children {
		if c.ID == target || childHasCP(c, target) {
			return true
		}
	}
	return false
}

func sunburstParent(n sunburst.SunburstNode, parentID, target string) string {
	for _, c := range n.Children {
		if c.ID == target {
			return parentID
		}
		if childHasSunburst(c, target) {
			return sunburstParent(c, c.ID, target)
		}
	}
	return ""
}

func childHasSunburst(n sunburst.SunburstNode, target string) bool {
	for _, c := range n.Children {
		if c.ID == target || childHasSunburst(c, target) {
			return true
		}
	}
	return false
}

// renderComponent renders a templ.Component to a string.
func renderComponent(c templ.Component) (string, error) {
	var b strings.Builder
	if err := c.Render(context.Background(), &b); err != nil {
		return "", err
	}
	return b.String(), nil
}

// errUnknownKind is returned by renderFull when the instance Kind is not one
// of the supported chart kinds.
var errUnknownKind = htmxError("htmx: unknown chart kind")

// errUnknownInstance is returned when an instance id is not registered.
var errUnknownInstance = htmxError("htmx: unknown instance")

// htmxError is a simple string error type.
type htmxError string

func (e htmxError) Error() string { return string(e) }
