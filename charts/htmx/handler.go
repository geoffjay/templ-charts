package htmx

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/geoffjay/templ-charts/charts/bar"
	"github.com/geoffjay/templ-charts/charts/core"
	"github.com/geoffjay/templ-charts/charts/heatmap"
	"github.com/geoffjay/templ-charts/charts/line"
	"github.com/geoffjay/templ-charts/charts/pie"
	"github.com/geoffjay/templ-charts/charts/tooltip"
)

// Handler is the http.Handler serving the htmx endpoints. Mount it at any
// prefix; routes are written as absolute paths under /charts/. If you mount
// it at a non-root prefix, strip the prefix before delegating (or use a
// sub-mount with http.StripPrefix).
//
//	Route                       Method  Action
//	/charts/{id}                GET     Full SVG render (hx-aware markup)
//	/charts/{id}/hover          GET     Bar/Pie hover tooltip HTML partial
//	/charts/{id}/slice          GET     Line slice tooltip HTML partial
//	/charts/{id}/click          POST    Bar activation toggle (re-render SVG)
//	/charts/{id}/toggle         POST    Series toggle (re-render SVG)
//
// Unknown instances return 404; unknown verbs return 400.
type Handler struct {
	registry *Registry
}

// NewHandler returns a Handler bound to the given registry.
func NewHandler(r *Registry) *Handler {
	return &Handler{registry: r}
}

// Registry returns the handler's underlying registry (handy for tests / demos
// that want to inspect or mutate state directly).
func (h *Handler) Registry() *Registry { return h.registry }

// RenderFull renders the complete SVG for the instance with the current state
// applied. It is the in-process equivalent of GET /charts/{id}: the page
// handlers use it to render the initial SVG inline without a round-trip.
func (h *Handler) RenderFull(id string) (string, error) {
	inst := h.registry.Get(id)
	if inst == nil {
		return "", errUnknownInstance
	}
	return renderFull(inst)
}

// ServeHTTP dispatches the request to the appropriate endpoint handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	// All routes start with "/charts/".
	if !strings.HasPrefix(path, "/charts/") {
		http.NotFound(w, r)
		return
	}
	rest := strings.TrimPrefix(path, "/charts/")
	// Split into id + optional action: "<id>" or "<id>/<action>".
	id, action, _ := strings.Cut(rest, "/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	inst := h.registry.Get(id)
	if inst == nil {
		http.NotFound(w, r)
		return
	}

	switch action {
	case "":
		// Full render: GET only.
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleFull(w, r, inst)
	case "hover":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleHover(w, r, inst)
	case "slice":
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleSlice(w, r, inst)
	case "click":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleClick(w, r, inst)
	case "toggle":
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.handleToggle(w, r, inst)
	default:
		http.NotFound(w, r)
	}
}

// handleFull re-renders the complete SVG for the instance and writes it with
// Content-Type image/svg+xml.
func (h *Handler) handleFull(w http.ResponseWriter, r *http.Request, inst *ChartInstance) {
	out, err := renderFull(inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	_, _ = w.Write([]byte(out))
}

// handleToggle flips the membership of the `series` query param in the
// instance's hidden-id list, then re-renders the full SVG (same body as
// handleFull — the legend swaps the chart container's innerHTML).
func (h *Handler) handleToggle(w http.ResponseWriter, r *http.Request, inst *ChartInstance) {
	series := r.URL.Query().Get("series")
	if series == "" {
		http.Error(w, "missing series", http.StatusBadRequest)
		return
	}
	inst.toggleHidden(series)
	out, err := renderFull(inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	_, _ = w.Write([]byte(out))
}

// handleClick toggles bar activation. v1 supports verb=activate which simply
// records the hovered key as the active key (no visual change in v1 bar);
// future verbs (select, etc.) can extend this. The response is the full SVG.
func (h *Handler) handleClick(w http.ResponseWriter, r *http.Request, inst *ChartInstance) {
	if inst.Kind != KindBar {
		http.Error(w, "click not supported for this chart kind", http.StatusBadRequest)
		return
	}
	key := r.URL.Query().Get("bar")
	verb := r.URL.Query().Get("verb")
	if verb == "" {
		verb = "activate"
	}
	switch verb {
	case "activate":
		// Toggle the hovered key: if it's the current hovered key, clear;
		// otherwise set it. This mirrors nivo's bar activation toggle.
		st := inst.State()
		if st.HoveredKey == key {
			inst.setHovered("")
		} else {
			inst.setHovered(key)
		}
	default:
		http.Error(w, "unknown verb: "+verb, http.StatusBadRequest)
		return
	}
	out, err := renderFull(inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	_, _ = w.Write([]byte(out))
}

// handleHover produces an HTML tooltip fragment for the bar/pie/line element
// under the cursor, plus an out-of-band swap of the full SVG (re-rendered with
// hover/active state). For pie, it also sets ActiveID so the arc pops. For line
// mesh, it sets HoverX/HoverY so the crosshair renders. The `?leave=1` query
// clears all hover state and returns an empty tooltip + clean SVG.
//
// The response body is the tooltip HTML (swapped into #tooltip-<id>), followed
// by an OOB <div hx-swap-oob="innerHTML:#chart-<id>"> containing the
// re-rendered SVG so a single request updates both the tooltip and the chart.
func (h *Handler) handleHover(w http.ResponseWriter, r *http.Request, inst *ChartInstance) {
	q := r.URL.Query()

	// leave=1: clear hover state, return empty tooltip + clean SVG OOB.
	if q.Get("leave") != "" {
		inst.clearHover()
		if inst.Kind == KindPie {
			inst.setActive("")
		}
		svg, err := renderFull(inst)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeHTML(w, oobWrapper(inst.ID, "", svg))
		return
	}

	switch inst.Kind {
	case KindBar:
		key := q.Get("bar")
		if key == "" {
			http.Error(w, "missing bar", http.StatusBadRequest)
			return
		}
		inst.setHovered(key)
		html, ok := barHoverTooltip(inst, key)
		if !ok {
			http.NotFound(w, r)
			return
		}
		svg, err := renderFull(inst)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeHTML(w, oobWrapper(inst.ID, html, svg))
	case KindPie:
		arcID := q.Get("arc")
		if arcID == "" {
			http.Error(w, "missing arc", http.StatusBadRequest)
			return
		}
		html, ok := pieHoverTooltip(inst, arcID, true)
		if !ok {
			http.NotFound(w, r)
			return
		}
		svg, err := renderFull(inst)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeHTML(w, oobWrapper(inst.ID, html, svg))
	case KindHeatmap:
		cellID := q.Get("cell")
		if cellID == "" {
			http.Error(w, "missing cell", http.StatusBadRequest)
			return
		}
		inst.setHovered(cellID)
		html, ok := heatmapHoverTooltip(inst, cellID)
		if !ok {
			http.NotFound(w, r)
			return
		}
		svg, err := renderFull(inst)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeHTML(w, oobWrapper(inst.ID, html, svg))
	case KindLine:
		// Line mesh hover: cursor x/y passed via hx-vals JS.
		x, y := q.Get("x"), q.Get("y")
		if x == "" || y == "" {
			http.Error(w, "missing cursor x/y", http.StatusBadRequest)
			return
		}
		xf, err := strconv.ParseFloat(x, 64)
		if err != nil {
			http.Error(w, "invalid x", http.StatusBadRequest)
			return
		}
		yf, err := strconv.ParseFloat(y, 64)
		if err != nil {
			http.Error(w, "invalid y", http.StatusBadRequest)
			return
		}
		html, ok := lineMeshHover(inst, xf, yf)
		if !ok {
			http.NotFound(w, r)
			return
		}
		svg, err := renderFull(inst)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeHTML(w, oobWrapper(inst.ID, html, svg))
	default:
		http.Error(w, "hover not supported for this chart kind", http.StatusBadRequest)
	}
}

// handleSlice produces an HTML table-tooltip fragment for a line slice, plus an
// OOB SVG swap with the crosshair rendered at the slice position. The `axis`
// param is "x" or "y"; `x` (or `y`) is the slice coordinate (float).
func (h *Handler) handleSlice(w http.ResponseWriter, r *http.Request, inst *ChartInstance) {
	if inst.Kind != KindLine {
		http.Error(w, "slice not supported for this chart kind", http.StatusBadRequest)
		return
	}
	q := r.URL.Query()
	axis := q.Get("axis")
	if axis == "" {
		axis = "x"
	}
	coordStr := q.Get("x")
	if axis == "y" {
		coordStr = q.Get("y")
	}
	coord, err := strconv.ParseFloat(coordStr, 64)
	if err != nil {
		http.Error(w, "invalid coordinate", http.StatusBadRequest)
		return
	}
	// Set hover coords so the crosshair renders at the slice position.
	if axis == "y" {
		inst.setHoverXY(0, coord)
	} else {
		inst.setHoverXY(coord, 0)
	}
	html, ok := lineSliceTooltip(inst, axis, coord)
	if !ok {
		http.NotFound(w, r)
		return
	}
	svg, err := renderFull(inst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeHTML(w, oobWrapper(inst.ID, html, svg))
}

// writeHTML writes an HTML fragment with the right content type.
func writeHTML(w http.ResponseWriter, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(html))
}

// oobWrapper bundles the tooltip HTML (swapped into #tooltip-<id>) with an
// out-of-band swap of the full SVG (re-rendered with hover state) targeting
// #chart-<id>. htmx processes the primary swap (tooltip) and the OOB swap
// (chart) from a single response.
func oobWrapper(id, tooltipHTML, svg string) string {
	return tooltipHTML +
		`<div hx-swap-oob="innerHTML:#chart-` + id + `">` + svg + `</div>`
}

// barHoverTooltip finds the bar with the given Key in the instance's computed
// bars and renders a BasicTooltip HTML fragment. Returns ("", false) if no
// bar matches.
func barHoverTooltip(inst *ChartInstance, key string) (string, bool) {
	props := inst.Props.(bar.BarProps)
	applyBarState(&props, inst.ID, inst.State())
	dims := core.UseDimensions(props.Width, props.Height, props.Margin)
	props.Width = dims.InnerWidth
	props.Height = dims.InnerHeight
	result := bar.UseBar(props)
	for _, b := range result.BarsWithValue {
		if b.Key == key {
			html, err := renderComponent(tooltip.BasicTooltip(tooltip.BasicTooltipProps{
				ID:             result.GetTooltipLabel(b.Data),
				FormattedValue: b.Data.FormattedValue,
				Color:          b.Color,
				EnableChip:     true,
			}))
			if err != nil {
				return "", false
			}
			return html, true
		}
	}
	return "", false
}

// pieHoverTooltip finds the arc with the given ID in the instance's computed
// data and renders a BasicTooltip HTML fragment. When setActive is true, it
// also updates the instance's ActiveID state (so the next full render pops
// the arc). Returns ("", false) if no arc matches.
func pieHoverTooltip(inst *ChartInstance, arcID string, setActive bool) (string, bool) {
	props := inst.Props.(pie.PieProps)
	st := inst.State()
	applyPieState(&props, inst.ID, st)
	result := pie.UsePie(props)
	for _, d := range result.DataWithArc {
		if d.ID == arcID {
			if setActive {
				// Toggle active off when hovering the same arc again is the
				// nivo mouseleave behaviour; v1 sets active on enter and the
				// demo clears it via a separate mouseleave htmx call. For the
				// simple case we just set it.
				inst.setActive(arcID)
			}
			html, err := renderComponent(tooltip.BasicTooltip(tooltip.BasicTooltipProps{
				ID:             d.Label,
				FormattedValue: d.FormattedValue,
				Color:          d.Color,
				EnableChip:     true,
			}))
			if err != nil {
				return "", false
			}
			return html, true
		}
	}
	return "", false
}

// heatmapHoverTooltip finds the cell with the given id in the instance's
// computed cells and renders a BasicTooltip HTML fragment ("serie - x" label,
// the formatted value, and the cell color). Returns ("", false) if no cell
// matches or the cell has no value.
func heatmapHoverTooltip(inst *ChartInstance, cellID string) (string, bool) {
	props := inst.Props.(heatmap.HeatMapProps)
	applyHeatmapState(&props, inst.ID, inst.State())
	dims := core.UseDimensions(props.Width, props.Height, props.Margin)
	props.Width = dims.InnerWidth
	props.Height = dims.InnerHeight
	result := heatmap.UseHeatMap(props)
	for _, c := range result.Cells {
		if c.ID == cellID && c.Value != nil {
			html, err := renderComponent(tooltip.BasicTooltip(tooltip.BasicTooltipProps{
				ID:             c.SerieID + " - " + c.X,
				FormattedValue: c.FormattedValue,
				Color:          c.Color,
				EnableChip:     true,
			}))
			if err != nil {
				return "", false
			}
			return html, true
		}
	}
	return "", false
}

// lineSliceTooltip finds the slice at the given axis coordinate in the
// instance's computed slices and renders a TableTooltip HTML fragment with
// one row per series point in the slice. Returns ("", false) if no slice
// matches.
func lineSliceTooltip(inst *ChartInstance, axis string, coord float64) (string, bool) {
	props := inst.Props.(line.LineProps)
	applyLineState(&props, inst.ID, inst.State())
	dims := core.UseDimensions(props.Width, props.Height, props.Margin)
	props.Width = dims.InnerWidth
	props.Height = dims.InnerHeight
	result := line.UseLine(props)
	slices := result.Slices
	if len(slices) == 0 {
		return "", false
	}
	for _, s := range slices {
		var sCoord float64
		if axis == "y" {
			sCoord = s.Y
		} else {
			sCoord = s.X
		}
		if sameFloat(sCoord, coord) {
			rows := make([]tooltip.TableTooltipRow, 0, len(s.Points))
			for _, p := range s.Points {
				rows = append(rows, tooltip.TableTooltipRow{
					Label: p.SeriesID,
					Value: p.Data.YFormatted,
					Color: p.SeriesColor,
				})
			}
			html, err := renderComponent(tooltip.TableTooltip(tooltip.TableTooltipProps{
				Rows: rows,
			}))
			if err != nil {
				return "", false
			}
			return html, true
		}
	}
	return "", false
}

// lineMeshHover finds the nearest point to (x, y) in chart units, sets the
// hover coordinates (for crosshair rendering), and returns a BasicTooltip HTML
// fragment. Returns ("", false) if no points exist.
func lineMeshHover(inst *ChartInstance, x, y float64) (string, bool) {
	props := inst.Props.(line.LineProps)
	applyLineState(&props, inst.ID, inst.State())
	dims := core.UseDimensions(props.Width, props.Height, props.Margin)
	props.Width = dims.InnerWidth
	props.Height = dims.InnerHeight
	result := line.UseLine(props)
	if len(result.Points) == 0 {
		return "", false
	}
	// Find nearest point by Euclidean distance.
	best := result.Points[0]
	bestDist := distSq(best.X, best.Y, x, y)
	for _, p := range result.Points[1:] {
		d := distSq(p.X, p.Y, x, y)
		if d < bestDist {
			bestDist = d
			best = p
		}
	}
	// Snap the crosshair to the resolved point (nivo's mesh behaviour) rather
	// than leaving it at the raw cursor position, so the crosshair and the
	// tooltip describe the same location.
	inst.setHoverXY(best.X, best.Y)
	html, err := renderComponent(tooltip.BasicTooltip(tooltip.BasicTooltipProps{
		ID:             best.SeriesID,
		FormattedValue: best.Data.YFormatted,
		Color:          best.SeriesColor,
		EnableChip:     true,
	}))
	if err != nil {
		return "", false
	}
	return html, true
}

func distSq(x1, y1, x2, y2 float64) float64 {
	dx := x1 - x2
	dy := y1 - y2
	return dx*dx + dy*dy
}

// sameFloat compares two floats with a small tolerance (the slice coord is
// formatted with %g in the hx-get URL, so round-trip may differ in the last
// digit).
func sameFloat(a, b float64) bool {
	return absFloat(a-b) < 1e-6
}

func absFloat(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
